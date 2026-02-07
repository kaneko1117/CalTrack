package usecase

import (
	"context"
	"fmt"
	"strings"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
)

type UserUsecase struct {
	userRepo  repository.UserRepository
	txManager repository.TransactionManager
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (u *UserUsecase) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.txManager.Execute(ctx, func(txCtx context.Context) error {
		exists, err := u.userRepo.ExistsByEmail(txCtx, user.Email())
		if err != nil {
			logError("Register", err, "email", user.Email().String())
			return err
		}
		if exists {
			logWarn("Register", "email already exists", "email", user.Email().String())
			return domainErrors.ErrEmailAlreadyExists
		}

		if err := u.userRepo.Save(txCtx, user); err != nil {
			logError("Register", err, "user_id", user.ID().String())
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfileInput はプロフィール更新の入力を表す
type UpdateProfileInput struct {
	Nickname      string
	Height        float64
	Weight        float64
	ActivityLevel string
}

// UpdateProfile は認証ユーザーのプロフィールを更新する
func (u *UserUsecase) UpdateProfile(ctx context.Context, userID vo.UserID, input UpdateProfileInput) (*entity.User, error) {
	var updatedUser *entity.User

	err := u.txManager.Execute(ctx, func(txCtx context.Context) error {
		user, err := u.userRepo.FindByID(txCtx, userID)
		if err != nil {
			logError("UpdateProfile", err, "user_id", userID.String())
			return err
		}
		if user == nil {
			logWarn("UpdateProfile", "user not found", "user_id", userID.String())
			return domainErrors.ErrUserNotFound
		}

		errs := user.UpdateProfile(input.Nickname, input.Height, input.Weight, input.ActivityLevel)
		if errs != nil {
			logWarn("UpdateProfile", "validation errors", "user_id", userID.String(), "errors", formatErrors(errs))
			return errs[0]
		}

		if err := u.userRepo.Update(txCtx, user); err != nil {
			logError("UpdateProfile", err, "user_id", userID.String())
			return err
		}

		updatedUser = user
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// GetProfile は認証ユーザーのプロフィールを取得する
func (u *UserUsecase) GetProfile(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetProfile", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetProfile", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	return user, nil
}

// formatErrors は複数のエラーをカンマ区切りの文字列に変換する
func formatErrors(errs []error) string {
	msgs := make([]string, len(errs))
	for i, err := range errs {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("[%s]", strings.Join(msgs, ", "))
}
