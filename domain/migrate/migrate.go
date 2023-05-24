package migrate

import (
	"github.com/MikhailR1337/task-sync-x/domain/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Teacher{},
		&models.Student{},
		&models.Homework{},
	)
	if err != nil {
		return err
	}
	return nil
}
