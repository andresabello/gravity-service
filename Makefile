# Define variables
COMPOSE_FILE = ./docker-compose.yml
PROJECT_NAME = backendify
PACKAGE = /app/internal/companies
TEST_FLAGS = 

# Define targets
dev:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up

build:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) build --no-cache

stop:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) stop

clean:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down -v

test:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) exec app go test -cover -v $(PACKAGE) $(TEST_FLAGS)

lint:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) exec app sh -c "cd /app && golangci-lint run -c .golangci.yml"