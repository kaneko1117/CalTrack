---
paths:
  - "backend/infrastructure/**/*.go"
---

# Infrastructure層規則

## Repository Implementation

### 定義場所
- `infrastructure/repository/{entity}_repository.go`

### 構造
- Usecase層のInterfaceを実装
- `*gorm.DB` を依存として持つ

### DBモデル
- `infrastructure/model/{entity}.go` に定義
- gormタグ使用可
- Entity ↔ Model 変換メソッドを持つ

### エラーマッピング
- DBエラー → ドメインエラーに変換
- `gorm.ErrRecordNotFound` → `ErrNotFound`
- duplicate entry → `ErrDuplicate`

## Migration

### 定義場所
- `infrastructure/migration/`

### ファイル命名
- `{version}_{description}.sql`
- 例: `001_create_users.sql`

### 構造
- `-- +migrate Up` セクション
- `-- +migrate Down` セクション

## 禁止事項

- ビジネスロジックの実装
- Domain層への依存以外の層への依存
- 環境変数のハードコード（configから取得）
