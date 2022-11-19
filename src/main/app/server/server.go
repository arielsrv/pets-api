package server

import (
	"log"
	"net/http"

	"github.com/src/main/app/config"

	"reflect"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html"
)

var routes = make(map[string]func(ctx *fiber.Ctx) error)

type Handler = func(ctx *fiber.Ctx) error

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
			ErrorHandler:          ErrorHandler,
			Views:                 html.New(config.String("views.folder"), ".html"),
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
			Format:     "${time} ${locals:requestid} ${status} - ${method} ${path}\n",
			TimeFormat: "2006/01/02 15:04:05",
		}))
	}

	if app.appConfig.Swagger {
		log.Println("Swagger enabled")
		app.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
	}

	return app
}

type AppConfig struct {
	Recovery  bool
	Swagger   bool
	RequestID bool
	Logger    bool
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
		registerHandler(handler)
	}
}

func registerHandler(action func(ctx *fiber.Ctx) error) {
	name := getFunctionName(action)
	routes[name] = action
}

func Use(action func(ctx *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	name := getFunctionName(action)
	return routes[name]
}

func getFunctionName(action func(ctx *fiber.Ctx) error) string {
	return runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
}

func SendString(ctx *fiber.Ctx, body string) error {
	if body == "" {
		ctx.Status(http.StatusNotFound)
	}

	return ctx.SendString(body)
}

func SendOk(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(data)
}

func SendCreated(ctx *fiber.Ctx, data interface{}) error {
	ctx.Status(http.StatusCreated)
	return ctx.JSON(data)
}
