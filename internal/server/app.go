package server

import (
	"log"
	"net/http"
	"os"
	reflect "reflect"
	"runtime"
	"sync"

	"github.com/internal/shared"
	"gopkg.in/yaml.v3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
)

var appConfigMtx sync.Once
var Configuration AppConfig
var routes = make(map[string]func(ctx *fiber.Ctx) error)

type App struct {
	*fiber.App
	config Config
}

func (app *App) Start(addr string) error {
	return app.Listen(addr)
}

func New(config ...Config) *App {
	app := &App{
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          shared.ErrorHandler,
		}),
		config: Config{
			Recovery:  true,
			Swagger:   true,
			RequestID: true,
			Logger:    true,
		},
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.Recovery {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if app.config.RequestID {
		app.Use(requestid.New())
	}

	if app.config.Logger {
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))
	}

	if app.config.Swagger {
		log.Println("Swagger enabled")
		app.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
	}

	return app
}

type Config struct {
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

func SendJSON(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(data)
}

type AppConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	GitLab struct {
		Token string `yaml:"token"`
	} `yaml:"gitlab"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"database"`
	GitLabClient struct {
		BaseURL           string `yaml:"base_url"`
		Timeout           int64  `yaml:"timeout"`
		ConnectionTimeout int64  `yaml:"connection_timeout"`
	} `yaml:"gitlab_client"`
}

func GetAppConfig() AppConfig {
	appConfigMtx.Do(func() {
		log.Println("Loading config file ...")
		f, err := os.Open("config.yml")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var appConfig AppConfig
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&appConfig)
		if err != nil {
			log.Fatal(err)
		}
		Configuration = appConfig
	})

	return Configuration
}

func RegisterHandler(action func(ctx *fiber.Ctx) error) {
	name := getFuncName(action)
	routes[name] = action
}

func Use(action func(ctx *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	name := getFuncName(action)
	return routes[name]
}

func getFuncName(action func(ctx *fiber.Ctx) error) string {
	return runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
}
