ARCH = amd64
OS = linux

VERSION = latest
COMMIT =
DATE = 
LDFLAGS = "-w -s -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.date=${DATE}'"
EXECUTABLE = ./.bin/go-moco

DOCKERFILE = ./build/go-moco-proxy.Dockerfile
CONTAINER_NAME = go-moco-proxy
IMAGE_NAME = go-moco-proxy
IMAGE_TAG = ${VERSION}
DOCKER_COMPOSE_FILE=./deploy/docker-compose/docker-compose.yml
DOCKER_COMPOSE_ENV=./deploy/docker-compose/config.env

.PHONY: run build
build:
	go generate ./...
	GOOS=${OS} GOARCH=${GOARCH} go build -ldflags ${LDFLAGS} -o ${EXECUTABLE} . 

run: build
	${EXECUTABLE}

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
	@if [ -f cov.out ]; then rm cov.out; fi;
	@echo "Running unit tests with coverage profile"
	@go test $$(go list ./... | grep -ve "/mock.*/") -coverprofile=cov.out -covermode=count 
	@go tool cover -func=cov.out

.PHONY: clean
clean:
	rm -f ./.bin/* ./cov.out
