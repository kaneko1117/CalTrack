package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/generative-ai-go/genai"

	"caltrack/pkg/logger"
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
	logger.Debug("Gemini PFC Estimator Request", "model", config.ModelName, "foodItems", input.FoodItems)

	// configで指定されたプロンプトを使用してリクエストを送信
	resp, err := model.GenerateContent(ctx, genai.Text(config.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// トークン使用量ログ
	if resp.UsageMetadata != nil {
		logger.Debug("Gemini PFC Estimator Token Usage",
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
	logger.Debug("Gemini PFC Estimator Response", "raw", responseText)

	// JSONをパース
	output, err := parsePfcResponse(responseText)
	if err != nil {
		logger.Error("Failed to parse PFC response", "error", err, "raw", responseText)
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
