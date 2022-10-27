package application

import (
	"net/http"

	"github.com/internal/handlers"
	"github.com/internal/server"
)

func RegisterRoutes(app *server.App) {
	app.Add(http.MethodGet, "/ping", server.Use(handlers.PingHandler{}.Ping))
	app.Add(http.MethodGet, "/repositories/groups", server.Use(handlers.RepositoriesHandler{}.GetGroups))
	app.Add(http.MethodPost, "/repositories", server.Use(handlers.RepositoriesHandler{}.CreateRepository))
}
