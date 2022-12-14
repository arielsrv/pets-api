package main

import (
	"log"

	"github.com/src/main/app"
	_ "github.com/src/resources/docs"
)

// @title Pets API.
// @description Backend for Pets Clients.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
