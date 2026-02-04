package dto

// AnalyzeImageRequest は画像解析リクエストDTO
type AnalyzeImageRequest struct {
	ImageData string `json:"imageData"` // Base64エンコードされた画像データ
	MimeType  string `json:"mimeType"`  // 画像のMIMEタイプ（例: image/jpeg, image/png）
}
