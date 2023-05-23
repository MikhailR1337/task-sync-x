package server

import (
	"github.com/MikhailR1337/task-sync-x/infrastructure/middlewares"
	"github.com/MikhailR1337/task-sync-x/infrastructure/routes"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, cfg *initializers.Config) {
	middlewares.AddCommonMiddleware(app)
	routes.PublicRoutes(app)

	middlewares.AddJwtMiddleware(app, cfg)
	routes.AuthorizedRoutes(app)
}
