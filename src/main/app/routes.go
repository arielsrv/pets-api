package app

import (
	"net/http"

	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
)

func Routes() func(*server.Routes) {
	routes := func(routes *server.Routes) {
		routes.Add(http.MethodGet, "/ping", handlers.PingHandler{}.Ping)
		routes.Add(http.MethodGet, "/apps/groups", handlers.AppHandler{}.GetGroups)
		routes.Add(http.MethodPost, "/apps", handlers.AppHandler{}.CreateApp)
		routes.Add(http.MethodGet, "/apps/types", handlers.AppHandler{}.GetAppTypes)
		routes.Add(http.MethodGet, "/apps/search", handlers.AppHandler{}.GetApp)
		routes.Add(http.MethodPost, "/apps/:appId/secrets", handlers.SecretHandler{}.CreateSecret)
		routes.Add(http.MethodGet, "/apps/:appId/secrets/:secretId/snippets", handlers.SnippetHandler{}.GetSnippet)
	}
	return routes
}
