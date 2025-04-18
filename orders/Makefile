.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version

# Migration variables
MIGRATIONS_PATH ?= migrations
DB_HOST ?= localhost
DB_PORT ?= 5221
DB_USER ?= user
DB_PASSWORD ?= password
DB_NAME ?= microservice_template_orders_db
MIGRATE_COMMAND = migrate -path $(MIGRATIONS_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Run all pending migrations
migrate-up:
	$(MIGRATE_COMMAND) up

# Run a specific number of migrations
migrate-up-one:
	$(MIGRATE_COMMAND) up 1

# Rollback all migrations
migrate-down:
	$(MIGRATE_COMMAND) down

# Rollback a specific number of migrations
migrate-down-one:
	$(MIGRATE_COMMAND) down 1

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $${name}

# Force set migration version
migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE_COMMAND) force $${version}

# Show current migration version
migrate-version:
	$(MIGRATE_COMMAND) version
