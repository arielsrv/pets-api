package module

import (
	. "net/http"

	. "github.com/internal/handlers"
	"github.com/internal/server"
)

func Routes() func(*server.Routes) {
	routes := func(routes *server.Routes) {
		routes.Add(MethodGet, "/ping", PingHandler{}.Ping)
		routes.Add(MethodGet, "/apps/groups", AppHandler{}.GetGroups)
		routes.Add(MethodPost, "/apps", AppHandler{}.CreateApp)
		routes.Add(MethodGet, "/apps/types", AppHandler{}.GetAppTypes)
		routes.Add(MethodGet, "/apps/search", AppHandler{}.GetApp)
		routes.Add(MethodGet, "/apps/conf", AppHandler{}.GetAppConf)
	}
	return routes
}
