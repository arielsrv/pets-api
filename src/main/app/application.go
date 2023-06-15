package app

import (
	"fmt"
	"log"

	"github.com/arielsrv/pets-api/src/main/app/config/env"

	"github.com/arielsrv/pets-api/src/main/app/config"
	"github.com/arielsrv/pets-api/src/main/app/server"
)

func Run() error {
	app := server.New()

	RegisterHandlers()
	RegisterRoutes()

	host := config.String("HOST")
	if env.IsEmpty(host) {
		host = "0.0.0.0"
	}

	port := config.String("PORT")
	if env.IsEmpty(port) {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on address: %s", address)

	return app.Start(address)
}
