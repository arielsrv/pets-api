package config

import (
	"log"
	. "os"
	. "strconv"
	. "strings"

	"github.com/beego/beego/v2/core/config"
)

//nolint:nolintlint,gochecknoinits
func init() {
	propertiesPath := "resources/application.properties"
	err := config.InitGlobalInstance("ini", propertiesPath)
	if err != nil {
		// fallback for test
		err = config.InitGlobalInstance("ini", "../../"+propertiesPath)
		if err != nil {
			log.Println("failed to load application properties")
		}
	}
}

// String function will try config key from config files,
// if the key is not found so will try
// fallback to environment variables
// String don't produce error.
func String(key string) string {
	value := config.DefaultString(key, "")
	if value == "" {
		return Getenv(ToUpper(Replace(key, ".", "_", -1)))
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
		env, err := Atoi(Getenv(ToUpper(Replace(key, ".", "_", -1))))
		if err != nil {
			return 0
		}
		return env
	}
	return value
}
