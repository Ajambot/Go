#!/bin/bash
sudo docker build --progress=plain -t httpserver -f loadbalancer/cmd/httpserver/Dockerfile .
