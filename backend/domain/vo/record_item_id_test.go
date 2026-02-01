package vo_test

import (
	"testing"

	"caltrack/domain/vo"

	"github.com/google/uuid"
)

func TestNewRecordItemID(t *testing.T) {
	recordItemID := vo.NewRecordItemID()

	if recordItemID.String() == "" {
		t.Error("NewRecordItemID() should return non-empty string")
	}
	if _, err := uuid.Parse(recordItemID.String()); err != nil {
		t.Errorf("NewRecordItemID() should return valid UUID, got: %s", recordItemID.String())
	}
}

func TestReconstructRecordItemID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("DBからRecordItemIDを復元できる", func(t *testing.T) {
		got := vo.ReconstructRecordItemID(validUUID)

		if got.String() != validUUID {
			t.Errorf("ReconstructRecordItemID(%q).String() = %v, want %v", validUUID, got.String(), validUUID)
		}
	})
}

func TestRecordItemID_Equals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1 := vo.ReconstructRecordItemID(validUUID)
	id2 := vo.ReconstructRecordItemID(validUUID)
	id3 := vo.NewRecordItemID()

	tests := []struct {
		name string
		id1  vo.RecordItemID
		id2  vo.RecordItemID
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
