COLOR_RED = \033[0;31m
COLOR_GREEN = \033[0;32m
COLOR_YELLOW = \033[0;33m
COLOR_BLUE = \033[0;34m
COLOR_RESET = \033[0m

dev:
	docker compose -f deployment/docker-compose.dev.yml up -d

migrate:
	migrate -source file://schema -database postgres://devuser:devpassword@localhost:8008/devdb?sslmode=disable up

swag:
	swag init -o docs/ -g main.go -d ./cmd/api,./internal/api/handlers --pd --pdl 1 -p pascalcase --st -q
	@echo -e "$(COLOR_GREEN)Swagger docs generated successfully!$(COLOR_RESET)"

run:
	go run cmd/api/main.go

generate-api:
	npx swagger-typescript-api generate -p ./docs/swagger.json -o ./generator/src -n api.ts --axios