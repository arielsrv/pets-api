package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/config"
)

//nolint:nolintlint,gochecknoinits
func init() {
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../..")
	err := os.Chdir(root)
	if err != nil {
		log.Fatalln(err)
	}
	propertiesPath := fmt.Sprintf("%s/resources/config/application.properties", root)
	err = config.InitGlobalInstance("ini", propertiesPath)
	if err != nil {
		log.Println("failed to load module properties")
	}
}

// String function will try config key from config files,
// if the key is not found so will try
// fallback to environment variables
// String don't produce error.
func String(key string) string {
	value := config.DefaultString(key, "")
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
	value := config.DefaultInt(key, 0)
	if value == 0 {
		env, err := strconv.Atoi(os.Getenv(strings.ToUpper(strings.ReplaceAll(key, ".", "_"))))
		if err != nil {
			return 0
		}
		return env
	}
	return value
}
