package files

import (
	"log"
	"os"
)

func GetFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	return string(content)
}
