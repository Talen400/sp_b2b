BIN_DIR := bin
BIN_NAME := api
DB_PATH := data/split.db

.PHONY: build run test migrate seed clean re sandbox help

build:
	go build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/api

run: build
	./$(BIN_DIR)/$(BIN_NAME)

test:
	go test ./... -v

migrate:
	go run ./cmd/api -migrate

seed: migrate
	go run ./cmd/api -seed

clean:
	rm -rf $(BIN_DIR) $(DB_PATH)

re: clean build

sandbox:
	@test -f cgibs/openapi-v0_0_10.json || { echo "Erro: cgibs/openapi-v0_0_10.json não encontrado. Coloque o arquivo OpenAPI oficial em cgibs/."; exit 1; }
	@echo "Subindo mock Prism da Plataforma Pública em http://localhost:4010"
	@echo "Pare com Ctrl+C quando quiser encerrar."
	npx @stoplight/prism-cli mock cgibs/openapi-v0_0_10.json -p 4010

help:
	@echo "Targets disponíveis:"
	@echo "  build     Compila o binário em $(BIN_DIR)/$(BIN_NAME)"
	@echo "  run       Sobe servidor (porta 8080)"
	@echo "  test      Roda todos os testes"
	@echo "  migrate   Aplica migrations no banco"
	@echo "  seed      Popula banco com cenário de demonstração"
	@echo "  clean     Remove bin/ e banco"
	@echo "  re        Clean + build"
	@echo "  sandbox   Sobe mock Prism da PP em http://localhost:4010"
