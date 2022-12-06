# Compile stage
FROM golang:1.19.3 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ADD . /app
WORKDIR /app

# Compile the application with the optimizations turned off
# This is important for the debugger to correctly work with the binary
RUN go build -gcflags "all=-N -l" -o server src/main/program.go

# Final stage
FROM debian:buster
#
EXPOSE 8080 40000
#
WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env app /app

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/server"]