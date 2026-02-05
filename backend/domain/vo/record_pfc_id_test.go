package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"

	"github.com/google/uuid"
)

func TestNewRecordPfcID(t *testing.T) {
	t.Run("新しいRecordPfcIDを生成できる", func(t *testing.T) {
		id := vo.NewRecordPfcID()

		if id.String() == "" {
			t.Error("NewRecordPfcID() should return non-empty string")
		}
		if _, err := uuid.Parse(id.String()); err != nil {
			t.Errorf("NewRecordPfcID() should return valid UUID, got: %s", id.String())
		}
	})

	t.Run("毎回異なるRecordPfcIDを生成する", func(t *testing.T) {
		id1 := vo.NewRecordPfcID()
		id2 := vo.NewRecordPfcID()

		if id1.Equals(id2) {
			t.Error("NewRecordPfcID() should generate different IDs")
		}
	})
}

func TestParseRecordPfcID(t *testing.T) {
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
			wantErr: domainErrors.ErrInvalidRecordPfcID,
		},
		{
			name:    "異常系_無効なフォーマット",
			input:   "invalid-uuid",
			wantErr: domainErrors.ErrInvalidRecordPfcID,
		},
		{
			name:    "異常系_短すぎる",
			input:   "550e8400-e29b-41d4",
			wantErr: domainErrors.ErrInvalidRecordPfcID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.ParseRecordPfcID(tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("ParseRecordPfcID() error = %v, want %v", err, tt.wantErr)
				}
				if !got.IsZero() {
					t.Error("ParseRecordPfcID() should return zero ID on error")
				}
			} else {
				if err != nil {
					t.Fatalf("ParseRecordPfcID() unexpected error = %v", err)
				}
				if got.String() != tt.input {
					t.Errorf("ParseRecordPfcID().String() = %v, want %v", got.String(), tt.input)
				}
			}
		})
	}
}

func TestReconstructRecordPfcID(t *testing.T) {
	t.Run("DBからRecordPfcIDを復元できる", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"

		got := vo.ReconstructRecordPfcID(validUUID)

		if got.String() != validUUID {
			t.Errorf("ReconstructRecordPfcID(%q).String() = %v, want %v", validUUID, got.String(), validUUID)
		}
	})
}

func TestRecordPfcID_IsZero(t *testing.T) {
	tests := []struct {
		name string
		id   vo.RecordPfcID
		want bool
	}{
		{
			name: "ゼロ値のRecordPfcID",
			id:   vo.ReconstructRecordPfcID("00000000-0000-0000-0000-000000000000"),
			want: true,
		},
		{
			name: "有効なRecordPfcID",
			id:   vo.ReconstructRecordPfcID("550e8400-e29b-41d4-a716-446655440000"),
			want: false,
		},
		{
			name: "新規生成したRecordPfcID",
			id:   vo.NewRecordPfcID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordPfcID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1 := vo.ReconstructRecordPfcID(validUUID)
	id2 := vo.ReconstructRecordPfcID(validUUID)
	id3 := vo.NewRecordPfcID()

	tests := []struct {
		name string
		id1  vo.RecordPfcID
		id2  vo.RecordPfcID
		want bool
	}{
		{
			name: "同じ値はtrue",
			id1:  id1,
			id2:  id2,
			want: true,
		},
		{
			name: "異なる値はfalse",
			id1:  id1,
			id2:  id3,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id1.Equals(tt.id2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
