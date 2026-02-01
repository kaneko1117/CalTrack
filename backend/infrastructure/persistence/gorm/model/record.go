package model

import "time"

// Record はカロリー記録を保持するGORMモデル
type Record struct {
	ID        string       `gorm:"primaryKey;size:36"`
	UserID    string       `gorm:"index;size:36;not null"`
	EatenAt   time.Time    `gorm:"index;not null"`
	CreatedAt time.Time
	User      User         `gorm:"foreignKey:UserID"`
	Items     []RecordItem `gorm:"foreignKey:RecordID"`
}
