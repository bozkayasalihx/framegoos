package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	// ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	list = []string{"pip", "pip3"}
    limit = 20
)

type Node struct {
	Prev *Node
	Next *Node
	Values []fs.DirEntry
}

type List struct {
	Head *Node
	Tail *Node
}

func NewList() *List {
	return &List{
		Head: &Node{},
		Tail: &Node{},
	}
}

func (L *List) Insert(vals []fs.DirEntry) {
	list := &Node{
		Next: L.Head,
		Values:  vals,
	}
	if L.Head != nil {
		L.Head.Prev= list
	}
	L.Head = list

	l := L.Head
	for l.Next != nil {
		l = l.Next
	}
	L.Tail= l
}

func (L *List) Display() {
	list := L.Head
	for list != nil {
		fmt.Printf("%+v\n", &list.Values)
		list = list.Next
	}
	fmt.Println()
}

func looper(list []string) (string, error) {
	var e error
	for idx, cmd := range list {
		err := exec.Command(cmd).Run()
		e = err
		if err == nil {
			return list[idx], e
		}
	}
	return "", e
}

func commandRunner(cmds ...string) ([]byte, error) {
	return exec.Command(cmds[0], cmds[1:]...).Output()
}

func remBg(cmd string) error {
	err := exec.Command(cmd, "install", "rembg").Run()
	if err != nil {
		return fmt.Errorf("couldn't install rembg -> %v", err)
	}
	return nil
}

func (list *List) Aggrator(path string) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("couldn't read dir %v", err)
	}
	chunk := limit
    a, b := 0, 0; 
    length := len(dirs);

	for a < length {
        b = a + chunk; 
        if b > length {
            b = length;
        }
		list.Insert(dirs[a:b])
        a = b;
	}
}

func main() {
	cmd, err := looper(list)
	if err != nil {
		log.Fatalf("looper error %v", err)

	}
	err = remBg(cmd)
	if err != nil {
		log.Fatalf("couldn't install rembg %v", err)
	}
	inputPath := path.Join(os.TempDir(), "test")
	outputPath := path.Join(os.TempDir(), "result")
	os.Mkdir(inputPath, 0777)
	os.Mkdir(outputPath, 0777)

	elems := []string{"ffmpeg", "-i", "test/test.mp4", inputPath + "/%04d.png"}
	commandRunner(elems...)
	list := NewList()
	list.Aggrator(inputPath)

    c := make(chan int)
    var total int
    var behind int
    
    for list != nil {
        head := list.Head 
        total += len(head.Values);
        for _, curDir := range head.Values {
            go func( ){
                _, e := exec.Command("rembg", "i", inputPath + "/" +  curDir.Name(), outputPath + "/" + curDir.Name()).Output()
                if e != nil {
                    panic(err);
                }
                c <- 1
            }()
        }
        fmt.Println("processsing...") 
        behind += <-c;
        if behind == total {
            head = head.Next; 
        }
    }


	// err = ffmpeg.Input(fmt.Sprintf("%s/%s", resultPath, "%04d.png"), ffmpeg.KwArgs{"r": "60"}).Output("./test/output.mp4", ffmpeg.KwArgs{"vcodec": "libx264", "crf": 15, "pix_fmt": "yuv420p"}).OverWriteOutput().ErrorToStdOut().Run()
	// if err != nil {
	//     log.Fatalf("couldn't ffmpeg build cmd %v" ,err);
	// }
}
