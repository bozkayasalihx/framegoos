package link

import (
	"fmt"
	"io/fs"
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


func (L *List) Head() *Node  {
    return L.head;
}

func (L *List)Tail()  *Node{
    return L.tail;
}
