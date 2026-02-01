package vo

import (
	domainErrors "caltrack/domain/errors"
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
