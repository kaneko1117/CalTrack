package usecase

import (
	"context"

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

// UpdateProfile は認証ユーザーのプロフィールを更新する
func (u *UserUsecase) UpdateProfile(ctx context.Context, userID vo.UserID, nickname vo.Nickname, height vo.Height, weight vo.Weight, activityLevel vo.ActivityLevel) (*entity.User, error) {
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

		user.ApplyProfile(nickname, height, weight, activityLevel)

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
