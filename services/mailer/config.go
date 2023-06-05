package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type Config struct {
	From         string `env:"FROM"`
	AuthPassword string `env:"AUTH_PASSWORD"`
	Host         string `env:"HOST"`
	Server       string `env:"SERVER"`
}

var (
	Cfg  Config
	once sync.Once
)

func InitConfig() {
	once.Do(func() {
		envType := os.Getenv("ENV")
		if envType == "" {
			envType = "dev"
		}
		if err := configor.New(&configor.Config{Environment: envType}).Load(&Cfg, "config.json"); err != nil {
			logrus.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(Cfg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
}
