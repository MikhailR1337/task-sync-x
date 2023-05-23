package server

import (
	"github.com/MikhailR1337/task-sync-x/infrastructure/middlewares"
	"github.com/MikhailR1337/task-sync-x/infrastructure/routes"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, cfg *initializers.Config, db *initializers.PgDb) {
	middlewares.AddCommonMiddleware(app)
	routes.PublicRoutes(app, db)

	middlewares.AddJwtMiddleware(app, cfg)
	routes.AuthorizedRoutes(app, db)
}
