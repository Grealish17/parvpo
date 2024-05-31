ifeq ($(POSTGRES_SETUP_PROD),)
	POSTGRES_SETUP_PROD := user=user password=pass dbname=prod host=localhost port=8002 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/
MOCKGEN_TAG=1.6.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/migrations

.PHONY: .migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: .prod-migration-up
prod-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" up

.PHONY: .prod-migration-down
prod-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" down

.PHONY: .prod-dc-up
prod-dc-up:
	sudo docker-compose up -d postgres zookeeper kafka1 kafka2 kafka3

.PHONY: .prod-dc-stop
prod-dc-stop:
	sudo docker-compose stop

.PHONY: .prod-dc-down
prod-dc-down:
	sudo docker-compose down

.PHONY: run-api
run-api:
	go build -o api.exe ./cmd/api
	sudo ./api.exe

.PHONY: run-app
run-app:
	go build -o app.exe ./cmd/app
	sudo ./app.exe