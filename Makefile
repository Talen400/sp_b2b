BIN_DIR := bin
BIN_NAME := api
DB_PATH := data/split.db
MIGRATIONS_DIR := migrations

.PHONY: build run test migrate seed clean re

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
