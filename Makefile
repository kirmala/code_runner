COMPOSE_BASE := compose.yml
COMPOSE_DEV := compose.dev.yml
COMPOSE_TEST := compose.test.yml

PHONY: dev-up dev-logs dev-down test-up test-logs test-down

dev-build:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) build

dev-up:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up -d

dev-logs:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) logs -f

dev-down:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) down

dev-restart-server:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) restart http_server

dev-restart-consumer:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) restart consumer


test-up:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) up --abort-on-container-exit --exit-code-from app_test --build

test-logs:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) logs -f

test-down:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) down -v

all-down: dev-down test-down

base-build:
	docker compose -f $(COMPOSE_BASE) build

build-runner:
	docker build -t runner ./consumer/service/docker/runner