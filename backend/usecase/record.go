package usecase

import (
	"context"

	"caltrack/domain/entity"
	"caltrack/domain/repository"
)

// RecordUsecase はカロリー記録に関するユースケースを提供する
type RecordUsecase struct {
	recordRepo repository.RecordRepository
	txManager  repository.TransactionManager
}

// NewRecordUsecase は RecordUsecase のインスタンスを生成する
func NewRecordUsecase(
	recordRepo repository.RecordRepository,
	txManager repository.TransactionManager,
) *RecordUsecase {
	return &RecordUsecase{
		recordRepo: recordRepo,
		txManager:  txManager,
	}
}

// Create は新しいカロリー記録を作成する
func (u *RecordUsecase) Create(ctx context.Context, record *entity.Record) error {
	err := u.txManager.Execute(ctx, func(txCtx context.Context) error {
		if err := u.recordRepo.Save(txCtx, record); err != nil {
			logError("Create", err, "record_id", record.ID().String())
			return err
		}
		return nil
	})

	return err
}
