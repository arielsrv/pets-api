package app

import (
	"fmt"
	"log"

	"github.com/src/main/app/config"
	"github.com/src/main/app/server"
)

func Run() error {
	app := server.New()
	app.Handlers(Handlers())
	app.Routing(Routes())

	host := config.String("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := config.String("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)

	return app.Start(address)
}
