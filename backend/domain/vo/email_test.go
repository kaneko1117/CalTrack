package vo_test

import (
	"strings"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewEmail_ValidEmail_ReturnsEmail(t *testing.T) {
	email, err := vo.NewEmail("user@example.com")

	if err != nil {
		t.Errorf("NewEmail should not return error for valid email, got: %v", err)
	}
	if email.String() != "user@example.com" {
		t.Errorf("Email.String() should return the email, got: %s", email.String())
	}
}

func TestNewEmail_ValidEmailWithSubdomain_ReturnsEmail(t *testing.T) {
	email, err := vo.NewEmail("user@mail.example.com")

	if err != nil {
		t.Errorf("NewEmail should not return error for valid email with subdomain, got: %v", err)
	}
	if email.String() != "user@mail.example.com" {
		t.Errorf("Email.String() should return the email, got: %s", email.String())
	}
}

func TestNewEmail_EmptyString_ReturnsError(t *testing.T) {
	_, err := vo.NewEmail("")

	if err == nil {
		t.Error("NewEmail should return error for empty string")
	}
	if err != domainErrors.ErrEmailRequired {
		t.Errorf("NewEmail should return ErrEmailRequired, got: %v", err)
	}
}

func TestNewEmail_NoAtSign_ReturnsError(t *testing.T) {
	_, err := vo.NewEmail("userexample.com")

	if err == nil {
		t.Error("NewEmail should return error for email without @")
	}
	if err != domainErrors.ErrInvalidEmailFormat {
		t.Errorf("NewEmail should return ErrInvalidEmailFormat, got: %v", err)
	}
}

func TestNewEmail_NoDomain_ReturnsError(t *testing.T) {
	_, err := vo.NewEmail("user@")

	if err == nil {
		t.Error("NewEmail should return error for email without domain")
	}
	if err != domainErrors.ErrInvalidEmailFormat {
		t.Errorf("NewEmail should return ErrInvalidEmailFormat, got: %v", err)
	}
}

func TestNewEmail_NoLocalPart_ReturnsError(t *testing.T) {
	_, err := vo.NewEmail("@example.com")

	if err == nil {
		t.Error("NewEmail should return error for email without local part")
	}
	if err != domainErrors.ErrInvalidEmailFormat {
		t.Errorf("NewEmail should return ErrInvalidEmailFormat, got: %v", err)
	}
}

func TestNewEmail_MaxLength254_ReturnsEmail(t *testing.T) {
	// Create email with exactly 254 characters
	// local@domain.com format: local (max 64) + @ + domain
	local := strings.Repeat("a", 64)
	domain := strings.Repeat("b", 254-64-1-4) + ".com" // 254 - 64 - 1(@) - 4(.com)
	email254 := local + "@" + domain

	if len(email254) != 254 {
		t.Fatalf("Test email should be 254 chars, got %d", len(email254))
	}

	_, err := vo.NewEmail(email254)

	if err != nil {
		t.Errorf("NewEmail should accept 254 character email, got: %v", err)
	}
}

func TestNewEmail_Exceeds254_ReturnsError(t *testing.T) {
	// Create email with 255 characters
	local := strings.Repeat("a", 64)
	domain := strings.Repeat("b", 255-64-1-4) + ".com"
	email255 := local + "@" + domain

	if len(email255) != 255 {
		t.Fatalf("Test email should be 255 chars, got %d", len(email255))
	}

	_, err := vo.NewEmail(email255)

	if err == nil {
		t.Error("NewEmail should return error for email exceeding 254 characters")
	}
	if err != domainErrors.ErrEmailTooLong {
		t.Errorf("NewEmail should return ErrEmailTooLong, got: %v", err)
	}
}

func TestEmail_Equals_SameValue_ReturnsTrue(t *testing.T) {
	email1, _ := vo.NewEmail("user@example.com")
	email2, _ := vo.NewEmail("user@example.com")

	if !email1.Equals(email2) {
		t.Error("Equals should return true for same email value")
	}
}

func TestEmail_Equals_DifferentValue_ReturnsFalse(t *testing.T) {
	email1, _ := vo.NewEmail("user1@example.com")
	email2, _ := vo.NewEmail("user2@example.com")

	if email1.Equals(email2) {
		t.Error("Equals should return false for different email values")
	}
}
