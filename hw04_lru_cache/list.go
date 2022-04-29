package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
	PrintList(f func(interface{}))
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.length
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) unlink(li *ListItem) {
	if li.Prev != nil {
		li.Prev.Next = li.Next
	} else {
		l.front = li.Next
	}
	if li.Next != nil {
		li.Next.Prev = li.Prev
	} else {
		l.back = li.Prev
	}
}

func (l *list) linkToFront(li *ListItem) {
	li.Next = l.front
	li.Prev = nil
	if l.front == nil {
		l.back = li
	} else {
		l.front.Prev = li
	}
	l.front = li
}

func (l *list) linkToBack(li *ListItem) {
	li.Prev = l.back
	li.Next = nil
	if l.back == nil {
		l.front = li
	} else {
		l.back.Next = li
	}
	l.back = li
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := new(ListItem)
	li.Value = v

	l.linkToFront(li)
	l.length++
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := new(ListItem)
	li.Value = v

	l.linkToBack(li)
	l.length++
	return li
}

func (l *list) Remove(li *ListItem) {
	l.unlink(li)
	l.length--
}

func (l *list) MoveToFront(li *ListItem) {
	l.unlink(li)
	l.linkToFront(li)
}

func (l *list) Clear() {
	l.front = nil
	l.back = nil
	l.length = 0

	// God save the Garbage Collector!
}

func (l *list) PrintList(f func(interface{})) {
	fmt.Printf("list length: %v: ", l.length)
	for i := l.front; i != nil; i = i.Next {
		f(i.Value)
	}
}
