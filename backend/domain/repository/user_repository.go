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
	// FindByID は指定IDのユーザーを取得する
	// 存在しない場合はnilとnilを返す
	FindByID(ctx context.Context, id vo.UserID) (*entity.User, error)
	// Update は既存ユーザーを更新する
	Update(ctx context.Context, user *entity.User) error
}
