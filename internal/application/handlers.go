package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/internal/clients"
	"github.com/internal/handlers"
	"github.com/internal/server"
	"github.com/internal/services"
)

func GetPingHandler() *handlers.PingHandler {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	return pingHandler
}

func GetRepositoriesHandler() *handlers.RepositoriesHandler {
	gitLabRequestBuilder := &rest.RequestBuilder{
		BaseURL: server.GetAppConfig().GitLabClient.BaseURL,
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", server.GetAppConfig().Credentials.GitlabToken)},
		},
		Timeout:        time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.Timeout),
		ConnectTimeout: time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.ConnectionTimeout),
	}

	gitLabClient := clients.NewGitLabClient(gitLabRequestBuilder)
	repositoriesService := services.NewRepositoriesService(gitLabClient)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)

	return repositoriesHandler
}
