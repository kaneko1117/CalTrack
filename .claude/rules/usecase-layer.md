---
paths:
  - "backend/usecase/**/*.go"
---

# Usecase層規則

## ファイル構成

**ドメイン単位で1ファイル**:
- `usecase/{domain}.go` - Usecase実装
- `usecase/{domain}_test.go` - テスト

```
usecase/
  user.go          # UserUsecase（Register, Login等のメソッドを持つ）
  user_test.go
  meal.go          # MealUsecase
  meal_test.go
```

## 命名規則

- 構造体: `{Domain}Usecase`（例: `UserUsecase`）
- メソッド: 動詞形（例: `Register`, `Login`, `Create`）
- コンストラクタ: `New{Domain}Usecase`

## 入出力

- **入力**: Handler層で変換された`*entity.{Domain}`を受け取る
- **出力**: `*entity.{Domain}`をそのまま返す（Output DTOは不要）

## Repository Interface

### 定義場所
- `domain/repository/{entity}_repository.go`

### メソッド規則
- 第一引数は `context.Context`
- 戻り値の最後は `error`
- Entity/VOを引数・戻り値に使用

### 標準メソッド
```go
Save(ctx context.Context, entity *Entity) error
FindByID(ctx context.Context, id ID) (*Entity, error)
Delete(ctx context.Context, id ID) error
```

## Output Port Interface

### 定義場所
- `usecase/port/{usecase}_output_port.go`

### メソッド
- `Success(ctx, output)`: 成功時
- `Error(ctx, err)`: エラー時

## Usecase

### 構造
- 依存はInterface経由で注入
- 1 Usecase = 1 ユーザーアクション

### 処理フロー
1. Inputバリデーション
2. Repository経由でデータ取得
3. Entityの振る舞い呼び出し
4. Repository経由でデータ保存
5. Output生成・返却

### 禁止事項
- 具体的なDB操作
- HTTP関連の処理
- 外部サービスの直接呼び出し
