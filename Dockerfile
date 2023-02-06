# Compile stage
FROM golang:1.20-bullseye AS build

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

ADD . /app
WORKDIR /app

RUN task

EXPOSE 8080 8080

CMD ["./build/program"]
