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

### 方針
- **基本的にドメイン（Entity/VO）を使用する**
- Usecase専用のInput/Output DTOは原則として定義しない
- どうしても必要な場合のみ、Usecase層内にDTOを定義してよい
  - 例: 複数Entityをまとめて返す必要がある場合
  - 例: ドメインに存在しない集計結果を返す場合

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

### 集計を含む処理フロー
- 基本的な集計はRepository経由でSQL集計結果を取得する（`architecture.md` の集計方針を参照）
- 取得した集計結果に対してビジネスルールを適用する場合は、Domainのメソッドを呼び出す
- Usecase層自体では集計ロジックを実装しない

### 禁止事項
- 具体的なDB操作
- HTTP関連の処理
- 外部サービスの直接呼び出し
- Usecase層での集計ロジック実装（SQLまたはDomainに委譲する）

## テスト

### モック方針
- **`backend/mock/` の自動生成モック（uber-go/mock）を使用すること**（手書きモック禁止）
- モック生成: `make mock-gen`
- Repository / Service のインターフェースに対するモックが `backend/mock/` に生成済み

## Service Interface（AI連携等）

### プロンプトの責務
- AIサービス（Gemini等）へのプロンプトはUsecase層で構築する
- Infrastructure層はプロンプトをConfigで受け取り、APIに送信するのみ
- ビジネスロジック（何を聞くか）はUsecase層、技術的実装（どう聞くか）はInfrastructure層
