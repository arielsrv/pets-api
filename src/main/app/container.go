package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/config"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/infrastructure"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"

	"github.com/arielsrv/golang-toolkit/rest"
	_ "github.com/go-sql-driver/mysql"
)

var secretStore = ProvideSecretStore()

func Handlers() []server.Handler {
	connectionString := getSecretValue("SECRETS_STORE_PROD_CONNECTION_STRING_KEY_NAME")
	mySqlDbClient := infrastructure.NewMySQLClient(connectionString)
	dbClient := infrastructure.NewDbClient(mySqlDbClient)

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	gitLabClient := gitlab.NewGitLabClient(&rest.RequestBuilder{
		BaseURL: config.String("gitlab.client.baseurl"),
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", getSecretValue("SECRETS_STORE_GITLAB_TOKEN_KEY_NAME"))},
		},
		Timeout:        time.Millisecond * time.Duration(config.Int("gitlab.client.pool.timeout")),
		ConnectTimeout: time.Millisecond * time.Duration(config.Int("gitlab.client.socket.connection-timeout")),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: config.Int("gitlab.client.pool.size"),
		},
	})

	appService := services.NewAppService(gitLabClient, dbClient)
	appHandler := handlers.NewAppHandler(appService)

	secretService := services.NewSecretService(dbClient, appService)
	secretHandler := handlers.NewSecretHandler(appService, secretService)

	snippetService := services.NewSnippetService(secretService)
	snippetHandler := handlers.NewSnippetHandler(snippetService)

	var h []server.Handler

	h = append(h, pingHandler.Ping)
	h = append(h, appHandler.CreateApp)
	h = append(h, appHandler.GetGroups)
	h = append(h, appHandler.GetAppTypes)
	h = append(h, appHandler.GetApp)
	h = append(h, secretHandler.CreateSecret)
	h = append(h, snippetHandler.GetSnippet)

	return h
}

func getSecretValue(key string) string {
	secret := secretStore.GetSecret(key)
	if secret.Err != nil {
		log.Fatalln(secret.Err)
	}
	return secret.Value
}

func ProvideSecretStore() infrastructure.ISecretStore {
	if config.GetEnv() != "dev" {
		return infrastructure.NewSecretStore()
	} else {
		return infrastructure.NewLocalSecretStore()
	}
}
