package vo

import (
	"github.com/google/uuid"
	domainErrors "caltrack/domain/errors"
)

// UUID はUUIDを表す値オブジェクト
type UUID struct {
	value uuid.UUID
}

// NewUUID は新しいUUIDを生成する
func NewUUID() UUID {
	return UUID{value: uuid.New()}
}

// ParseUUID は文字列からUUIDを生成する
func ParseUUID(value string) (UUID, error) {
	if value == "" {
		return UUID{}, domainErrors.ErrUUIDRequired
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return UUID{}, domainErrors.ErrInvalidUUIDFormat
	}
	return UUID{value: parsed}, nil
}

// ReconstructUUID はDBからUUIDを復元する（バリデーションなし）
func ReconstructUUID(value string) UUID {
	parsed, _ := uuid.Parse(value)
	return UUID{value: parsed}
}

// String はUUIDの文字列表現を返す
func (u UUID) String() string {
	return u.value.String()
}

// IsZero はUUIDがゼロ値かを判定する
func (u UUID) IsZero() bool {
	return u.value == uuid.Nil
}

// Equals は2つのUUIDが等しいかを比較する
func (u UUID) Equals(other UUID) bool {
	return u.value == other.value
}
