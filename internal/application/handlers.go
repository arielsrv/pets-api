package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/internal/infrastructure"

	"github.com/arielsrv/golang-toolkit/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/internal/clients"
	"github.com/internal/handlers"
	"github.com/internal/server"
	"github.com/internal/services"
)

func RegisterHandlers() {
	dataAccess := infrastructure.NewDataAccessService()
	dataAccess.Open()

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	gitLabRb := &rest.RequestBuilder{
		BaseURL: server.GetAppConfig().GitLabClient.BaseURL,
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", server.GetAppConfig().GitLab.Token)},
		},
		Timeout:        time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.Timeout),
		ConnectTimeout: time.Millisecond * time.Duration(server.GetAppConfig().GitLabClient.ConnectionTimeout),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: 20,
		},
	}

	gitLabClient := clients.NewGitLabClient(gitLabRb)
	repositoriesService := services.NewAppService(gitLabClient, dataAccess)
	appHandler := handlers.NewAppHandler(repositoriesService)

	server.RegisterHandler(pingHandler.Ping)
	server.RegisterHandler(appHandler.CreateApp)
	server.RegisterHandler(appHandler.GetGroups)
	server.RegisterHandler(appHandler.GetAppTypes)
	server.RegisterHandler(appHandler.GetApp)
}