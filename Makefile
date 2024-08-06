# Directories
API_DIR=./api
CLIENT_DIR=./client

# Docker settings
DOCKER_IMAGE=mbrav/slack-bots:latest
# DOCKER_IMAGE_API=myregistry/myapi:latest
# DOCKER_IMAGE_CLIENT=myregistry/myclient:latest

# Go settings
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_RUN=$(GO_CMD) run

# Docker commands
# DOCKER_BUILD=docker build
DOCKER_BUILD=docker buildx build --platform linux/amd64
DOCKER_PUSH=docker push

# Targets
.PHONY: all test build run docker-build docker-push

all: test mkdir build docker-build docker-push

# help:
# 	@echo "Makefile commands:"
# 	@echo "run          - Run Go package"
# 	@echo "build        - Package for all code and deps into a zip archive"
# 	@echo ""
# 	@echo "Current Version: ${TAG}"
#
## Run tests for both api and client
test:
	@echo "Running tests..."
	$(GO_TEST) ./...

mkdir:
	mkdir -p ./bin

## Build api and client
build: build-api build-client

build-api:
	@echo "Building API..."
	cd $(API_DIR) && $(GO_BUILD) -o ../bin/api
	@command -v upx >/dev/null 2>&1 && upx ./bin/api || echo "Skipping compression"

build-client:
	@echo "Building Client..."
	cd $(CLIENT_DIR) && $(GO_BUILD) -o ../bin/client
	@command -v upx >/dev/null 2>&1 && upx ./bin/client || echo "Skipping compression"

run-api:
	@echo "Running API..."
	cd $(API_DIR) && $(GO_RUN) .

run-client:
	@echo "Running Client..."
	cd $(CLIENT_DIR) && $(GO_RUN) .

## Build Docker images for api and client
# docker-build: docker-build-api docker-build-client
docker-build:
	@echo "Building Docker image ${DOCKER_IMAGE} ..."
	$(DOCKER_BUILD) -t ${DOCKER_IMAGE} .
	
# docker-build-api:
# 	@echo "Building Docker image for API..."
# 	cd $(API_DIR) && $(DOCKER_BUILD) -t $(DOCKER_IMAGE_API) .
#
# docker-build-client:
# 	@echo "Building Docker image for Client..."
# 	cd $(CLIENT_DIR) && $(DOCKER_BUILD) -t $(DOCKER_IMAGE_CLIENT) .
#
## Push Docker images to the registry
# docker-push: docker-push-api docker-push-client
docker-push: docker-build
	@echo "Pushing Docker image ${DOCKER_IMAGE} ..."
	$(DOCKER_PUSH) ${DOCKER_IMAGE}

# docker-push-api:
# 	@echo "Pushing Docker image for API..."
# 	$(DOCKER_PUSH) $(DOCKER_IMAGE_API)
#
# docker-push-client:
# 	@echo "Pushing Docker image for Client..."
# 	$(DOCKER_PUSH) $(DOCKER_IMAGE_CLIENT)

deploy:
	@echo "Deploying to k8s with Kustomize ..."
	kubectl apply -k ./kustomize

deploy-down:
	@echo "Removing deploy from  k8s with Kustomize ..."
	kubectl delete -k ./kustomize

