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
	"time"
)

func processCPUTime() (user, system time.Duration, err error) {
	var r syscall.Rusage
	err = syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	if err != nil {
		return
	}

	user = time.Duration(r.Utime.Sec)*time.Second +
		time.Duration(r.Utime.Usec)*time.Microsecond

	system = time.Duration(r.Stime.Sec)*time.Second +
		time.Duration(r.Stime.Usec)*time.Microsecond

	return
}

func getCPUUtilization() (float64, error) {
	u1, s1, err := processCPUTime()
	if err != nil {
		return float64(-1.0), err
	}

	t1 := time.Now()

	time.Sleep(1 * time.Second)

	u2, s2, err := processCPUTime()
	if err != nil {
		return float64(-1.0), err
	}
	t2 := time.Now()

	cpuTime := (u2 + s2) - (u1 + s1)
	wall := t2.Sub(t1)

	cpuPercent := float64(cpuTime) / float64(wall) * 100
	return cpuPercent, nil
}

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

		if req.RequestLine.RequestTarget == "/status" {
			cpuUtil, cpuErr := getCPUUtilization()
			var resp []byte
			if cpuErr != nil {
				err := w.WriteStatusLine(500)
				if err != nil {
					log.Fatal("Error", err)
					return
				}
				resp = fmt.Appendf(nil, "{ Error: %s }", cpuErr.Error())
			} else {
				err := w.WriteStatusLine(200)
				if err != nil {
					log.Fatal("Error", err)
					return
				}
				resp = fmt.Appendf(nil, `{ "CPUUsage": %f }`, cpuUtil)
			}

			h := response.GetDefaultHeaders(len(resp))
			h.Overwrite("Content-Type", "application/json")
			err = w.WriteHeaders(h)
			if err != nil {
				log.Println("Error", err)
				return
			}

			_, err = w.WriteBody(resp)
			if err != nil {
				log.Println("Error", err)
				return
			}

		} else {
			log.Println("Executing expensive calculation")
			time.Sleep(5 * time.Second)
			log.Println("Finished expensive calculation")
			h := response.GetDefaultHeaders(len(ok))
			err := w.WriteStatusLine(200)
			if err != nil {
				log.Fatal("Error", err)
				return
			}
			h.Overwrite("Content-Type", "text/html")
			err = w.WriteHeaders(h)
			if err != nil {
				log.Println("Error", err)
				return
			}
			_, err = w.WriteBody([]byte(ok))
			if err != nil {
				log.Println("Error", err)
				return
			}

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
