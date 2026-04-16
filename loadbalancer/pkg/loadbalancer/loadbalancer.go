package loadbalancer

import (
	"bytes"
	"encoding/json"
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
	"time"
)

type Scheduler interface {
	Next([]*server.Server) (int, error)
}

type StatsResponse struct {
	CPUUsage float64 `json:"CPUUsage"`
}

type LoadBalancer struct {
	Servers              []*server.Server
	statusCheckFrequency time.Duration
	scheduler            Scheduler
}

func MakeLB(algo string, statusCheckFrequency time.Duration) (*LoadBalancer, error) {
	var scheduler Scheduler
	switch algo {
	case "rr":
		scheduler = algorithm.NewRoundRobin()
	case "wrr":
		scheduler = algorithm.NewWeightedRoundRobin()
	case "lc":
		scheduler = algorithm.NewLeastConnections()
	case "rb":
		scheduler = algorithm.NewResourceBased()
	default:
		return nil, errors.New("Error: selected scheduling algorithm is not valid")
	}

	return &LoadBalancer{make([]*server.Server, 0), statusCheckFrequency, scheduler}, nil
}

func (lb *LoadBalancer) Register(server *server.Server) {
	lb.Servers = append(lb.Servers, server)
}

func (lb *LoadBalancer) Remove(index int) {
	lb.Servers = append(lb.Servers[:index], lb.Servers[index+1:]...)
}

func (lb *LoadBalancer) handler(w *response.Writer, req *request.Request) {
	log.Println("Received a request")
	nextServer, err := lb.scheduler.Next(lb.Servers)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	for lb.Servers[nextServer].Healthy == false {
		nextServer, err = lb.scheduler.Next(lb.Servers)
		if err != nil {
			log.Fatal("Error", err)
			return
		}
	}
	targetUrl := lb.Servers[nextServer].Url
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
	lb.Servers[nextServer].AddConnection()
	resp, err := client.Do(newReq)
	if err != nil {
		log.Fatal("Error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	lb.Servers[nextServer].RemoveConnection()
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
		log.Println("Error", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		log.Println("Error", err)
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
	go lb.healthCheckRoutine()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func statusCheck(url string) (StatsResponse, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url + "/status")
	if err != nil {
		log.Println()
		return StatsResponse{}, errors.New(fmt.Sprint(url, " unreachable. Error: ", err))
	}

	var status StatsResponse
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return StatsResponse{}, err
	}

	return status, nil
}

func (lb *LoadBalancer) checkHealth() {
	for _, r := range lb.Servers {
		status, err := statusCheck(r.Url)
		if err != nil {
			log.Println("Error checking status of", r.Url, err)
			r.SetHealthy(false)
		}
		r.Stats.CPUUsage = status.CPUUsage
	}
}

func (lb *LoadBalancer) healthCheckRoutine() {
	t := time.NewTicker(lb.statusCheckFrequency)
	for ; true; <-t.C { // Starts a health check immediately and then every 20 seconds
		log.Println("Starting health check...")
		lb.checkHealth()
		log.Println("Health check completed")
		//for i, r := range lb.Servers {
		//	fmt.Println("Server ", i, " usage: ", r.Stats.CPUUsage)
		//}
	}
}
