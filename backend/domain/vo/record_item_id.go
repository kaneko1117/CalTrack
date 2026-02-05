package vo

import (
	domainErrors "caltrack/domain/errors"
)

// RecordItemID はカロリー記録明細の識別子を表す値オブジェクト
type RecordItemID struct {
	value UUID
}

// NewRecordItemID は新しいRecordItemIDを生成する
func NewRecordItemID() RecordItemID {
	return RecordItemID{value: NewUUID()}
}

// ParseRecordItemID は文字列からRecordItemIDを生成する
func ParseRecordItemID(value string) (RecordItemID, error) {
	parsed, err := ParseUUID(value)
	if err != nil {
		return RecordItemID{}, domainErrors.ErrInvalidRecordItemID
	}
	return RecordItemID{value: parsed}, nil
}

// ReconstructRecordItemID はDBからRecordItemIDを復元する
func ReconstructRecordItemID(value string) RecordItemID {
	return RecordItemID{value: ReconstructUUID(value)}
}

// String はRecordItemIDの文字列表現を返す
func (r RecordItemID) String() string {
	return r.value.String()
}

// IsZero はRecordItemIDがゼロ値かを判定する
func (r RecordItemID) IsZero() bool {
	return r.value.IsZero()
}

// Equals は2つのRecordItemIDが等しいかを比較する
func (r RecordItemID) Equals(other RecordItemID) bool {
	return r.value.Equals(other.value)
}
