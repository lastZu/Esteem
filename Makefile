.PHONY: build
build:
	go build -v ./cmd/esteem

.DEFAULT_GOAL := build