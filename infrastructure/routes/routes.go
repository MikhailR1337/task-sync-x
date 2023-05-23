package routes

import (
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App, db *initializers.PgDb) {
	app.Get("/", MainPageHandler.Get)

	app.Get("/registration", RegistrationHandler.Get)
	app.Post("registration", RegistrationHandler.Registrate)

	app.Get("/login", LoginHandler.Get)
	app.Post("/login", LoginHandler.Login)
}

func AuthorizedRoutes(app *fiber.App, db *initializers.PgDb) {
	app.Get("/profile")
	app.Patch("profile")
	app.Delete("profile")

	app.Get("/homeworks")

	app.Get("/homeworks/:id")
	app.Post("/homeworks/:id")
	app.Patch("/homeworks/:id")
	app.Delete("/homeworks/:id")
}
