package configuration

import (
	"log"
	"os"
)

func ReadVersion() string {
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		log.Panicln(err)
	}
	return string(versionBytes)
}
