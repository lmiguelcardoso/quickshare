.PHONY: help up down restart logs test fmt build clean

.DEFAULT_GOAL := help

help:
	@echo "Comandos disponíveis:"
	@echo "  make up          - Sobe os containers"
	@echo "  make down        - Para os containers"
	@echo "  make restart     - Reinicia os containers"
	@echo "  make logs        - Mostra logs"
	@echo "  make test        - Roda os testes"
	@echo "  make test-cov    - Roda testes com coverage"
	@echo "  make fmt         - Formata o código"
	@echo "  make build       - Rebuilda as imagens"
	@echo "  make clean       - Remove tudo"

# Docker
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

logs-api:
	docker-compose logs -f api

logs-db:
	docker-compose logs -f postgres

build:
	docker-compose build

clean:
	docker-compose down -v --rmi local

# Tests
test:
	go test ./...

test-v:
	go test -v ./...

test-cov:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Dev
run:
	go run cmd/main.go

fmt:
	go fmt ./...

tidy:
	go mod tidy

# Database
db-shell:
	docker-compose exec postgres psql -U postgres -d quickshare

