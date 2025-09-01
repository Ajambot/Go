package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
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
		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Connection has been closed.")
	}

}
