package model

import "time"

// Session はセッション情報を保持するGORMモデル
type Session struct {
	ID        string    `gorm:"primaryKey;size:44"` // Base64エンコードされたSessionID（44文字）
	UserID    string    `gorm:"index;size:36;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}
