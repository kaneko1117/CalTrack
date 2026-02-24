package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const defaultGeminiModelName = "gemini-3-flash-preview"

// GeminiConfig はGemini APIクライアントと設定を保持する構造体
type GeminiConfig struct {
	Client    *genai.Client
	modelName string
}

// NewGeminiConfig はGeminiクライアントを初期化し、GeminiConfig構造体を返す
// APIキーが未設定の場合、Client=nilのGeminiConfigを返す（画像解析機能は無効）
func NewGeminiConfig() (*GeminiConfig, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("GEMINI_API_KEY is not set. Image analysis feature will be disabled.")
		return &GeminiConfig{Client: nil, modelName: defaultGeminiModelName}, nil
	}

	modelName := os.Getenv("GEMINI_MODEL_NAME")
	if modelName == "" {
		modelName = defaultGeminiModelName
	}
	log.Printf("Gemini model: %s", modelName)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	log.Println("Gemini client initialized")
	return &GeminiConfig{Client: client, modelName: modelName}, nil
}

// Close はGeminiクライアントをクローズする
func (g *GeminiConfig) Close() {
	if g.Client != nil {
		if err := g.Client.Close(); err != nil {
			log.Printf("Failed to close Gemini client: %v", err)
		}
		log.Println("Gemini client closed")
	}
}

// GeminiModelName は使用するGeminiモデル名を返す
// usecase.AIConfig インターフェースを満たす
func (g *GeminiConfig) GeminiModelName() string {
	return g.modelName
}
