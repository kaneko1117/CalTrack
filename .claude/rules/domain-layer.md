---
paths:
  - "backend/domain/**/*.go"
---

# Domain層規則

## Value Object (VO)

### 構造
- フィールドは非公開（小文字始まり）
- ファクトリ関数 `New{VO名}()` で生成
- 不変（immutable）であること

### バリデーション
- ファクトリ関数内でバリデーション
- 無効な値でのインスタンス生成を防ぐ
- エラーは `domain/errors` で定義

### 禁止事項
- setter メソッド
- フレームワーク依存のタグ（`gorm:`, `json:` 等）
- 外部パッケージへの依存

## Entity

### 構造
- 識別子（ID）を持つ
- フィールドは非公開
- Getterで値を公開

### ファクトリ関数
- `New{Entity}()`: 新規作成（バリデーションあり）
- `Reconstruct{Entity}()`: DB復元用（バリデーションなし）

### 振る舞い
- ドメインロジックはEntityのメソッドとして実装
- 状態変更は専用メソッド経由
- 不変条件を常に維持

### 禁止事項
- Repository呼び出し
- 他Entityの直接生成（IDのみ保持）
- フレームワーク依存

## エラー定義

- `domain/errors/errors.go` に集約
- `var Err{Name} = errors.New("{message}")` 形式
- ドメイン固有のエラーのみ定義
