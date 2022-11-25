package app

import (
	"log"

	"github.com/src/main/app/handlers/apps"
	secrets2 "github.com/src/main/app/handlers/secrets"
	"github.com/src/main/app/handlers/snippets"
	apps2 "github.com/src/main/app/services/apps"
	secrets3 "github.com/src/main/app/services/secrets"
	snippets2 "github.com/src/main/app/services/snippets"

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
var dbClient = ProvideDBClient()
var restClients = config.ProvideRestClients()

func RegisterHandlers() {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)
	server.RegisterHandler(pingHandler)

	gitLabClient := gitlab.NewGitLabClient(restClients.Get("gitlab"), secretStore)

	appService := apps2.NewAppService(gitLabClient, dbClient, secretStore)
	appHandler := apps.NewAppHandler(appService)
	server.RegisterHandler(appHandler)

	secretService := secrets3.NewSecretService(dbClient, appService)
	secretHandler := secrets2.NewSecretHandler(appService, secretService)
	server.RegisterHandler(secretHandler)

	snippetService := snippets2.NewSnippetService(secretService)
	snippetHandler := snippets.NewSnippetHandler(snippetService)
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

func ProvideDBClient() database.IDbClient {
	connectionString := getSecretValue("SECRETS_STORE_PETS-API_PROD_CONNECTION_STRING_KEY_NAME")
	mySQLClient := database.NewMySQLClient(connectionString)
	return database.NewDBClient(mySQLClient)
}
