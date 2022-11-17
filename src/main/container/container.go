package container

import (
	"fmt"
	"net/http"
	"time"

	"github.com/src/main/clients/gitlab"
	"github.com/src/main/config"
	handlers2 "github.com/src/main/handlers"
	"github.com/src/main/infrastructure"
	"github.com/src/main/server"
	services2 "github.com/src/main/services"

	"github.com/arielsrv/golang-toolkit/rest"
	_ "github.com/go-sql-driver/mysql"
)

func Handlers() []server.Handler {
	dbClient := infrastructure.NewDbClient()
	dbClient.Open()

	pingService := services2.NewPingService()
	pingHandler := handlers2.NewPingHandler(pingService)
	gitLabRb := &rest.RequestBuilder{
		BaseURL: config.String("gitlab.client.baseurl"),
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", config.String("gitlab.token"))},
		},
		Timeout:        time.Millisecond * time.Duration(config.Int("gitlab.client.pool.timeout")),
		ConnectTimeout: time.Millisecond * time.Duration(config.Int("gitlab.client.socket.connection-timeout")),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: config.Int("gitlab.client.pool.size"),
		},
	}

	gitLabClient := gitlab.NewGitLabClient(gitLabRb)
	appService := services2.NewAppService(gitLabClient, dbClient)
	appHandler := handlers2.NewAppHandler(appService)

	secretService := services2.NewSecretService(dbClient, appService)
	secretHandler := handlers2.NewSecretHandler(appService, secretService)

	snippetService := services2.NewSnippetService(secretService)
	snippetHandler := handlers2.NewSnippetHandler(snippetService)

	var handlers []server.Handler
	handlers = append(handlers, pingHandler.Ping)
	handlers = append(handlers, appHandler.CreateApp)
	handlers = append(handlers, appHandler.GetGroups)
	handlers = append(handlers, appHandler.GetAppTypes)
	handlers = append(handlers, appHandler.GetApp)
	handlers = append(handlers, secretHandler.CreateSecret)
	handlers = append(handlers, snippetHandler.GetSnippet)

	return handlers
}
