package main

import (
	"net/smtp"
	"time"

	"github.com/MikhailR1337/task-sync-x/mailer/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

type request struct {
	Email    string `json:"email"`
	Template string `json:"template"`
	Subject  string `json:"subject"`
}

func run() error {
	initializers.InitConfig()
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))

	app.Post("/email", func(c *fiber.Ctx) error {
		req := request{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return err
		}
		auth := smtp.PlainAuth("", initializers.Cfg.From, initializers.Cfg.AuthPassword, initializers.Cfg.Host)
		headers := "MIME-version: 1.0;\nContent-type: text/html; charset=\"UTF-8\";"
		msg := "Subject: " + req.Subject + "\n" + headers + "\n\n" + req.Template
		err := smtp.SendMail(initializers.Cfg.Server, auth, initializers.Cfg.From, []string{req.Email}, []byte(msg))
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	})

	port := ":3001"
	logrus.Fatal(app.Listen(port))
	return nil

}

func main() {
	if err := run(); err != nil {
		logrus.WithError(err)
	}
}
