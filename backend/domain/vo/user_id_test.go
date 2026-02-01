package vo_test

import (
	"testing"

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

func TestReconstructUserID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	got := vo.ReconstructUserID(validUUID)

	if got.String() != validUUID {
		t.Errorf("ReconstructUserID(%q).String() = %v, want %v", validUUID, got.String(), validUUID)
	}
}

func TestUserID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1 := vo.ReconstructUserID(validUUID)
	id2 := vo.ReconstructUserID(validUUID)
	id3 := vo.NewUserID()

	tests := []struct {
		name string
		id1  vo.UserID
		id2  vo.UserID
		want bool
	}{
		{"同じ値はtrue", id1, id2, true},
		{"異なる値はfalse", id1, id3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id1.Equals(tt.id2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
