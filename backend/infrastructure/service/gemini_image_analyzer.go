package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"

	"caltrack/domain/vo"
	usecaseService "caltrack/usecase/service"
)

// GeminiImageAnalyzer はGemini APIを使用した画像解析サービスの実装
type GeminiImageAnalyzer struct {
	client *genai.Client
}

// NewGeminiImageAnalyzer はGeminiImageAnalyzerを生成する
func NewGeminiImageAnalyzer(client *genai.Client) *GeminiImageAnalyzer {
	return &GeminiImageAnalyzer{
		client: client,
	}
}

// geminiResponse はGeminiのレスポンスをパースするための構造体
type geminiResponse struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

// Analyze は画像データを解析し、認識した食品のリストを返す
// configからモデル名・プロンプトを受け取る（ビジネスロジックはUsecase層で管理）
func (g *GeminiImageAnalyzer) Analyze(ctx context.Context, config usecaseService.ImageAnalyzerConfig, imageData string, mimeType string) ([]usecaseService.AnalyzedItem, error) {
	if g.client == nil {
		return nil, errors.New("Gemini client is not initialized")
	}

	// configで指定されたモデルを使用
	model := g.client.GenerativeModel(config.ModelName)

	// 画像データをBlobとして作成
	imageBlob := genai.Blob{
		MIMEType: mimeType,
		Data:     []byte(imageData),
	}

	// configで指定されたプロンプトを使用してリクエストを送信
	resp, err := model.GenerateContent(ctx, genai.Text(config.Prompt), imageBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// レスポンスからテキストを抽出
	responseText, err := extractResponseText(resp)
	if err != nil {
		log.Printf("[ERROR] Failed to extract response text: %v", err)
		return nil, err
	}

	// JSONをパース
	items, err := parseGeminiResponse(responseText)
	if err != nil {
		log.Printf("[ERROR] Failed to parse Gemini response: %v, response: %s", err, responseText)
		return nil, err
	}

	return items, nil
}

// extractResponseText はGeminiのレスポンスからテキスト部分を抽出する
func extractResponseText(resp *genai.GenerateContentResponse) (string, error) {
	if resp == nil || len(resp.Candidates) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	var textParts []string
	for _, part := range candidate.Content.Parts {
		if text, ok := part.(genai.Text); ok {
			textParts = append(textParts, string(text))
		}
	}

	if len(textParts) == 0 {
		return "", fmt.Errorf("no text parts in response")
	}

	return strings.Join(textParts, ""), nil
}

// parseGeminiResponse はGeminiのレスポンステキストをパースしてAnalyzedItemのスライスに変換する
func parseGeminiResponse(responseText string) ([]usecaseService.AnalyzedItem, error) {
	// 前後の空白を除去
	responseText = strings.TrimSpace(responseText)

	// マークダウンのコードブロックを除去（念のため）
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// JSONをパース
	var responses []geminiResponse
	if err := json.Unmarshal([]byte(responseText), &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// ドメインオブジェクトに変換
	items := make([]usecaseService.AnalyzedItem, 0, len(responses))
	for _, r := range responses {
		// 食品名のVO作成
		itemName, err := vo.NewItemName(r.Name)
		if err != nil {
			log.Printf("[WARN] Invalid item name skipped: %s, error: %v", r.Name, err)
			continue
		}

		// カロリーのVO作成
		calories, err := vo.NewCalories(r.Calories)
		if err != nil {
			log.Printf("[WARN] Invalid calories skipped: %d, error: %v", r.Calories, err)
			continue
		}

		items = append(items, usecaseService.AnalyzedItem{
			Name:     itemName,
			Calories: calories,
		})
	}

	return items, nil
}
