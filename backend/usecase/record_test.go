package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
	"caltrack/usecase"
)

// mockRecordRepository はRecordRepositoryのモック実装
type mockRecordRepository struct {
	save func(ctx context.Context, record *entity.Record) error
}

func (m *mockRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	return m.save(ctx, record)
}

// mockRecordTransactionManager はTransactionManagerのモック実装
type mockRecordTransactionManager struct{}

func (m *mockRecordTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// validRecord はテスト用の有効なRecordを生成する
func validRecord(t *testing.T) *entity.Record {
	t.Helper()
	userID := vo.NewUserID()
	record, err := entity.NewRecord(userID, time.Now())
	if err != nil {
		t.Fatalf("failed to create valid record: %v", err)
	}
	return record
}

func TestRecordUsecase_Create(t *testing.T) {
	t.Run("正常系_記録が保存される", func(t *testing.T) {
		var savedRecord *entity.Record
		repo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				savedRecord = record
				return nil
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(repo, txManager)
		record := validRecord(t)
		err := uc.Create(context.Background(), record)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if savedRecord == nil {
			t.Error("record should be saved")
		}
		if savedRecord.ID().String() != record.ID().String() {
			t.Errorf("saved record ID = %s, want %s", savedRecord.ID().String(), record.ID().String())
		}
	})

	t.Run("異常系_保存時にエラーが発生", func(t *testing.T) {
		saveErr := errors.New("save error")
		repo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				return saveErr
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(repo, txManager)
		err := uc.Create(context.Background(), validRecord(t))

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}
