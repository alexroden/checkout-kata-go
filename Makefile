.PHONY: setup start down cleanup logs

setup:
	@docker compose -p kata up -d dynamodb
	@sleep 5
	@docker compose -p kata up -d dynamo-init

start:
	DOCKER_BUILDKIT=1 docker compose -p kata up -d api
	@sleep 20
	@make logs

down:
	@docker compose -p kata down

cleanup:
	@docker compose -p kata down -v --rmi all

logs:
	@docker logs -f kata