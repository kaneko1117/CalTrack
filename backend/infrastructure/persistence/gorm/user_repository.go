package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
	"caltrack/infrastructure/persistence/gorm/model"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) error {
	tx := GetTx(ctx, r.db)
	m := toUserModel(user)
	if err := tx.Create(&m).Error; err != nil {
		logError("Save", err, "user_id", user.ID().String())
		return err
	}
	return nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	tx := GetTx(ctx, r.db)
	var m model.User
	err := tx.Where("email = ?", email.String()).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logError("FindByEmail", err, "email", email.String())
		return nil, err
	}
	return toUserEntity(&m)
}

func (r *GormUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	tx := GetTx(ctx, r.db)
	var count int64
	err := tx.Model(&model.User{}).Where("email = ?", email.String()).Count(&count).Error
	if err != nil {
		logError("ExistsByEmail", err, "email", email.String())
		return false, err
	}
	return count > 0, nil
}

func toUserModel(user *entity.User) model.User {
	return model.User{
		ID:             user.ID().String(),
		Email:          user.Email().String(),
		HashedPassword: user.HashedPassword().String(),
		Nickname:       user.Nickname().String(),
		Weight:         user.Weight().Kg(),
		Height:         user.Height().Cm(),
		BirthDate:      user.BirthDate().Time(),
		Gender:         user.Gender().String(),
		ActivityLevel:  user.ActivityLevel().String(),
		CreatedAt:      user.CreatedAt(),
		UpdatedAt:      user.UpdatedAt(),
	}
}

func toUserEntity(m *model.User) (*entity.User, error) {
	return entity.ReconstructUser(
		m.ID,
		m.Email,
		m.HashedPassword,
		m.Nickname,
		m.Weight,
		m.Height,
		m.BirthDate,
		m.Gender,
		m.ActivityLevel,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
