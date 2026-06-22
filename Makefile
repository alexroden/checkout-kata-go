.PHONY: setup start cleanup

setup:
	@docker compose -p kata up -d dynamodb
	@sleep 5
	@docker compose -p kata up -d dynamo-init

start:
	DOCKER_BUILDKIT=1 docker compose -p kata up -d api

cleanup:
	@docker compose -p kata down -v --rmi all
