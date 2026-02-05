package model

import "time"

// AdviceCache はAIアドバイスキャッシュを保持するGORMモデル
type AdviceCache struct {
	ID        string    `gorm:"primaryKey;size:36"`
	UserID    string    `gorm:"size:36;not null;index:uk_user_date,unique"`
	CacheDate time.Time `gorm:"type:date;not null;index:uk_user_date,unique"`
	Advice    string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (AdviceCache) TableName() string {
	return "advice_caches"
}
