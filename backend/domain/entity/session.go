package entity

import (
	"time"

	"caltrack/domain/vo"
)

// Session はユーザーセッションを表すエンティティ
type Session struct {
	id        vo.SessionID
	userID    vo.UserID
	expiresAt vo.ExpiresAt
	createdAt time.Time
}

// NewSession は新しいセッションを生成する
// userIDStr: ユーザーID文字列
// 戻り値: 生成されたSession, エラー
func NewSession(userIDStr string) (*Session, error) {
	userID := parseUserID(userIDStr)

	sessionID, err := parseSessionID()
	if err != nil {
		return nil, err
	}

	return &Session{
		id:        sessionID,
		userID:    userID,
		expiresAt: vo.NewExpiresAt(),
		createdAt: time.Now(),
	}, nil
}

// NewSessionWithUserID はUserID VOから新しいセッションを生成する
// 既にUserID VOを持っている場合に使用（認証後など）
func NewSessionWithUserID(userID vo.UserID) (*Session, error) {
	sessionID, err := vo.NewSessionID()
	if err != nil {
		return nil, err
	}

	return &Session{
		id:        sessionID,
		userID:    userID,
		expiresAt: vo.NewExpiresAt(),
		createdAt: time.Now(),
	}, nil
}

// ReconstructSession はDBからの復元用
// NOTE: DBデータは保存時にバリデーション済みなのでVO変換は本来不要だが、
// SessionIDはデータ破損検知のためバリデーションあり。
func ReconstructSession(
	sessionIDStr string,
	userIDStr string,
	expiresAt time.Time,
	createdAt time.Time,
) (*Session, error) {
	sessionID, err := vo.ParseSessionID(sessionIDStr)
	if err != nil {
		return nil, err
	}

	userID := vo.ReconstructUserID(userIDStr)

	return &Session{
		id:        sessionID,
		userID:    userID,
		expiresAt: vo.ParseExpiresAt(expiresAt),
		createdAt: createdAt,
	}, nil
}

// parseUserID はUserID文字列をVOに変換する
func parseUserID(s string) vo.UserID {
	return vo.ReconstructUserID(s)
}

// parseSessionID は新しいSessionIDを生成する
func parseSessionID() (vo.SessionID, error) {
	return vo.NewSessionID()
}

// ID はセッションIDを返す
func (s *Session) ID() vo.SessionID {
	return s.id
}

// UserID はユーザーIDを返す
func (s *Session) UserID() vo.UserID {
	return s.userID
}

// ExpiresAt は有効期限を返す
func (s *Session) ExpiresAt() vo.ExpiresAt {
	return s.expiresAt
}

// CreatedAt は作成日時を返す
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// IsExpired はセッションが有効期限切れかどうかを返す
func (s *Session) IsExpired() bool {
	return s.expiresAt.IsExpired()
}

// ValidateNotExpired はセッションが有効期限内かどうかを検証する
// 有効期限切れの場合は ErrSessionExpired を返す
func (s *Session) ValidateNotExpired() error {
	return s.expiresAt.ValidateNotExpired()
}
