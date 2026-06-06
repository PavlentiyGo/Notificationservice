include .env
export

export PROJECT_ROOT = ${shell pwd}


sub-up:
	go run services/subscription/cmd/main.go

env-up:
	docker compose up -d --build

migrate-up:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make migrate-up SERVICE=subscription"; \
		echo "Available services: subscription"; \
		exit 1; \
	fi
	@migrate \
	-path ${PROJECT_ROOT}/services/$(SERVICE)/migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/$(SERVICE)?sslmode=disable \
	up

migrate-down:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make migrate-up SERVICE=subscription"; \
		echo "Available services: subscription"; \
		exit 1; \
	fi
	@migrate \
	-path ${PROJECT_ROOT}/services/$(SERVICE)/migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/$(SERVICE)?sslmode=disable \
	down