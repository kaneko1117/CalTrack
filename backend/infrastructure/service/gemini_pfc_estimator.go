package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/google/generative-ai-go/genai"

	usecaseService "caltrack/usecase/service"
)

// GeminiPfcEstimator はGemini APIを使用したPFC推定サービスの実装
type GeminiPfcEstimator struct {
	client *genai.Client
}

// NewGeminiPfcEstimator はGeminiPfcEstimatorを生成する
func NewGeminiPfcEstimator(client *genai.Client) *GeminiPfcEstimator {
	return &GeminiPfcEstimator{client: client}
}

// pfcResponse はGeminiのPFC推定レスポンスをパースするための構造体
type pfcResponse struct {
	Protein float64 `json:"protein"`
	Fat     float64 `json:"fat"`
	Carbs   float64 `json:"carbs"`
}

// Estimate は食品名リストからPFC値を推定する
// configからモデル名・プロンプトを受け取る（ビジネスロジックはUsecase層で管理）
func (g *GeminiPfcEstimator) Estimate(ctx context.Context, config usecaseService.PfcEstimatorConfig, input usecaseService.PfcEstimateInput) (*usecaseService.PfcEstimateOutput, error) {
	if g.client == nil {
		return nil, errors.New("Gemini client is not initialized")
	}

	// 食品が0件の場合は0を返す
	if len(input.FoodItems) == 0 {
		return &usecaseService.PfcEstimateOutput{Protein: 0, Fat: 0, Carbs: 0}, nil
	}

	// configで指定されたモデルを使用
	model := g.client.GenerativeModel(config.ModelName)

	// リクエストログ
	if config.Log.EnableRequestLog {
		log.Printf("[INFO] Gemini PFC Estimator Request - Model: %s, FoodItems: %v", config.ModelName, input.FoodItems)
	}

	// configで指定されたプロンプトを使用してリクエストを送信
	resp, err := model.GenerateContent(ctx, genai.Text(config.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// トークン使用量ログ
	if config.Log.EnableTokenLog && resp.UsageMetadata != nil {
		log.Printf("[INFO] Gemini PFC Estimator Token Usage - PromptTokens: %d, CandidatesTokens: %d, TotalTokens: %d",
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
		log.Printf("[INFO] Gemini PFC Estimator Response - Raw: %s", responseText)
	}

	// JSONをパース
	output, err := parsePfcResponse(responseText)
	if err != nil {
		log.Printf("[ERROR] Failed to parse PFC response: %v, raw: %s", err, responseText)
		return nil, fmt.Errorf("failed to parse PFC response: %w", err)
	}

	return output, nil
}

// parsePfcResponse はGeminiのレスポンステキストをパースしてPfcEstimateOutputに変換する
func parsePfcResponse(text string) (*usecaseService.PfcEstimateOutput, error) {
	// JSON部分を抽出
	jsonStr := extractJSON(text)

	// JSONをパース
	var resp pfcResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, err
	}

	// マイナス値を0に補正
	if resp.Protein < 0 {
		resp.Protein = 0
	}
	if resp.Fat < 0 {
		resp.Fat = 0
	}
	if resp.Carbs < 0 {
		resp.Carbs = 0
	}

	return &usecaseService.PfcEstimateOutput{
		Protein: resp.Protein,
		Fat:     resp.Fat,
		Carbs:   resp.Carbs,
	}, nil
}

// extractJSON はテキストからJSON部分を抽出する
func extractJSON(text string) string {
	// マークダウンのコードブロックを除去
	re := regexp.MustCompile("```json\\s*([\\s\\S]*?)\\s*```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}

	// JSON部分のみ抽出
	re = regexp.MustCompile(`\{[\s\S]*\}`)
	match := re.FindString(text)
	if match != "" {
		return match
	}

	// 見つからなければそのまま返す
	return text
}
