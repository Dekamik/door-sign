package configuration

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Weather struct {
		Lat float32 `yaml:"lat"`
		Lon float32 `yaml:"lon"`
	} `yaml:"weather"`
}

func (c *Config) ReadYamlConfig() *Config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Panicln(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Panicln(err)
	}
	return c
}
