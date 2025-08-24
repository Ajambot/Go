package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	add, err := net.ResolveUDPAddr("udp", ":42069")

	if err != nil {
		log.Fatal("error ", err)
	}

	conn, err := net.DialUDP("udp", nil, add)

	if err != nil {
		log.Fatal("error ", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error ", err)
		}

		n, err := conn.Write([]byte(str))

		if err != nil {
			log.Fatal("error ", err, "wrote ", n, "bytes")
		}
	}
}
