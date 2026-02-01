package vo

import (
	"time"

	domainErrors "caltrack/domain/errors"
)

// EatenAt はカロリー記録の食事日時を表す値オブジェクト
type EatenAt struct {
	value time.Time
}

// NewEatenAt は指定された時刻からEatenAtを生成する
// 未来の日時の場合はエラーを返す
func NewEatenAt(t time.Time) (EatenAt, error) {
	now := nowFunc()
	if t.After(now) {
		return EatenAt{}, domainErrors.ErrEatenAtMustNotBeFuture
	}
	return EatenAt{value: t}, nil
}

// ReconstructEatenAt はDBからEatenAtを復元する（バリデーションなし）
func ReconstructEatenAt(t time.Time) EatenAt {
	return EatenAt{value: t}
}

// Time はEatenAtのtime.Time表現を返す
func (e EatenAt) Time() time.Time {
	return e.value
}

// Equals は2つのEatenAtが等しいかを比較する
func (e EatenAt) Equals(other EatenAt) bool {
	return e.value.Equal(other.value)
}
