# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  seeder FILENAME     Process the specified seeder file name"

# application repository and binary file name
NAME=Bougette

# application repository path
REPOSITORY=github.com/projects/${NAME}


install:
	go mod download

run-dev:
	echo "Starting Application In Development Mode"
	#gow run ./cmd/api
	go run ./cmd/api


migrate:
	echo "Running migrations up"
	go run ./internal/database/migrate_up.go

#E.G make seeder FILENAME=category
.PHONY: seeder
seeder:
ifdef FILENAME
	echo "Seedinng : $(FILENAME).go"
	go run "./internal/database/seeders/$(FILENAME)_seeder.go"
else
	echo "Error: FILENAME is not specified. Please provide the filename using 'make seeder FILENAME=<filename>'"
	exit 1
endif

migrate_fresh:
	echo "Running fresh migrations"
	go run ./internal/database/migrate_fresh.go

migrate_down:
	go run ./migrations/drop_migrations.go

test:
	 gotest ./tests/... -v