package vo

import (
	domainErrors "caltrack/domain/errors"
)

// RecordPfcID はカロリー記録のPFC情報の識別子を表す値オブジェクト
type RecordPfcID struct {
	value UUID
}

// NewRecordPfcID は新しいRecordPfcIDを生成する
func NewRecordPfcID() RecordPfcID {
	return RecordPfcID{value: NewUUID()}
}

// ParseRecordPfcID は文字列からRecordPfcIDを生成する
func ParseRecordPfcID(value string) (RecordPfcID, error) {
	parsed, err := ParseUUID(value)
	if err != nil {
		return RecordPfcID{}, domainErrors.ErrInvalidRecordPfcID
	}
	return RecordPfcID{value: parsed}, nil
}

// ReconstructRecordPfcID はDBからRecordPfcIDを復元する
func ReconstructRecordPfcID(value string) RecordPfcID {
	return RecordPfcID{value: ReconstructUUID(value)}
}

// String はRecordPfcIDの文字列表現を返す
func (r RecordPfcID) String() string {
	return r.value.String()
}

// IsZero はRecordPfcIDがゼロ値かを判定する
func (r RecordPfcID) IsZero() bool {
	return r.value.IsZero()
}

// Equals は2つのRecordPfcIDが等しいかを比較する
func (r RecordPfcID) Equals(other RecordPfcID) bool {
	return r.value.Equals(other.value)
}
