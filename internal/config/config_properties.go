package config

import (
	"log"

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
	log.Println("loaded config " + propertiesPath)
}

func String(key string) string {
	value := config.DefaultString(key, "")
	return value
}

func Int(key string) int {
	value := config.DefaultInt(key, 0)
	return value
}
