.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version

# Migration variables
MIGRATIONS_PATH ?= migrations
DB_HOST ?= localhost
DB_PORT ?= 5221
DB_USER ?= user
DB_PASSWORD ?= password
DB_NAME ?= microservice_template_orders_db
MIGRATE_COMMAND = migrate -path $(MIGRATIONS_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# run all pending migrations
migrate-up:
	$(MIGRATE_COMMAND) up

# run a specific number of migrations
migrate-up-one:
	$(MIGRATE_COMMAND) up 1

# rollback all migrations
migrate-down:
	$(MIGRATE_COMMAND) down

# rollback a specific number of migrations
migrate-down-one:
	$(MIGRATE_COMMAND) down 1

# create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $${name}

# force set migration version
migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE_COMMAND) force $${version}

# show current migration version
migrate-version:
	$(MIGRATE_COMMAND) version

# WARN: development/testing only - dangerous for production!
migrate-reset:
	PGPASSWORD=$(DB_PASSWORD) dropdb --host=$(DB_HOST) --port=$(DB_PORT) --username=$(DB_USER) $(DB_NAME) --if-exists
	PGPASSWORD=$(DB_PASSWORD) createdb --host=$(DB_HOST) --port=$(DB_PORT) --username=$(DB_USER) $(DB_NAME)	
	$(MIGRATE_COMMAND) up


