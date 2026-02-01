package vo_test

import (
	"testing"

	"caltrack/domain/vo"

	"github.com/google/uuid"
)

func TestNewRecordID(t *testing.T) {
	recordID := vo.NewRecordID()

	if recordID.String() == "" {
		t.Error("NewRecordID() should return non-empty string")
	}
	if _, err := uuid.Parse(recordID.String()); err != nil {
		t.Errorf("NewRecordID() should return valid UUID, got: %s", recordID.String())
	}
}

func TestReconstructRecordID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("DBからRecordIDを復元できる", func(t *testing.T) {
		got := vo.ReconstructRecordID(validUUID)

		if got.String() != validUUID {
			t.Errorf("ReconstructRecordID(%q).String() = %v, want %v", validUUID, got.String(), validUUID)
		}
	})
}

func TestRecordID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1 := vo.ReconstructRecordID(validUUID)
	id2 := vo.ReconstructRecordID(validUUID)
	id3 := vo.NewRecordID()

	tests := []struct {
		name string
		id1  vo.RecordID
		id2  vo.RecordID
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
