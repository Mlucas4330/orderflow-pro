.PHONY: help up down logs build db-migrate-dev db-migrate-test lint test build-images push-images

# --- Variáveis de Configuração ---
DOCKER_USER := mlucas4330
VERSION     := v1.0.0
SERVICES    := order-service inventory-service notification-service product-service

COMPOSE_BASE := docker-compose.yml
COMPOSE_DEV  := docker-compose.dev.yml
COMPOSE_TEST := docker-compose.test.yml

# ==============================================================================
# Comandos do Projeto OrderFlow Pro
# ==============================================================================

## help: Mostra esta ajuda.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## up: Sobe o ambiente de desenvolvimento em background (base + dev).
up:
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up -d --build

## down: Para e remove todos os contêineres e volumes (reset completo).
down:
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) down -v 2>/dev/null || true
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) down -v 2>/dev/null || true

## logs: Mostra os logs do ambiente de desenvolvimento em tempo real.
logs:
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) logs -f

## build: Força a reconstrução das imagens dos serviços de desenvolvimento.
build:
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) build

## db-migrate-dev: Aplica as migrações na base de dados de DESENVOLVIMENTO.
db-migrate-dev:
	@goose -dir "db/migrations" postgres "postgres://gopher:mysecretpassword@localhost:5432/orderflow_dev_db?sslmode=disable" up

## db-migrate-test: Aplica as migrações na base de dados de TESTE.
db-migrate-test:
	@goose -dir "db/migrations" postgres "postgres://gopher:mysecretpassword@localhost:5432/orderflow_test_db?sslmode=disable" up

## lint: Roda o linter para verificar a qualidade do código.
lint:
	@golangci-lint run

## test: Roda a suíte de testes completa num ambiente Docker isolado (base + teste).
test:
	@docker-compose -f $(COMPOSE_BASE) -f $(COMPOSE_TEST) up --build --abort-on-container-exit

## build-images: Constrói as imagens Docker para todos os microsserviços.
build-images:
	@$(foreach service,$(SERVICES), \
		docker build -t $(DOCKER_USER)/orderflow-pro-$(service):$(VERSION) -f cmd/$(service)/Dockerfile . ;\
	)

## push-images: Envia todas as imagens construídas para o Docker Hub.
push-images:
	@docker login
	@$(foreach service,$(SERVICES), \
		docker push $(DOCKER_USER)/orderflow-pro-$(service):$(VERSION) ;\
	)