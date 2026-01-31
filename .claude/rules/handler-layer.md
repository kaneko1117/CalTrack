---
paths:
  - "backend/handler/**/*.go"
---

# Handler層規則

## ファイル構成

**ドメイン単位でディレクトリ**:
```
handler/
  common/
    error_code.go    # エラーコード定義
    response.go      # レスポンスヘルパー関数
  {domain}/
    handler.go       # {Domain}Handler
    handler_test.go
    dto/
      request.go     # リクエストDTO + ToDomain()
      response.go    # レスポンスDTO
```

## 命名規則

- 構造体: `{Domain}Handler`（例: `UserHandler`）
- メソッド: 動詞形（例: `Register`, `Login`, `Create`）
- コンストラクタ: `New{Domain}Handler`

## Request DTO

### 構造
- JSONタグ使用（`json:`）
- `ToDomain()` メソッドでEntityに変換

```go
func (r RegisterUserRequest) ToDomain() (*entity.User, error, []error) {
    // パース処理
    // entity.NewUser() を呼び出し
}
```

## Response DTO

### 構造
- JSONタグ使用
- Entityから必要なフィールドのみ抽出

## 共通ヘルパー（handler/common）

共通化できる処理はヘルパー関数として定義:

```go
// response.go
func RespondError(c echo.Context, status int, code, message string) error
func RespondValidationError(c echo.Context, details []string) error
```

## Handler

### 処理フロー
1. リクエストのバインド
2. `req.ToDomain()` でEntity変換
3. Usecase実行
4. レスポンス返却

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
