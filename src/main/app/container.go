package app

import (
	"log"

	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"

	"github.com/src/main/app/infrastructure/database"
	"github.com/src/main/app/infrastructure/secrets"

	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"

	_ "github.com/go-sql-driver/mysql"
)

var secretStore = ProvideSecretStore()
var dbClient = ProvideDbClient()
var restClients = config.ProvideRestClients()

func RegisterHandlers() {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)
	server.RegisterHandler(pingHandler)

	gitLabClient := gitlab.NewGitLabClient(restClients.Get("gitlab"), secretStore)

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
