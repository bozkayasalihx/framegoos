package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)


var (
    list = []string{"pip","pip3"}
)

type Node struct {
	prev *Node
	next *Node
	key  []fs.DirEntry
}

type List struct {
	head *Node
	tail *Node
}


func NewList() *List {
    return &List{
        head: &Node{},
        tail: &Node{},
    }
}

func (L *List) Insert(key []fs.DirEntry) {
	list := &Node{
		next: L.head,
		key:  key,
	}
	if L.head != nil {
		L.head.prev = list
	}
	L.head = list

	l := L.head
	for l.next != nil {
		l = l.next
	}
	L.tail = l
}

func (L *List) Display() {
	list := L.head
	for list != nil {
		fmt.Printf("%+v\n", &list.key)
		list = list.next
	}
	fmt.Println()
}




func looper(list []string) (string ,error)  {
    var e error
    for idx, cmd := range list {
        err := exec.Command(cmd).Run()
        e = err 
        if err == nil {
            return list[idx], e
        }
    } 
    return "", e;
}


func commandRunner(cmds ...string) error {
    err := exec.Command(cmds[0], cmds[1:]...).Run()
    return err;
}

func remBg(cmd string) error {
    err := exec.Command(cmd, "install", "rembg").Run()
    if err != nil {
        return fmt.Errorf("couldn't install rembg -> %v" ,err);
    }

    return nil;
}


func checkFFMPEG(){
    err := exec.Command("ffmpeg", "--version").Run()
    if err != err {
        log.Fatal("install ffmpeg")
    }

}

func (list *List) Aggrator(path string) {
    dirs, err := os.ReadDir(path);
    if err != nil {
        panic(err);
    }
    chunk := 20
    var total int

    for {
        list.Insert(dirs[total:total+chunk])
        total += chunk;        
        if total >= len(dirs) {
            break
        }
    }
}


func (list *List) mainExecutor(filepath, resultPath string){
    l := list.head
    wg := sync.WaitGroup{}
    for l != nil {
        wg.Add(len(l.key))
        for _, dir := range l.key {
            go func() {
                defer wg.Done()
                input := fmt.Sprintf("%s/%s", filepath, dir.Name()) 
                output := fmt.Sprintf("%s/%s",resultPath,dir.Name())
                cmd := []string{"rembg", "i", input, output};
                err := commandRunner(cmd...)
                if err != nil {
                    fmt.Printf("got an error %v", err);
                }
            }()
        }
        fmt.Println("done..")
        wg.Wait()
        l = l.next
    }

    fmt.Println("done");
}


func main() {
    cmd, err := looper(list);
    if err != nil {
        panic(err)
    }
    err = remBg(cmd);
    if err != nil {
        panic(err);
    }
    checkFFMPEG()
    filePath := path.Join(os.TempDir(), "test");
    resultPath := path.Join(os.TempDir(), "result");
    os.Mkdir(filePath,0777)
    os.Mkdir(resultPath, 0777);

    elems := []string{"ffmpeg", "-i", "test.mp4", filePath + "/%04d.png"}
    err = commandRunner(elems...);
    if err != nil {
        panic(err);
    }

    list := NewList();
    list.Aggrator(filePath);
    list.mainExecutor(filePath, resultPath);
}
