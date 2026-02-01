package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewCalories(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		wantValue int
		wantErr   error
	}{
		// 正常系
		{"正常なカロリー100kcal", 100, 100, nil},
		{"正常なカロリー1kcal", 1, 1, nil},
		{"大きなカロリー10000kcal", 10000, 10000, nil},
		// 異常系
		{"0はエラー", 0, 0, domainErrors.ErrCaloriesMustBePositive},
		{"負の値はエラー", -10, 0, domainErrors.ErrCaloriesMustBePositive},
		{"負の値-1はエラー", -1, 0, domainErrors.ErrCaloriesMustBePositive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewCalories(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewCalories(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.Value() != tt.wantValue {
				t.Errorf("NewCalories(%v).Value() = %v, want %v", tt.input, got.Value(), tt.wantValue)
			}
		})
	}
}

func TestReconstructCalories(t *testing.T) {
	t.Run("DBからの復元", func(t *testing.T) {
		calories := vo.ReconstructCalories(250)
		if calories.Value() != 250 {
			t.Errorf("ReconstructCalories(250).Value() = %v, want 250", calories.Value())
		}
	})

	t.Run("不正な値でも復元できる", func(t *testing.T) {
		// DBから取得した値はバリデーションをスキップ
		calories := vo.ReconstructCalories(0)
		if calories.Value() != 0 {
			t.Errorf("ReconstructCalories(0).Value() = %v, want 0", calories.Value())
		}
	})
}
