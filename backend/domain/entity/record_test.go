package entity_test

import (
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewRecord_正常系_有効なパラメータでRecord作成(t *testing.T) {
	userID := vo.NewUserID()
	eatenAt := time.Now().Add(-1 * time.Hour)
	items := []entity.RecordItemInput{
		{Name: "おにぎり", Calories: 180},
		{Name: "味噌汁", Calories: 50},
	}

	record, errs := entity.NewRecord(userID, eatenAt, items)

	if len(errs) > 0 {
		t.Fatalf("NewRecord() unexpected errors = %v", errs)
	}
	if record == nil {
		t.Fatal("NewRecord() returned nil record")
	}
	if !record.UserID().Equals(userID) {
		t.Errorf("NewRecord().UserID() = %v, want %v", record.UserID(), userID)
	}
	if len(record.Items()) != 2 {
		t.Errorf("NewRecord().Items() length = %d, want %d", len(record.Items()), 2)
	}
	if record.ID().String() == "" {
		t.Error("NewRecord().ID() should not be empty")
	}
}

func TestNewRecord_異常系_itemsが空の場合エラー(t *testing.T) {
	userID := vo.NewUserID()
	eatenAt := time.Now().Add(-1 * time.Hour)
	items := []entity.RecordItemInput{}

	record, errs := entity.NewRecord(userID, eatenAt, items)

	if record != nil {
		t.Error("NewRecord() should return nil when items is empty")
	}
	if len(errs) != 1 {
		t.Fatalf("NewRecord() errors count = %d, want 1", len(errs))
	}
	if errs[0] != domainErrors.ErrRecordItemsRequired {
		t.Errorf("NewRecord() error = %v, want %v", errs[0], domainErrors.ErrRecordItemsRequired)
	}
}

func TestNewRecord_異常系_未来の日時の場合エラー(t *testing.T) {
	userID := vo.NewUserID()
	eatenAt := time.Now().Add(1 * time.Hour) // 未来の日時
	items := []entity.RecordItemInput{
		{Name: "おにぎり", Calories: 180},
	}

	record, errs := entity.NewRecord(userID, eatenAt, items)

	if record != nil {
		t.Error("NewRecord() should return nil for future eaten_at")
	}
	if len(errs) != 1 {
		t.Fatalf("NewRecord() errors count = %d, want 1", len(errs))
	}
	if errs[0] != domainErrors.ErrEatenAtMustNotBeFuture {
		t.Errorf("NewRecord() error = %v, want %v", errs[0], domainErrors.ErrEatenAtMustNotBeFuture)
	}
}

func TestNewRecord_異常系_無効なRecordItemの場合エラー(t *testing.T) {
	userID := vo.NewUserID()
	eatenAt := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name      string
		items     []entity.RecordItemInput
		wantErr   error
		errCount  int
	}{
		{
			name:      "食品名が空の場合",
			items:     []entity.RecordItemInput{{Name: "", Calories: 180}},
			wantErr:   domainErrors.ErrItemNameRequired,
			errCount:  1,
		},
		{
			name:      "カロリーが0以下の場合",
			items:     []entity.RecordItemInput{{Name: "おにぎり", Calories: 0}},
			wantErr:   domainErrors.ErrCaloriesMustBePositive,
			errCount:  1,
		},
		{
			name:      "食品名とカロリー両方無効の場合",
			items:     []entity.RecordItemInput{{Name: "", Calories: -1}},
			wantErr:   nil, // 複数エラー
			errCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, errs := entity.NewRecord(userID, eatenAt, tt.items)

			if record != nil {
				t.Error("NewRecord() should return nil for invalid items")
			}
			if len(errs) != tt.errCount {
				t.Errorf("NewRecord() errors count = %d, want %d", len(errs), tt.errCount)
			}
			if tt.wantErr != nil && errs[0] != tt.wantErr {
				t.Errorf("NewRecord() error = %v, want %v", errs[0], tt.wantErr)
			}
		})
	}
}

func TestReconstructRecord_DB復元(t *testing.T) {
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
}

func TestRecord_TotalCalories_合計カロリー計算(t *testing.T) {
	idStr := "record-123"
	userIDStr := "user-456"
	eatenAtTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		items    []entity.RecordItem
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

func TestReconstructRecordItem_DB復元(t *testing.T) {
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
}
