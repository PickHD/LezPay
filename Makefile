# Load .env configurations
ifneq (,$(wildcard ./auth/cmd/v1/.env))
    export
    MIGRATIONS_SOURCE = file://db/migrations
	DATABASE_URL = postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
endif

# Run & Build entire services 
run:
	docker compose up -d --build --force-recreate
	
# Stop entire services
stop:
	docker compose down

# Remove entire services
remove:
	docker compose down -v

# Run golang-migrate up
migrate-up:
	migrate -source $(MIGRATIONS_SOURCE) -database "$(DATABASE_URL)" up

# Run golang-migrate down
migrate-down:
	migrate -source $(MIGRATIONS_SOURCE) -database "$(DATABASE_URL)" down -all

# Run golang-migrate create
migrate-create:
	migrate create -ext sql -dir ./db/migrations $(filter-out $@,$(MAKECMDGOALS))

# Run linters
lint:
	golangci-lint run
	
.PHONY: run stop remove