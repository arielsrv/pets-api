package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/config"

	"github.com/internal/clients/gitlab"

	"github.com/internal/infrastructure"

	"github.com/arielsrv/golang-toolkit/rest"
	_ "github.com/go-sql-driver/mysql"
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
		BaseURL: config.DefaultString("gitlab.client.baseurl", ""),
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", config.DefaultString("gitlab.token", ""))},
		},
		Timeout:        time.Millisecond * time.Duration(config.DefaultInt("gitlab.client.pool.timeout", 500)),
		ConnectTimeout: time.Millisecond * time.Duration(config.DefaultInt("gitlab.client.socket.connection-timeout", 1000)),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: config.DefaultInt("gitlab.client.pool.size", 20),
		},
	}

	gitLabClient := gitlab.NewGitLabClient(gitLabRb)
	appService := services.NewAppService(gitLabClient, dataAccess)
	appHandler := handlers.NewAppHandler(appService)

	server.RegisterHandler(pingHandler.Ping)
	server.RegisterHandler(appHandler.CreateApp)
	server.RegisterHandler(appHandler.GetGroups)
	server.RegisterHandler(appHandler.GetAppTypes)
	server.RegisterHandler(appHandler.GetApp)
}
