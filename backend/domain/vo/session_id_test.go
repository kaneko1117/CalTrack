package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewSessionID(t *testing.T) {
	sessionID, err := vo.NewSessionID()

	if err != nil {
		t.Fatalf("NewSessionID() error = %v, want nil", err)
	}
	if sessionID.String() == "" {
		t.Error("NewSessionID() should return non-empty string")
	}
}

func TestNewSessionID_Uniqueness(t *testing.T) {
	// 複数回生成して重複しないことを確認
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := vo.NewSessionID()
		if err != nil {
			t.Fatalf("NewSessionID() error = %v", err)
		}
		if ids[id.String()] {
			t.Errorf("NewSessionID() generated duplicate ID: %s", id.String())
		}
		ids[id.String()] = true
	}
}

func TestParseSessionID(t *testing.T) {
	// 有効なセッションIDを生成
	validID, _ := vo.NewSessionID()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"有効なセッションID", validID.String(), nil},
		// 異常系
		{"空文字はエラー", "", domainErrors.ErrInvalidSessionID},
		{"無効なbase64はエラー", "not-valid-base64!!!", domainErrors.ErrInvalidSessionID},
		{"短すぎるとエラー", "YWJjZA==", domainErrors.ErrInvalidSessionID}, // "abcd" in base64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.ParseSessionID(tt.input)

			if err != tt.wantErr {
				t.Errorf("ParseSessionID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.input {
				t.Errorf("ParseSessionID(%q).String() = %v, want %v", tt.input, got.String(), tt.input)
			}
		})
	}
}

func TestSessionID_Equals(t *testing.T) {
	id1, _ := vo.NewSessionID()
	id2, _ := vo.ParseSessionID(id1.String())
	id3, _ := vo.NewSessionID()

	tests := []struct {
		name string
		id1  vo.SessionID
		id2  vo.SessionID
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
