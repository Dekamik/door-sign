package config

import (
	"os"
)

func ReadVersion() string {
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		panic(err)
	}
	return string(versionBytes)
}
