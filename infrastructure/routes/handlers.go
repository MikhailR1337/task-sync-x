package routes

import (
	"errors"
	"strconv"

	"github.com/MikhailR1337/task-sync-x/domain/models"
	"github.com/MikhailR1337/task-sync-x/infrastructure/utilities"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	MainPageHandler     = &mainPageHandler{&initializers.DB}
	RegistrationHandler = &registrationHandler{&initializers.DB}
	LoginHandler        = &loginHandler{&initializers.DB}
	ProfileHandler      = &profileHandler{&initializers.DB}
	HomeworkHandler     = &homeworksHandler{&initializers.DB}
)

const (
	UniqueViolationErr = "23505"
)

type (
	mainPageHandler struct {
		storage *initializers.PgDb
	}
	registrationHandler struct {
		storage *initializers.PgDb
	}
	loginHandler struct {
		storage *initializers.PgDb
	}
	profileHandler struct {
		storage *initializers.PgDb
	}
	homeworksHandler struct {
		storage *initializers.PgDb
	}
)

func (h *mainPageHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *registrationHandler) Get(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

type RegistrateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegistrateResponse struct {
	Id string `json:"id"`
}

func (h *registrationHandler) Registrate(c *fiber.Ctx) error {
	req := RegistrateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	password, err := utilities.HashPassword(req.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	newUser := models.User{
		Name:     req.Name,
		Email:    &req.Email,
		Password: password,
		Role:     req.Password,
	}

	if err := h.storage.Create(&newUser).Error; err != nil {
		if IsDuplicatedKeyError(err) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	id := strconv.FormatUint(uint64(newUser.Id), 10)

	return c.Status(fiber.StatusCreated).JSON(RegistrateResponse{Id: id})
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

func IsDuplicatedKeyError(err error) bool {
	var perr *pgconn.PgError
	if errors.As(err, &perr) {
		return perr.Code == UniqueViolationErr
	}
	return false
}
