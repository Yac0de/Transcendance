# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_FILE = ./srcs/docker-compose.yaml
DB_CONTAINER = postgres
DUMP_FILE = ./srcs/data/postgresql/db_backup.sql

include ./srcs/.env
export $(shell grep -v '^#' ./srcs/.env | xargs)

# Phony targets
.PHONY: all build up down restart clean  re prune

# Default target: build and start containers
all: build up

# Build Docker containers
build:
	@echo "Building Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build
	@echo "Docker containers built successfully."

# Start Docker containers
up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Containers are up and running."

# Stop and remove Docker containers without removing volumes
down:
	@echo "Stopping and removing Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "Containers stopped and removed."

# Exporter la base de données PostgreSQL dans un fichier SQL
dump:
	@echo "Exporting database to $(DUMP_FILE)..."
	@docker exec -t $(DB_CONTAINER) pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) > $(DUMP_FILE)
	@echo "Database exported to $(DUMP_FILE)."

# Restaurer la base de données à partir d'un fichier SQL
restore:
	@echo "Restoring database from $(DUMP_FILE)..."
	@cat $(DUMP_FILE) | docker exec -i $(DB_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)
	@echo "Database restored from $(DUMP_FILE)."

# Restart Docker containers
restart: down up

# Clean Docker resources
clean:
	@echo "Cleaning Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down --remove-orphans
	@docker system prune --force
	@echo "Cleanup complete."

# Rebuild Docker containers from scratch without removing volumes
re: down build up

# Clean up all dangling Docker resources
prune:
	@echo "Pruning all unused Docker resources..."
	@docker system prune --all --volumes --force
	@echo "Prune complete."

