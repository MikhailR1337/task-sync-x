package repository

import (
	"errors"

	"github.com/MikhailR1337/task-sync-x/infrastructure/models"
	"github.com/MikhailR1337/task-sync-x/initializers"
)

var (
	errTeacherNotFound   = errors.New("teacher is not found")
	errTeacherNotCreated = errors.New("teacher is not created")
)

var Teacher = &teacher{&initializers.DB}

type teacher struct {
	storage *initializers.PgDb
}

func (h *teacher) GetList() (*[]models.Teacher, error) {
	teachers := &[]models.Teacher{}
	result := h.storage.Find(teachers)
	if result.Error != nil {
		return nil, errTeacherNotFound
	}
	return teachers, nil
}

func (h *teacher) GetById(id uint) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	result := h.storage.Where("id = ?", id).Take(teacher)
	if result.Error != nil {
		return nil, errTeacherNotFound
	}
	return teacher, nil
}

func (h *teacher) GetByEmail(email string) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	result := h.storage.Where("email = ?", email).Take(teacher)
	if result.Error != nil {
		return nil, errTeacherNotFound
	}
	return teacher, nil
}

func (h *teacher) Create(model *models.Teacher) error {
	if err := h.storage.Create(model).Error; err != nil {
		return errTeacherNotCreated
	}
	return nil
}