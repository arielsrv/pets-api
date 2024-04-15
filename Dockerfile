FROM golang:1.22.2 AS build

ADD . /app
WORKDIR /app

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN task build

ENTRYPOINT ["./build/program"]
