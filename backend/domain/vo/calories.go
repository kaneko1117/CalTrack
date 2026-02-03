package vo

import (
	domainErrors "caltrack/domain/errors"
)

// カロリー達成率の閾値定数
const (
	CaloriesAchievementLowerThreshold = 0.8 // 達成下限（80%）
	CaloriesAchievementUpperThreshold = 1.0 // 達成上限（100%）
)

// Calories はカロリーを表すValue Object
type Calories struct {
	value int
}

// NewCalories は新しいCaloriesを生成する
// 1以上の正の整数のみ許可する
func NewCalories(value int) (Calories, error) {
	if value < 1 {
		return Calories{}, domainErrors.ErrCaloriesMustBePositive
	}
	return Calories{value: value}, nil
}

// ReconstructCalories はDBから復元する際に使用する
// バリデーションをスキップする
func ReconstructCalories(value int) Calories {
	return Calories{value: value}
}

// Value はカロリー値を返す
func (c Calories) Value() int {
	return c.value
}

// ZeroCalories は0カロリーを返す（記録がない日用）
func ZeroCalories() Calories {
	return Calories{value: 0}
}

// IsAchieved は目標に対して達成（80%以上100%以下）しているか判定
func (c Calories) IsAchieved(target Calories) bool {
	if target.value == 0 {
		return false
	}
	ratio := float64(c.value) / float64(target.value)
	return ratio >= CaloriesAchievementLowerThreshold && ratio <= CaloriesAchievementUpperThreshold
}

// IsOver は目標を超過（100%超）しているか判定
func (c Calories) IsOver(target Calories) bool {
	if target.value == 0 {
		return c.value > 0
	}
	return float64(c.value)/float64(target.value) > CaloriesAchievementUpperThreshold
}

// Add は加算した新しいCaloriesを返す
func (c Calories) Add(other Calories) Calories {
	return Calories{value: c.value + other.value}
}
