package routes

import (
	"net/http"

	handlers2 "github.com/src/main/handlers"
	"github.com/src/main/server"
)

func Routes() func(*server.Routes) {
	routes := func(routes *server.Routes) {
		routes.Add(http.MethodGet, "/ping", handlers2.PingHandler{}.Ping)
		routes.Add(http.MethodGet, "/apps/groups", handlers2.AppHandler{}.GetGroups)
		routes.Add(http.MethodPost, "/apps", handlers2.AppHandler{}.CreateApp)
		routes.Add(http.MethodGet, "/apps/types", handlers2.AppHandler{}.GetAppTypes)
		routes.Add(http.MethodGet, "/apps/search", handlers2.AppHandler{}.GetApp)
		routes.Add(http.MethodPost, "/apps/:appId/secrets", handlers2.SecretHandler{}.CreateSecret)
		routes.Add(http.MethodGet, "/apps/:appId/secrets/:secretId/snippets", handlers2.SnippetHandler{}.GetSnippet)
	}
	return routes
}
