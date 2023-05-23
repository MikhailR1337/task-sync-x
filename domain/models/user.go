package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint `gorm:"primaryKey"`
	Name      string
	Role      string
	Email     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
