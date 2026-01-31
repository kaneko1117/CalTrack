# CalTrack

カロリー管理アプリケーション

## 必須ルール

**マルチエージェント構成を使用すること**

メイン（Claude）がチームメンバーを呼び出して作業を進める:

```
メイン（Claude）
  │
  ├─→ PM「このタスクお願い」
  │     ↓
  │   ← 「了解。まずプランナーに設計出してもらおう」
  │
  ├─→ プランナー呼び出し
  │     ↓
  │   ← 「設計まとまりました」
  │
  ├─→ PM「プランナーの結果これです」
  │     ↓
  │   ← 「いいね。ユーザーに確認取ってからエンジニアに振ろう」
  │
  └─→ ... 繰り返し
```

**重要:**
- メンバーは**他のメンバーを呼び出さない**
- メンバーは**結果を返すだけ**
- **メインがフロー制御を行う**

**禁止事項**:
- メンバーが他のメンバーを呼び出すこと
- PMを経由せずに個別の層を実装すること

## チームメンバー

| メンバー | 役割 | キャラクター |
|---------|------|-------------|
| PM | 司令塔、次の担当を判断 | 冷静で的確 |
| プランナー | 設計作成 | 慎重派、ドキュメント重視 |
| エンジニア | コード実装 | 素直、一生懸命 |
| QA | テスト実行 | 厳格、品質重視 |
| DevOps | PR作成・マージ | 手堅い、慎重 |
| 技術リード | リファクタリング | 職人気質 |

### タスク振り分け

| タスク種別 | キーワード例 | 最初に呼ぶメンバー |
|-----------|-------------|------------------|
| 新機能実装 | `#123`, `Issue` | PM |
| 設計 | `設計して`, `構成を考えて` | プランナー |
| 実装 | `実装して`, `追加して` | エンジニア |
| テスト | `テストして`, `ビルド確認` | QA |
| PR作成 | `PR作って`, `マージして` | DevOps |
| リファクタリング | `リファクタ`, `改善して` | 技術リード |

### 基本フロー

```
PM → プランナー → [承認] → エンジニア → QA → [承認] → DevOps
 ↑        │                                         │
 └────────┴────── 結果を返して次を判断 ─────────────┘
```

### 機能実装フロー（GitHub Issue）

```
Issue #{n}
    ↓
メイン → PM（次は？） → 「プランナーに設計出してもらおう」
    ↓
メイン → プランナー（設計作成） → 「設計まとまりました」
    ↓
メイン → PM（次は？） → 「いいね。ユーザーに確認取ろう」
    ↓
[ユーザー承認]
    ↓
メイン → PM（承認された） → 「エンジニアに実装振ろう」
    ↓
メイン → エンジニア（実装） → 「実装できました！」
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
    ├── agents/           # チームメンバー定義
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
- **テストケース名は日本語で記述**
- マジックナンバー禁止（定数化）
- エラーコード・メッセージは定数化
- Usecase, Handler はドメイン単位で1ファイル

### Backend (Go)
- `go fmt` でフォーマット
- テストファイル: `{file}_test.go`
- テストパッケージ: `{package}_test`
- **テストケース名は日本語**: `t.Run` の第1引数、テーブル駆動テストのnameは日本語

### Frontend (TypeScript/React)
- ESLint / Prettier でフォーマット
- `any` 禁止
- **`interface` ではなく `type` を使用**
- テストファイル: `{file}.test.ts(x)`
- 定数: `ERROR_CODE_XXX`, `ERROR_MESSAGE_XXX`
- **テストケース名は日本語**: `describe`, `it` の第1引数は日本語
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
