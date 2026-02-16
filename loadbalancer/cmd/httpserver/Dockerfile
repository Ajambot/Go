FROM golang:1.24-alpine

WORKDIR /app

# Copy BOTH modules first so replace works
COPY loadbalancer /app/loadbalancer
COPY httpfromtcp /app/httpfromtcp

WORKDIR /app/loadbalancer

RUN go mod download

RUN go build -o httpserver ./cmd/httpserver

ENTRYPOINT ["./httpserver"]
