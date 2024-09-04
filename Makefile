VERSION = v0.0.1
COMMIT = $$(git rev-parse HEAD)
DATE = $$(date --rfc-3339=seconds)
LDFLAGS = "-w -s -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.date=${DATE}'"
EXECUTABLE = ./.bin/moco-proxy
MAINDIR= .

GOARCH = amd64
GOOS = linux

.PHONY: generate build run
generate:
	@go generate ./...

build: generate
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags ${LDFLAGS} -o ${EXECUTABLE} ${MAINDIR} 

run: build
	@${EXECUTABLE}

DOCKERFILE = ./build/dev-moco-proxy.Dockerfile
IMAGE_NAME = moco-proxy
IMAGE_TAG = ${VERSION}
CONTAINER_NAME = moco-proxy
DOCKER_COMPOSE_FILE=./deploy/docker-compose/docker-compose.yml
DOCKER_COMPOSE_ENV=./deploy/docker-compose/config.env

.PHONY: docker-image docker-container
docker-image:
	docker build --tag ${IMAGE_NAME}:${IMAGE_TAG} --target=${target} -f ${DOCKERFILE} .

docker-container:
	docker run --rm -d --network="host" --name=${CONTAINER_NAME} ${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: docker-up docker-down
docker-up:
	docker compose --env-file=${DOCKER_COMPOSE_ENV} --project-directory=. -f=${DOCKER_COMPOSE_FILE} up -d

docker-down:
	docker compose --env-file=${DOCKER_COMPOSE_ENV} --project-directory=. -f=${DOCKER_COMPOSE_FILE} down

TEST_DIRS = $$(go list ./... | grep -ve "/mock.*/")

.PHONY: lint test test-race bench coverage
lint:
	@echo "Running golangci-lint"
	@golangci-lint run --config ./.golangci.yml ./...

test:
	@echo "Running unit tests"
	@go test -count=1 ./...

bench:
	@echo "Running bench tests"
	@go test -bench=.  ${TEST_DIRS}

test-race:
	@echo "Running unit tests with race detector"
	@go test -v -count=1 ${TEST_DIRS} -race

coverage:
	@if [ -f cov.out ]; then rm cov.out; fi;
	@echo "Running unit tests with coverage profile"
	@go test ${TEST_DIRS} -coverprofile=cov.out -covermode=count 
	@go tool cover -func=cov.out

.PHONY: clean
clean:
	rm -f ./.bin/* ./cov.out
