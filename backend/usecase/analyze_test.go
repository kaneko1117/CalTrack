package usecase_test

import (
	"context"
	"errors"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase"
	"caltrack/usecase/service"
)

// mockImageAnalyzer は画像解析サービスのモック実装
type mockImageAnalyzer struct {
	analyze func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error)
}

func (m *mockImageAnalyzer) Analyze(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
	return m.analyze(ctx, imageData, mimeType)
}

// TestAnalyzeUsecase_AnalyzeImage は画像解析機能のテスト
func TestAnalyzeUsecase_AnalyzeImage(t *testing.T) {
	t.Run("正常系_画像解析が成功し、結果が返る", func(t *testing.T) {
		// 解析結果のモック
		itemName1, err := vo.NewItemName("ハンバーガー")
		if err != nil {
			t.Fatalf("failed to create ItemName: %v", err)
		}
		calories1, err := vo.NewCalories(500)
		if err != nil {
			t.Fatalf("failed to create Calories: %v", err)
		}
		itemName2, err := vo.NewItemName("ポテト")
		if err != nil {
			t.Fatalf("failed to create ItemName: %v", err)
		}
		calories2, err := vo.NewCalories(300)
		if err != nil {
			t.Fatalf("failed to create Calories: %v", err)
		}

		mockAnalyzer := &mockImageAnalyzer{
			analyze: func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
				return []service.AnalyzedItem{
					{Name: itemName1, Calories: calories1},
					{Name: itemName2, Calories: calories2},
				}, nil
			},
		}

		uc := usecase.NewAnalyzeUsecase(mockAnalyzer)
		result, err := uc.AnalyzeImage(context.Background(), "base64encodedimage", "image/jpeg")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("result should not be nil")
		}
		if len(result.Items) != 2 {
			t.Errorf("got %d items, want 2", len(result.Items))
		}
		if result.Items[0].Name.String() != "ハンバーガー" {
			t.Errorf("got %s, want ハンバーガー", result.Items[0].Name.String())
		}
		if result.Items[0].Calories.Value() != 500 {
			t.Errorf("got %d, want 500", result.Items[0].Calories.Value())
		}
		if result.Items[1].Name.String() != "ポテト" {
			t.Errorf("got %s, want ポテト", result.Items[1].Name.String())
		}
		if result.Items[1].Calories.Value() != 300 {
			t.Errorf("got %d, want 300", result.Items[1].Calories.Value())
		}
	})

	t.Run("異常系_画像データが空の場合、ErrImageDataRequiredを返す", func(t *testing.T) {
		mockAnalyzer := &mockImageAnalyzer{
			analyze: func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
				return nil, nil
			},
		}

		uc := usecase.NewAnalyzeUsecase(mockAnalyzer)
		_, err := uc.AnalyzeImage(context.Background(), "", "image/jpeg")

		if err != domainErrors.ErrImageDataRequired {
			t.Errorf("got %v, want ErrImageDataRequired", err)
		}
	})

	t.Run("異常系_MIMEタイプが空の場合、ErrMimeTypeRequiredを返す", func(t *testing.T) {
		mockAnalyzer := &mockImageAnalyzer{
			analyze: func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
				return nil, nil
			},
		}

		uc := usecase.NewAnalyzeUsecase(mockAnalyzer)
		_, err := uc.AnalyzeImage(context.Background(), "base64encodedimage", "")

		if err != domainErrors.ErrMimeTypeRequired {
			t.Errorf("got %v, want ErrMimeTypeRequired", err)
		}
	})

	t.Run("異常系_解析結果が空の場合、ErrNoFoodDetectedを返す", func(t *testing.T) {
		mockAnalyzer := &mockImageAnalyzer{
			analyze: func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
				return []service.AnalyzedItem{}, nil
			},
		}

		uc := usecase.NewAnalyzeUsecase(mockAnalyzer)
		_, err := uc.AnalyzeImage(context.Background(), "base64encodedimage", "image/jpeg")

		if err != domainErrors.ErrNoFoodDetected {
			t.Errorf("got %v, want ErrNoFoodDetected", err)
		}
	})

	t.Run("異常系_画像解析サービスがエラーを返す場合", func(t *testing.T) {
		analyzeErr := errors.New("analysis service error")
		mockAnalyzer := &mockImageAnalyzer{
			analyze: func(ctx context.Context, imageData string, mimeType string) ([]service.AnalyzedItem, error) {
				return nil, analyzeErr
			},
		}

		uc := usecase.NewAnalyzeUsecase(mockAnalyzer)
		_, err := uc.AnalyzeImage(context.Background(), "base64encodedimage", "image/jpeg")

		if err != analyzeErr {
			t.Errorf("got %v, want analyzeErr", err)
		}
	})
}
