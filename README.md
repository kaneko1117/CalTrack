# CalTrack

カロリー追跡アプリケーション

## 技術スタック

- **フロントエンド**: React + Vite + TypeScript + shadcn/ui + Tailwind CSS
- **バックエンド**: Go + Gin + GORM
- **データベース**: MySQL 8.0
- **インフラ**: Docker Compose

## 必要条件

- Docker
- Docker Compose

## セットアップ

### 1. コンテナの起動

```bash
docker compose up --build
```

### 2. アクセス

| サービス | URL |
|---------|-----|
| フロントエンド | http://localhost:5173 |
| バックエンド API | http://localhost:8080 |
| ヘルスチェック | http://localhost:8080/health |

## 開発

### ホットリロード

- **フロントエンド**: Vite の HMR により、コードを変更すると自動的にブラウザが更新されます
- **バックエンド**: Air により、Go コードを変更すると自動的にサーバーが再起動されます

### ディレクトリ構成

```
CalTrack/
├── docker-compose.yml
├── frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── components.json          # shadcn/ui設定
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   └── ui/              # shadcn/uiコンポーネント
│   │   └── lib/
│   │       └── utils.ts
│   └── index.html
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── config/
│   │   └── database.go
│   ├── models/
│   │   └── user.go              # サンプルモデル
│   └── handlers/
│       └── health.go
└── README.md
```

### 環境変数

| 変数 | 説明 | デフォルト値 |
|-----|------|------------|
| MYSQL_ROOT_PASSWORD | MySQL root パスワード | rootpassword |
| MYSQL_DATABASE | データベース名 | caltrack |
| MYSQL_USER | MySQL ユーザー名 | caltrack |
| MYSQL_PASSWORD | MySQL パスワード | caltrack |
| DB_HOST | データベースホスト | mysql |
| DB_PORT | データベースポート | 3306 |

## コマンド

### コンテナの起動
```bash
docker compose up --build
```

### コンテナの停止
```bash
docker compose down
```

### コンテナとボリュームの削除
```bash
docker compose down -v
```

### ログの確認
```bash
# 全サービス
docker compose logs -f

# 特定のサービス
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
```
