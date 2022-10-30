package application

import (
	"net/http"

	"github.com/internal/handlers"
	"github.com/internal/server"
)

func RegisterRoutes(app *server.App) {
	app.Add(http.MethodGet, "/ping", server.Use(handlers.PingHandler{}.Ping))
	app.Add(http.MethodGet, "/apps/groups", server.Use(handlers.AppHandler{}.GetGroups))
	app.Add(http.MethodPost, "/apps", server.Use(handlers.AppHandler{}.CreateApp))
	app.Add(http.MethodGet, "/apps/types", server.Use(handlers.AppHandler{}.GetAppTypes))
	app.Add(http.MethodGet, "/search/apps", server.Use(handlers.AppHandler{}.GetApp))
}
