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

	curString := ""
	for {
		buf := make([]byte, 8)

		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		buf = buf[:n]
		slice := strings.Split(string(buf), "\n")
		if len(slice) > 1 {
			curString += slice[0]
			fmt.Printf("read: %s\n", curString)
			curString = slice[1]
		} else {
			curString += string(buf)
		}
	}
	if curString != "" {
		fmt.Printf("read: %s\n", curString)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

}
