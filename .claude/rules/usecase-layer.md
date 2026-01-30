---
paths:
  - "backend/usecase/**/*.go"
---

# Usecase層規則

## Input/Output DTO

### Input
- Usecaseへの入力データ構造
- バリデーションタグ使用可（`validate:`）
- Handler層からのみ生成

### Output
- Usecaseからの出力データ構造
- Entity/VOから変換して生成
- プリミティブ型またはDTO型のみ

## Repository Interface

### 定義場所
- `usecase/repository/{entity}_repository.go`

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
