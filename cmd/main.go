package main

import (
	"fmt"
	"log"

	_ "github.com/docs"
	"github.com/internal/application"
	"github.com/internal/server"
)

// @title       Golang Template API
// @version     1.0
// @description This is a sample swagger for Golang Template API
// @BasePath    /.
func main() {
	app := server.New()

	application.RegisterHandlers()
	application.RegisterRoutes(app)

	host := server.GetAppConfig().Server.Host
	port := server.GetAppConfig().Server.Port

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)
	log.Fatal(app.Start(address))
}
