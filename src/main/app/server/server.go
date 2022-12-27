package server

import (
	"net/http"

	"github.com/src/main/app/config"

	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html"
)

type App struct {
	*fiber.App
	appConfig AppConfig
}

var handlers = make(map[string]any)
var routes []Route

const (
	GET  = http.MethodGet
	POST = http.MethodPost
)

type Route struct {
	Verb   string
	Path   string
	Action func(ctx *fiber.Ctx) error
}

func (app *App) Start(addr string) error {
	for _, route := range routes {
		app.Add(route.Verb, route.Path, route.Action)
	}
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
			Format:     "${cyan}${time} ${locals:requestid} ${yellow}${status} - ${green}${method} ${path}\n",
			TimeFormat: "2006/01/02 15:04:05",
		}))
	}

	if app.appConfig.Swagger {
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

func RegisterHandler(handler any) {
	key := getType(handler)
	handlers[key] = handler
}

func Resolve[T any](_ ...T) *T {
	args := make([]T, 1)
	key := getType(args[0])
	return handlers[key].(*T)
}

func getType(value any) string {
	name := reflect.TypeOf(value)
	if name.Kind() == reflect.Ptr {
		name = name.Elem()
	}
	return name.String()
}

func Register(verb, path string, action func(ctx *fiber.Ctx) error) {
	route := &Route{
		Verb:   verb,
		Path:   path,
		Action: action,
	}
	routes = append(routes, *route)
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
