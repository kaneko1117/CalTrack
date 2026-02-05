package vo

import (
	domainErrors "caltrack/domain/errors"
)

// UserID はユーザーの識別子を表す値オブジェクト
type UserID struct {
	value UUID
}

// NewUserID は新しいUserIDを生成する
func NewUserID() UserID {
	return UserID{value: NewUUID()}
}

// ParseUserID は文字列からUserIDを生成する
func ParseUserID(value string) (UserID, error) {
	parsed, err := ParseUUID(value)
	if err != nil {
		return UserID{}, domainErrors.ErrInvalidUserID
	}
	return UserID{value: parsed}, nil
}

// ReconstructUserID はDBからUserIDを復元する
func ReconstructUserID(value string) UserID {
	return UserID{value: ReconstructUUID(value)}
}

// String はUserIDの文字列表現を返す
func (u UserID) String() string {
	return u.value.String()
}

// IsZero はUserIDがゼロ値かを判定する
func (u UserID) IsZero() bool {
	return u.value.IsZero()
}

// Equals は2つのUserIDが等しいかを比較する
func (u UserID) Equals(other UserID) bool {
	return u.value.Equals(other.value)
}
