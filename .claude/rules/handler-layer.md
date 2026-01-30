---
paths:
  - "backend/handler/**/*.go"
---

# Handler層規則

## Request DTO

### 定義場所
- `handler/request/{handler}_request.go`

### 構造
- JSONタグ使用（`json:`）
- バリデーションタグ使用（`binding:`）
- path/query/body/headerからのマッピング

## Response DTO

### 定義場所
- `handler/response/{handler}_response.go`

### 構造
- JSONタグ使用
- Usecase Outputから変換

### エラーレスポンス
```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

## Handler

### 定義場所
- `handler/handler/{handler}.go`

### 構造
- Usecaseを依存として持つ
- Output Portを実装

### 処理フロー
1. リクエストのバインド
2. バリデーション
3. Usecase Input生成
4. Usecase実行
5. レスポンス返却

### エラーハンドリング
| ドメインエラー | HTTPステータス |
|--------------|---------------|
| ErrNotFound | 404 |
| ErrValidation | 400 |
| ErrUnauthorized | 401 |
| ErrForbidden | 403 |
| ErrConflict | 409 |
| その他 | 500 |

## Router

### 定義場所
- `handler/router/router.go`

### 規則
- RESTful設計
- バージョニング: `/api/v1/`
- リソース名は複数形

## 禁止事項

- ビジネスロジックの実装
- Repository直接呼び出し
- Infrastructure層への依存
