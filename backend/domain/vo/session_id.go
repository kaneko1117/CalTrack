package vo

import (
	"crypto/rand"
	"encoding/base64"

	domainErrors "caltrack/domain/errors"
)

// セッションIDの長さ（バイト数）
// 32バイト = 256ビットのエントロピー
const sessionIDBytes = 32

// SessionID はセッションを識別するランダムな文字列
type SessionID struct {
	value string
}

// NewSessionID は新しいセッションIDを生成する
func NewSessionID() (SessionID, error) {
	bytes := make([]byte, sessionIDBytes)
	if _, err := rand.Read(bytes); err != nil {
		return SessionID{}, domainErrors.ErrSessionIDGenerationFailed
	}
	// URL-safe Base64エンコード
	value := base64.URLEncoding.EncodeToString(bytes)
	return SessionID{value: value}, nil
}

// ParseSessionID は文字列からSessionIDを復元する
func ParseSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, domainErrors.ErrInvalidSessionID
	}
	// Base64デコードで妥当性を検証
	decoded, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return SessionID{}, domainErrors.ErrInvalidSessionID
	}
	// 長さの検証
	if len(decoded) != sessionIDBytes {
		return SessionID{}, domainErrors.ErrInvalidSessionID
	}
	return SessionID{value: value}, nil
}

// String はセッションIDの文字列表現を返す
func (s SessionID) String() string {
	return s.value
}

// Equals は2つのセッションIDが等しいか比較する
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}
