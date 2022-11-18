package app

import (
	"fmt"
	"log"
	"os"

	"github.com/src/main/app/server"
)

func Run() error {
	app := server.New()
	app.Handlers(Handlers())
	app.Routing(Routes())

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

	return app.Start(address)
}
