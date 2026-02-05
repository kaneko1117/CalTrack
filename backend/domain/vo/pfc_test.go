package vo_test

import (
	"testing"

	"caltrack/domain/vo"
)

func TestNewPfc(t *testing.T) {
	tests := []struct {
		name        string
		protein     float64
		fat         float64
		carbs       float64
		wantProtein float64
		wantFat     float64
		wantCarbs   float64
	}{
		{"標準的なPFC", 50.0, 30.0, 100.0, 50.0, 30.0, 100.0},
		{"ゼロ", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{"小数点あり", 25.5, 15.3, 80.7, 25.5, 15.3, 80.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfc := vo.NewPfc(tt.protein, tt.fat, tt.carbs)

			if pfc.Protein() != tt.wantProtein {
				t.Errorf("Protein() = %v, want %v", pfc.Protein(), tt.wantProtein)
			}
			if pfc.Fat() != tt.wantFat {
				t.Errorf("Fat() = %v, want %v", pfc.Fat(), tt.wantFat)
			}
			if pfc.Carbs() != tt.wantCarbs {
				t.Errorf("Carbs() = %v, want %v", pfc.Carbs(), tt.wantCarbs)
			}
		})
	}
}
