package routes

import (
	"errors"
	"strconv"
	"time"

	"github.com/MikhailR1337/task-sync-x/app/application/forms"
	"github.com/MikhailR1337/task-sync-x/app/application/utilities"
	"github.com/MikhailR1337/task-sync-x/app/infrastructure/models"
	"github.com/MikhailR1337/task-sync-x/app/infrastructure/repository"
	"github.com/MikhailR1337/task-sync-x/app/initializers"
	"github.com/MikhailR1337/task-sync-x/app/services/mailer"
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
	errNotFound       = errors.New("404 Oops, page not found")
	errSomethingWrong = errors.New("something wrong. try again")
	errConflict       = errors.New("oops... we already have this email")
	errBadCredentials = errors.New("email or password is incorrect")
	errValidation     = errors.New("something wrong with your data. change something and try again")
	errPoints         = errors.New("points should be a number")
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
	return c.Render("main", nil)
}

func (h *registrationHandler) Get(c *fiber.Ctx) error {
	return c.Render("registration", nil)
}

func (h *registrationHandler) Registrate(c *fiber.Ctx) error {
	req := forms.RegistrateRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.Render("registration", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	err := initializers.Validator.Struct(req)
	if err != nil {
		logrus.WithError(err)
		return c.Render("registration", fiber.Map{
			"error": errValidation,
		})
	}
	password, err := utilities.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err)
		return c.Render("registration", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	if req.Role == Roles.Teacher {
		_, err := repository.Teacher.GetByEmail(req.Email)
		if err == nil {
			logrus.WithError(err)
			return c.Render("registration", fiber.Map{
				"error": errConflict,
			})
		}
		newTeacher := &models.Teacher{
			Name:     req.Name,
			Email:    req.Email,
			Password: password,
		}
		err = repository.Teacher.Create(newTeacher)
		if err != nil {
			logrus.WithError(err)
			return c.Render("registration", fiber.Map{
				"error": errSomethingWrong,
			})
		}
	} else if req.Role == Roles.Student {
		_, err := repository.Student.GetByEmail(req.Email)
		if err == nil {
			logrus.WithError(err)
			return c.Render("registration", fiber.Map{
				"error": errConflict,
			})
		}
		newStudent := &models.Student{
			Name:     req.Name,
			Email:    req.Email,
			Password: password,
		}
		err = repository.Student.Create(newStudent)
		if err != nil {
			logrus.WithError(err)
			return c.Render("registration", fiber.Map{
				"error": errSomethingWrong,
			})
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
		return c.Render("login", fiber.Map{
			"error": errSomethingWrong,
		})
	}

	err := initializers.Validator.Struct(req)
	if err != nil {
		logrus.WithError(err)
		return c.Render("login", fiber.Map{
			"error": errValidation,
		})
	}
	if req.Role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(req.Email)
		if err != nil {
			logrus.WithError(err)
			return c.Render("login", fiber.Map{
				"error": errBadCredentials,
			})
		}
		if !utilities.CheckPasswordHash(req.Password, teacher.Password) {
			logrus.WithError(err)
			return c.Render("login", fiber.Map{
				"error": errBadCredentials,
			})
		}
	} else if req.Role == Roles.Student {
		student, err := repository.Student.GetByEmail(req.Email)
		if err != nil {
			logrus.WithError(err)
			return c.Render("login", fiber.Map{
				"error": errBadCredentials,
			})
		}
		if !utilities.CheckPasswordHash(req.Password, student.Password) {
			logrus.WithError(err)
			return c.Render("login", fiber.Map{
				"error": errBadCredentials,
			})
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
		return c.Render("login", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = initializers.Cfg.JwtCookieKey
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HTTPOnly = true

	c.Cookie(cookie)
	return c.Redirect("/profile")
}

func (h *loginHandler) SignOut(c *fiber.Ctx) error {
	c.ClearCookie(initializers.Cfg.JwtCookieKey)
	return c.SendStatus(fiber.StatusOK)
}

func (h *profileHandler) Get(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Redirect("/login")
		}
		students, err := repository.Student.GetByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
			return c.Render("profileTeacher", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		return c.Render("profileTeacher", fiber.Map{
			"email":    teacher.Email,
			"name":     teacher.Name,
			"role":     Roles.Teacher,
			"students": *students,
		})
	}
	student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	if student.TeacherId != 0 {
		teacher, err := repository.Teacher.GetById(student.TeacherId)
		if err != nil {
			logrus.WithError(err)
			return c.Render("/profileStudent", fiber.Map{
				"error": errSomethingWrong,
			})
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
		"teachers": *teachers,
	})
}

func (h *profileHandler) Update(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	role := jwtPayload["roles"].(string)
	req := forms.StudentUpdateRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		if role == Roles.Student {
			return c.Render("profileStudent", fiber.Map{
				"error": errSomethingWrong,
			})
		} else {
			return c.Render("profileTeacher", fiber.Map{
				"error": errSomethingWrong,
			})
		}
	}
	if role == Roles.Student {
		student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Redirect("/login")
		}
		teacherId, err := strconv.ParseUint(req.Teacher, 10, 32)
		if err != nil {
			logrus.WithError(err)
			return c.Render("profileStudent", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		student.TeacherId = uint(teacherId)
		repository.Student.Update(student)
		return c.SendStatus(fiber.StatusOK)
	}
	return c.Status(fiber.StatusNotFound).Redirect("/")
}

func (h *profileHandler) Delete(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Redirect("/login")
		}
		err = repository.Teacher.Delete(teacher)
		if err != nil {
			logrus.WithError(err)
			return c.Render("profileTeacher", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		err = repository.Homework.DeleteByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
		}
	} else {
		student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Redirect("/login")
		}
		err = repository.Student.Delete(student)
		if err != nil {
			logrus.WithError(err)
			return c.Render("profileStudent", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		err = repository.Homework.DeleteByStudentId(student.Id)
		if err != nil {
			logrus.WithError(err)
		}
	}
	c.ClearCookie(initializers.Cfg.JwtCookieKey)
	return c.SendStatus(fiber.StatusOK)
}

func (h *homeworksHandler) GetList(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
		if err != nil {
			logrus.WithError(err)
			return c.Redirect("/login")
		}
		students, err := repository.Student.GetByTeacherId(teacher.Id)
		if err != nil {
			logrus.WithError(err)
			return c.Render("homeworks", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		homeworks, err := repository.Homework.GetByTeacherId(teacher.Id)
		if err != nil {
			return c.Render("homeworks", fiber.Map{})
		}
		return c.Render("homeworks", fiber.Map{
			"homeworks": *homeworks,
			"students":  *students,
			"isTeacher": true,
		})
	}
	student, err := repository.Student.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	homeworks, err := repository.Homework.GetByStudentId(student.Id)
	if err != nil {
		return c.Render("homeworks", fiber.Map{})
	}
	return c.Render("homeworks", fiber.Map{
		"homeworks": *homeworks,
	})
}

func (h *homeworksHandler) Create(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	req := forms.CreateHomeworkRequest{}
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	err = initializers.Validator.Struct(req)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errValidation,
		})
	}
	teacher, err := repository.Teacher.GetByEmail(jwtPayload["sub"].(string))
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	currentPoints, err := strconv.ParseUint(req.CurrentPoints, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errPoints,
		})
	}
	maxPoints, err := strconv.ParseUint(req.MaxPoints, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errPoints,
		})
	}
	studentId, err := strconv.ParseUint(req.Student, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errSomethingWrong,
		})
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
		return c.Render("homeworks", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	student, err := repository.Student.GetById(uint(studentId))
	if err != nil {
		logrus.WithError(err)
		return c.Render("homeworks", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	mailer.NewHomework(student.Email, student.Name, newHomework.Name)
	return c.Redirect("/homeworks")
}

func (h *homeworksHandler) Get(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errNotFound,
		})
	}
	homework, err := repository.Homework.GetById(uint(homeworkId))
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/homeworks")
	}
	teacher, err := repository.Teacher.GetById(homework.TeacherId)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	student, err := repository.Student.GetById(homework.StudentId)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	role := jwtPayload["roles"].(string)
	return c.Render("homework", fiber.Map{
		"id":                 homework.Id,
		"name":               homework.Name,
		"description":        homework.Description,
		"currentPoints":      homework.CurrentPoints,
		"maxPoints":          homework.MaxPoints,
		"type":               homework.Type,
		"status":             homework.Status,
		"teacher":            teacher.Name,
		"student":            student.Name,
		"isTeacher":          role == Roles.Teacher,
		"isStudentCanStart":  homework.Status == "new",
		"isStudentCanFinish": homework.Status == "processing",
		"isTeacherCanCheck":  homework.Status == "finished",
		"isChecked":          homework.Status == "checked",
	})
}

func (h *homeworksHandler) Update(c *fiber.Ctx) error {
	jwtPayload, err := utilities.GetJwtPayload(c)
	if err != nil {
		logrus.WithError(err)
		return c.Redirect("/login")
	}
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errNotFound,
		})
	}
	homework, err := repository.Homework.GetById(uint(homeworkId))
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errNotFound,
		})
	}
	role := jwtPayload["roles"].(string)
	if role == Roles.Teacher {
		req := forms.UpdateHomeworkTeacherRequest{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return c.Render("homework", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		currentPoints, err := strconv.ParseUint(req.CurrentPoints, 10, 32)
		if err != nil {
			logrus.WithError(err)
			return c.Render("homework", fiber.Map{
				"error": errPoints,
			})
		}
		homework.CurrentPoints = uint8(currentPoints)
		homework.Status = req.Status
		student, err := repository.Student.GetById(homework.StudentId)
		if err != nil {
			logrus.WithError(err)
			return c.Render("homeworks", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		mailer.CheckedHomework(student.Email, student.Name, homework.Name)
	} else if role == Roles.Student {
		req := forms.UpdateHomeworkStudentRequest{}
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err)
			return c.Render("homework", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		homework.Status = req.Status
		teacher, err := repository.Teacher.GetById(homework.TeacherId)
		if err != nil {
			logrus.WithError(err)
			return c.Render("homeworks", fiber.Map{
				"error": errSomethingWrong,
			})
		}
		mailer.CheckedHomework(teacher.Email, teacher.Name, homework.Name)
	}
	err = repository.Homework.Update(homework)
	if err != nil {
		logrus.WithError(err)
		return c.Render("homework", fiber.Map{
			"error": errSomethingWrong,
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *homeworksHandler) Delete(c *fiber.Ctx) error {
	homeworkParam := c.Params("id")
	homeworkId, err := strconv.ParseUint(homeworkParam, 10, 32)
	if err != nil {
		logrus.WithError(err)
	}
	err = repository.Homework.Delete(uint(homeworkId))
	if err != nil {
		logrus.WithError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}
