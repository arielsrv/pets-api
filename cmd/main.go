package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/internal/clients"

	_ "github.com/docs"
	"github.com/internal/handlers"
	"github.com/internal/server"
	"github.com/internal/services"
)

// @title       Golang Template API
// @version     1.0
// @description This is a sample swagger for Golang Template API
// @BasePath    /.
func main() {
	app := server.New()

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		gitlabToken = "glpat-BWcsGBLXz-1yQxzd3BG3"
	}
	gitLabRequestBuilder := &rest.RequestBuilder{
		BaseURL: "https://gitlab.tiendanimal.com:8088/api/v4",
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", gitlabToken)},
		},
		Timeout:        time.Millisecond * 1000,
		ConnectTimeout: time.Millisecond * 2000,
	}

	gitLabClient := clients.NewGitLabClient(gitLabRequestBuilder)
	repositoriesService := services.NewRepositoriesService(gitLabClient)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)

	app.Add(http.MethodGet, "/ping", pingHandler.Ping)
	app.Add(http.MethodGet, "/repositories/groups", repositoriesHandler.GetGroups)
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateRepository)

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)
	log.Fatal(app.Start(address))
}
