package vo

import (
	"time"

	domainErrors "caltrack/domain/errors"
)

// MealType は食事タイプを表す
type MealType int

const (
	MealTypeBreakfast MealType = iota + 1 // 朝食 (5:00 - 11:00)
	MealTypeLunch                         // 昼食 (11:00 - 14:00)
	MealTypeSnack                         // 間食 (14:00 - 17:00)
	MealTypeDinner                        // 夕食 (17:00 - 21:00)
	MealTypeLateNight                     // 夜食 (21:00 - 5:00)
)

// String は食事タイプの日本語名を返す
func (m MealType) String() string {
	switch m {
	case MealTypeBreakfast:
		return "朝食"
	case MealTypeLunch:
		return "昼食"
	case MealTypeSnack:
		return "間食"
	case MealTypeDinner:
		return "夕食"
	case MealTypeLateNight:
		return "夜食"
	default:
		return "不明"
	}
}

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

// MealType は食事日時から食事タイプを判定して返す
func (e EatenAt) MealType() MealType {
	hour := e.value.Hour()

	switch {
	case hour >= 5 && hour < 11:
		return MealTypeBreakfast
	case hour >= 11 && hour < 14:
		return MealTypeLunch
	case hour >= 14 && hour < 17:
		return MealTypeSnack
	case hour >= 17 && hour < 21:
		return MealTypeDinner
	default:
		return MealTypeLateNight
	}
}
