package main

import (
	"fmt"
	"log"
	"net/http"

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

	pingHandler := application.GetPingHandler()
	repositoriesHandler := application.GetRepositoriesHandler()

	app.Add(http.MethodGet, "/ping", pingHandler.Ping)
	app.Add(http.MethodGet, "/repositories/groups", repositoriesHandler.GetGroups)
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateRepository)

	host := server.GetAppConfig().Server.Host
	port := server.GetAppConfig().Server.Port

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)
	log.Fatal(app.Start(address))
}
