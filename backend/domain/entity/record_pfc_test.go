package entity_test

import (
	"testing"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

func TestNewRecordPfc(t *testing.T) {
	t.Run("正常系_有効なパラメータでRecordPfc作成", func(t *testing.T) {
		recordID := vo.NewRecordID()
		protein := 20.5
		fat := 10.3
		carbs := 50.8

		recordPfc := entity.NewRecordPfc(recordID, protein, fat, carbs)

		if recordPfc == nil {
			t.Fatal("NewRecordPfc() returned nil")
		}
		if !recordPfc.RecordID().Equals(recordID) {
			t.Errorf("NewRecordPfc().RecordID() = %v, want %v", recordPfc.RecordID(), recordID)
		}
		if recordPfc.ID().String() == "" {
			t.Error("NewRecordPfc().ID() should not be empty")
		}
		if recordPfc.Protein() != protein {
			t.Errorf("NewRecordPfc().Protein() = %v, want %v", recordPfc.Protein(), protein)
		}
		if recordPfc.Fat() != fat {
			t.Errorf("NewRecordPfc().Fat() = %v, want %v", recordPfc.Fat(), fat)
		}
		if recordPfc.Carbs() != carbs {
			t.Errorf("NewRecordPfc().Carbs() = %v, want %v", recordPfc.Carbs(), carbs)
		}
	})

	t.Run("正常系_0の値も許容される", func(t *testing.T) {
		recordID := vo.NewRecordID()
		protein := 0.0
		fat := 0.0
		carbs := 0.0

		recordPfc := entity.NewRecordPfc(recordID, protein, fat, carbs)

		if recordPfc == nil {
			t.Fatal("NewRecordPfc() returned nil")
		}
		if recordPfc.Protein() != 0.0 {
			t.Errorf("NewRecordPfc().Protein() = %v, want 0.0", recordPfc.Protein())
		}
		if recordPfc.Fat() != 0.0 {
			t.Errorf("NewRecordPfc().Fat() = %v, want 0.0", recordPfc.Fat())
		}
		if recordPfc.Carbs() != 0.0 {
			t.Errorf("NewRecordPfc().Carbs() = %v, want 0.0", recordPfc.Carbs())
		}
	})

	t.Run("正常系_小数点以下も正しく保持される", func(t *testing.T) {
		recordID := vo.NewRecordID()
		protein := 20.123
		fat := 10.456
		carbs := 50.789

		recordPfc := entity.NewRecordPfc(recordID, protein, fat, carbs)

		if recordPfc.Protein() != protein {
			t.Errorf("NewRecordPfc().Protein() = %v, want %v", recordPfc.Protein(), protein)
		}
		if recordPfc.Fat() != fat {
			t.Errorf("NewRecordPfc().Fat() = %v, want %v", recordPfc.Fat(), fat)
		}
		if recordPfc.Carbs() != carbs {
			t.Errorf("NewRecordPfc().Carbs() = %v, want %v", recordPfc.Carbs(), carbs)
		}
	})
}

func TestReconstructRecordPfc(t *testing.T) {
	t.Run("DB復元", func(t *testing.T) {
		idStr := "550e8400-e29b-41d4-a716-446655440000"
		recordIDStr := "660e8400-e29b-41d4-a716-446655440001"
		protein := 20.5
		fat := 10.3
		carbs := 50.8

		recordPfc := entity.ReconstructRecordPfc(idStr, recordIDStr, protein, fat, carbs)

		if recordPfc.ID().String() != idStr {
			t.Errorf("ReconstructRecordPfc().ID() = %v, want %v", recordPfc.ID().String(), idStr)
		}
		if recordPfc.RecordID().String() != recordIDStr {
			t.Errorf("ReconstructRecordPfc().RecordID() = %v, want %v", recordPfc.RecordID().String(), recordIDStr)
		}
		if recordPfc.Protein() != protein {
			t.Errorf("ReconstructRecordPfc().Protein() = %v, want %v", recordPfc.Protein(), protein)
		}
		if recordPfc.Fat() != fat {
			t.Errorf("ReconstructRecordPfc().Fat() = %v, want %v", recordPfc.Fat(), fat)
		}
		if recordPfc.Carbs() != carbs {
			t.Errorf("ReconstructRecordPfc().Carbs() = %v, want %v", recordPfc.Carbs(), carbs)
		}
	})

	t.Run("DB復元_0の値も正しく復元される", func(t *testing.T) {
		idStr := "550e8400-e29b-41d4-a716-446655440000"
		recordIDStr := "660e8400-e29b-41d4-a716-446655440001"
		protein := 0.0
		fat := 0.0
		carbs := 0.0

		recordPfc := entity.ReconstructRecordPfc(idStr, recordIDStr, protein, fat, carbs)

		if recordPfc.Protein() != 0.0 {
			t.Errorf("ReconstructRecordPfc().Protein() = %v, want 0.0", recordPfc.Protein())
		}
		if recordPfc.Fat() != 0.0 {
			t.Errorf("ReconstructRecordPfc().Fat() = %v, want 0.0", recordPfc.Fat())
		}
		if recordPfc.Carbs() != 0.0 {
			t.Errorf("ReconstructRecordPfc().Carbs() = %v, want 0.0", recordPfc.Carbs())
		}
	})
}

func TestRecordPfc_Pfc(t *testing.T) {
	t.Run("Pfc VOを取得できる", func(t *testing.T) {
		recordID := vo.NewRecordID()
		protein := 20.5
		fat := 10.3
		carbs := 50.8

		recordPfc := entity.NewRecordPfc(recordID, protein, fat, carbs)

		pfc := recordPfc.Pfc()

		if pfc.Protein() != protein {
			t.Errorf("Pfc().Protein() = %v, want %v", pfc.Protein(), protein)
		}
		if pfc.Fat() != fat {
			t.Errorf("Pfc().Fat() = %v, want %v", pfc.Fat(), fat)
		}
		if pfc.Carbs() != carbs {
			t.Errorf("Pfc().Carbs() = %v, want %v", pfc.Carbs(), carbs)
		}
	})
}
