package repository

import (
	"errors"

	"github.com/MikhailR1337/task-sync-x/infrastructure/models"
	"github.com/MikhailR1337/task-sync-x/initializers"
	"gorm.io/gorm/clause"
)

var (
	errHomeworkNotFound   = errors.New("homework is not found")
	errHomeworkNotCreated = errors.New("homework is not created")
	errHomeworkNotUpdated = errors.New("homework is not updated")
	errHomeworkNotDeleted = errors.New("homework is not deleted")
)

const homeworkOrder = "(case status when 'new' then 1 when 'processing' then 2 when 'finished' then 3 when 'checked' then 4 end)"

var Homework = &homework{&initializers.DB}

type homework struct {
	storage *initializers.PgDb
}

func (h *homework) GetById(id uint) (*models.Homework, error) {
	homework := &models.Homework{}
	result := h.storage.Where("id = ?", id).Take(homework)
	if result.Error != nil {
		return nil, errHomeworkNotFound
	}
	return homework, nil
}

func (h *homework) GetByTeacherId(id uint) (*[]models.Homework, error) {
	homeworks := &[]models.Homework{}
	result := h.storage.Where("teacher_id", id).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: homeworkOrder},
	}).Find(homeworks)
	if result.Error != nil {
		return nil, errHomeworkNotFound
	}
	return homeworks, nil
}

func (h *homework) GetByStudentId(id uint) (*[]models.Homework, error) {
	homeworks := &[]models.Homework{}
	result := h.storage.Where("student_id", id).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: homeworkOrder},
	}).Find(homeworks)
	if result.Error != nil {
		return nil, errHomeworkNotFound
	}
	return homeworks, nil
}

func (h *homework) Create(model *models.Homework) error {
	if err := h.storage.Create(model).Error; err != nil {
		return errHomeworkNotCreated
	}
	return nil
}

func (h *homework) Update(model *models.Homework) error {
	if err := h.storage.Save(model).Error; err != nil {
		return errHomeworkNotUpdated
	}
	return nil
}

func (h *homework) DeleteByTeacherId(id uint) error {
	if err := h.storage.Where("teacher_id", id).Delete(&models.Homework{}).Error; err != nil {
		return errHomeworkNotDeleted
	}
	return nil
}

func (h *homework) DeleteByStudentId(id uint) error {
	if err := h.storage.Where("student_id", id).Delete(&models.Homework{}).Error; err != nil {
		return errHomeworkNotDeleted
	}
	return nil
}

func (h *homework) Delete(id uint) error {
	if err := h.storage.Where("id", id).Delete(&models.Homework{}).Error; err != nil {
		return errHomeworkNotDeleted
	}
	return nil
}
