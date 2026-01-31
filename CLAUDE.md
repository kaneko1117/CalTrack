# CalTrack

カロリー管理アプリケーション

## 必須ルール

**機能実装時は必ずサブエージェントを使用すること**

Issue番号を受け取ったら、直接実装せず必ず以下のフローで進める:

```
1. orchestrator エージェントを起動
2. 各層の設計エージェント（domain-design, usecase-design等）で設計
3. ユーザー承認を得る
4. impl エージェントで実装
5. test-pr エージェントでテスト・PR作成
```

**禁止事項**:
- サブエージェントを使わずに直接設計・実装すること
- orchestratorを経由せずに個別の層を実装すること

## 技術スタック

### Backend
- **言語**: Go 1.24
- **フレームワーク**: Gin
- **ORM**: GORM
- **DB**: MySQL 8.0
- **ホットリロード**: Air

### Frontend
- **言語**: TypeScript
- **フレームワーク**: React + Vite
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
- **状態管理**: Jotai
- **HTTP**: Axios
- **テスト**: Vitest

## ディレクトリ構成

```
CalTrack/
├── backend/
│   ├── domain/           # VO, Entity, Domain Errors
│   │   ├── vo/
│   │   ├── entity/
│   │   └── errors/
│   ├── usecase/          # DTO, Interfaces, Usecases
│   │   ├── dto/
│   │   ├── repository/   # Repository Interface
│   │   ├── service/      # Service Interface
│   │   ├── port/         # Output Port Interface
│   │   └── {usecase}/
│   ├── infrastructure/   # 実装
│   │   ├── database/
│   │   ├── migration/
│   │   ├── model/
│   │   ├── repository/   # Repository Implementation
│   │   └── service/      # Service Implementation
│   └── handler/          # HTTP層
│       ├── request/
│       ├── response/
│       ├── handler/
│       ├── middleware/
│       └── router/
├── frontend/src/
│   ├── features/         # 機能単位
│   │   └── {feature}/
│   │       ├── types/
│   │       ├── api/
│   │       ├── hooks/
│   │       └── components/
│   ├── components/       # 共通コンポーネント
│   ├── hooks/            # 共通Hooks
│   ├── lib/              # ユーティリティ
│   └── types/            # 共通型定義
└── .claude/
    ├── agents/           # サブエージェント定義
    └── rules/            # コーディングルール
```

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

## サブエージェントシステム

### 使用方法

GitHub Issueに機能要件を作成し、Issue番号を指定:

```
#2 を実装して
```

### エージェントフロー

```
GitHub Issue → orchestrator → 各層設計 → 承認 → 子Issue → impl → test-pr → PR
```

### 利用可能エージェント

| エージェント | 用途 |
|-------------|------|
| orchestrator | 全体フロー管理 |
| domain-design | Domain層設計 |
| usecase-design | Usecase層設計 |
| infrastructure-design | Infrastructure層設計 |
| handler-design | Handler層設計 |
| frontend-data-design | Frontend Data Layer設計 |
| frontend-ui-design | Frontend UI Layer設計 |
| impl | コード実装 |
| test-pr | テスト・PR作成 |

## コーディング規約

### Backend (Go)

- `go fmt` でフォーマット
- テストファイル: `{file}_test.go`
- テストパッケージ: `{package}_test`
- テスト命名: `Test{対象}_{条件}_{期待結果}`

### Frontend (TypeScript/React)

- ESLint / Prettier でフォーマット
- `any` 禁止
- テストファイル: `{file}.test.ts(x)`
- Vitest 使用
- Jotai Atom命名: `{name}Atom`

### コンポーネント設計

| 種別 | 特徴 |
|-----|------|
| Container | Hooks使用、ロジック担当 |
| Presentational | propsのみ、再利用可能 |

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

## Issue / PR 運用

### 機能要件Issue
- ラベル: `feature`
- 内容: 機能仕様

### 設計Issue（子Issue）
- タイトル: `[設計] {機能名}: {層名}`
- ラベル: `design`
- 内容: 詳細設計
- 参照: `Parent: #{親Issue番号}`

### PR
- 形式: `Closes #{設計Issue番号}`
- マージ時: 設計Issueが自動クローズ
