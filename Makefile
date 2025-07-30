# Makefile
GO_FILE := main.go
NAME := $(shell grep "appName" $(GO_FILE) | head -1 | awk -F'"' '{print $$2}')
VERSION := "v$(shell grep "appVersion" $(GO_FILE) | head -1 | awk -F'"' '{print $$2}')"

build: clean
	@echo "Building $(NAME) version: $(VERSION)"
	go build -ldflags "-X main.AppVersion=$(VERSION)" -o build/$(NAME)

clean:
	rm -rf build

run: build
	./build/$(NAME)

release:
	goreleaser release --clean