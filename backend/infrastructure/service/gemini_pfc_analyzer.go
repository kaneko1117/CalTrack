package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/generative-ai-go/genai"

	"caltrack/pkg/logger"
	usecaseService "caltrack/usecase/service"
)

// GeminiPfcAnalyzer はGemini APIを使用したPFC分析サービスの実装
type GeminiPfcAnalyzer struct {
	client *genai.Client
}

// NewGeminiPfcAnalyzer はGeminiPfcAnalyzerを生成する
func NewGeminiPfcAnalyzer(client *genai.Client) *GeminiPfcAnalyzer {
	return &GeminiPfcAnalyzer{
		client: client,
	}
}

// Analyze はNutritionAdviceInputを受け取り、Gemini APIでアドバイスを生成する
func (g *GeminiPfcAnalyzer) Analyze(ctx context.Context, config usecaseService.PfcAnalyzerConfig, input usecaseService.NutritionAdviceInput) (*usecaseService.NutritionAdviceOutput, error) {
	if g.client == nil {
		return nil, errors.New("Gemini client is not initialized")
	}

	// モデルの取得
	model := g.client.GenerativeModel(config.ModelName)

	// リクエストログ
	logger.Debug("Gemini PFC Analyzer Request", "model", config.ModelName)

	// Gemini APIにリクエスト送信
	resp, err := model.GenerateContent(ctx, genai.Text(config.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// トークン使用量ログ
	if resp.UsageMetadata != nil {
		logger.Debug("Gemini API Token Usage",
			"promptTokens", resp.UsageMetadata.PromptTokenCount,
			"candidatesTokens", resp.UsageMetadata.CandidatesTokenCount,
			"totalTokens", resp.UsageMetadata.TotalTokenCount)
	}

	// レスポンスからテキストを抽出
	responseText, err := extractResponseText(resp)
	if err != nil {
		logger.Error("Failed to extract response text", "error", err)
		return nil, err
	}

	// レスポンスログ
	logger.Debug("Gemini PFC Analyzer Response", "adviceLength", len(responseText))

	return &usecaseService.NutritionAdviceOutput{
		Advice: responseText,
	}, nil
}
