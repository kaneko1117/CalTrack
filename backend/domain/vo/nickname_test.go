package vo_test

import (
	"strings"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewNickname_ValidNickname_ReturnsNickname(t *testing.T) {
	nickname, err := vo.NewNickname("John")

	if err != nil {
		t.Errorf("NewNickname should not return error for valid nickname, got: %v", err)
	}
	if nickname.String() != "John" {
		t.Errorf("Nickname.String() should return the nickname, got: %s", nickname.String())
	}
}

func TestNewNickname_EmptyString_ReturnsError(t *testing.T) {
	_, err := vo.NewNickname("")

	if err == nil {
		t.Error("NewNickname should return error for empty string")
	}
	if err != domainErrors.ErrNicknameRequired {
		t.Errorf("NewNickname should return ErrNicknameRequired, got: %v", err)
	}
}

func TestNewNickname_TooLong51_ReturnsError(t *testing.T) {
	longNickname := strings.Repeat("a", 51)
	_, err := vo.NewNickname(longNickname)

	if err == nil {
		t.Error("NewNickname should return error for 51 char nickname")
	}
	if err != domainErrors.ErrNicknameTooLong {
		t.Errorf("NewNickname should return ErrNicknameTooLong, got: %v", err)
	}
}

func TestNewNickname_MinLength1_ReturnsNickname(t *testing.T) {
	_, err := vo.NewNickname("A")

	if err != nil {
		t.Errorf("NewNickname should accept 1 char nickname, got: %v", err)
	}
}

func TestNewNickname_MaxLength50_ReturnsNickname(t *testing.T) {
	nickname50 := strings.Repeat("a", 50)
	_, err := vo.NewNickname(nickname50)

	if err != nil {
		t.Errorf("NewNickname should accept 50 char nickname, got: %v", err)
	}
}

func TestNewNickname_Exceeds50_ReturnsError(t *testing.T) {
	nickname51 := strings.Repeat("a", 51)
	_, err := vo.NewNickname(nickname51)

	if err == nil {
		t.Error("NewNickname should reject 51 char nickname")
	}
}
