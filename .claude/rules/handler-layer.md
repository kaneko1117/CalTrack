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
- `ToDomain()` メソッドで **Entity または VO** に変換
- プリミティブ → VO/Entity 変換は **Handler層（DTOのToDomain()）** で行う
- Usecase層にプリミティブやInput DTOを渡さない
- バリデーションエラーはHandler層で処理してHTTPレスポンスを返す

## Response DTO

### 構造
- JSONタグ使用
- Entityから必要なフィールドのみ抽出

## 共通ヘルパー（handler/common）

共通化できる処理はヘルパー関数として定義（`RespondError`, `RespondValidationError` 等）

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

## テスト

### モック方針
- **Usecaseインターフェースに対するモックを使用すること**（Repository層のモックは使わない）
- 各handler.goに定義された `{Domain}UsecaseInterface` に対してモックを作成
- テストの関心事はHandler層のみ（HTTPリクエスト/レスポンスのマッピング）

## 禁止事項

- ビジネスロジックの実装
- Repository直接呼び出し
- Infrastructure層への依存
