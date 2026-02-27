package loadbalancer

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/pkg/request"
	"httpfromtcp/pkg/response"
	httpserver "httpfromtcp/pkg/server"
	"io"
	"loadbalancer/pkg/algorithm"
	"loadbalancer/pkg/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Scheduler interface {
	Next([]server.Server) (int, error)
}

type LoadBalancer struct {
	Servers   []server.Server
	scheduler Scheduler
}

func MakeLB(algo string) (*LoadBalancer, error) {
	var scheduler Scheduler
	switch algo {
	case "rr":
		scheduler = algorithm.NewRoundRobin()
	case "wrr":
		scheduler = algorithm.NewWeightedRoundRobin()
	default:
		return nil, errors.New("Error: selected scheduling algorithm is not valid")
	}

	return &LoadBalancer{make([]server.Server, 0), scheduler}, nil
}

func (lb *LoadBalancer) Register(server server.Server) {
	lb.Servers = append(lb.Servers, server)
}

func (lb *LoadBalancer) Remove(id int) {
	for i, server := range lb.Servers {
		if server.Id == id {
			lb.Servers = append(lb.Servers[:i], lb.Servers[i+1:]...)
			break
		}
	}
}

func (lb *LoadBalancer) handler(w *response.Writer, req *request.Request) {
	fmt.Println("Received a request")
	nextServer, err := lb.scheduler.Next(lb.Servers)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	targetUrl := fmt.Sprintf("http://localhost:%d", lb.Servers[nextServer].Id) // different port for VMs
	newReq, err := http.NewRequest(req.RequestLine.Method, targetUrl+req.RequestLine.RequestTarget, bytes.NewReader(req.Body))
	if err != nil {
		log.Fatal("Error", err)
		return
	}

	ch := req.Headers.Range()
	for h := range ch {
		newReq.Header.Add(h[0], h[1])
	}

	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	fmt.Println(string(body))

	newH := response.GetDefaultHeaders(len(body))
	for k, v := range resp.Header {
		newH.Overwrite(k, strings.Join(v, ""))
	}
	err = w.WriteStatusLine(response.StatusCode((resp.StatusCode)))
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	err = w.WriteHeaders(newH)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
}

func (lb *LoadBalancer) Start() {
	port := 42069
	lbServer, err := httpserver.Serve(port, lb.handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer lbServer.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
