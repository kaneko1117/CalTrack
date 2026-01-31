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
			return err
		}
		if exists {
			return domainErrors.ErrEmailAlreadyExists
		}

		return u.userRepo.Save(txCtx, user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
