package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const badReq = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const intSer = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const ok = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

const port = 42069

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(400)
			h := response.GetDefaultHeaders(len(badReq))
			h.Overwrite("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(badReq))
			return
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			w.WriteStatusLine(500)
			h := response.GetDefaultHeaders(len(intSer))
			h.Overwrite("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(intSer))
			return
		}

		w.WriteStatusLine(200)
		h := response.GetDefaultHeaders(len(ok))
		h.Overwrite("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(ok))
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
