package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewWeight(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantKg  float64
		wantErr error
	}{
		// 正常系
		{"valid weight 70.5kg", 70.5, 70.5, nil},
		{"valid weight 1kg", 1.0, 1.0, nil},
		// 異常系
		{"zero", 0, 0, domainErrors.ErrWeightMustBePositive},
		{"negative", -10, 0, domainErrors.ErrWeightMustBePositive},
		{"too heavy 501kg", 501, 0, domainErrors.ErrWeightTooHeavy},
		// 境界値
		{"min positive 0.1kg", 0.1, 0.1, nil},
		{"max weight 500kg", 500, 500, nil},
		{"exceeds max 500.1kg", 500.1, 0, domainErrors.ErrWeightTooHeavy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewWeight(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewWeight(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.Kg() != tt.wantKg {
				t.Errorf("NewWeight(%v).Kg() = %v, want %v", tt.input, got.Kg(), tt.wantKg)
			}
		})
	}
}
