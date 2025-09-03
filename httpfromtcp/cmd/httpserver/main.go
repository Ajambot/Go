package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
		h := response.GetDefaultHeaders(len(intSer))
		if req.RequestLine.RequestTarget == "/yourproblem" {
			err := w.WriteStatusLine(400)
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
			_, err = w.WriteBody([]byte(badReq))
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			return
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			err := w.WriteStatusLine(500)
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
			_, err = w.WriteBody([]byte(intSer))
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			return
		}
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			h.Remove("Content-Length")
			h.Set("Transfer-Encoding", "chunked")
			h.Set("Trailer", "X-Content-SHA256")
			h.Set("Trailer", "X-Content-Length")
			err := w.WriteStatusLine(200)
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			err = w.WriteHeaders(h)
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			resp, err := http.Get("http://httpbin.org/" + target)
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			body := new(bytes.Buffer)
			for {
				buf := make([]byte, 1024)
				n, err := resp.Body.Read(buf)
				body.Write(buf[:n])
				fmt.Println("Read:", n, "bytes from httpbin")
				if n > 0 {
					_, err2 := w.WriteChunkedBody(buf[:n])
					if err2 != nil {
						log.Fatal("Error", err)
						return
					}
				}
				if err == io.EOF {
					_, err = w.WriteChunkedBodyDone()
					if err != nil {
						log.Fatal("Error", err)
						return
					}
					hash := sha256.Sum256(body.Bytes())
					t := headers.NewHeaders()
					t.Set("X-Content-SHA256", hex.EncodeToString(hash[:]))
					t.Set("X-Content-Length", fmt.Sprint(body.Len()))
					err = w.WriteTrailers(t)
					if err != nil {
						log.Fatal("Error", err)
						return
					}
					break
				}
			}
			return
		}

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
