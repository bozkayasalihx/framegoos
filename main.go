package main

import (
	"fmt"
	"github.com/bozkayasalihx/framegoos/util"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)

var (
	list          = []string{"pip", "pip3"}
	limit         = 20
	outdir        = "./test/output.mp4"
	backgroundDir = "./test/background.mp4"
	latestDir     = "./test/latest.mp4"
)

type Node struct {
	Prev   *Node
	Next   *Node
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
		Next:   L.Head,
		Values: vals,
	}
	if L.Head != nil {
		L.Head.Prev = list
	}
	L.Head = list

	l := L.Head
	for l.Next != nil {
		l = l.Next
	}
	L.Tail = l
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
	a, b := 0, 0
	length := len(dirs)

	for a < length {
		b = a + chunk
		if b > length {
			b = length
		}
		fmt.Println(a)
		fmt.Println(b)
		list.Insert(dirs[a:b])
		a = b
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

	elems := []string{"ffmpeg", "-i", "test/test.mp4", "-vf", "fps=30", inputPath + "/%04d.png"}
	out, e := commandRunner(elems...)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(out))
	list := NewList()
	list.Aggrator(inputPath)

	wg := sync.WaitGroup{}

	fmt.Println(inputPath)

	for list != nil {
		head := list.Head
		for _, curDir := range head.Values {
			wg.Add(1)
			go func(c fs.DirEntry) {
				fmt.Println("command running")
				_, e := exec.Command("rembg", "i", inputPath+"/"+c.Name(), outputPath+"/"+c.Name()).Output()
				if e != nil {
					panic(err)
				}
				fmt.Println("command done...")
				defer wg.Done()
			}(curDir)
		}
		head = head.Next
		// break
	}
	wg.Wait()
	fmt.Println("all done")

	err = ffmpeg.Input(fmt.Sprintf("%s/%s", outputPath, "%04d.png"), ffmpeg.KwArgs{"r": "30"}).
		Output(outdir, ffmpeg.KwArgs{"vcodec": "libx264", "crf": 15, "pix_fmt": "yuv420p"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Fatalf("couldn't ffmpeg build cmd %v", err)
	}

	overlay := ffmpeg.Input(outdir).Filter("scale", ffmpeg.Args{"64:-1"})
	err = ffmpeg.Filter(
		[]*ffmpeg.Stream{
			ffmpeg.Input(backgroundDir),
			overlay,
		}, "overlay", ffmpeg.Args{"10:10"}, ffmpeg.KwArgs{"enable": "gte(t,1)"}).
		Output(latestDir).OverWriteOutput().ErrorToStdOut().Run()

	err = util.Cleanup(inputPath)
	err = util.Cleanup(outputPath)
}
