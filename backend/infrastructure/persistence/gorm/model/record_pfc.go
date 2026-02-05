package model

// RecordPfc は食事記録のPFC情報を保持するGORMモデル
type RecordPfc struct {
	ID       string  `gorm:"primaryKey;size:36"`
	RecordID string  `gorm:"uniqueIndex;size:36;not null"`
	Protein  float64 `gorm:"not null"`
	Fat      float64 `gorm:"not null"`
	Carbs    float64 `gorm:"not null"`
}

// TableName はテーブル名を明示的に指定する
func (RecordPfc) TableName() string {
	return "record_pfcs"
}
