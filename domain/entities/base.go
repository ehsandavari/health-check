package entities

import (
	"gorm.io/gorm"
	"time"
)

type Base1 struct {
	CreatedAt time.Time `gorm:"not null"`
}

type Base2 struct {
	CreatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Base3 struct {
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
