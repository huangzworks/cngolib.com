package main

import (
	"container/list"
	"fmt"
)

func main() {
    // 创建一个新的链表，并向链表里面添加几个数字
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

    // 遍历链表并打印它包含的元素
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

}
