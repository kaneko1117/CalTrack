package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase"
)

// mockRecordRepository はRecordRepositoryのモック実装
type mockRecordRepository struct {
	save                     func(ctx context.Context, record *entity.Record) error
	findByUserIDAndDateRange func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
}

func (m *mockRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	return m.save(ctx, record)
}

func (m *mockRecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
	if m.findByUserIDAndDateRange != nil {
		return m.findByUserIDAndDateRange(ctx, userID, startTime, endTime)
	}
	return nil, nil
}

// mockRecordUserRepository はUserRepositoryのモック実装（Record用）
type mockRecordUserRepository struct {
	findByID func(ctx context.Context, id vo.UserID) (*entity.User, error)
}

func (m *mockRecordUserRepository) Save(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockRecordUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

func (m *mockRecordUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return false, nil
}

func (m *mockRecordUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
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
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				savedRecord = record
				return nil
			},
		}
		userRepo := &mockRecordUserRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
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
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				return saveErr
			},
		}
		userRepo := &mockRecordUserRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		err := uc.Create(context.Background(), validRecord(t))

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

// validUserForRecord はRecord用テストのための有効なUserを生成する
func validUserForRecord(t *testing.T, userID vo.UserID) *entity.User {
	t.Helper()
	// ReconstructUserを使用してuserIDを指定できるようにする
	user, err := entity.ReconstructUser(
		userID.String(),
		"test@example.com",
		"hashedpassword",
		"testuser",
		70.5,
		175.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create valid user: %v", err)
	}
	return user
}

func TestRecordUsecase_GetTodayCalories(t *testing.T) {
	t.Run("正常系_今日のカロリー情報を取得", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		record2, _ := entity.NewRecord(userID, time.Now())
		_ = record2.AddItem("昼食", 500)
		_ = record2.AddItem("おやつ", 100)

		records := []*entity.Record{record1, record2}

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return records, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		output, err := uc.GetTodayCalories(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 合計カロリーの検証（300 + 500 + 100 = 900）
		if output.TotalCalories != 900 {
			t.Errorf("TotalCalories = %d, want 900", output.TotalCalories)
		}

		// 目標カロリーの検証
		expectedTargetCalories := user.CalculateTargetCalories()
		if output.TargetCalories != expectedTargetCalories {
			t.Errorf("TargetCalories = %d, want %d", output.TargetCalories, expectedTargetCalories)
		}

		// 差分の検証（目標 - 実績）
		expectedDifference := expectedTargetCalories - 900
		if output.Difference != expectedDifference {
			t.Errorf("Difference = %d, want %d", output.Difference, expectedDifference)
		}

		// Record数の検証
		if len(output.Records) != 2 {
			t.Errorf("len(Records) = %d, want 2", len(output.Records))
		}
	})

	t.Run("正常系_記録が0件の場合", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{}, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		output, err := uc.GetTodayCalories(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 合計カロリーは0
		if output.TotalCalories != 0 {
			t.Errorf("TotalCalories = %d, want 0", output.TotalCalories)
		}

		// Record数は0
		if len(output.Records) != 0 {
			t.Errorf("len(Records) = %d, want 0", len(output.Records))
		}
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userID := vo.NewUserID()

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil // ユーザーが存在しない
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Record取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return nil, repoErr
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, userRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
