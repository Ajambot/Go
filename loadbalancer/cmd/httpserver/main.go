package main

import (
	"fmt"
	"httpfromtcp/pkg/request"
	"httpfromtcp/pkg/response"
	"httpfromtcp/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const ok = `<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>`

const port = 6967

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		fmt.Println("Received request")
		h := response.GetDefaultHeaders(len(ok))
		err := w.WriteStatusLine(200)
		if err != nil {
			log.Fatal("Error", err)
			return
		}
		h.Overwrite("Content-Type", "text/html")
		err = w.WriteHeaders(h)
		if err != nil {
			log.Fatal("Error", err)
			return
		}
		_, err = w.WriteBody([]byte(ok))
		if err != nil {
			log.Fatal("Error", err)
			return
		}
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
