package vo

import (
	"github.com/google/uuid"
	domainErrors "caltrack/domain/errors"
)

type UserID struct {
	value string
}

func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

func ParseUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, domainErrors.ErrInvalidUserID
	}
	if _, err := uuid.Parse(value); err != nil {
		return UserID{}, domainErrors.ErrInvalidUserID
	}
	return UserID{value: value}, nil
}

func (u UserID) String() string {
	return u.value
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}
