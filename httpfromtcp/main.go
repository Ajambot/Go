package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error reading file.")
		return
	}

	ch := getLinesChannel(file)

	for l := range ch {
		fmt.Println("read:", l)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()
		curString := ""
		for {
			buf := make([]byte, 8)

			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			buf = buf[:n]
			slice := strings.Split(string(buf), "\n")
			if len(slice) > 1 {
				curString += slice[0]
				ch <- curString
				curString = slice[1]
			} else {
				curString += string(buf)
			}
		}
		if curString != "" {
			ch <- curString
		}
	}()
	return ch
}
