package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
)

type Env string

const (
	Dev Env = "dev"
)

func init() {
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../..")
	err := os.Chdir(root)
	if err != nil {
		log.Fatalln(err)
	}

	env := os.Getenv("env")
	if env == "" {
		os.Setenv("env", string(Dev))
		env = string(Dev)
	}

	propertiesPath := fmt.Sprintf("%s/resources/config", root)
	err = archaius.Init(archaius.WithRequiredFiles([]string{
		fmt.Sprintf("%s/application.yml", propertiesPath),
		fmt.Sprintf("%s/%s/application.yml", propertiesPath, env),
	}))
	if err != nil {
		openlog.Error("Error:" + err.Error())
	}

	log.Println(archaius.GetString("app.env", ""))
}

// String function will try config key from config files,
// if the key is not found so will try
// fallback to environment variables
// String don't produce error.
func String(key string) string {
	value := archaius.GetString(key, "")
	if value == "" {
		return os.Getenv(strings.ToUpper(strings.ReplaceAll(key, ".", "_")))
	}
	return value
}

// Int  function will try config key from config files,
// if the key is not found so will try
// fallback to environment variables
// Int don't produce error.
func Int(key string) int {
	value := archaius.GetInt(key, 0)
	if value == 0 {
		env, err := strconv.Atoi(os.Getenv(strings.ToUpper(strings.ReplaceAll(key, ".", "_"))))
		if err != nil {
			return 0
		}
		return env
	}
	return value
}
