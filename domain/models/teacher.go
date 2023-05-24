package models

import (
	"time"

	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Id        uint       `gorm:"primaryKey"`
	Email     string     `gorm:"uniqueIndex;notnull"`
	Name      string     `gorm:"not null"`
	Password  string     `gorm:"not null"`
	Students  []Student  `gorm:"foreignKey:TeacherId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Homeworks []Homework `gorm:"foreignKey:TeacherId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
