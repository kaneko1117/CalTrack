package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Email     string         `gorm:"uniqueIndex;size:255" json:"email"`
	Name      string         `gorm:"size:255" json:"name"`
}

func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
