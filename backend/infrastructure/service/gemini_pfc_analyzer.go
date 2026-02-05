package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"

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
	if config.Log.EnableRequestLog {
		log.Printf("[INFO] Gemini PFC Analyzer Request - Model: %s", config.ModelName)
	}

	// Gemini APIにリクエスト送信
	resp, err := model.GenerateContent(ctx, genai.Text(config.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// トークン使用量ログ
	if config.Log.EnableTokenLog && resp.UsageMetadata != nil {
		log.Printf("[INFO] Gemini API Token Usage - PromptTokens: %d, CandidatesTokens: %d, TotalTokens: %d",
			resp.UsageMetadata.PromptTokenCount,
			resp.UsageMetadata.CandidatesTokenCount,
			resp.UsageMetadata.TotalTokenCount)
	}

	// レスポンスからテキストを抽出
	responseText, err := extractResponseText(resp)
	if err != nil {
		log.Printf("[ERROR] Failed to extract response text: %v", err)
		return nil, err
	}

	// レスポンスログ
	if config.Log.EnableResponseLog {
		log.Printf("[INFO] Gemini PFC Analyzer Response - Advice length: %d chars", len(responseText))
	}

	return &usecaseService.NutritionAdviceOutput{
		Advice: responseText,
	}, nil
}
