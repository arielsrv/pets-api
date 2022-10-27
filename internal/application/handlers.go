package application

import (
	"fmt"
	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/internal/clients"
	"github.com/internal/handlers"
	"github.com/internal/services"
	"net/http"
	"os"
	"time"
)

func GetPingHandler() *handlers.PingHandler {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	return pingHandler
}

func GetRepositoriesHandler() *handlers.RepositoriesHandler {
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

	return repositoriesHandler
}
