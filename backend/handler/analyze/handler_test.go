package analyze_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/analyze"
	"caltrack/handler/analyze/dto"
	"caltrack/usecase"
)

// MockAnalyzeUsecase はAnalyzeUsecaseのモック
type MockAnalyzeUsecase struct {
	AnalyzeImageFunc func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error)
}

func (m *MockAnalyzeUsecase) AnalyzeImage(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
	return m.AnalyzeImageFunc(ctx, imageData, mimeType)
}

func setupTestRouter(handler *analyze.AnalyzeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 認証済みユーザーのミドルウェアを設定
	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	r.POST("/analyze-image", handler.AnalyzeImage)
	return r
}

func TestAnalyzeHandler_AnalyzeImage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAnalyzeUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "正常系_画像解析が成功する",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					name, _ := vo.NewItemName("唐揚げ")
					calories, _ := vo.NewCalories(250)
					return &usecase.AnalyzeOutput{
						Items: []usecase.AnalyzedItemOutput{
							{Name: name, Calories: calories},
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var res dto.AnalyzeImageResponse
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Len(t, res.Items, 1)
				assert.Equal(t, "唐揚げ", res.Items[0].Name)
				assert.Equal(t, 250, res.Items[0].Calories)
			},
		},
		{
			name: "正常系_複数の食品が検出される",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					name1, _ := vo.NewItemName("ご飯")
					calories1, _ := vo.NewCalories(250)
					name2, _ := vo.NewItemName("味噌汁")
					calories2, _ := vo.NewCalories(50)
					return &usecase.AnalyzeOutput{
						Items: []usecase.AnalyzedItemOutput{
							{Name: name1, Calories: calories1},
							{Name: name2, Calories: calories2},
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var res dto.AnalyzeImageResponse
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Len(t, res.Items, 2)
				assert.Equal(t, "ご飯", res.Items[0].Name)
				assert.Equal(t, 250, res.Items[0].Calories)
				assert.Equal(t, "味噌汁", res.Items[1].Name)
				assert.Equal(t, 50, res.Items[1].Calories)
			},
		},
		{
			name:           "異常系_リクエストボディが不正",
			requestBody:    "invalid json",
			setupMock:      func(m *MockAnalyzeUsecase) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Invalid request body")
			},
		},
		{
			name: "異常系_ImageDataが空",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					return nil, domainErrors.ErrImageDataRequired
				}
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), domainErrors.ErrImageDataRequired.Error())
			},
		},
		{
			name: "異常系_MimeTypeが空",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					return nil, domainErrors.ErrMimeTypeRequired
				}
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), domainErrors.ErrMimeTypeRequired.Error())
			},
		},
		{
			name: "異常系_食品が検出されなかった",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					return nil, domainErrors.ErrNoFoodDetected
				}
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), domainErrors.ErrNoFoodDetected.Error())
			},
		},
		{
			name: "異常系_画像解析に失敗",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					return nil, domainErrors.ErrImageAnalysisFailed
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Image analysis failed")
			},
		},
		{
			name: "異常系_その他のエラー",
			requestBody: dto.AnalyzeImageRequest{
				ImageData: "base64encodedimage",
				MimeType:  "image/jpeg",
			},
			setupMock: func(m *MockAnalyzeUsecase) {
				m.AnalyzeImageFunc = func(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error) {
					return nil, errors.New("unexpected error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			mockUsecase := &MockAnalyzeUsecase{}
			tt.setupMock(mockUsecase)

			// ハンドラの作成
			handler := analyze.NewAnalyzeHandler(mockUsecase)
			router := setupTestRouter(handler)

			// リクエストボディの作成
			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			// リクエストの作成と実行
			req := httptest.NewRequest(http.MethodPost, "/analyze-image", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスの確認
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestAnalyzeHandler_AnalyzeImage_Unauthorized(t *testing.T) {
	// 未認証のルーターをセットアップ
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockUsecase := &MockAnalyzeUsecase{}
	handler := analyze.NewAnalyzeHandler(mockUsecase)

	r.POST("/analyze-image", handler.AnalyzeImage)

	reqBody := dto.AnalyzeImageRequest{
		ImageData: "base64encodedimage",
		MimeType:  "image/jpeg",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/analyze-image", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not authenticated")
}
