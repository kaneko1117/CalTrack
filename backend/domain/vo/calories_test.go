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

func TestZeroCalories(t *testing.T) {
	t.Run("0カロリーが返される", func(t *testing.T) {
		zero := vo.ZeroCalories()
		if zero.Value() != 0 {
			t.Errorf("ZeroCalories().Value() = %v, want 0", zero.Value())
		}
	})
}

func TestCalories_IsAchieved(t *testing.T) {
	tests := []struct {
		name     string
		actual   int
		target   int
		expected bool
	}{
		// 80%以上100%以下が達成
		{"80%ちょうどは達成", 80, 100, true},
		{"90%は達成", 90, 100, true},
		{"100%ちょうどは達成", 100, 100, true},
		{"79%は未達成", 79, 100, false},
		{"101%は超過なので未達成", 101, 100, false},
		{"50%は未達成", 50, 100, false},
		{"0%は未達成", 0, 100, false},
		// 目標が0の場合は常にfalse
		{"目標0の場合はfalse", 100, 0, false},
		// 大きな数値でのテスト
		{"1600/2000は80%で達成", 1600, 2000, true},
		{"2000/2000は100%で達成", 2000, 2000, true},
		{"1599/2000は79.95%で未達成", 1599, 2000, false},
		{"2001/2000は100.05%で超過", 2001, 2000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := vo.ReconstructCalories(tt.actual)
			target := vo.ReconstructCalories(tt.target)
			got := actual.IsAchieved(target)
			if got != tt.expected {
				t.Errorf("Calories(%d).IsAchieved(Calories(%d)) = %v, want %v", tt.actual, tt.target, got, tt.expected)
			}
		})
	}
}

func TestCalories_IsOver(t *testing.T) {
	tests := []struct {
		name     string
		actual   int
		target   int
		expected bool
	}{
		// 100%超が超過
		{"100%ちょうどは超過ではない", 100, 100, false},
		{"101%は超過", 101, 100, true},
		{"99%は超過ではない", 99, 100, false},
		{"80%は超過ではない", 80, 100, false},
		{"200%は超過", 200, 100, true},
		{"0%は超過ではない", 0, 100, false},
		// 目標が0の場合
		{"目標0で実績ありは超過", 100, 0, true},
		{"目標0で実績0は超過ではない", 0, 0, false},
		// 大きな数値でのテスト
		{"2000/2000は超過ではない", 2000, 2000, false},
		{"2001/2000は超過", 2001, 2000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := vo.ReconstructCalories(tt.actual)
			target := vo.ReconstructCalories(tt.target)
			got := actual.IsOver(target)
			if got != tt.expected {
				t.Errorf("Calories(%d).IsOver(Calories(%d)) = %v, want %v", tt.actual, tt.target, got, tt.expected)
			}
		})
	}
}

func TestCalories_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"100 + 200 = 300", 100, 200, 300},
		{"0 + 100 = 100", 0, 100, 100},
		{"100 + 0 = 100", 100, 0, 100},
		{"0 + 0 = 0", 0, 0, 0},
		{"1500 + 500 = 2000", 1500, 500, 2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := vo.ReconstructCalories(tt.a)
			b := vo.ReconstructCalories(tt.b)
			got := a.Add(b)
			if got.Value() != tt.expected {
				t.Errorf("Calories(%d).Add(Calories(%d)).Value() = %v, want %v", tt.a, tt.b, got.Value(), tt.expected)
			}
		})
	}
}
