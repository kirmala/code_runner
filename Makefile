.PHONY: launch_services launch_with_tests stop_services build_services

launch_services: build_services
	docker compose up

launch_with_tests: build_services
	docker compose --profile test up --abort-on-container-exit --exit-code-from app_test

stop_services:
	docker compose down

build_services:
	docker compose build
	docker build -t code_processor consumer/usecases/services/docker_code_processor
	
