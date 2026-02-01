.PHONY: help up down restart logs logs-backend logs-frontend logs-mysql \
        build build-backend build-frontend \
        test test-backend test-frontend \
        lint lint-backend lint-frontend fmt fmt-backend fmt-frontend \
        migrate migrate-status migrate-down migrate-new \
        shell-backend shell-frontend shell-mysql clean

# =============================================================================
# ヘルプ
# =============================================================================

help:
	@echo "CalTrack - 利用可能なコマンド"
	@echo ""
	@echo "起動・停止:"
	@echo "  make up              - コンテナ起動"
	@echo "  make down            - コンテナ停止"
	@echo "  make restart         - コンテナ再起動"
	@echo "  make clean           - コンテナとボリューム削除"
	@echo ""
	@echo "ログ:"
	@echo "  make logs            - 全サービスのログ"
	@echo "  make logs-backend    - バックエンドのログ"
	@echo "  make logs-frontend   - フロントエンドのログ"
	@echo "  make logs-mysql      - MySQLのログ"
	@echo ""
	@echo "ビルド:"
	@echo "  make build           - 全サービスをビルド"
	@echo "  make build-backend   - バックエンドをビルド"
	@echo "  make build-frontend  - フロントエンドをビルド"
	@echo ""
	@echo "テスト:"
	@echo "  make test            - 全テスト実行"
	@echo "  make test-backend    - バックエンドテスト"
	@echo "  make test-frontend   - フロントエンドテスト"
	@echo ""
	@echo "Lint・フォーマット:"
	@echo "  make lint            - 全Lint実行"
	@echo "  make lint-backend    - バックエンドLint"
	@echo "  make lint-frontend   - フロントエンドLint"
	@echo "  make fmt             - 全フォーマット実行"
	@echo "  make fmt-backend     - バックエンドフォーマット"
	@echo "  make fmt-frontend    - フロントエンドフォーマット"
	@echo ""
	@echo "マイグレーション:"
	@echo "  make migrate         - マイグレーション実行"
	@echo "  make migrate-status  - マイグレーション状態確認"
	@echo "  make migrate-down    - ロールバック（1つ戻す）"
	@echo "  make migrate-new NAME=xxx - 新規マイグレーション作成"
	@echo ""
	@echo "シェル:"
	@echo "  make shell-backend   - バックエンドコンテナに入る"
	@echo "  make shell-frontend  - フロントエンドコンテナに入る"
	@echo "  make shell-mysql     - MySQLコンテナに入る"

# =============================================================================
# 起動・停止
# =============================================================================

up:
	docker compose up --build -d

down:
	docker compose down

restart: down up

clean:
	docker compose down -v

# =============================================================================
# ログ
# =============================================================================

logs:
	docker compose logs -f

logs-backend:
	docker compose logs -f backend

logs-frontend:
	docker compose logs -f frontend

logs-mysql:
	docker compose logs -f mysql

# =============================================================================
# ビルド
# =============================================================================

build:
	docker compose build

build-backend:
	docker compose build backend

build-frontend:
	docker compose build frontend

# =============================================================================
# テスト
# =============================================================================

test: test-backend test-frontend

test-backend:
	docker compose exec -e ENV=test backend gotestsum --format testname -- ./...

test-frontend:
	docker compose exec frontend npm run test

# =============================================================================
# Lint・フォーマット
# =============================================================================

lint: lint-backend lint-frontend

lint-backend:
	docker compose exec backend go vet ./...

lint-frontend:
	docker compose exec frontend npm run lint

fmt: fmt-backend fmt-frontend

fmt-backend:
	docker compose exec backend go fmt ./...

fmt-frontend:
	docker compose exec frontend npm run format

# =============================================================================
# マイグレーション
# =============================================================================

migrate:
	docker compose exec backend sql-migrate up -env=development

migrate-status:
	docker compose exec backend sql-migrate status -env=development

migrate-down:
	docker compose exec backend sql-migrate down -env=development -limit=1

migrate-new:
	cd backend && sql-migrate new $(NAME)

# =============================================================================
# シェル
# =============================================================================

shell-backend:
	docker compose exec backend sh

shell-frontend:
	docker compose exec frontend sh

shell-mysql:
	docker compose exec mysql mysql -u caltrack -pcaltrack caltrack
