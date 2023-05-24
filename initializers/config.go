package initializers

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
	PgHost         string `env:"PG_HOST"`
	PgUser         string `env:"PG_USER"`
	PgPassword     string `env:"PG_PASSWORD"`
	PgDb           string `env:"PG_DB"`
	PgPort         string `env:"PG_PORT"`
	PgTZ           string `env:"PG_TZ"`
	JwtSecretKey   string `env:"JWT_SECRET_KEY"`
	ContextKeyUser string `env:"CONTEXT_KEY_USER"`
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
