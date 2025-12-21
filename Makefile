.PHONY: up down restart build logs ps clean help init test test-unit test-e2e test-coverage test-unit-coverage

# .envファイルの初期化
init:
	@if [ ! -f .env ]; then \
		cp .env.sample .env; \
		echo ".env ファイルを作成しました"; \
	else \
		echo ".env ファイルは既に存在します"; \
	fi

# Docker Composeの起動
up:
	docker compose up -d

# Docker Composeの停止
down:
	docker compose down

# Docker Composeの再起動
restart:
	docker compose restart

# Docker Composeのビルド
build:
	docker compose build

# Docker Composeのビルドして起動
up-build:
	docker compose up -d --build

# ログの表示
logs:
	docker compose logs -f

# コンテナの状態確認
ps:
	docker compose ps

# コンテナとボリュームの削除
clean:
	docker compose down -v

# APIコンテナに接続
api-shell:
	docker compose exec api sh

# DBコンテナに接続
db-shell:
	docker compose exec db mysql -uroot -p$$(grep MYSQL_ROOT_PASSWORD .env | cut -d '=' -f2) $$(grep MYSQL_DATABASE .env | cut -d '=' -f2)

# テストデータ作成
seed:
	docker compose exec api go run tool/seed/main.go

# Mockファイル作成
mock:
	mockery

# 単体テスト実行
test:
	go test -v $$(go list ./... | grep -v /e2e)

# E2Eテスト実行
test-e2e:
	go test -v ./e2e/...

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo ""
	@echo "【Docker関連】"
	@echo "  make init             - .env ファイルを .env.sample から作成"
	@echo "  make up               - Docker Composeを起動"
	@echo "  make down             - Docker Composeを停止"
	@echo "  make restart          - Docker Composeを再起動"
	@echo "  make build            - Docker Composeをビルド"
	@echo "  make up-build         - ビルドしてから起動"
	@echo "  make logs             - ログを表示"
	@echo "  make ps               - コンテナの状態確認"
	@echo "  make clean            - コンテナとボリュームを削除"
	@echo "  make api-shell        - APIコンテナに接続"
	@echo "  make db-shell         - DBコンテナに接続"
	@echo ""
	@echo "【開発関連】"
	@echo "  make seed             - テストデータ作成"
	@echo "  make mock             - Mockファイル作成"
	@echo ""
	@echo "【テスト関連】"
	@echo "  make test             - 単体テスト実行（e2e以外）"
	@echo "  make test-e2e         - E2Eテスト実行"
	@echo ""
	@echo "【その他】"
	@echo "  make help             - このヘルプを表示"
