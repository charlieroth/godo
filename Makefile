SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ===
# Define Dependencies
GOLANG   := golang:1.23
ALPINE   := alpine:3.21
POSTGRES := postgres:17.3

GODO_APP        := godo
BASE_IMAGE_NAME := localhost/charlieroth
VERSION         := 0.1.0
GODO_IMAGE      := $(BASE_IMAGE_NAME)/$(GODO_APP):$(VERSION)

# ===
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	docker pull $(POSTGRES) & \
	wait;

# ===
# Build Containers

build: godo

godo:
	docker build \
		-f zarf/docker/Dockerfile \
		-t $(GODO_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
		.

# ===
# Docker Compose

compose-up:
	cd ./zarf/compose && docker compose -f docker-compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./zarf/compose && docker compose -f docker-compose.yaml down

compose-logs:
	cd ./zarf/compose && docker compose -f docker-compose.yaml logs

# ===
# Administration

pgcli:
	pgcli postgresql://postgres:postgres@localhost

liveness:
	curl -s http://localhost:8080/liveness

readiness:
	curl -s http://localhost:8080/readiness

# ===
# Testing

test-down:
	docker stop godotest
	docker rm godotest -v

test-r:
	CGO_ENABLED=0 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

vuln-check:
	govulncheck ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -check=all ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check

# ===
# Modules Support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ===
# Local

run:
	go run cmd/server/main.go

ready:
	curl -i http://localhost:8080/readiness

live:
	curl -i http://localhost:8080/liveness

# ==============================================================================
# Help command
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  dev-gotooling           Install Go tooling"
	@echo "  dev-brew                Install brew dependencies"
	@echo "  dev-docker              Pull Docker images"
	@echo "  build                   Build all the containers"
	@echo "  godo                    Build the godo container"
	@echo "  compose-up              Start the Docker Compose cluster"
	@echo "  compose-build-up        Build and start the Docker Compose cluster"
	@echo "  compose-down            Stop the Docker Compose cluster"
	@echo "  compose-logs            Show the logs for the Docker Compose cluster"
	@echo "  pgcli                   Connect to the database"
	@echo "  liveness                Check the liveness of the server"
	@echo "  readiness               Check the readiness of the server"
	@echo "  test                    Run the tests"
	@echo "  test-race               Run the tests with race detection"
	@echo "  test-only               Run the tests without race detection"
	@echo "  test-down               Stop the test Docker container"
	@echo "  vuln-check              Check the vulnerabilities in the code"
	@echo "  tidy                    Tidy the go modules"
	@echo "  deps-reset              Reset the go modules"
	@echo "  deps-list               List the go modules"
	@echo "  deps-upgrade            Upgrade the go modules"
	@echo "  deps-cleancache         Clean the go module cache"
	@echo "  list                    List the go modules"
	@echo "  run                     Run the server"
	@echo "  ready                   Check the readiness of the server"
	@echo "  live                    Check the liveness of the server"
