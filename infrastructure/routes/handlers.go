package routes

import "github.com/gofiber/fiber/v2"

var (
	MainPageHandler     = &mainPageHandler{}
	RegistrationHandler = &registrationHandler{}
	LoginHandler        = &loginHandler{}
)

type (
	mainPageHandler     struct{}
	registrationHandler struct{}
	loginHandler        struct{}
	// profileHandler      struct{}
	// homeworksHandler    struct{}
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
