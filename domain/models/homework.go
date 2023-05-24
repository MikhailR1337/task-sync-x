package models

import (
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	Id            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Description   string
	CurrentPoints uint8  `gorm:"not null;default:0"`
	MaxPoints     uint8  `gorm:"not null;default:40"`
	Type          string `gorm:"not null"`
	Status        string `gorm:"not null"`
	TeacherId     uint   `gorm:"not null"`
	StudentId     uint   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
