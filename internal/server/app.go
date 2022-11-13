package server

import (
	"log"
	"net/http"

	"github.com/internal/services"

	"reflect"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html"
	"github.com/internal/shared"
)

var routes = make(map[string]func(ctx *fiber.Ctx) error)

type App struct {
	*fiber.App
	appConfig AppConfig
}

func (app *App) Start(addr string) error {
	return app.Listen(addr)
}

func New(appConfig ...AppConfig) *App {
	app := &App{
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          shared.ErrorHandler,
			Views:                 html.New("./views", ".html"),
		}),
		appConfig: AppConfig{
			Recovery:  true,
			Swagger:   true,
			RequestID: true,
			Logger:    true,
		},
	}

	if len(appConfig) > 0 {
		app.appConfig = appConfig[0]
	}

	if app.appConfig.Recovery {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if app.appConfig.RequestID {
		app.Use(requestid.New())
	}

	if app.appConfig.Logger {
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))
	}

	if app.appConfig.Swagger {
		log.Println("Swagger enabled")
		app.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
	}

	app.Get("/snippets", func(c *fiber.Ctx) error {
		snippetService := services.NewSnippetService()
		return c.Render("snippets/index", fiber.Map{
			"Snippets": snippetService.GetSecrets(),
		})
	})

	return app
}

type AppConfig struct {
	Recovery  bool
	Swagger   bool
	RequestID bool
	Logger    bool
}

func SendString(ctx *fiber.Ctx, body string) error {
	if body == "" {
		ctx.Status(http.StatusNotFound)
	}

	return ctx.SendString(body)
}

type Routes struct {
	routes []Route
}

type Route struct {
	Verb   string
	Path   string
	Action func(ctx *fiber.Ctx) error
}

func (r *Routes) Add(verb string, path string, action func(ctx *fiber.Ctx) error) {
	route := &Route{
		Verb:   verb,
		Path:   path,
		Action: Use(action),
	}
	r.routes = append(r.routes, *route)
}

func (app *App) Routing(f func(*Routes)) {
	r := new(Routes)
	f(r)

	for _, route := range r.routes {
		app.Add(route.Verb, route.Path, route.Action)
	}
}

func (app *App) Handlers(handlers []Handler) {
	for _, handler := range handlers {
		RegisterHandler(handler)
	}
}

func SendJSON(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(data)
}

func RegisterHandler(action func(ctx *fiber.Ctx) error) {
	name := getFunctionName(action)
	routes[name] = action
}

type Handler = func(ctx *fiber.Ctx) error

func Use(action func(ctx *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	name := getFunctionName(action)
	return routes[name]
}

func getFunctionName(action func(ctx *fiber.Ctx) error) string {
	return runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
}
