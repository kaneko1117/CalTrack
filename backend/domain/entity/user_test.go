package entity_test

import (
	"testing"
	"time"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/entity"
	"caltrack/domain/vo"
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
		{"invalid email", "invalid", domainErrors.ErrInvalidEmailFormat},
		{"empty email", "", domainErrors.ErrEmailRequired},
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

func TestReconstructUser(t *testing.T) {
	userID := vo.NewUserID()
	email, _ := vo.NewEmail("test@example.com")
	hashedPassword := vo.NewHashedPassword("$2a$10$hashedpassword")
	nickname, _ := vo.NewNickname("testuser")
	weight, _ := vo.NewWeight(70.5)
	height, _ := vo.NewHeight(175.0)
	birthDate, _ := vo.NewBirthDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	gender, _ := vo.NewGender("male")
	activityLevel, _ := vo.NewActivityLevel("moderate")
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	user := entity.ReconstructUser(
		userID, email, hashedPassword, nickname,
		weight, height, birthDate, gender, activityLevel,
		createdAt, updatedAt,
	)

	if !user.ID().Equals(userID) {
		t.Errorf("ID mismatch")
	}
	if !user.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt(), createdAt)
	}
	if !user.UpdatedAt().Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", user.UpdatedAt(), updatedAt)
	}
}
