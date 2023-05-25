package routes

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MikhailR1337/task-sync-x/application/utilities"
	"github.com/MikhailR1337/task-sync-x/infrastructure/models"
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
	students := []models.Student{}
	result = h.storage.Where("teacher_id", teacher.Id).Find(&students)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Render("profileTeacher", fiber.Map{
		"email":    teacher.Email,
		"name":     teacher.Name,
		"role":     Roles.Teacher,
		"students": students,
	})
}

func (h *profileHandler) GetStudent(c *fiber.Ctx, jwtPayload jwt.MapClaims) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&student)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if student.TeacherId != 0 {
		teacher := models.Teacher{}
		h.storage.Where("id = ?", student.TeacherId).Take(&teacher)
		return c.Render("profileStudent", fiber.Map{
			"email":       student.Email,
			"name":        student.Name,
			"role":        Roles.Student,
			"teacherName": teacher.Name,
		})
	}
	teachers := []models.Teacher{}
	result = h.storage.Find(&teachers)
	if result.Error != nil {
		return c.Render("profileStudent", fiber.Map{
			"email": student.Email,
			"name":  student.Name,
			"role":  Roles.Student,
		})
	}
	return c.Render("profileStudent", fiber.Map{
		"email":    student.Email,
		"name":     student.Name,
		"role":     Roles.Student,
		"teachers": teachers,
	})
}

type ProfileUpdateRequest struct {
	Teacher string `json:"teacher"`
}

func (h *profileHandler) Update(c *fiber.Ctx) error {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	role := jwtPayload["roles"].(string)
	req := ProfileUpdateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if role == Roles.Teacher {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	} else if role == Roles.Student {
		return h.UpdateStudent(c, jwtPayload, &req)
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
}

func (h *profileHandler) UpdateStudent(c *fiber.Ctx, jwtPayload jwt.MapClaims, req *ProfileUpdateRequest) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&student)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	teacherId, err := strconv.ParseUint(req.Teacher, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	student.TeacherId = uint(teacherId)
	h.storage.Save(&student)
	return c.Redirect("/profile")
}

func (h *profileHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) GetList(c *fiber.Ctx) error {
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
		return h.GetTeachersHomeworks(c, jwtPayload)
	} else if role == Roles.Student {
		return h.GetStudentHomeworks(c, jwtPayload)
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
}

func (h *homeworksHandler) GetTeachersHomeworks(c *fiber.Ctx, jwtPayload jwt.MapClaims) error {
	teacher := models.Teacher{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&teacher)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	homeworks := []models.Homework{}
	result = h.storage.Where("teacher_id", teacher.Id).Find(&homeworks)
	if result.Error != nil {
		return c.Render("homeworks", fiber.Map{})
	}
	students := []models.Student{}
	result = h.storage.Where("teacher_id", teacher.Id).Find(&students)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Render("homeworks", fiber.Map{
		"homeworks": homeworks,
		"students":  students,
		"isTeacher": true,
	})
}

func (h *homeworksHandler) GetStudentHomeworks(c *fiber.Ctx, jwtPayload jwt.MapClaims) error {
	student := models.Student{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&student)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	homeworks := []models.Homework{}
	result = h.storage.Where("student_id", student.Id).Find(&homeworks)
	if result.Error != nil {
		return c.Render("homeworks", fiber.Map{})
	}
	return c.Render("homeworks", fiber.Map{
		"homeworks": homeworks,
	})
}

type CreateHomeworkRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	CurrentPoints string `json:"currentPoints"`
	MaxPoints     string `json:"maxPoints"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Student       string `json:"student"`
}

func (h *homeworksHandler) Create(c *fiber.Ctx) error {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	teacher := models.Teacher{}
	result := h.storage.Where("email = ?", jwtPayload["sub"].(string)).Take(&teacher)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	req := CreateHomeworkRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	currentPoints, err := strconv.ParseUint(req.CurrentPoints, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	maxPoints, err := strconv.ParseUint(req.MaxPoints, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	studentId, err := strconv.ParseUint(req.Student, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	newHomework := models.Homework{
		Name:          req.Name,
		Description:   req.Description,
		CurrentPoints: uint8(currentPoints),
		MaxPoints:     uint8(maxPoints),
		Type:          req.Type,
		Status:        req.Status,
		TeacherId:     teacher.Id,
		StudentId:     uint(studentId),
	}
	if err := h.storage.Create(&newHomework).Error; err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Redirect("/homeworks")
}

func (h *homeworksHandler) Get(c *fiber.Ctx) error {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	homework := models.Homework{}
	result := h.storage.Where("id = ?", uint(homeworkId)).Take(&homework)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	teacher := models.Teacher{}
	result = h.storage.Where("id = ?", homework.TeacherId).Take(&teacher)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	student := models.Student{}
	result = h.storage.Where("id = ?", homework.StudentId).Take(&student)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	role := jwtPayload["roles"].(string)
	return c.Render("homework", fiber.Map{
		"id":            homework.Id,
		"name":          homework.Name,
		"description":   homework.Description,
		"currentPoints": homework.CurrentPoints,
		"maxPoints":     homework.MaxPoints,
		"type":          homework.Type,
		"status":        homework.Status,
		"teacher":       teacher.Name,
		"student":       student.Name,
		"isTeacher":     role == Roles.Teacher,
	})
}

type UpdateHomeworkStudentRequest struct {
	Status string `json:"status"`
}

type UpdateHomeworkTeacherRequest struct {
	Status        string `json:"status"`
	CurrentPoints string `json:"currentPoints"`
}

func (h *homeworksHandler) Update(c *fiber.Ctx) error {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	homework := models.Homework{}
	result := h.storage.Where("id = ?", uint(homeworkId)).Take(&homework)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		req := UpdateHomeworkTeacherRequest{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		currentPoints, err := strconv.ParseUint(req.CurrentPoints, 10, 32)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		homework.CurrentPoints = uint8(currentPoints)
		homework.Status = req.Status
	} else if role == Roles.Student {
		req := UpdateHomeworkStudentRequest{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		homework.Status = req.Status
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
	h.storage.Save(&homework)
	return c.Redirect(fmt.Sprintf("/homeworks/%s", homeworkParam))
}

func (h *homeworksHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}
