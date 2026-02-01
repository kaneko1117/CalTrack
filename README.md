# CalTrack

カロリー管理アプリケーション

## 技術スタック

### Backend
- **言語**: Go 1.24
- **フレームワーク**: Gin
- **ORM**: GORM
- **DB**: MySQL 8.0
- **マイグレーション**: sql-migrate
- **ホットリロード**: Air
- **APIドキュメント**: Swagger (swaggo)

### Frontend
- **言語**: TypeScript
- **フレームワーク**: React + Vite
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
- **ルーティング**: React Router
- **HTTP**: Axios
- **テスト**: Vitest + React Testing Library
- **コンポーネントカタログ**: Storybook

### インフラ
- Docker Compose

## 必要条件

- Docker
- Docker Compose

## セットアップ

```bash
# コンテナ起動
make up

# コンテナ停止
make down
```

## アクセス

| サービス | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| Backend API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |
| Storybook | http://localhost:6006 |
| ヘルスチェック | http://localhost:8080/health |

## コマンド

`make help` で全コマンドを確認できます。

### 起動・停止

```bash
make up              # コンテナ起動
make down            # コンテナ停止
make restart         # コンテナ再起動
make clean           # コンテナとボリューム削除
```

### ログ

```bash
make logs            # 全サービスのログ
make logs-backend    # バックエンドのログ
make logs-frontend   # フロントエンドのログ
make logs-mysql      # MySQLのログ
```

### ビルド・テスト

```bash
make build           # 全サービスをビルド
make test            # 全テスト実行
make test-backend    # バックエンドテスト
make test-frontend   # フロントエンドテスト
make lint            # 全Lint実行
make fmt             # 全フォーマット実行
```

### ドキュメント

```bash
make swagger         # Swaggerドキュメント生成
make storybook       # Storybook起動（ポート6006）
make build-storybook # Storybookビルド
```

### マイグレーション

```bash
make migrate         # マイグレーション実行
make migrate-status  # マイグレーション状態確認
make migrate-down    # ロールバック（1つ戻す）
make migrate-new NAME=xxx  # 新規マイグレーション作成
```

### シェル

```bash
make shell-backend   # バックエンドコンテナに入る
make shell-frontend  # フロントエンドコンテナに入る
make shell-mysql     # MySQLコンテナに入る
```

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
│   ├── docs/             # Swagger自動生成
│   └── migrations/       # sql-migrate
├── frontend/
│   ├── .storybook/       # Storybook設定
│   └── src/
│       ├── features/     # 機能単位
│       │   └── {feature}/
│       │       ├── types/
│       │       ├── api/
│       │       ├── hooks/
│       │       └── components/
│       ├── components/ui/ # shadcn/ui
│       ├── pages/        # ページコンポーネント
│       ├── routes/       # React Router設定
│       ├── hooks/        # 共通Hooks
│       └── lib/          # ユーティリティ
└── Makefile
```

## 環境変数

| 変数 | 説明 | デフォルト値 |
|-----|------|------------|
| MYSQL_DATABASE | データベース名 | caltrack |
| MYSQL_USER | MySQL ユーザー名 | caltrack |
| MYSQL_PASSWORD | MySQL パスワード | caltrack |
| DB_HOST | データベースホスト | mysql |
| DB_PORT | データベースポート | 3306 |
| ENV | 環境 (production/test/development) | development |

## ポート

| サービス | ポート |
|---------|--------|
| Frontend | 5173 |
| Backend | 8080 |
| MySQL | 3307 |
| Storybook | 6006 |
