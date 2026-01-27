BINARY=url_shorter
MIGRATIONS_DIR=./migrations

ifeq ($(OS),Windows_NT)
MIGRATE=.\tools\migrate.exe
else
MIGRATE=./tools/migrate
endif

include .env

//DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

build:
	go build -o $(BINARY) .

run:
	go run .

test:
	go test ./...

fmt:
	go fmt ./...
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

migrate-status:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

