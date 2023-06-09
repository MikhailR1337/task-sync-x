package server

import (
	"github.com/MikhailR1337/task-sync-x/app/application/middlewares"
	"github.com/MikhailR1337/task-sync-x/app/application/routes"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	middlewares.AddCommonMiddleware(app)
	routes.PublicRoutes(app)

	middlewares.AddJwtMiddleware(app)
	routes.AuthorizedRoutes(app)
}
