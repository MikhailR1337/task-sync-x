package routes

import (
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
)

var (
	MainPageHandler     = &mainPageHandler{initializers.DB}
	RegistrationHandler = &registrationHandler{initializers.DB}
	LoginHandler        = &loginHandler{initializers.DB}
	ProfileHandler      = &profileHandler{initializers.DB}
	HomeworkHandler     = &homeworksHandler{initializers.DB}
)

type (
	mainPageHandler struct {
		storage initializers.PgDb
	}
	registrationHandler struct {
		storage initializers.PgDb
	}
	loginHandler struct {
		storage initializers.PgDb
	}
	profileHandler struct {
		storage initializers.PgDb
	}
	homeworksHandler struct {
		storage initializers.PgDb
	}
)

func (h *mainPageHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *registrationHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *registrationHandler) Registrate(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *loginHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *loginHandler) Login(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *profileHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *profileHandler) Update(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *profileHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) GetList(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) Create(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) Update(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}
