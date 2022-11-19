package main

import (
	"log"

	"github.com/src/main/app"
	_ "github.com/src/resources/docs"
)

// @title       Pets API
// @version     1.0
// @description Create apps, services and infrastructure.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
