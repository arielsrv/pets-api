package main

import (
	"log"

	"github.com/arielsrv/pets-api/src/main/app"
	_ "github.com/arielsrv/pets-api/src/resources/docs"
)

// @title Pets API.
// @description Backend for Pets Clients.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
