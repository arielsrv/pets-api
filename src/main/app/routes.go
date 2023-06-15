package app

import (
	"github.com/arielsrv/pets-api/src/main/app/handlers"
	"github.com/arielsrv/pets-api/src/main/app/handlers/apps"
	"github.com/arielsrv/pets-api/src/main/app/handlers/secrets"
	"github.com/arielsrv/pets-api/src/main/app/handlers/snippets"
	"github.com/arielsrv/pets-api/src/main/app/server"
)

func RegisterRoutes() {
	server.Register(server.GET, "/ping", server.Resolve[handlers.PingHandler]().Ping)
	server.Register(server.GET, "/apps/groups", server.Resolve[apps.AppHandler]().GetGroups)
	server.Register(server.POST, "/apps", server.Resolve[apps.AppHandler]().CreateApp)
	server.Register(server.GET, "/apps/types", server.Resolve[apps.AppHandler]().GetAppTypes)
	server.Register(server.GET, "/apps/search", server.Resolve[apps.AppHandler]().GetApp)
	server.Register(server.POST, "/apps/:appId/secrets", server.Resolve[secrets.SecretHandler]().CreateSecret)
	server.Register(server.GET, "/apps/:appId/secrets/:secretId/snippets", server.Resolve[snippets.SnippetHandler]().GetSnippet)
}
