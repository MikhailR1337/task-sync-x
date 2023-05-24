package routes

import (
	"errors"
	"time"

	"github.com/MikhailR1337/task-sync-x/domain/models"
	"github.com/MikhailR1337/task-sync-x/infrastructure/utilities"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var (
	MainPageHandler     = &mainPageHandler{&initializers.DB}
	RegistrationHandler = &registrationHandler{&initializers.DB}
	LoginHandler        = &loginHandler{&initializers.DB}
	ProfileHandler      = &profileHandler{&initializers.DB}
	HomeworkHandler     = &homeworksHandler{&initializers.DB}
	Roles               = roles{
		Teacher: "teacher",
		Student: "student",
	}
)

var (
	errConflict       = errors.New("Conflict")
	errBadCredentials = errors.New("email or password is incorrect")
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
	roles = struct {
		Teacher string
		Student string
	}
)

func (h *mainPageHandler) Get(c *fiber.Ctx) error {
	return c.Render("main", fiber.Map{})
}

func (h *registrationHandler) Get(c *fiber.Ctx) error {
	return c.Render("registration", fiber.Map{})
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
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if req.Role == Roles.Teacher {
		err := h.CreateTeacher(req)
		if err != nil {
			if errors.Is(err, errConflict) {
				return c.SendStatus(fiber.StatusConflict)
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if req.Role == Roles.Student {
		err := h.CreateStudent(req)
		if err != nil {
			if errors.Is(err, errConflict) {
				return c.SendStatus(fiber.StatusConflict)
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
	return c.Status(fiber.StatusCreated).Redirect("/login")
}

func (h *registrationHandler) CreateTeacher(req RegistrateRequest) error {
	teacher := models.Teacher{}
	result := h.storage.Where("email = ?", req.Email).Take(&teacher)
	if result.Error == nil {
		return errConflict
	}
	password, err := utilities.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err)
		return err
	}
	newTeacher := models.Teacher{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	if err := h.storage.Create(&newTeacher).Error; err != nil {
		logrus.WithError(err)
		return err
	}
	return nil
}

func (h *registrationHandler) CreateStudent(req RegistrateRequest) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", req.Email).Take(&student)
	if result.Error == nil {
		return errConflict
	}
	password, err := utilities.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err)
		return err
	}
	newStudent := models.Student{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	if err := h.storage.Create(&newStudent).Error; err != nil {
		logrus.WithError(err)
		return err
	}
	return nil
}

func (h *loginHandler) Get(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *loginHandler) Login(c *fiber.Ctx) error {
	req := LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if req.Role == Roles.Teacher {
		err := h.LoginTeacher(req)
		if err != nil {
			if errors.Is(err, errBadCredentials) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if req.Role == Roles.Student {
		err := h.LoginStudent(req)
		if err != nil {
			if errors.Is(err, errBadCredentials) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}

	payload := jwt.MapClaims{
		"sub":   req.Email,
		"roles": req.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(initializers.Cfg.JwtSecretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	cookie := new(fiber.Cookie)
	cookie.Name = initializers.Cfg.JwtCookieKey
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HTTPOnly = true

	c.Cookie(cookie)
	return c.Redirect("/profile")
}

func (h *loginHandler) LoginTeacher(req LoginRequest) error {
	teacher := models.Teacher{}
	result := h.storage.Where("email = ?", req.Email).Take(&teacher)
	if result.Error != nil {
		return errBadCredentials
	}
	if !utilities.CheckPasswordHash(req.Password, teacher.Password) {
		return errBadCredentials
	}
	return nil
}

func (h *loginHandler) LoginStudent(req LoginRequest) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", req.Email).Take(&student)
	if result.Error != nil {
		return errBadCredentials
	}
	if !utilities.CheckPasswordHash(req.Password, student.Password) {
		return errBadCredentials
	}
	return nil
}

func (h *profileHandler) Get(c *fiber.Ctx) error {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		return h.GetTeacher(c, jwtPayload)
	} else if role == Roles.Student {
		return h.GetStudent(c, jwtPayload)
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
}

func (h *profileHandler) GetTeacher(c *fiber.Ctx, jwtPayload jwt.MapClaims) error {
	teacher := models.Teacher{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&teacher)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Render("profile", fiber.Map{
		"email": teacher.Email,
		"name":  teacher.Name,
		"role":  Roles.Teacher,
	})
}

func (h *profileHandler) GetStudent(c *fiber.Ctx, jwtPayload jwt.MapClaims) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&student)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Render("profile", fiber.Map{
		"email": student.Email,
		"name":  student.Name,
		"role":  Roles.Student,
	})
}

func (h *profileHandler) Update(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *profileHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) GetList(c *fiber.Ctx) error {
	return c.Render("error", fiber.Map{
		"Error": "404 not found",
	})
}

func (h *homeworksHandler) Get(c *fiber.Ctx) error {
	return c.Render("error", fiber.Map{
		"Error": "404 not found",
	})
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
