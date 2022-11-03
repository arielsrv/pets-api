package shared

import (
	"log"
	"strconv"

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

func GetProperty(key string) string {
	value, err := config.String(key)
	if err != nil {
		log.Printf("missing conf key: %s", key)
	}
	return value
}

func GetIntProperty(key string) int {
	value, err := config.String(key)
	if err != nil {
		log.Printf("missing conf key: %s", key)
		return 0
	}
	parseInt, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("failed to parse conf value: %s", key)
		return 0
	}
	return parseInt
}
