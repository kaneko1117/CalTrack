package vo

import (
	domainErrors "caltrack/domain/errors"
)

// RecordID はカロリー記録の識別子を表す値オブジェクト
type RecordID struct {
	value UUID
}

// NewRecordID は新しいRecordIDを生成する
func NewRecordID() RecordID {
	return RecordID{value: NewUUID()}
}

// ParseRecordID は文字列からRecordIDを生成する
func ParseRecordID(value string) (RecordID, error) {
	parsed, err := ParseUUID(value)
	if err != nil {
		return RecordID{}, domainErrors.ErrInvalidRecordID
	}
	return RecordID{value: parsed}, nil
}

// ReconstructRecordID はDBからRecordIDを復元する
func ReconstructRecordID(value string) RecordID {
	return RecordID{value: ReconstructUUID(value)}
}

// String はRecordIDの文字列表現を返す
func (r RecordID) String() string {
	return r.value.String()
}

// IsZero はRecordIDがゼロ値かを判定する
func (r RecordID) IsZero() bool {
	return r.value.IsZero()
}

// Equals は2つのRecordIDが等しいかを比較する
func (r RecordID) Equals(other RecordID) bool {
	return r.value.Equals(other.value)
}
