package main

import (
	"time"

	"github.com/MikhailR1337/task-sync-x/application/server"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
)

func main() {
	initializers.InitConfig()
	initializers.InitValidator()
	err := initializers.InitDb()
	if err != nil {
		logrus.Fatal(err)
	}
	engine := html.New("public/template", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ViewsLayout:  "index",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Static("/", "./public")
	server.Init(app)

	port := ":3000"
	logrus.Fatal(app.Listen(port))
}
