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
		{"正常な体重70.5kg", 70.5, 70.5, nil},
		{"正常な体重1kg", 1.0, 1.0, nil},
		// 異常系
		{"0はエラー", 0, 0, domainErrors.ErrWeightMustBePositive},
		{"負の値はエラー", -10, 0, domainErrors.ErrWeightMustBePositive},
		{"501kg以上はエラー", 501, 0, domainErrors.ErrWeightTooHeavy},
		// 境界値
		{"最小値0.1kgは有効", 0.1, 0.1, nil},
		{"最大値500kgは有効", 500, 500, nil},
		{"500.1kgは無効", 500.1, 0, domainErrors.ErrWeightTooHeavy},
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
