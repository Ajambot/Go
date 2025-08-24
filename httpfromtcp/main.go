package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	port := ":42069"
	listener, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatal("error", e)
	}
	fmt.Println("Listening on port:", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Println("Connection has been accepted.")
		ch := getLinesChannel(conn)

		for l := range ch {
			fmt.Println(l)
		}

		fmt.Println("Connection has been closed.")
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
