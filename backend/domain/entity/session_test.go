package entity_test

import (
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewSession_Success(t *testing.T) {
	validUserID := "550e8400-e29b-41d4-a716-446655440000"

	session, errs := entity.NewSession(validUserID)

	if errs != nil {
		t.Fatalf("NewSession() unexpected errors: %v", errs)
	}
	if session.ID().String() == "" {
		t.Error("ID should not be empty")
	}
	if session.UserID().String() != validUserID {
		t.Errorf("UserID = %v, want %v", session.UserID().String(), validUserID)
	}
	if session.ExpiresAt().Time().IsZero() {
		t.Error("ExpiresAt should not be zero")
	}
	if session.CreatedAt().IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestNewSession_InvalidUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr error
	}{
		{"空のユーザーIDはエラー", "", domainErrors.ErrInvalidUserID},
		{"無効なUUID形式はエラー", "invalid-uuid", domainErrors.ErrInvalidUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := entity.NewSession(tt.userID)

			if len(errs) != 1 {
				t.Fatalf("got %d errors, want 1: %v", len(errs), errs)
			}
			if errs[0] != tt.wantErr {
				t.Errorf("got err = %v, want %v", errs[0], tt.wantErr)
			}
		})
	}
}

func TestNewSession_GeneratesUniqueIDs(t *testing.T) {
	validUserID := "550e8400-e29b-41d4-a716-446655440000"

	session1, _ := entity.NewSession(validUserID)
	session2, _ := entity.NewSession(validUserID)

	if session1.ID().Equals(session2.ID()) {
		t.Error("should generate unique session IDs")
	}
}

func TestNewSessionWithUserID_Success(t *testing.T) {
	userID, _ := vo.ParseUserID("550e8400-e29b-41d4-a716-446655440000")

	session, err := entity.NewSessionWithUserID(userID)

	if err != nil {
		t.Fatalf("NewSessionWithUserID() unexpected error: %v", err)
	}
	if session.ID().String() == "" {
		t.Error("ID should not be empty")
	}
	if !session.UserID().Equals(userID) {
		t.Errorf("UserID = %v, want %v", session.UserID().String(), userID.String())
	}
	if session.ExpiresAt().Time().IsZero() {
		t.Error("ExpiresAt should not be zero")
	}
	if session.CreatedAt().IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestNewSessionWithUserID_GeneratesUniqueIDs(t *testing.T) {
	userID, _ := vo.ParseUserID("550e8400-e29b-41d4-a716-446655440000")

	session1, _ := entity.NewSessionWithUserID(userID)
	session2, _ := entity.NewSessionWithUserID(userID)

	if session1.ID().Equals(session2.ID()) {
		t.Error("should generate unique session IDs")
	}
}

func TestReconstructSession_Success(t *testing.T) {
	// 有効なセッションIDを生成してそのString表現を使用
	validSessionID, _ := vo.NewSessionID()
	sessionIDStr := validSessionID.String()
	userIDStr := "550e8400-e29b-41d4-a716-446655440000"
	expiresAt := time.Date(2024, 6, 22, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	session, err := entity.ReconstructSession(
		sessionIDStr,
		userIDStr,
		expiresAt,
		createdAt,
	)

	if err != nil {
		t.Fatalf("ReconstructSession() unexpected error: %v", err)
	}
	if session.ID().String() != sessionIDStr {
		t.Errorf("ID = %v, want %v", session.ID().String(), sessionIDStr)
	}
	if session.UserID().String() != userIDStr {
		t.Errorf("UserID = %v, want %v", session.UserID().String(), userIDStr)
	}
	if !session.ExpiresAt().Time().Equal(expiresAt) {
		t.Errorf("ExpiresAt = %v, want %v", session.ExpiresAt().Time(), expiresAt)
	}
	if !session.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", session.CreatedAt(), createdAt)
	}
}

func TestReconstructSession_InvalidSessionID(t *testing.T) {
	userIDStr := "550e8400-e29b-41d4-a716-446655440000"
	expiresAt := time.Date(2024, 6, 22, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		sessionID string
	}{
		{"空のセッションIDはエラー", ""},
		{"無効なbase64はエラー", "not-valid-base64!!!"},
		{"短すぎるbase64はエラー", "YWJjZA=="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.ReconstructSession(
				tt.sessionID,
				userIDStr,
				expiresAt,
				createdAt,
			)

			if err != domainErrors.ErrInvalidSessionID {
				t.Errorf("got err = %v, want %v", err, domainErrors.ErrInvalidSessionID)
			}
		})
	}
}

func TestReconstructSession_InvalidUserID(t *testing.T) {
	validSessionID, _ := vo.NewSessionID()
	sessionIDStr := validSessionID.String()
	expiresAt := time.Date(2024, 6, 22, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		userID string
	}{
		{"空のユーザーIDはエラー", ""},
		{"無効なUUID形式はエラー", "invalid-uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.ReconstructSession(
				sessionIDStr,
				tt.userID,
				expiresAt,
				createdAt,
			)

			if err != domainErrors.ErrInvalidUserID {
				t.Errorf("got err = %v, want %v", err, domainErrors.ErrInvalidUserID)
			}
		})
	}
}

func TestSession_IsExpired(t *testing.T) {
	validSessionID, _ := vo.NewSessionID()
	sessionIDStr := validSessionID.String()
	userIDStr := "550e8400-e29b-41d4-a716-446655440000"
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		// 有効期限内（未来）
		{"期限内は期限切れでない", time.Now().Add(24 * time.Hour), false},
		// 有効期限切れ（過去）
		{"期限切れは期限切れ", time.Now().Add(-24 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, _ := entity.ReconstructSession(
				sessionIDStr,
				userIDStr,
				tt.expiresAt,
				createdAt,
			)

			if got := session.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_ValidateNotExpired(t *testing.T) {
	validSessionID, _ := vo.NewSessionID()
	sessionIDStr := validSessionID.String()
	userIDStr := "550e8400-e29b-41d4-a716-446655440000"
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		expiresAt time.Time
		wantErr   error
	}{
		// 有効期限内（未来）
		{"有効期限内はエラーなし", time.Now().Add(24 * time.Hour), nil},
		// 有効期限切れ（過去）
		{"有効期限切れはエラー", time.Now().Add(-24 * time.Hour), domainErrors.ErrSessionExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, _ := entity.ReconstructSession(
				sessionIDStr,
				userIDStr,
				tt.expiresAt,
				createdAt,
			)

			if err := session.ValidateNotExpired(); err != tt.wantErr {
				t.Errorf("ValidateNotExpired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
