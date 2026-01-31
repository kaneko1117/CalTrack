package vo

import (
	"time"

	domainErrors "caltrack/domain/errors"
)

// セッションの有効期間（7日間）
const SessionDurationDays = 7

// ExpiresAt はセッションの有効期限を表す
type ExpiresAt struct {
	value time.Time
}

// NewExpiresAt は現在時刻から7日後の有効期限を生成する
func NewExpiresAt() ExpiresAt {
	expiresAt := nowFunc().AddDate(0, 0, SessionDurationDays)
	return ExpiresAt{value: expiresAt}
}

// ParseExpiresAt は時刻からExpiresAtを復元する（DBからの復元用）
func ParseExpiresAt(t time.Time) ExpiresAt {
	return ExpiresAt{value: t}
}

// Time は有効期限のtime.Time表現を返す
func (e ExpiresAt) Time() time.Time {
	return e.value
}

// IsExpired は有効期限が過ぎているか判定する
func (e ExpiresAt) IsExpired() bool {
	return nowFunc().After(e.value)
}

// ValidateNotExpired は有効期限内かどうかを検証する
func (e ExpiresAt) ValidateNotExpired() error {
	if e.IsExpired() {
		return domainErrors.ErrSessionExpired
	}
	return nil
}
