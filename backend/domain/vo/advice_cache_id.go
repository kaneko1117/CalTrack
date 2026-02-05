package vo

import (
	domainErrors "caltrack/domain/errors"
)

// AdviceCacheID はアドバイスキャッシュの識別子を表す値オブジェクト
type AdviceCacheID struct {
	value UUID
}

// NewAdviceCacheID は新しいAdviceCacheIDを生成する
func NewAdviceCacheID() AdviceCacheID {
	return AdviceCacheID{value: NewUUID()}
}

// ParseAdviceCacheID は文字列からAdviceCacheIDを生成する
func ParseAdviceCacheID(value string) (AdviceCacheID, error) {
	parsed, err := ParseUUID(value)
	if err != nil {
		return AdviceCacheID{}, domainErrors.ErrInvalidAdviceCacheID
	}
	return AdviceCacheID{value: parsed}, nil
}

// ReconstructAdviceCacheID はDBからAdviceCacheIDを復元する
func ReconstructAdviceCacheID(value string) AdviceCacheID {
	return AdviceCacheID{value: ReconstructUUID(value)}
}

// String はAdviceCacheIDの文字列表現を返す
func (a AdviceCacheID) String() string {
	return a.value.String()
}

// IsZero はAdviceCacheIDがゼロ値かを判定する
func (a AdviceCacheID) IsZero() bool {
	return a.value.IsZero()
}

// Equals は2つのAdviceCacheIDが等しいかを比較する
func (a AdviceCacheID) Equals(other AdviceCacheID) bool {
	return a.value.Equals(other.value)
}
