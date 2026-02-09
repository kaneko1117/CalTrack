# =============================================================================
# 変数定義
# =============================================================================

# Docker Compose ファイル
COMPOSE_DEV = docker compose -f docker-compose.yml -f docker-compose.dev.yml
COMPOSE_PROD = docker compose -f docker-compose.yml -f docker-compose.prod.yml

.PHONY: help up down restart logs logs-backend logs-frontend logs-mysql \
        build build-backend build-frontend \
        test test-backend test-frontend \
        lint lint-backend lint-frontend fmt fmt-backend fmt-frontend \
        migrate migrate-status migrate-down migrate-new \
        shell-backend shell-frontend shell-mysql clean \
        swagger mock-gen mock-clean storybook build-storybook \
        up-prod down-prod

# =============================================================================
# ヘルプ
# =============================================================================

help:
	@echo "CalTrack - 利用可能なコマンド"
	@echo ""
	@echo "起動・停止（開発環境）:"
	@echo "  make up              - 開発用コンテナ起動"
	@echo "  make down            - 開発用コンテナ停止"
	@echo "  make restart         - 開発用コンテナ再起動"
	@echo "  make clean           - コンテナとボリューム削除"
	@echo ""
	@echo "起動・停止（本番環境）:"
	@echo "  make up-prod         - 本番用コンテナ起動"
	@echo "  make down-prod       - 本番用コンテナ停止"
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
	@echo ""
	@echo "Swagger:"
	@echo "  make swagger         - Swaggerドキュメント生成"
	@echo ""
	@echo "モック:"
	@echo "  make mock-gen        - モックファイル生成"
	@echo "  make mock-clean      - モックファイル削除"
	@echo ""
	@echo "Storybook:"
	@echo "  make storybook       - Storybook起動（ポート6006）"
	@echo "  make build-storybook - Storybookビルド"

# =============================================================================
# 起動・停止（開発環境）
# =============================================================================

up:
	$(COMPOSE_DEV) up --build -d

down:
	$(COMPOSE_DEV) down

restart: down up

clean:
	$(COMPOSE_DEV) down -v

# =============================================================================
# 起動・停止（本番環境）
# =============================================================================

up-prod:
	$(COMPOSE_PROD) up -d --build

down-prod:
	$(COMPOSE_PROD) down

# =============================================================================
# ログ
# =============================================================================

logs:
	$(COMPOSE_DEV) logs -f

logs-backend:
	$(COMPOSE_DEV) logs -f backend

logs-frontend:
	$(COMPOSE_DEV) logs -f frontend

logs-mysql:
	$(COMPOSE_DEV) logs -f mysql

# =============================================================================
# ビルド
# =============================================================================

build:
	$(COMPOSE_DEV) build

build-backend:
	$(COMPOSE_DEV) build backend

build-frontend:
	$(COMPOSE_DEV) build frontend

# =============================================================================
# テスト
# =============================================================================

test: test-backend test-frontend

test-backend:
	$(COMPOSE_DEV) exec -e ENV=test backend gotestsum --format testname -- ./...

test-frontend:
	$(COMPOSE_DEV) exec frontend npm run test

# =============================================================================
# Lint・フォーマット
# =============================================================================

lint: lint-backend lint-frontend

lint-backend:
	$(COMPOSE_DEV) exec backend go vet ./...

lint-frontend:
	$(COMPOSE_DEV) exec frontend npm run lint

fmt: fmt-backend fmt-frontend

fmt-backend:
	$(COMPOSE_DEV) exec backend go fmt ./...

fmt-frontend:
	$(COMPOSE_DEV) exec frontend npm run format

# =============================================================================
# マイグレーション
# =============================================================================

migrate:
	$(COMPOSE_DEV) exec backend sql-migrate up -env=development

migrate-status:
	$(COMPOSE_DEV) exec backend sql-migrate status -env=development

migrate-down:
	$(COMPOSE_DEV) exec backend sql-migrate down -env=development -limit=1

migrate-new:
	cd backend && sql-migrate new $(NAME)

# =============================================================================
# シェル
# =============================================================================

shell-backend:
	$(COMPOSE_DEV) exec backend sh

shell-frontend:
	$(COMPOSE_DEV) exec frontend sh

shell-mysql:
	$(COMPOSE_DEV) exec mysql mysql -u caltrack -pcaltrack caltrack

# =============================================================================
# Swagger
# =============================================================================

swagger:
	$(COMPOSE_DEV) exec backend swag init

# =============================================================================
# モック生成
# =============================================================================

MOCKGEN := $(shell go env GOPATH)/bin/mockgen

mock-gen:
	cd backend && mkdir -p mock
	cd backend && $(MOCKGEN) -source=domain/repository/user_repository.go -destination=mock/mock_user_repository.go -package=mock
	cd backend && $(MOCKGEN) -source=domain/repository/session_repository.go -destination=mock/mock_session_repository.go -package=mock
	cd backend && $(MOCKGEN) -source=domain/repository/record_repository.go -destination=mock/mock_record_repository.go -package=mock
	cd backend && $(MOCKGEN) -source=domain/repository/record_pfc_repository.go -destination=mock/mock_record_pfc_repository.go -package=mock
	cd backend && $(MOCKGEN) -source=domain/repository/advice_cache_repository.go -destination=mock/mock_advice_cache_repository.go -package=mock
	cd backend && $(MOCKGEN) -source=domain/repository/transaction.go -destination=mock/mock_transaction_manager.go -package=mock
	cd backend && $(MOCKGEN) -source=usecase/service/image_analyzer.go -destination=mock/mock_image_analyzer.go -package=mock
	cd backend && $(MOCKGEN) -source=usecase/service/pfc_analyzer.go -destination=mock/mock_pfc_analyzer.go -package=mock
	cd backend && $(MOCKGEN) -source=usecase/service/pfc_estimator.go -destination=mock/mock_pfc_estimator.go -package=mock
	cd backend && $(MOCKGEN) -source=usecase/ai_config.go -destination=mock/mock_ai_config.go -package=mock
	@echo "Mock generation completed."

mock-clean:
	rm -rf backend/mock/
	@echo "Mock files cleaned."

# =============================================================================
# Storybook
# =============================================================================

storybook:
	$(COMPOSE_DEV) exec frontend npm run storybook

build-storybook:
	$(COMPOSE_DEV) exec frontend npm run build-storybook
