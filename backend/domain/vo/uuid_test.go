package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"

	"github.com/google/uuid"
)

func TestNewUUID(t *testing.T) {
	t.Run("新しいUUIDを生成できる", func(t *testing.T) {
		u := vo.NewUUID()

		if u.String() == "" {
			t.Error("NewUUID() should return non-empty string")
		}
		if _, err := uuid.Parse(u.String()); err != nil {
			t.Errorf("NewUUID() should return valid UUID, got: %s", u.String())
		}
	})

	t.Run("毎回異なるUUIDを生成する", func(t *testing.T) {
		u1 := vo.NewUUID()
		u2 := vo.NewUUID()

		if u1.Equals(u2) {
			t.Error("NewUUID() should generate different UUIDs")
		}
	})
}

func TestParseUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "正常系_有効なUUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: nil,
		},
		{
			name:    "異常系_空文字",
			input:   "",
			wantErr: domainErrors.ErrUUIDRequired,
		},
		{
			name:    "異常系_無効なフォーマット",
			input:   "invalid-uuid",
			wantErr: domainErrors.ErrInvalidUUIDFormat,
		},
		{
			name:    "異常系_短すぎる",
			input:   "550e8400-e29b-41d4",
			wantErr: domainErrors.ErrInvalidUUIDFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.ParseUUID(tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("ParseUUID() error = %v, want %v", err, tt.wantErr)
				}
				if !got.IsZero() {
					t.Error("ParseUUID() should return zero UUID on error")
				}
			} else {
				if err != nil {
					t.Fatalf("ParseUUID() unexpected error = %v", err)
				}
				if got.String() != tt.input {
					t.Errorf("ParseUUID().String() = %v, want %v", got.String(), tt.input)
				}
			}
		})
	}
}

func TestReconstructUUID(t *testing.T) {
	t.Run("正常系_有効なUUID", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"

		got := vo.ReconstructUUID(validUUID)

		if got.String() != validUUID {
			t.Errorf("ReconstructUUID(%q).String() = %v, want %v", validUUID, got.String(), validUUID)
		}
	})

	t.Run("DB復元時はエラーにならない", func(t *testing.T) {
		// DB復元はバリデーションしないのでエラーにならない
		invalidUUID := "invalid-uuid"

		got := vo.ReconstructUUID(invalidUUID)

		// 無効なUUIDはゼロ値になる
		if !got.IsZero() {
			t.Error("ReconstructUUID() with invalid UUID should return zero UUID")
		}
	})
}

func TestUUID_IsZero(t *testing.T) {
	tests := []struct {
		name string
		uuid vo.UUID
		want bool
	}{
		{
			name: "ゼロ値のUUID",
			uuid: vo.ReconstructUUID("00000000-0000-0000-0000-000000000000"),
			want: true,
		},
		{
			name: "有効なUUID",
			uuid: vo.ReconstructUUID("550e8400-e29b-41d4-a716-446655440000"),
			want: false,
		},
		{
			name: "新規生成したUUID",
			uuid: vo.NewUUID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	u1 := vo.ReconstructUUID(validUUID)
	u2 := vo.ReconstructUUID(validUUID)
	u3 := vo.NewUUID()

	tests := []struct {
		name string
		u1   vo.UUID
		u2   vo.UUID
		want bool
	}{
		{
			name: "同じ値はtrue",
			u1:   u1,
			u2:   u2,
			want: true,
		},
		{
			name: "異なる値はfalse",
			u1:   u1,
			u2:   u3,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u1.Equals(tt.u2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
