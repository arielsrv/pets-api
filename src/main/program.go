package main

import (
	"log"

	"github.com/src/main/application"

	_ "github.com/docs"
)

// @title       Pets API
// @version     1.0
// @description Create apps, services and infrastructure.
func main() {
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
