package vo_test

import (
	"testing"

	"github.com/google/uuid"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewUserID_ReturnsValidUUID(t *testing.T) {
	userID := vo.NewUserID()

	if userID.String() == "" {
		t.Error("NewUserID should return non-empty string")
	}

	if _, err := uuid.Parse(userID.String()); err != nil {
		t.Errorf("NewUserID should return valid UUID, got: %s", userID.String())
	}
}

func TestParseUserID_ValidUUID_ReturnsUserID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	userID, err := vo.ParseUserID(validUUID)

	if err != nil {
		t.Errorf("ParseUserID should not return error for valid UUID, got: %v", err)
	}
	if userID.String() != validUUID {
		t.Errorf("ParseUserID should return UserID with same value, expected: %s, got: %s", validUUID, userID.String())
	}
}

func TestParseUserID_EmptyString_ReturnsError(t *testing.T) {
	_, err := vo.ParseUserID("")

	if err == nil {
		t.Error("ParseUserID should return error for empty string")
	}
	if err != domainErrors.ErrInvalidUserID {
		t.Errorf("ParseUserID should return ErrInvalidUserID, got: %v", err)
	}
}

func TestParseUserID_InvalidFormat_ReturnsError(t *testing.T) {
	_, err := vo.ParseUserID("invalid")

	if err == nil {
		t.Error("ParseUserID should return error for invalid format")
	}
	if err != domainErrors.ErrInvalidUserID {
		t.Errorf("ParseUserID should return ErrInvalidUserID, got: %v", err)
	}
}

func TestUserID_Equals_SameValue_ReturnsTrue(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	userID1, _ := vo.ParseUserID(validUUID)
	userID2, _ := vo.ParseUserID(validUUID)

	if !userID1.Equals(userID2) {
		t.Error("Equals should return true for same UUID value")
	}
}

func TestUserID_Equals_DifferentValue_ReturnsFalse(t *testing.T) {
	userID1 := vo.NewUserID()
	userID2 := vo.NewUserID()

	if userID1.Equals(userID2) {
		t.Error("Equals should return false for different UUID values")
	}
}
