COMPOSE_BASE := compose.yml
COMPOSE_DEV := compose.dev.yml
COMPOSE_TEST := compose.test.yml
COMPOSE_PROD := compose.prod.yml
CLUSTER := ./redis_cluster/compose.yml

PHONY: network-create cluster-build cluster-up cluster-down dev-build dev-up dev-logs dev-down dev-restart-server dev-restart-consumer runner-build test-up test-logs test-down base-build prod-up prod-down all-down

network-create:
	docker network inspect code_runner_network >/dev/null 2>&1 || \
	docker network create code_runner_network

cluster-build:
	docker compose -f $(CLUSTER) build

cluster-up: network-create
	docker compose -f $(CLUSTER) up -d

cluster-down:
	docker compose -f $(CLUSTER) down


dev-build:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) build

dev-up: network-create cluster-up
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up -d

dev-logs:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) logs -f

dev-down:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) down

dev-restart-server:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) restart http_server

dev-restart-consumer:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) restart consumer

runner-build:
	docker build -t runner ./consumer/service/docker

test-up: runner-build network-create cluster-up
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) up --abort-on-container-exit --exit-code-from app_test --build

test-logs:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) logs -f

test-down:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) down -v


base-build:
	docker compose -f $(COMPOSE_BASE) build

prod-up: network-create
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_PROD) up --build
prod-down:
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_PROD) down


all-down: dev-down test-down prod-down cluster-down