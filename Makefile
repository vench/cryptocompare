V := @

# Build
OUT_DIR = ./bin
DOCKER_IMAGE_NAME = github.com/vench/cryptocompare
DOCKER_FILE = deployments/docker/Dockerfile
ACTION ?= build

.PHONY: vendor
vendor:
	$(V)go mod tidy -compat=1.17
	$(V)go mod vendor

.PHONY: build
build:
	$(V)CGO_ENABLED=1 go build -o ${OUT_DIR}/cryptocompare ./cmd/cryptocompare/main.go

.PHONY: test
test: GO_TEST_FLAGS += -race
test:
	$(V)go test -mod=vendor $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...

.PHONY: docker-build-local
docker-build-local:
	$(V)docker build --network host -t ${DOCKER_IMAGE_NAME}:local -f ${DOCKER_FILE} --build-arg ACTION=${ACTION} .

.PHONY: init
init:
	$(V)make docker-build-local
	$(V)make up
	$(V)echo done...

.PHONY: down
down:
	$(V)docker-compose  -f deployments/docker-compose.yaml down

.PHONY: up
up:
	$(V)docker-compose  -f deployments/docker-compose.yaml up -d




