package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/go-chassis/go-archaius"
	"github.com/src/main/app/config/env"
)

const (
	File = "config.yml"
)

func init() {
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../../..")
	err := os.Chdir(root)
	if err != nil {
		log.Fatalln(err)
	}

	var compositeConfig []string
	propertiesPath := fmt.Sprintf("%s/resources/config", root)
	environment, scope := env.GetEnv(), env.GetScope()

	scopeConfig := fmt.Sprintf("%s/%s/%s.%s", propertiesPath, environment, scope, File)
	if _, err = os.Stat(scopeConfig); err == nil {
		compositeConfig = append(compositeConfig, scopeConfig)
	}

	envConfig := fmt.Sprintf("%s/%s/%s", propertiesPath, environment, File)
	if _, err = os.Stat(envConfig); err == nil {
		compositeConfig = append(compositeConfig, envConfig)
	}

	sharedConfig := fmt.Sprintf("%s/%s", propertiesPath, File)
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

	log.Printf("INFO: ENV: %s, SCOPE: %s", environment, scope)
}

func String(key string) string {
	return archaius.GetString(key, "")
}

func Int(key string) int {
	return archaius.GetInt(key, 0)
}

func TryInt(key string, defaultValue int) int {
	return archaius.GetInt(key, defaultValue)
}
