package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	userID := vo.NewUserID()

	if userID.String() == "" {
		t.Error("NewUserID() should return non-empty string")
	}
	if _, err := uuid.Parse(userID.String()); err != nil {
		t.Errorf("NewUserID() should return valid UUID, got: %s", userID.String())
	}
}

func TestParseUserID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name    string
		input   string
		wantID  string
		wantErr error
	}{
		// 正常系
		{"valid UUID", validUUID, validUUID, nil},
		// 異常系
		{"empty string", "", "", domainErrors.ErrInvalidUserID},
		{"invalid format", "invalid", "", domainErrors.ErrInvalidUserID},
		{"partial UUID", "550e8400-e29b", "", domainErrors.ErrInvalidUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.ParseUserID(tt.input)

			if err != tt.wantErr {
				t.Errorf("ParseUserID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantID {
				t.Errorf("ParseUserID(%q).String() = %v, want %v", tt.input, got.String(), tt.wantID)
			}
		})
	}
}

func TestUserID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1, _ := vo.ParseUserID(validUUID)
	id2, _ := vo.ParseUserID(validUUID)
	id3 := vo.NewUserID()

	tests := []struct {
		name string
		id1  vo.UserID
		id2  vo.UserID
		want bool
	}{
		{"same value", id1, id2, true},
		{"different value", id1, id3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id1.Equals(tt.id2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
