package main

import (
	"log"

	_ "github.com/docs"
	"github.com/internal/application"
)

// @title       Pets API
// @version     1.0
// @description Create apps, services and infrastructure.
func main() {
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
