package main

import (
	"time"

	"github.com/MikhailR1337/task-sync-x/infrastructure/server"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := initializers.InitConfig()

	err := initializers.InitDb(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	server.Init(app, cfg)

	port := ":3000"
	logrus.Fatal(app.Listen(port))
}
