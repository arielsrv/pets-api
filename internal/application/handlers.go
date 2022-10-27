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

func RegisterHandlers() {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	gitLabRb := &rest.RequestBuilder{
		BaseURL: server.GetAppConfig().GitLabClient.BaseURL,
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", server.GetAppConfig().Credentials.GitlabToken)},
		},
		Timeout:        time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.Timeout),
		ConnectTimeout: time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.ConnectionTimeout),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: 20,
		},
	}

	gitLabClient := clients.NewGitLabClient(gitLabRb)
	repositoriesService := services.NewRepositoriesService(gitLabClient)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)

	server.RegisterHandler(pingHandler.Ping)
	server.RegisterHandler(repositoriesHandler.CreateRepository)
	server.RegisterHandler(repositoriesHandler.GetGroups)
}
