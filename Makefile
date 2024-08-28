ARCH = amd64
OS = linux

VERSION = latest
COMMIT =
DATE = 
LDFLAGS = "-w -s -X main.version='${VERSION}' -X main.commit='${COMMIT}' -X main.date='${DATE}'"
EXECUTABLE = ./.bin/go-moco

DOCKERFILE = ./deploy/docker/go-moco-proxy.Dockerfile
CONTAINER_NAME = go-moco-proxy
IMAGE_NAME = go-moco-proxy
IMAGE_TAG = ${VERSION}
DOCKER_COMPOSE_CONFIG=./deploy/docker/config.env

.PHONY: run build
build:
	go generate ./...
	GOOS=${OS} GOARCH=${GOARCH} go build -ldflags ${LDFLAGS} -o ${EXECUTABLE} ./cmd/moco-proxy/

run: build
	${EXECUTABLE}

.PHONY: docker-image docker-container
docker-image:
	docker build --tag ${IMAGE_NAME}:${IMAGE_TAG} --target=${target} -f ${DOCKERFILE} .

docker-container:
	docker run --rm -d --network="host" --name=${CONTAINER_NAME} ${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: docker-up docker-down
docker-up:
	docker compose --env-file=${DOCKER_COMPOSE_CONFIG} up -d 

docker-down:
	docker compose --env-file=${DOCKER_COMPOSE_CONFIG} down 

.PHONY: lint test test-race bench coverage
lint:
	@echo "Running golangci-lint"
	@golangci-lint run --config ./.golangci.yml ./...

test:
	@echo "Running unit tests"
	@go test -count=1 ./...

bench:
	@echo "Running bench tests"
	@go test -bench=.  ./...

test-race:
	@echo "Running unit tests with race detector"
	@go test -v -count=1 ./... -race

coverage:
	@if [ -f coverage.out ]; then rm coverage.out; fi;
	@echo "Running unit tests with coverage profile"
	@go test ./... -coverprofile=coverage.out -covermode=count
	@go tool cover -func=cov.out

.PHONY: clean
clean:
	rm -f ./.bin/* ./cov.out
