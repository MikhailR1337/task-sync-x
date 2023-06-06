package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Id        uint       `gorm:"primaryKey"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Name      string     `gorm:"not null"`
	Password  string     `gorm:"not null"`
	TeacherId uint       `gorm:"default:null"`
	Homeworks []Homework `gorm:"foreignKey:StudentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
