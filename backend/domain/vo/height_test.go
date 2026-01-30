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
		{"valid height 170.5cm", 170.5, 170.5, nil},
		{"valid height 100cm", 100, 100, nil},
		// 異常系
		{"zero", 0, 0, domainErrors.ErrHeightMustBePositive},
		{"negative", -10, 0, domainErrors.ErrHeightMustBePositive},
		{"too tall 301cm", 301, 0, domainErrors.ErrHeightTooTall},
		// 境界値
		{"min positive 0.1cm", 0.1, 0.1, nil},
		{"max height 300cm", 300, 300, nil},
		{"exceeds max 300.1cm", 300.1, 0, domainErrors.ErrHeightTooTall},
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
		{"170cm to 1.7m", 170, 1.7},
		{"100cm to 1m", 100, 1.0},
		{"185.5cm to 1.855m", 185.5, 1.855},
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
