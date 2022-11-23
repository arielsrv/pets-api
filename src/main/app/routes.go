package app

import (
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
)

func RegisterRoutes() {
	server.Register(server.GET, "/ping", server.Use[handlers.PingHandler]().Ping)
	server.Register(server.GET, "/apps/groups", server.Use[handlers.AppHandler]().GetGroups)
	server.Register(server.POST, "/apps", server.Use[handlers.AppHandler]().CreateApp)
	server.Register(server.GET, "/apps/types", server.Use[handlers.AppHandler]().GetAppTypes)
	server.Register(server.GET, "/apps/search", server.Use[handlers.AppHandler]().GetApp)
	server.Register(server.POST, "/apps/:appId/secrets", server.Use[handlers.SecretHandler]().CreateSecret)
	server.Register(server.GET, "/apps/:appId/secrets/:secretId/snippets", server.Use[handlers.SnippetHandler]().GetSnippet)
}
