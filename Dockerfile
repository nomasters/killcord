FROM golang:1.10.0
WORKDIR /go/src/github.com/nomasters/killcord
RUN apt-get update && apt-get upgrade -y && apt-get install python python-pip zip unzip -y
RUN pip install awscli
COPY . .