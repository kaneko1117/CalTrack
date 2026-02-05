package service

import (
	"context"

	"caltrack/domain/vo"
)

type NutritionAdviceInput struct {
	TargetCalories  int
	TargetPfc       vo.Pfc
	CurrentCalories int
	CurrentPfc      vo.Pfc
	FoodItems       []string
}

type PfcAnalyzerLogConfig struct {
	EnableRequestLog  bool
	EnableResponseLog bool
	EnableTokenLog    bool
}

type PfcAnalyzerConfig struct {
	ModelName string
	Prompt    string
	Log       PfcAnalyzerLogConfig
}

type NutritionAdviceOutput struct {
	Advice string
}

type PfcAnalyzer interface {
	Analyze(ctx context.Context, config PfcAnalyzerConfig, input NutritionAdviceInput) (*NutritionAdviceOutput, error)
}
