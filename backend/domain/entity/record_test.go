package entity_test

import (
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewRecord(t *testing.T) {
	t.Run("正常系_有効なパラメータでRecord作成", func(t *testing.T) {
		userID := vo.NewUserID()
		eatenAt := time.Now().Add(-1 * time.Hour)

		record, err := entity.NewRecord(userID, eatenAt)

		if err != nil {
			t.Fatalf("NewRecord() unexpected error = %v", err)
		}
		if record == nil {
			t.Fatal("NewRecord() returned nil record")
		}
		if !record.UserID().Equals(userID) {
			t.Errorf("NewRecord().UserID() = %v, want %v", record.UserID(), userID)
		}
		if len(record.Items()) != 0 {
			t.Errorf("NewRecord().Items() length = %d, want %d", len(record.Items()), 0)
		}
		if record.ID().String() == "" {
			t.Error("NewRecord().ID() should not be empty")
		}
	})

	t.Run("異常系_未来の日時の場合エラー", func(t *testing.T) {
		userID := vo.NewUserID()
		eatenAt := time.Now().Add(1 * time.Hour) // 未来の日時

		record, err := entity.NewRecord(userID, eatenAt)

		if record != nil {
			t.Error("NewRecord() should return nil for future eaten_at")
		}
		if err != domainErrors.ErrEatenAtMustNotBeFuture {
			t.Errorf("NewRecord() error = %v, want %v", err, domainErrors.ErrEatenAtMustNotBeFuture)
		}
	})
}

func TestRecord_AddItem(t *testing.T) {
	t.Run("正常系_アイテム追加", func(t *testing.T) {
		userID := vo.NewUserID()
		eatenAt := time.Now().Add(-1 * time.Hour)
		record, _ := entity.NewRecord(userID, eatenAt)

		err := record.AddItem("おにぎり", 180)

		if err != nil {
			t.Fatalf("AddItem() unexpected error = %v", err)
		}
		if len(record.Items()) != 1 {
			t.Errorf("AddItem() items count = %d, want %d", len(record.Items()), 1)
		}
		if record.Items()[0].Name().String() != "おにぎり" {
			t.Errorf("AddItem() name = %v, want %v", record.Items()[0].Name().String(), "おにぎり")
		}
		if record.Items()[0].Calories().Value() != 180 {
			t.Errorf("AddItem() calories = %v, want %v", record.Items()[0].Calories().Value(), 180)
		}
	})

	t.Run("異常系_無効なアイテム", func(t *testing.T) {
		userID := vo.NewUserID()
		eatenAt := time.Now().Add(-1 * time.Hour)

		tests := []struct {
			name     string
			itemName string
			calories int
			wantErr  error
		}{
			{
				name:     "食品名が空の場合",
				itemName: "",
				calories: 180,
				wantErr:  domainErrors.ErrItemNameRequired,
			},
			{
				name:     "カロリーが0以下の場合",
				itemName: "おにぎり",
				calories: 0,
				wantErr:  domainErrors.ErrCaloriesMustBePositive,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				record, _ := entity.NewRecord(userID, eatenAt)
				err := record.AddItem(tt.itemName, tt.calories)

				if err != tt.wantErr {
					t.Errorf("AddItem() error = %v, want %v", err, tt.wantErr)
				}
				if len(record.Items()) != 0 {
					t.Errorf("AddItem() should not add item on error, items count = %d", len(record.Items()))
				}
			})
		}
	})
}

func TestReconstructRecord(t *testing.T) {
	t.Run("DB復元", func(t *testing.T) {
		idStr := "record-123"
		userIDStr := "user-456"
		eatenAtTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		createdAt := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
		items := []entity.RecordItem{
			*entity.ReconstructRecordItem("item-1", idStr, "おにぎり", 180),
			*entity.ReconstructRecordItem("item-2", idStr, "味噌汁", 50),
		}

		record := entity.ReconstructRecord(idStr, userIDStr, eatenAtTime, createdAt, items)

		if record.ID().String() != idStr {
			t.Errorf("ReconstructRecord().ID() = %v, want %v", record.ID().String(), idStr)
		}
		if record.UserID().String() != userIDStr {
			t.Errorf("ReconstructRecord().UserID() = %v, want %v", record.UserID().String(), userIDStr)
		}
		if !record.EatenAt().Time().Equal(eatenAtTime) {
			t.Errorf("ReconstructRecord().EatenAt() = %v, want %v", record.EatenAt().Time(), eatenAtTime)
		}
		if !record.CreatedAt().Equal(createdAt) {
			t.Errorf("ReconstructRecord().CreatedAt() = %v, want %v", record.CreatedAt(), createdAt)
		}
		if len(record.Items()) != 2 {
			t.Errorf("ReconstructRecord().Items() length = %d, want %d", len(record.Items()), 2)
		}
	})
}

func TestRecord_TotalCalories(t *testing.T) {
	idStr := "record-123"
	userIDStr := "user-456"
	eatenAtTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		items     []entity.RecordItem
		wantTotal int
	}{
		{
			name: "単一アイテム",
			items: []entity.RecordItem{
				*entity.ReconstructRecordItem("item-1", idStr, "おにぎり", 180),
			},
			wantTotal: 180,
		},
		{
			name: "複数アイテム",
			items: []entity.RecordItem{
				*entity.ReconstructRecordItem("item-1", idStr, "おにぎり", 180),
				*entity.ReconstructRecordItem("item-2", idStr, "味噌汁", 50),
				*entity.ReconstructRecordItem("item-3", idStr, "焼き鮭", 200),
			},
			wantTotal: 430,
		},
		{
			name:      "アイテムなし",
			items:     []entity.RecordItem{},
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := entity.ReconstructRecord(idStr, userIDStr, eatenAtTime, createdAt, tt.items)

			if got := record.TotalCalories(); got != tt.wantTotal {
				t.Errorf("TotalCalories() = %d, want %d", got, tt.wantTotal)
			}
		})
	}
}

func TestReconstructRecordItem(t *testing.T) {
	t.Run("DB復元", func(t *testing.T) {
		idStr := "item-123"
		recordIDStr := "record-456"
		nameStr := "おにぎり"
		caloriesVal := 180

		item := entity.ReconstructRecordItem(idStr, recordIDStr, nameStr, caloriesVal)

		if item.ID().String() != idStr {
			t.Errorf("ReconstructRecordItem().ID() = %v, want %v", item.ID().String(), idStr)
		}
		if item.RecordID().String() != recordIDStr {
			t.Errorf("ReconstructRecordItem().RecordID() = %v, want %v", item.RecordID().String(), recordIDStr)
		}
		if item.Name().String() != nameStr {
			t.Errorf("ReconstructRecordItem().Name() = %v, want %v", item.Name().String(), nameStr)
		}
		if item.Calories().Value() != caloriesVal {
			t.Errorf("ReconstructRecordItem().Calories() = %v, want %v", item.Calories().Value(), caloriesVal)
		}
	})
}
