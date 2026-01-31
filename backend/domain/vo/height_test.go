package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewHeight(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantCm  float64
		wantErr error
	}{
		// 正常系
		{"正常な身長170.5cm", 170.5, 170.5, nil},
		{"正常な身長100cm", 100, 100, nil},
		// 異常系
		{"0はエラー", 0, 0, domainErrors.ErrHeightMustBePositive},
		{"負の値はエラー", -10, 0, domainErrors.ErrHeightMustBePositive},
		{"301cm以上はエラー", 301, 0, domainErrors.ErrHeightTooTall},
		// 境界値
		{"最小値0.1cmは有効", 0.1, 0.1, nil},
		{"最大値300cmは有効", 300, 300, nil},
		{"300.1cmは無効", 300.1, 0, domainErrors.ErrHeightTooTall},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewHeight(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewHeight(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.Cm() != tt.wantCm {
				t.Errorf("NewHeight(%v).Cm() = %v, want %v", tt.input, got.Cm(), tt.wantCm)
			}
		})
	}
}

func TestHeight_Meters(t *testing.T) {
	tests := []struct {
		name       string
		inputCm    float64
		wantMeters float64
	}{
		{"170cmは1.7mに変換", 170, 1.7},
		{"100cmは1mに変換", 100, 1.0},
		{"185.5cmは1.855mに変換", 185.5, 1.855},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := vo.NewHeight(tt.inputCm)
			if got := h.Meters(); got != tt.wantMeters {
				t.Errorf("Height.Meters() = %v, want %v", got, tt.wantMeters)
			}
		})
	}
}
