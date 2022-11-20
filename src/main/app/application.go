package app

import (
	"fmt"
	"log"

	"github.com/src/main/app/config/env"

	"github.com/src/main/app/config"
	"github.com/src/main/app/server"
)

func Run() error {
	app := server.New()

	app.SetHandlers(Handlers())
	app.SetRoutes(Routes())

	host := config.String("HOST")
	if env.IsEmpty(host) {
		host = "127.0.0.1"
	}

	port := config.String("PORT")
	if env.IsEmpty(host) {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)

	return app.Start(address)
}
