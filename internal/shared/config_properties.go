package shared

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
