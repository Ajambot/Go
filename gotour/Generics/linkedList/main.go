package main

import "fmt"

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func add[T comparable](h *List[T], v T) *List[T] {
	if h == nil {
		n := List[T]{next: nil, val: v}
		return &n
	}

	p := h

	for p.next != nil {
		p = p.next
	}
	n := List[T]{next: nil, val: v}
	p.next = &n
	return h
}

func Walk[T comparable](h *List[T]) {
	for h != nil {
		fmt.Println(h.val)
		h = h.next
	}

}

func main() {
	h := &List[string]{val: "hello"}
	h = add(h, "world")
	Walk(h)
	h2 := &List[int]{val: 1}
	h2 = add(h2, 2)
	h2 = add(h2, 3)
	Walk(h2)
}
