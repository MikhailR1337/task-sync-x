package middlewares

import (
	"fmt"

	"github.com/MikhailR1337/task-sync-x/app/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/sirupsen/logrus"
)

func AddCommonMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))
}

func AddJwtMiddleware(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		TokenLookup: fmt.Sprintf("cookie:%s", initializers.Cfg.JwtCookieKey),
		SigningKey:  []byte(initializers.Cfg.JwtSecretKey),
		ContextKey:  initializers.Cfg.ContextKeyUser,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logrus.WithError(err)
			return c.Redirect("/login")
		},
	}))
}
