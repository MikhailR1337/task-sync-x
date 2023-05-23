package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint    `gorm:"primaryKey"`
	Email     *string `gorm:"uniqueIndex"`
	Name      string
	Role      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
