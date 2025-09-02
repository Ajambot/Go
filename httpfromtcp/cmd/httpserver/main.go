package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
		}

		w.Write([]byte("All good, frfr\n"))
		return nil
	}

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
