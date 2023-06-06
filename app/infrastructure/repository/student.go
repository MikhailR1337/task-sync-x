package repository

import (
	"errors"

	"github.com/MikhailR1337/task-sync-x/app/infrastructure/models"
	"github.com/MikhailR1337/task-sync-x/app/initializers"
)

var (
	errStudentNotFound   = errors.New("student is not found")
	errStudentNotCreated = errors.New("student is not created")
	errStudentNotUpdated = errors.New("student is not updated")
	errStudentNotDeleted = errors.New("student is not deleted")
)

var Student = &student{&initializers.DB}

type student struct {
	storage *initializers.PgDb
}

func (h *student) GetById(id uint) (*models.Student, error) {
	student := &models.Student{}
	result := h.storage.Where("id = ?", id).Take(student)
	if result.Error != nil {
		return nil, errStudentNotFound
	}
	return student, nil
}

func (h *student) GetByTeacherId(id uint) (*[]models.Student, error) {
	students := &[]models.Student{}
	result := h.storage.Where("teacher_id = ?", id).Find(students)
	if result.Error != nil {
		return nil, errStudentNotFound
	}
	return students, nil
}

func (h *student) GetByEmail(email string) (*models.Student, error) {
	student := &models.Student{}
	result := h.storage.Where("email = ?", email).Take(student)
	if result.Error != nil {
		return nil, errStudentNotFound
	}
	return student, nil
}

func (h *student) Create(model *models.Student) error {
	if err := h.storage.Create(model).Error; err != nil {
		return errStudentNotCreated
	}
	return nil
}

func (h *student) Update(model *models.Student) error {
	if err := h.storage.Save(model).Error; err != nil {
		return errStudentNotUpdated
	}
	return nil
}

func (h *student) Delete(model *models.Student) error {
	if err := h.storage.Delete(model).Error; err != nil {
		return errStudentNotDeleted
	}
	return nil
}
