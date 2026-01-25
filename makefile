CONFIG_PATH=./config/local.yml

MIGRATION_PATH=migrations
DB_USER=pics
DB_PASS=12345678
DB_HOST=localhost
DB_NAME=pics
DB_PORT=8002

up:
	CONFIG_PATH=$(CONFIG_PATH) go run ./cmd/picCompres/main.go

migrate:
	MIGRATIONS_PATH=$(MIGRATION_PATH) \
    PG_USER=$(DB_USER) \
    PG_PASS=$(DB_PASS) \
    PG_HOST=$(DB_HOST) \
    PG_PORT=$(DB_PORT) \
    PG_DB=$(DB_NAME) \
    go run ./cmd/migrate/main.go

