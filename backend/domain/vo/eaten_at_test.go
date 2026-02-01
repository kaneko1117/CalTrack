package vo

import (
	"testing"
	"time"

	domainErrors "caltrack/domain/errors"
)

func TestNewEatenAt(t *testing.T) {
	// 現在時刻を固定
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	oneSecondAgo := fixedNow.Add(-1 * time.Second)
	oneHourAgo := fixedNow.Add(-1 * time.Hour)
	oneDayAgo := fixedNow.AddDate(0, 0, -1)
	oneSecondLater := fixedNow.Add(1 * time.Second)
	oneHourLater := fixedNow.Add(1 * time.Hour)
	oneDayLater := fixedNow.AddDate(0, 0, 1)

	tests := []struct {
		name    string
		input   time.Time
		wantErr error
	}{
		// 正常系
		{"現在時刻は有効", fixedNow, nil},
		{"1秒前は有効", oneSecondAgo, nil},
		{"1時間前は有効", oneHourAgo, nil},
		{"1日前は有効", oneDayAgo, nil},
		// 異常系
		{"1秒後はエラー", oneSecondLater, domainErrors.ErrEatenAtMustNotBeFuture},
		{"1時間後はエラー", oneHourLater, domainErrors.ErrEatenAtMustNotBeFuture},
		{"1日後はエラー", oneDayLater, domainErrors.ErrEatenAtMustNotBeFuture},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eatenAt, err := NewEatenAt(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewEatenAt(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if err == nil && !eatenAt.Time().Equal(tt.input) {
				t.Errorf("NewEatenAt(%v).Time() = %v, want %v", tt.input, eatenAt.Time(), tt.input)
			}
		})
	}
}

func TestReconstructEatenAt(t *testing.T) {
	input := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	eatenAt := ReconstructEatenAt(input)

	if !eatenAt.Time().Equal(input) {
		t.Errorf("ReconstructEatenAt(%v).Time() = %v, want %v", input, eatenAt.Time(), input)
	}
}

func TestEatenAt_Equals(t *testing.T) {
	time1 := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	time3 := time.Date(2024, 6, 15, 13, 0, 0, 0, time.UTC)

	eatenAt1 := ReconstructEatenAt(time1)
	eatenAt2 := ReconstructEatenAt(time2)
	eatenAt3 := ReconstructEatenAt(time3)

	tests := []struct {
		name string
		a    EatenAt
		b    EatenAt
		want bool
	}{
		{"同じ時刻は等しい", eatenAt1, eatenAt2, true},
		{"異なる時刻は等しくない", eatenAt1, eatenAt3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.want {
				t.Errorf("EatenAt.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEatenAt_MealType(t *testing.T) {
	tests := []struct {
		name     string
		hour     int
		wantType MealType
		wantName string
	}{
		// 朝食 (5:00 - 11:00)
		{"5時は朝食", 5, MealTypeBreakfast, "朝食"},
		{"10時は朝食", 10, MealTypeBreakfast, "朝食"},
		// 昼食 (11:00 - 14:00)
		{"11時は昼食", 11, MealTypeLunch, "昼食"},
		{"13時は昼食", 13, MealTypeLunch, "昼食"},
		// 間食 (14:00 - 17:00)
		{"14時は間食", 14, MealTypeSnack, "間食"},
		{"16時は間食", 16, MealTypeSnack, "間食"},
		// 夕食 (17:00 - 21:00)
		{"17時は夕食", 17, MealTypeDinner, "夕食"},
		{"20時は夕食", 20, MealTypeDinner, "夕食"},
		// 夜食 (21:00 - 5:00)
		{"21時は夜食", 21, MealTypeLateNight, "夜食"},
		{"0時は夜食", 0, MealTypeLateNight, "夜食"},
		{"4時は夜食", 4, MealTypeLateNight, "夜食"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eatenAt := ReconstructEatenAt(time.Date(2024, 6, 15, tt.hour, 0, 0, 0, time.UTC))

			if got := eatenAt.MealType(); got != tt.wantType {
				t.Errorf("MealType() = %v, want %v", got, tt.wantType)
			}
			if got := eatenAt.MealType().String(); got != tt.wantName {
				t.Errorf("MealType().String() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMealType_String(t *testing.T) {
	t.Run("不明な値", func(t *testing.T) {
		// 定義されていないMealType値の場合
		invalidType := MealType(99)
		if got := invalidType.String(); got != "不明" {
			t.Errorf("MealType(99).String() = %v, want %v", got, "不明")
		}
	})
}
