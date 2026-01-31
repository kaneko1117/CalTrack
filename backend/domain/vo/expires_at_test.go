package vo

import (
	"testing"
	"time"

	domainErrors "caltrack/domain/errors"
)

func TestNewExpiresAt(t *testing.T) {
	// 現在時刻を固定
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	expiresAt := NewExpiresAt()

	// 7日後であることを確認
	expected := fixedNow.AddDate(0, 0, 7)
	if !expiresAt.Time().Equal(expected) {
		t.Errorf("NewExpiresAt().Time() = %v, want %v", expiresAt.Time(), expected)
	}
}

func TestParseExpiresAt(t *testing.T) {
	input := time.Date(2024, 6, 22, 12, 0, 0, 0, time.UTC)
	expiresAt := ParseExpiresAt(input)

	if !expiresAt.Time().Equal(input) {
		t.Errorf("ParseExpiresAt(%v).Time() = %v, want %v", input, expiresAt.Time(), input)
	}
}

func TestExpiresAt_IsExpired(t *testing.T) {
	// 現在時刻を固定
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		// 有効期限内
		{"future", fixedNow.Add(1 * time.Hour), false},
		{"exactly now", fixedNow, false},
		// 有効期限切れ
		{"past by 1 second", fixedNow.Add(-1 * time.Second), true},
		{"past by 1 day", fixedNow.AddDate(0, 0, -1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseExpiresAt(tt.expiresAt)
			if got := e.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpiresAt_ValidateNotExpired(t *testing.T) {
	// 現在時刻を固定
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	tests := []struct {
		name      string
		expiresAt time.Time
		wantErr   error
	}{
		// 有効期限内
		{"valid", fixedNow.Add(1 * time.Hour), nil},
		// 有効期限切れ
		{"expired", fixedNow.Add(-1 * time.Second), domainErrors.ErrSessionExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseExpiresAt(tt.expiresAt)
			if err := e.ValidateNotExpired(); err != tt.wantErr {
				t.Errorf("ValidateNotExpired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
