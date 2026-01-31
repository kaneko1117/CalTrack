package usecase

import (
	"context"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
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
