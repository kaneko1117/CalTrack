# CalTrack

カロリー管理アプリケーション

## 必須ルール

**全ての開発タスクはサブエージェントを使用すること**

どんなプロンプトを受け取っても、直接実装せず必ず以下のフローで進める:

```
1. orchestrator エージェントを起動（全プロンプトを受け付ける司令塔）
2. タスク種別を判断し、適切なサブエージェントに振り分け
3. ユーザー承認を得る
4. 実装・テスト・PR作成
```

**禁止事項**:
- サブエージェントを使わずに直接設計・実装すること
- orchestratorを経由せずに個別の層を実装すること

## サブエージェントシステム

### エージェント一覧

| エージェント | 役割 | 参照rules |
|-------------|------|----------|
| orchestrator | 全プロンプト受付・タスク振り分け | clean-architecture, coding |
| design | 各層の設計 | clean-architecture, {layer}-layer, coding |
| impl | コード実装 | {layer}-layer, coding |
| test | Build・Test実行 | coding |
| pr | PR作成・マージ | coding |
| refactor | リファクタリング | clean-architecture, {layer}-layer, coding |

### タスク振り分け

| タスク種別 | キーワード例 | 振り分け先 |
|-----------|-------------|-----------|
| 新機能実装 | `#123`, `Issue` | 機能実装フロー（全層） |
| 設計 | `設計して`, `構成を考えて` | design |
| 実装 | `実装して`, `追加して` | impl → test → pr |
| テスト | `テストして`, `ビルド確認` | test |
| PR作成 | `PR作って`, `マージして` | pr |
| リファクタリング | `リファクタ`, `改善して` | refactor → test → pr |

### 基本フロー

```
design → 【設計承認】→ impl → test → 【テスト結果承認】→ pr → マージ
```

### 機能実装フロー（GitHub Issue）

```
Issue #{n} → orchestrator → Backend各層 → Frontend各層 → 完了
                              ↓
                        Domain → Usecase → Infrastructure → Handler
                              ↓
                        Data Layer → UI Layer
```

---

## 技術スタック

### Backend
- **言語**: Go 1.24
- **フレームワーク**: Gin
- **ORM**: GORM
- **DB**: MySQL 8.0
- **マイグレーション**: sql-migrate
- **ホットリロード**: Air

### Frontend
- **言語**: TypeScript
- **フレームワーク**: React + Vite
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
- **ルーティング**: React Router
- **HTTP**: Axios
- **テスト**: Vitest + React Testing Library

---

## ディレクトリ構成

```
CalTrack/
├── backend/
│   ├── domain/           # VO, Entity, Domain Errors
│   │   ├── vo/
│   │   ├── entity/
│   │   └── errors/
│   ├── usecase/          # Usecase（ドメイン単位）
│   ├── infrastructure/   # 実装
│   │   ├── persistence/  # Repository実装
│   │   └── service/      # Service実装
│   ├── handler/          # HTTP層（ドメイン単位）
│   │   ├── common/       # 共通エラーコード・レスポンス
│   │   └── {domain}/     # dto/, handler.go
│   ├── config/           # DB設定
│   └── migrations/       # sql-migrate
├── frontend/src/
│   ├── features/         # 機能単位
│   │   └── {feature}/
│   │       ├── types/    # 型定義・定数
│   │       ├── api/      # API関数
│   │       ├── hooks/    # カスタムフック
│   │       └── components/
│   ├── components/ui/    # shadcn/ui
│   ├── pages/            # ページコンポーネント
│   ├── routes/           # React Router設定
│   ├── hooks/            # 共通Hooks
│   └── lib/              # ユーティリティ
└── .claude/
    ├── agents/           # サブエージェント定義
    └── rules/            # コーディングルール
```

---

## 開発環境

### 起動
```bash
docker compose up --build
```

### ポート
| サービス | ポート |
|---------|--------|
| Frontend | 5173 |
| Backend | 8080 |
| MySQL | 3307 |

### ヘルスチェック
- Frontend: http://localhost:5173
- Backend: http://localhost:8080/health

---

## アーキテクチャ

### Clean Architecture

依存方向: Handler → Usecase → Domain ← Infrastructure

```
┌─────────────────────────────────────────┐
│              Handler                     │
├─────────────────────────────────────────┤
│              Usecase                     │
├─────────────────────────────────────────┤
│              Domain                      │
├─────────────────────────────────────────┤
│           Infrastructure                 │
└─────────────────────────────────────────┘
```

### 層間ルール

| 層 | 依存可能 | 依存禁止 |
|----|---------|---------|
| Domain | なし | 全ての外部層 |
| Usecase | Domain | Infrastructure, Handler |
| Infrastructure | Domain, Usecase (Interface) | Handler |
| Handler | Usecase (Interface) | Infrastructure |

---

## コーディング規約

### 共通
- **コメントは日本語で記述**
- マジックナンバー禁止（定数化）
- エラーコード・メッセージは定数化
- Usecase, Handler はドメイン単位で1ファイル

### Backend (Go)
- `go fmt` でフォーマット
- テストファイル: `{file}_test.go`
- テストパッケージ: `{package}_test`

### Frontend (TypeScript/React)
- ESLint / Prettier でフォーマット
- `any` 禁止
- **`interface` ではなく `type` を使用**
- テストファイル: `{file}.test.ts(x)`
- 定数: `ERROR_CODE_XXX`, `ERROR_MESSAGE_XXX`
- **必須スキル参照**: `.claude/skills/` 配下のスキルを必ず読む
  - `vercel-react-best-practices` - パフォーマンス最適化
  - `vercel-composition-patterns` - コンポーネント構成
  - `web-design-guidelines` - Webデザイン

---

## コマンド

### Backend
```bash
cd backend && go build ./...     # ビルド
cd backend && go test ./... -v   # テスト
```

### Frontend
```bash
cd frontend && npm run build     # ビルド
cd frontend && npm run test      # テスト
cd frontend && npm run lint      # Lint
```

---

## Issue / PR 運用

### 機能要件Issue
- ラベル: `feature`
- 内容: 機能仕様

### 設計Issue（子Issue）
- タイトル: `feat({layer}): {要約}`
- 内容: 詳細設計
- 参照: `Closes #{親Issue番号}`

### PR
- 形式: `Closes #{設計Issue番号}`
- マージ時: 設計Issueが自動クローズ
