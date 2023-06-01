package routes

import (
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App) {
	app.Get("/", MainPageHandler.Get)

	app.Get("/registration", RegistrationHandler.Get)
	app.Post("/registration", RegistrationHandler.Registrate)

	app.Get("/login", LoginHandler.Get)
	app.Post("/login", LoginHandler.Login)
}

func AuthorizedRoutes(app *fiber.App) {
	app.Get("/profile", ProfileHandler.Get)
	app.Patch("/profile", ProfileHandler.Update)
	app.Delete("/profile", ProfileHandler.Delete)

	app.Get("/homeworks", HomeworkHandler.GetList)
	app.Post("/homeworks", HomeworkHandler.Create)

	app.Get("/homeworks/:id", HomeworkHandler.Get)
	app.Patch("/homeworks/:id", HomeworkHandler.Update)
	app.Delete("/homeworks/:id", HomeworkHandler.Delete)
}
