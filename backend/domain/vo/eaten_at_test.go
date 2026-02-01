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
