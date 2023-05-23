package models

import (
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	Id            uint `gorm:"primaryKey"`
	Name          string
	Description   string
	CurrentPoints uint8
	MaxPoints     uint8
	Type          string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
