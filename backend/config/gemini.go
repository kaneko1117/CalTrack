package config

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient はGemini APIクライアントのグローバルインスタンス
var GeminiClient *genai.Client

// InitGemini はGeminiクライアントを初期化する
// アプリケーション起動時に1回だけ呼び出す
func InitGemini() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("GEMINI_API_KEY is not set. Image analysis feature will be disabled.")
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		return
	}

	GeminiClient = client
	log.Println("Gemini client initialized")
}

// CloseGemini はGeminiクライアントをクローズする
// アプリケーション終了時に呼び出す
func CloseGemini() {
	if GeminiClient != nil {
		if err := GeminiClient.Close(); err != nil {
			log.Printf("Failed to close Gemini client: %v", err)
		}
		log.Println("Gemini client closed")
	}
}
