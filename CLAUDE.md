# CalTrack

カロリー管理アプリケーション

## 必須ルール

**マルチエージェント構成を使用すること**

メインエージェント（Claude）が全てのサブエージェント呼び出しを管理する:

```
メインエージェント（Claude）
  │
  ├─→ orchestrator「タスクを処理したい」
  │     ↓
  │   ← 返却「まず design を呼んで」
  │
  ├─→ design 呼び出し
  │     ↓
  │   ← 返却「設計完了」
  │
  ├─→ orchestrator「design の結果はこれ。次は？」
  │     ↓
  │   ← 返却「次は impl を呼んで」
  │
  └─→ ... 繰り返し
```

**重要:**
- サブエージェントは**他のサブエージェントを呼び出さない**
- サブエージェントは**結果を返すだけ**
- **メインがフロー制御を行う**

**禁止事項**:
- サブエージェントが他のサブエージェントを呼び出すこと
- orchestratorを経由せずに個別の層を実装すること

## マルチエージェントシステム

### エージェント一覧

| 種別 | エージェント | 役割 |
|-----|-------------|------|
| 司令塔 | orchestrator | タスク分析・次のエージェント判断（**指示を返すだけ**） |
| ワーカー | design | 各層の設計作成 |
| ワーカー | impl | コード実装 |
| ワーカー | test | Build・Test実行 |
| ワーカー | pr | PR作成・マージ |
| ワーカー | refactor | リファクタリング |

### タスク振り分け

| タスク種別 | キーワード例 | 最初に呼ぶエージェント |
|-----------|-------------|---------------------|
| 新機能実装 | `#123`, `Issue` | orchestrator |
| 設計 | `設計して`, `構成を考えて` | design |
| 実装 | `実装して`, `追加して` | impl |
| テスト | `テストして`, `ビルド確認` | test |
| PR作成 | `PR作って`, `マージして` | pr |
| リファクタリング | `リファクタ`, `改善して` | refactor |

### 基本フロー

```
orchestrator → design → [承認] → impl → test → [承認] → pr
     ↑            │                                    │
     └────────────┴────── 結果を返して次を判断 ────────┘
```

### 機能実装フロー（GitHub Issue）

```
Issue #{n}
    ↓
メイン → orchestrator（次は？） → 「design を呼んで」
    ↓
メイン → design（設計作成） → 設計完了
    ↓
メイン → orchestrator（次は？） → 「承認後 impl を呼んで」
    ↓
[ユーザー承認]
    ↓
メイン → impl（実装） → 実装完了
    ↓
... 繰り返し（Backend各層 → Frontend各層）
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
