package model

import "time"

type User struct {
	ID             string    `gorm:"primaryKey;size:36"`
	Email          string    `gorm:"uniqueIndex;size:254;not null"`
	HashedPassword string    `gorm:"size:60;not null"`
	Nickname       string    `gorm:"size:50;not null"`
	Weight         float64   `gorm:"not null"`
	Height         float64   `gorm:"not null"`
	BirthDate      time.Time `gorm:"not null"`
	Gender         string    `gorm:"size:10;not null"`
	ActivityLevel  string    `gorm:"size:20;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Records        []Record `gorm:"foreignKey:UserID"`
}
