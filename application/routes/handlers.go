package routes

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MikhailR1337/task-sync-x/application/forms"
	"github.com/MikhailR1337/task-sync-x/application/utilities"
	"github.com/MikhailR1337/task-sync-x/infrastructure/models"
	"github.com/MikhailR1337/task-sync-x/infrastructure/repository"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var (
	MainPageHandler     = &mainPageHandler{}
	RegistrationHandler = &registrationHandler{}
	LoginHandler        = &loginHandler{}
	ProfileHandler      = &profileHandler{}
	HomeworkHandler     = &homeworksHandler{}
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
	mainPageHandler     struct{}
	registrationHandler struct{}
	loginHandler        struct{}
	profileHandler      struct{}
	homeworksHandler    struct{}
	roles               struct {
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

func (h *registrationHandler) Registrate(c *fiber.Ctx) error {
	req := forms.RegistrateRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	err := initializers.Validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}
	password, err := utilities.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err)
		return err
	}
	if req.Role == Roles.Teacher {
		_, err := repository.Teacher.GetByEmail(req.Email)
		if err == nil {
			logrus.WithError(err)
			return errConflict
		}
		newTeacher := &models.Teacher{
			Name:     req.Name,
			Email:    req.Email,
			Password: password,
		}
		err = repository.Teacher.Create(newTeacher)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if req.Role == Roles.Student {
		_, err := repository.Student.GetByEmail(req.Email)
		if err == nil {
			logrus.WithError(err)
			return errConflict
		}
		newStudent := &models.Student{
			Name:     req.Name,
			Email:    req.Email,
			Password: password,
		}
		err = repository.Student.Create(newStudent)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.Status(fiber.StatusCreated).Redirect("/login")
}

func (h *loginHandler) Get(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func (h *loginHandler) Login(c *fiber.Ctx) error {
	req := forms.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	err := initializers.Validator.Struct(req)
	if err != nil {
		logrus.WithError(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}
	if req.Role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(req.Email)
		if err != nil {
			logrus.WithError(err)
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		if !utilities.CheckPasswordHash(req.Password, teacher.Password) {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(errBadCredentials.Error())
		}
	} else if req.Role == Roles.Student {
		student, err := repository.Student.GetByEmail(req.Email)
		if err != nil {
			logrus.WithError(err)
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		if !utilities.CheckPasswordHash(req.Password, student.Password) {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(errBadCredentials.Error())
		}
	}

	payload := jwt.MapClaims{
		"sub":   req.Email,
		"roles": req.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(initializers.Cfg.JwtSecretKey))
	if err != nil {
		logrus.WithError(err)
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

func (h *profileHandler) Get(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		students, err := repository.Student.GetByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		return c.Render("profileTeacher", fiber.Map{
			"email":    teacher.Email,
			"name":     teacher.Name,
			"role":     Roles.Teacher,
			"students": students,
		})
	}
	student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}
	if student.TeacherId != 0 {
		teacher, err := repository.Teacher.GetById(student.TeacherId)
		if err != nil {
			logrus.WithError(err)
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.Render("profileStudent", fiber.Map{
			"email":       student.Email,
			"name":        student.Name,
			"role":        Roles.Student,
			"teacherName": teacher.Name,
		})
	}
	teachers, err := repository.Teacher.GetList()
	if err != nil {
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

func (h *profileHandler) Update(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	role := jwtPayload["roles"].(string)
	req := forms.StudentUpdateRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if role == Roles.Student {
		student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		teacherId, err := strconv.ParseUint(req.Teacher, 10, 32)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		student.TeacherId = uint(teacherId)
		repository.Student.Update(student)
		return c.Redirect("/profile")
	}
	return c.Status(fiber.StatusNotFound).Redirect("/")
}

func (h *profileHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}

func (h *homeworksHandler) GetList(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		students, err := repository.Student.GetByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		homeworks, err := repository.Homework.GetByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Render("homeworks", fiber.Map{
			"homeworks": homeworks,
			"students":  students,
			"isTeacher": true,
		})
	}
	student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	homeworks, err := repository.Homework.GetByStudentId(student.Id)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{})
	}
	return c.Render("homeworks", fiber.Map{
		"homeworks": homeworks,
	})
}

func (h *homeworksHandler) Create(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	req := forms.CreateHomeworkRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	err = initializers.Validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}
	teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
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
	err = repository.Homework.Create(&newHomework)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Redirect("/homeworks")
}

func (h *homeworksHandler) Get(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	homework, err := repository.Homework.GetById(uint(homeworkId))
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	teacher, err := repository.Teacher.GetById(homework.TeacherId)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	student, err := repository.Student.GetById(homework.StudentId)
	if err != nil {
		logrus.WithError(err)
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

func (h *homeworksHandler) Update(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	homework, err := repository.Homework.GetById(uint(homeworkId))
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		req := forms.UpdateHomeworkTeacherRequest{}
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
		req := forms.UpdateHomeworkStudentRequest{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		homework.Status = req.Status
	} else {
		return c.Status(fiber.StatusNotFound).Redirect("/")
	}
	err = repository.Homework.Update(homework)
	if err != nil {
		logrus.WithError(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Redirect(fmt.Sprintf("/homeworks/%s", homeworkParam))
}

func (h *homeworksHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("HELLO")
}
