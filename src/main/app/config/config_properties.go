package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-chassis/go-archaius"
)

type Env string

const (
	Dev  Env = "dev"
	Prod Env = "prod"
)

func init() {
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../../..")
	err := os.Chdir(root)
	if err != nil {
		log.Fatalln(err)
	}

	env := GetEnv()
	scope := GetScope()

	propertiesPath := fmt.Sprintf("%s/resources/config", root)

	var compositeConfig []string

	scopeConfig := fmt.Sprintf("%s/%s/%s.config.yml", propertiesPath, env, scope)
	if _, err = os.Stat(scopeConfig); err == nil {
		compositeConfig = append(compositeConfig, scopeConfig)
	}

	envConfig := fmt.Sprintf("%s/%s/config.yml", propertiesPath, env)
	if _, err = os.Stat(envConfig); err == nil {
		compositeConfig = append(compositeConfig, envConfig)
	}

	sharedConfig := fmt.Sprintf("%s/config.yml", propertiesPath)
	if _, err = os.Stat(sharedConfig); err == nil {
		compositeConfig = append(compositeConfig, sharedConfig)
	}

	err = archaius.Init(
		archaius.WithENVSource(),
		archaius.WithRequiredFiles(compositeConfig),
	)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("INFO: env mode: %s", archaius.GetString("app.env", ""))
}

func GetScope() string {
	return strings.ToLower(os.Getenv("SCOPE"))
}

func GetEnv() string {
	env := os.Getenv("app.env")
	if env != "" {
		return env
	}
	scope := GetScope()
	if scope == "" {
		return string(Dev)
	}
	return string(Prod)
}

func String(key string) string {
	return archaius.GetString(key, "")
}

func Int(key string) int {
	return archaius.GetInt(key, 0)
}
