package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/src/main/app/config/env"

	"github.com/src/main/app/infrastructure/database"
	"github.com/src/main/app/infrastructure/secrets"

	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/config"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"

	"github.com/arielsrv/golang-toolkit/rest"
	_ "github.com/go-sql-driver/mysql"
)

var secretStore = ProvideSecretStore()
var dbClient = ProvideDbClient()

func RegisterHandlers() {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)
	server.RegisterHandler(pingHandler)

	gitLabClient := gitlab.NewGitLabClient(&rest.RequestBuilder{
		BaseURL: config.String("gitlab.client.baseurl"),
		Headers: http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", getSecretValue("SECRETS_STORE_PETS-API_GITLAB_TOKEN_KEY_NAME"))},
		},
		Timeout:        time.Millisecond * time.Duration(config.Int("gitlab.client.pool.timeout")),
		ConnectTimeout: time.Millisecond * time.Duration(config.Int("gitlab.client.socket.connection-timeout")),
		CustomPool: &rest.CustomPool{
			MaxIdleConnsPerHost: config.Int("gitlab.client.pool.size"),
		},
	})

	appService := services.NewAppService(gitLabClient, dbClient)
	appHandler := handlers.NewAppHandler(appService)
	server.RegisterHandler(appHandler)

	secretService := services.NewSecretService(dbClient, appService)
	secretHandler := handlers.NewSecretHandler(appService, secretService)
	server.RegisterHandler(secretHandler)

	snippetService := services.NewSnippetService(secretService)
	snippetHandler := handlers.NewSnippetHandler(snippetService)
	server.RegisterHandler(snippetHandler)
}

func getSecretValue(key string) string {
	secret := secretStore.GetSecret(key)
	if secret.Err != nil {
		log.Fatalln(secret.Err)
	}
	return secret.Value
}

func ProvideSecretStore() secrets.ISecretStore {
	if !env.IsDev() {
		return secrets.NewSecretStore()
	} else {
		return secrets.NewLocalSecretStore()
	}
}

func ProvideDbClient() database.IDbClient {
	connectionString := getSecretValue("SECRETS_STORE_PETS-API_PROD_CONNECTION_STRING_KEY_NAME")
	mySqlDbClient := database.NewMySQLClient(connectionString)
	return database.NewDbClient(mySqlDbClient)
}
