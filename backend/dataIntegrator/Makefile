.PHONY:  \
	build \
    checks \
    compush \
    help \
    init \
    lint \
    migrator-clickhouse-create \
    migrator-clickhouse-down \
    migrator-clickhouse-down-all \
    migrator-clickhouse-status \
    migrator-clickhouse-up \
    migrator-postgres-create \
    migrator-postgres-down \
    migrator-postgres-down-all \
    migrator-postgres-status \
    mocks \
    proto-gen \
    stage \
    start \
    stop \
    test

CHECK_EMOJI := ✅
INFO_EMOJI := ℹ️

MIGRATION_NAME ?= migration

MIGRATIONS_PATH_POSTGRES := migrations/postgres
DB_URL_POSTGRES := postgres://admin:password@localhost:5432/preorders?sslmode=disable

lint:
	@echo "[INFO] Running golangci-lint..."
	@golangci-lint run
	@echo "[OK] Lint completed successfully!"

migrator-postgres-create:
	@echo "$(INFO_EMOJI) Creating new migration for postgres '$(MIGRATION_NAME)'..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) create $(MIGRATION_NAME) sql
	@echo "$(CHECK_EMOJI) Migration files created!"

migrator-postgres-up:
	@echo "$(INFO_EMOJI) Applying migrations for postgres UP..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" up
	@echo "$(CHECK_EMOJI) Migrations applied!"

migrator-postgres-down:
	@echo "$(INFO_EMOJI) Rolling back LAST migration for postgres..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" down
	@echo "$(CHECK_EMOJI) Rolled back one migration!"

migrator-postgres-down-all:
	@echo "$(INFO_EMOJI) Rolling back ALL migrations for postgres..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" reset
	@echo "$(CHECK_EMOJI) Rolled back all migrations!"

migrator-postgres-status:
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" status


proto-gen:
	@echo "[INFO] Generating protobuf code..."
	@protoc --version
	@protoc \
		--experimental_allow_proto3_optional \
		--proto_path= \
		--go_out=./api/grpc/gen \
		--go_opt=module=gl.iteco.com/technology/go_masterdata/preorderservice/api/grpc/gen \
		--go-grpc_out=./api/grpc/gen \
		--go-grpc_opt=module=gl.iteco.com/technology/go_masterdata/preorderservice/api/grpc/gen \
		./api/grpc/proto/preorders/preorders.proto \
		./api/grpc/proto/offers/offers.proto \
		./api/grpc/proto/trades/trades.proto
	@echo "[OK] Protobuf code generated successfully!"