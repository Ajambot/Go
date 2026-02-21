package main

import (
	"fmt"
	"httpfromtcp/pkg/request"
	"httpfromtcp/pkg/response"
	"httpfromtcp/pkg/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Port is not a valid integer")
		return
	}

	ok := fmt.Sprintf(`
		<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Request returned from server at port %d</p>
			</body>
		</html>
		`, port)

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

	server, err := server.Serve(int(port), handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}
	defer server.Close()
	log.Printf("Server started on port %d", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
