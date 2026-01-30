package vo

import (
	"golang.org/x/crypto/bcrypt"

	domainErrors "caltrack/domain/errors"
)

const minPasswordLength = 8

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {
	if value == "" {
		return Password{}, domainErrors.ErrPasswordRequired
	}
	if len(value) < minPasswordLength {
		return Password{}, domainErrors.ErrPasswordTooShort
	}
	return Password{value: value}, nil
}

func (p Password) Hash() (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, err
	}
	return HashedPassword{value: string(hash)}, nil
}

type HashedPassword struct {
	value string
}

func NewHashedPassword(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func (h HashedPassword) String() string {
	return h.value
}

func (h HashedPassword) Compare(password Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h.value), []byte(password.value))
	return err == nil
}
