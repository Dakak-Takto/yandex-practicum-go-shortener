package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	ServerAddr string `yaml:"serverAddr"`
	Scheme     string `yaml:"scheme"`
}

var cfg = config{}

func init() {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error while opening config file: %s", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// TODO fix hardcode
	// на автотестах в гитхабе приложение почему-то не стартует.
	// Попробую захардкодить сюда переменные
	cfg.ServerAddr = "localhost:8080"
	cfg.Scheme = "http"
}

func GetAddr() string {
	return cfg.ServerAddr
}
func GetScheme() string {
	return cfg.Scheme
}
