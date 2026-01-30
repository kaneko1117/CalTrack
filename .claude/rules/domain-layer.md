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
- `New{Entity}()`: 新規作成（プリミティブ型を受け取り、内部でVO変換、エラーをまとめて返す）
- `Reconstruct{Entity}()`: DB復元用（VOを直接受け取る、バリデーションなし）

### VO変換パターン
- Entity内でプリミティブ→VO変換は専用のparse関数に分離
- 条件分岐のネストを避け、フラットに保つ
- エラーは`appendIfErr`等で集約し、最後にまとめて返す

```go
// 良い例
func NewUser(emailStr string, ...) (*User, []error) {
    var errs []error
    email, err := parseEmail(emailStr)
    errs = appendIfErr(errs, err)
    // ...
    if len(errs) > 0 {
        return nil, errs
    }
    return &User{...}, nil
}

func parseEmail(s string) (vo.Email, error) {
    return vo.NewEmail(s)
}
```

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
