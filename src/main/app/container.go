package app

import (
	"log"

	pingHandler "github.com/src/main/app/handlers"
	appHandler "github.com/src/main/app/handlers/apps"
	secretHandler "github.com/src/main/app/handlers/secrets"
	snippetHandler "github.com/src/main/app/handlers/snippets"
	pingService "github.com/src/main/app/services"
	appService "github.com/src/main/app/services/apps"
	secretService "github.com/src/main/app/services/secrets"
	snippetService "github.com/src/main/app/services/snippets"

	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"

	"github.com/src/main/app/infrastructure/database"
	"github.com/src/main/app/infrastructure/secrets"

	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/server"

	_ "github.com/go-sql-driver/mysql"
)

var secretStore = ProvideSecretStore()
var dbClient = ProvideDBClient()
var restClients = config.ProvideRestClients()

func RegisterHandlers() {
	newPingService := pingService.NewPingService()
	newPingHandler := pingHandler.NewPingHandler(newPingService)
	server.RegisterHandler(newPingHandler)

	gitLabClient := gitlab.NewGitLabClient(restClients.Get("gitlab"), secretStore)

	newAppService := appService.NewAppService(gitLabClient, dbClient, secretStore)
	newAppHandler := appHandler.NewAppHandler(newAppService)
	server.RegisterHandler(newAppHandler)

	newSecretService := secretService.NewSecretService(dbClient, newAppService)
	newSecretHandler := secretHandler.NewSecretHandler(newAppService, newSecretService)
	server.RegisterHandler(newSecretHandler)

	newSnippetService := snippetService.NewSnippetService(newSecretService)
	newSnippetHandler := snippetHandler.NewSnippetHandler(newSnippetService)
	server.RegisterHandler(newSnippetHandler)
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
	}

	return secrets.NewLocalSecretStore()
}

func ProvideDBClient() database.IDbClient {
	connectionString := getSecretValue("PROD_CONNECTION_STRING")
	mySQLClient := database.NewMySQLClient(connectionString)

	return database.NewDBClient(mySQLClient)
}
