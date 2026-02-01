package model

// RecordItem はカロリー記録明細を保持するGORMモデル
type RecordItem struct {
	ID       string `gorm:"primaryKey;size:36"`
	RecordID string `gorm:"index;size:36;not null"`
	Name     string `gorm:"size:100;not null"`
	Calories int    `gorm:"not null"`
}
