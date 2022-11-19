package config

import (
	"fmt"
	"github.com/go-chassis/go-archaius"
	"log"
	"os"
	"path"
	"runtime"
)

type Env string

const (
	Dev Env = "dev"
)

func init() {
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../../..")
	err := os.Chdir(root)
	if err != nil {
		log.Fatalln(err)
	}

	env := os.Getenv("env")
	if env == "" {
		log.Println("INFO: missing env variable, setting default env ...")
		os.Setenv("env", string(Dev))
		env = string(Dev)
	}

	propertiesPath := fmt.Sprintf("%s/resources/config", root)
	err = archaius.Init(
		archaius.WithRequiredFiles([]string{
			fmt.Sprintf("%s/application.yml", propertiesPath),
			fmt.Sprintf("%s/%s/application.yml", propertiesPath, env),
		}),
		archaius.WithENVSource(),
	)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("INFO: env mode: %s", archaius.GetString("app.env", ""))
}

func String(key string) string {
	return archaius.GetString(key, "")
}

func Int(key string) int {
	return archaius.GetInt(key, 0)
}
