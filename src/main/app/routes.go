package app

import (
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/handlers/apps"
	"github.com/src/main/app/handlers/secrets"
	"github.com/src/main/app/handlers/snippets"
	"github.com/src/main/app/server"
)

func RegisterRoutes() {
	server.Register(server.GET, "/ping", server.Use[handlers.PingHandler]().Ping)
	server.Register(server.GET, "/apps/groups", server.Use[apps.AppHandler]().GetGroups)
	server.Register(server.POST, "/apps", server.Use[apps.AppHandler]().CreateApp)
	server.Register(server.GET, "/apps/types", server.Use[apps.AppHandler]().GetAppTypes)
	server.Register(server.GET, "/apps/search", server.Use[apps.AppHandler]().GetApp)
	server.Register(server.POST, "/apps/:appId/secrets", server.Use[secrets.SecretHandler]().CreateSecret)
	server.Register(server.GET, "/apps/:appId/secrets/:secretId/snippets", server.Use[snippets.SnippetHandler]().GetSnippet)
}
