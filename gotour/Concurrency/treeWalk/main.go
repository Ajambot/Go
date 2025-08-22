package main

import "golang.org/x/tour/tree"

import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		Walk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	calc := func(t *tree.Tree, ch chan int) {
		Walk(t, ch)
		close(ch)
	}
	ch1 := make(chan int)
	ch2 := make(chan int)
	go calc(t1, ch1)
	go calc(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if ok1 == false || ok2 == false {
			break
		}
		fmt.Println(v1, v2)
		if v1 != v2 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
}
