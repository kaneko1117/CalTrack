package vo

import (
	domainErrors "caltrack/domain/errors"
)

const maxNicknameLength = 50

type Nickname struct {
	value string
}

func NewNickname(value string) (Nickname, error) {
	if value == "" {
		return Nickname{}, domainErrors.ErrNicknameRequired
	}
	if len(value) > maxNicknameLength {
		return Nickname{}, domainErrors.ErrNicknameTooLong
	}
	return Nickname{value: value}, nil
}

func (n Nickname) String() string {
	return n.value
}
