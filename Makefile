# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_FILE = ./srcs/docker-compose.yaml

# Phony targets
.PHONY: all build up down restart clean fclean re prune

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

# Stop and remove Docker containers
down:
	@echo "Stopping and removing Docker containers..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "Containers stopped and removed."

# Restart Docker containers
restart: down up

# Clean Docker resources including volumes
clean:
	@echo "Cleaning Docker containers and volumes..."
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
	@docker system prune --volumes --force
	@echo "Cleanup complete."

# Full clean including data directory
fclean: clean
	@echo "Removing data directory..."
	@rm -rf $(DATA_DIR)
	@echo "Data directory removed."

# Rebuild Docker containers from scratch without removing data
re: clean all

# Clean up all dangling Docker resources
prune:
	@echo "Pruning all unused Docker resources..."
	@docker system prune --all --volumes --force
	@echo "Prune complete."