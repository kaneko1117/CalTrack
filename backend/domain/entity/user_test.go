package entity_test

import (
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
)

func TestNewUser_Success(t *testing.T) {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user, errs := entity.NewUser(
		"test@example.com",
		"$2a$10$hashedpassword",
		"testuser",
		70.5,
		175.0,
		birthDate,
		"male",
		"moderate",
	)

	if errs != nil {
		t.Fatalf("NewUser() unexpected errors: %v", errs)
	}
	if user.ID().String() == "" {
		t.Error("ID should not be empty")
	}
	if user.Email().String() != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", user.Email().String())
	}
	if user.Nickname().String() != "testuser" {
		t.Errorf("Nickname = %v, want testuser", user.Nickname().String())
	}
	if user.Weight().Kg() != 70.5 {
		t.Errorf("Weight = %v, want 70.5", user.Weight().Kg())
	}
	if user.Height().Cm() != 175.0 {
		t.Errorf("Height = %v, want 175.0", user.Height().Cm())
	}
	if user.Gender().String() != "male" {
		t.Errorf("Gender = %v, want male", user.Gender().String())
	}
	if user.ActivityLevel().String() != "moderate" {
		t.Errorf("ActivityLevel = %v, want moderate", user.ActivityLevel().String())
	}
}

func TestNewUser_ValidationErrors(t *testing.T) {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"無効なメールアドレス形式", "invalid", domainErrors.ErrInvalidEmailFormat},
		{"空のメールアドレス", "", domainErrors.ErrEmailRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := entity.NewUser(tt.email, "$2a$10$hash", "nick", 70, 170, birthDate, "male", "moderate")
			if len(errs) != 1 || errs[0] != tt.wantErr {
				t.Errorf("got errs = %v, want [%v]", errs, tt.wantErr)
			}
		})
	}
}

func TestNewUser_MultipleErrors(t *testing.T) {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	_, errs := entity.NewUser("", "$2a$10$hash", "", -1, 400, birthDate, "invalid", "unknown")

	if len(errs) != 6 {
		t.Fatalf("got %d errors, want 6: %v", len(errs), errs)
	}
}

func TestNewUser_GeneratesUniqueIDs(t *testing.T) {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user1, _ := entity.NewUser("test@example.com", "$2a$10$hash", "nick", 70, 170, birthDate, "male", "moderate")
	user2, _ := entity.NewUser("test@example.com", "$2a$10$hash", "nick", 70, 170, birthDate, "male", "moderate")

	if user1.ID().Equals(user2.ID()) {
		t.Error("should generate unique IDs")
	}
}

func TestReconstructUser_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user, err := entity.ReconstructUser(
		userID,
		"test@example.com",
		"$2a$10$hashedpassword",
		"testuser",
		70.5,
		175.0,
		birthDate,
		"male",
		"moderate",
		createdAt,
		updatedAt,
	)

	if err != nil {
		t.Fatalf("ReconstructUser() unexpected error: %v", err)
	}
	if user.ID().String() != userID {
		t.Errorf("ID = %v, want %v", user.ID().String(), userID)
	}
	if user.Email().String() != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", user.Email().String())
	}
	if !user.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt(), createdAt)
	}
	if !user.UpdatedAt().Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", user.UpdatedAt(), updatedAt)
	}
}
