package repository

import (
	"context"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error)
	ExistsByEmail(ctx context.Context, email vo.Email) (bool, error)
}
