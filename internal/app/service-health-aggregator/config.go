package service_health_aggregator

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout_ms"` // this will be ms from config
}

func NewConfig() *Config {
	var cfg Config
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err    #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	log.Printf("%+v", cfg)
	return &cfg
}
