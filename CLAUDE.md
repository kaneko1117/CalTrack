# CalTrack

カロリー管理アプリケーション

## 必須ルール

**マルチエージェント構成を使用すること**

詳細は `.claude/agents/workflow.md` を参照。

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

## 詳細ルール

- **アーキテクチャ**: `.claude/rules/architecture.md`
- **コーディング規約**: `.claude/rules/common.md`
- **各層規則**: `.claude/rules/{layer}-layer.md`
- **環境変数ポリシー**: `.claude/rules/env-file-policy.md`
- **マルチエージェント**: `.claude/agents/workflow.md`
