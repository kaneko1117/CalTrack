package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase/user"
)

type mockUserRepository struct {
	existsByEmail func(ctx context.Context, email vo.Email) (bool, error)
	save          func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return m.existsByEmail(ctx, email)
}

func (m *mockUserRepository) Save(ctx context.Context, u *entity.User) error {
	return m.save(ctx, u)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func validInput() user.RegisterUserInput {
	return user.RegisterUserInput{
		Email:         "test@example.com",
		Password:      "password123",
		Nickname:      "testuser",
		Weight:        70.5,
		Height:        175.0,
		BirthDate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:        "male",
		ActivityLevel: "moderate",
	}
}

func TestRegisterUserUsecase_Success(t *testing.T) {
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	output, err := uc.Execute(context.Background(), validInput())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.UserID == "" {
		t.Error("UserID should not be empty")
	}
}

func TestRegisterUserUsecase_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return true, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	_, err := uc.Execute(context.Background(), validInput())

	if err != domainErrors.ErrEmailAlreadyExists {
		t.Errorf("got %v, want ErrEmailAlreadyExists", err)
	}
}

func TestRegisterUserUsecase_InvalidEmail(t *testing.T) {
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	input := validInput()
	input.Email = "invalid"

	_, err := uc.Execute(context.Background(), input)

	if !errors.Is(err, domainErrors.ErrInvalidEmailFormat) {
		t.Errorf("got %v, want ErrInvalidEmailFormat", err)
	}
}

func TestRegisterUserUsecase_InvalidPassword(t *testing.T) {
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	input := validInput()
	input.Password = "short"

	_, err := uc.Execute(context.Background(), input)

	if err != domainErrors.ErrPasswordTooShort {
		t.Errorf("got %v, want ErrPasswordTooShort", err)
	}
}

func TestRegisterUserUsecase_RepositoryError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, repoErr
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	_, err := uc.Execute(context.Background(), validInput())

	if err != repoErr {
		t.Errorf("got %v, want repoErr", err)
	}
}

func TestRegisterUserUsecase_SaveError(t *testing.T) {
	saveErr := errors.New("save error")
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return saveErr
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	_, err := uc.Execute(context.Background(), validInput())

	if !errors.Is(err, saveErr) {
		t.Errorf("got %v, want saveErr", err)
	}
}

func TestRegisterUserUsecase_MultipleValidationErrors(t *testing.T) {
	repo := &mockUserRepository{
		existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
			return false, nil
		},
		save: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	txManager := &mockTransactionManager{}

	uc := user.NewRegisterUserUsecase(repo, txManager)
	input := validInput()
	input.Nickname = ""
	input.Weight = -1
	input.Gender = "invalid"

	_, err := uc.Execute(context.Background(), input)

	if !errors.Is(err, domainErrors.ErrNicknameRequired) {
		t.Errorf("should contain ErrNicknameRequired: %v", err)
	}
	if !errors.Is(err, domainErrors.ErrWeightMustBePositive) {
		t.Errorf("should contain ErrWeightMustBePositive: %v", err)
	}
	if !errors.Is(err, domainErrors.ErrInvalidGender) {
		t.Errorf("should contain ErrInvalidGender: %v", err)
	}
}
