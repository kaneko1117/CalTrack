package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
	"caltrack/infrastructure/persistence/gorm/model"
)

// GormSessionRepository はSessionRepositoryのGORM実装
type GormSessionRepository struct {
	db *gorm.DB
}

// NewGormSessionRepository は新しいGormSessionRepositoryを生成する
func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

// Save はセッションを保存する
func (r *GormSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	tx := GetTx(ctx, r.db)
	m := toSessionModel(session)
	if err := tx.Create(&m).Error; err != nil {
		logError("Save", err, "session_id", session.ID().String())
		return err
	}
	return nil
}

// FindByID はセッションIDでセッションを取得する
func (r *GormSessionRepository) FindByID(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
	tx := GetTx(ctx, r.db)
	var m model.Session
	err := tx.Where("id = ?", sessionID.String()).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logError("FindByID", err, "session_id", sessionID.String())
		return nil, err
	}
	return toSessionEntity(&m)
}

// DeleteByID はセッションを削除する
func (r *GormSessionRepository) DeleteByID(ctx context.Context, sessionID vo.SessionID) error {
	tx := GetTx(ctx, r.db)
	result := tx.Where("id = ?", sessionID.String()).Delete(&model.Session{})
	if result.Error != nil {
		logError("DeleteByID", result.Error, "session_id", sessionID.String())
		return result.Error
	}
	return nil
}

// DeleteByUserID はユーザーの全セッションを削除する
func (r *GormSessionRepository) DeleteByUserID(ctx context.Context, userID vo.UserID) error {
	tx := GetTx(ctx, r.db)
	result := tx.Where("user_id = ?", userID.String()).Delete(&model.Session{})
	if result.Error != nil {
		logError("DeleteByUserID", result.Error, "user_id", userID.String())
		return result.Error
	}
	return nil
}

// toSessionModel はEntityからGORMモデルに変換する
func toSessionModel(session *entity.Session) model.Session {
	return model.Session{
		ID:        session.ID().String(),
		UserID:    session.UserID().String(),
		ExpiresAt: session.ExpiresAt().Time(),
		CreatedAt: session.CreatedAt(),
	}
}

// toSessionEntity はGORMモデルからEntityに変換する
func toSessionEntity(m *model.Session) (*entity.Session, error) {
	return entity.ReconstructSession(
		m.ID,
		m.UserID,
		m.ExpiresAt,
		m.CreatedAt,
	)
}
