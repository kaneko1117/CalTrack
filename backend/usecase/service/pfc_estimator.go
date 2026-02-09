package service

import (
	"context"
)

// PfcEstimateInput は食品名からPFCを推定するための入力
type PfcEstimateInput struct {
	FoodItems []string
}

// PfcEstimateOutput はPFC推定結果
type PfcEstimateOutput struct {
	Protein float64
	Fat     float64
	Carbs   float64
}

// PfcEstimatorConfig はPfcEstimatorの設定
type PfcEstimatorConfig struct {
	ModelName string
	Prompt    string
}

// PfcEstimator は食品名からPFCを推定するサービスインターフェース
type PfcEstimator interface {
	Estimate(ctx context.Context, config PfcEstimatorConfig, input PfcEstimateInput) (*PfcEstimateOutput, error)
}
