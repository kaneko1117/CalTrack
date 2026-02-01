package vo

import (
	"github.com/google/uuid"
)

type UserID struct {
	value string
}

func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// ReconstructUserID はDBからUserIDを復元する
func ReconstructUserID(value string) UserID {
	return UserID{value: value}
}

func (u UserID) String() string {
	return u.value
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}
