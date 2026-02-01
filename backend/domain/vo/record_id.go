package vo

import (
	"github.com/google/uuid"
)

// RecordID はカロリー記録の識別子を表す値オブジェクト
type RecordID struct {
	value string
}

// NewRecordID は新しいRecordIDを生成する
func NewRecordID() RecordID {
	return RecordID{value: uuid.New().String()}
}

// ReconstructRecordID はDBからRecordIDを復元する
func ReconstructRecordID(value string) RecordID {
	return RecordID{value: value}
}

// String はRecordIDの文字列表現を返す
func (r RecordID) String() string {
	return r.value
}

// Equals は2つのRecordIDが等しいかを比較する
func (r RecordID) Equals(other RecordID) bool {
	return r.value == other.value
}
