package app

import (
	"github.com/src/main/app/handlers"
	. "github.com/src/main/app/server"
)

func RegisterRoutes() {
	Register(GET, "/ping", Use[handlers.PingHandler]().Ping)
	Register(GET, "/apps/groups", Use[handlers.AppHandler]().GetGroups)
	Register(POST, "/apps", Use[handlers.AppHandler]().CreateApp)
	Register(GET, "/apps/types", Use[handlers.AppHandler]().GetAppTypes)
	Register(GET, "/apps/search", Use[handlers.AppHandler]().GetApp)
	Register(POST, "/apps/:appId/secrets", Use[handlers.SecretHandler]().CreateSecret)
	Register(GET, "/apps/:appId/secrets/:secretId/snippets", Use[handlers.SnippetHandler]().GetSnippet)
}
