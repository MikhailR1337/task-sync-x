package middlewares

import (
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
)

func AddCommonMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))
}

func AddJwtMiddleware(group fiber.Router, cfg *initializers.Config) {
	group.Use(jwtware.New(jwtware.Config{
		SigningKey: cfg.JwtSecretKey,
		ContextKey: cfg.ContextKeyUser,
	}))
}
