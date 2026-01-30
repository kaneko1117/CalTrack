package user

import (
	"context"
	"errors"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
)

type RegisterUserInput struct {
	Email         string
	Password      string
	Nickname      string
	Weight        float64
	Height        float64
	BirthDate     time.Time
	Gender        string
	ActivityLevel string
}

type RegisterUserOutput struct {
	UserID string
}

type RegisterUserUsecase struct {
	userRepo  repository.UserRepository
	txManager repository.TransactionManager
}

func NewRegisterUserUsecase(
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
) *RegisterUserUsecase {
	return &RegisterUserUsecase{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (u *RegisterUserUsecase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	password, err := vo.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	var output *RegisterUserOutput

	err = u.txManager.Execute(ctx, func(txCtx context.Context) error {
		exists, err := u.userRepo.ExistsByEmail(txCtx, email)
		if err != nil {
			return err
		}
		if exists {
			return domainErrors.ErrEmailAlreadyExists
		}

		user, errs := entity.NewUser(
			input.Email,
			hashedPassword.String(),
			input.Nickname,
			input.Weight,
			input.Height,
			input.BirthDate,
			input.Gender,
			input.ActivityLevel,
		)
		if errs != nil {
			return errors.Join(errs...)
		}

		if err := u.userRepo.Save(txCtx, user); err != nil {
			return err
		}

		output = &RegisterUserOutput{
			UserID: user.ID().String(),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
