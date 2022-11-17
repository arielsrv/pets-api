package main

import (
	"fmt"
	"log"
	"os"

	"github.com/src/container"
	"github.com/src/routes"

	_ "github.com/docs"
	"github.com/src/server"
)

// @title       Pets API
// @version     1.0
// @description Create apps, services and infrastructure.
func main() {
	app := server.New()
	app.Handlers(container.Handlers())
	app.Routing(routes.Routes())

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)
	log.Fatal(app.Start(address))
}
