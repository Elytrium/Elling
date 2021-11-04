package elling

import (
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type DBModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}