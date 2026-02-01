package vo

import (
	"github.com/google/uuid"
)

// RecordItemID はカロリー記録明細の識別子を表す値オブジェクト
type RecordItemID struct {
	value string
}

// NewRecordItemID は新しいRecordItemIDを生成する
func NewRecordItemID() RecordItemID {
	return RecordItemID{value: uuid.New().String()}
}

// ReconstructRecordItemID はDBからRecordItemIDを復元する
func ReconstructRecordItemID(value string) RecordItemID {
	return RecordItemID{value: value}
}

// String はRecordItemIDの文字列表現を返す
func (r RecordItemID) String() string {
	return r.value
}

// Equals は2つのRecordItemIDが等しいかを比較する
func (r RecordItemID) Equals(other RecordItemID) bool {
	return r.value == other.value
}
