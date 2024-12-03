# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_FILE = ./srcs/docker-compose.yaml

# Phony targets
.PHONY: all build up down restart clean re prune prod

# Default: development environment (everything except frontend_prod)
all: build up

# Development build (excludes frontend_prod)
build:
	@echo "Building Docker containers for development..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build frontend backend db adminer websocket alertmanager prometheus node-exporter postgres-exporter grafana

# Start containers for development
up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d frontend backend db adminer websocket alertmanager prometheus node-exporter postgres-exporter grafana
	@echo "Development environment is up on port 3000"

# Production environment (excludes frontend)
prod:
	@echo "Building and starting production environment..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d --build frontend_prod backend db adminer websocket alertmanager prometheus node-exporter postgres-exporter grafana nginx
	@echo "Production environment is up on port 8000"

# Stop and remove containers
down:
	@echo "Stopping and removing Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "Containers stopped and removed."

# Clean Docker resources
clean:
	@echo "Cleaning Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down --remove-orphans
	@docker volume rm -f srcs_frontend_build
	@docker volume rm -f srcs_grafana-data
	@docker volume rm -f srcs_prometheus_data
	@docker system prune --force
	@echo "Cleanup complete."

# Deep clean
prune:
	@echo "Pruning all unused Docker resources..."
	@docker system prune --all --volumes --force
	@echo "Prune complete."

# Rebuild containers
re: down build up

# Restart containers
restart: down up
