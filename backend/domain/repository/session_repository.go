package repository

import (
	"context"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// SessionRepository はセッションの永続化を担当するインターフェース
type SessionRepository interface {
	// Save はセッションを保存する
	Save(ctx context.Context, session *entity.Session) error

	// FindByID はセッションIDでセッションを取得する
	FindByID(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error)

	// DeleteByID はセッションを削除する
	DeleteByID(ctx context.Context, sessionID vo.SessionID) error

	// DeleteByUserID はユーザーの全セッションを削除する（ログアウト時等）
	DeleteByUserID(ctx context.Context, userID vo.UserID) error
}
