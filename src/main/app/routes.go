package app

import (
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
	"net/http"
)

func RegisterRoutes() {
	server.RegisterRoute(http.MethodGet, "/ping", server.Action[handlers.PingHandler]().Ping)
	server.RegisterRoute(http.MethodGet, "/apps/groups", server.Action[handlers.AppHandler]().GetGroups)
	server.RegisterRoute(http.MethodPost, "/apps", server.Action[handlers.AppHandler]().CreateApp)
	server.RegisterRoute(http.MethodGet, "/apps/types", server.Action[handlers.AppHandler]().GetAppTypes)
	server.RegisterRoute(http.MethodGet, "/apps/search", server.Action[handlers.AppHandler]().GetApp)
	server.RegisterRoute(http.MethodPost, "/apps/:appId/secrets", server.Action[handlers.SecretHandler]().CreateSecret)
	server.RegisterRoute(http.MethodGet, "/apps/:appId/secrets/:secretId/snippets", server.Action[handlers.SnippetHandler]().GetSnippet)
}
