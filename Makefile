IMAGE_REPO ?= $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/baseimage-ci
IMAGE_TAG ?= ds-resource-injector

ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
    GOARCH ?= amd64
else ifeq ($(ARCH),aarch64)
    GOARCH ?= arm64
else ifeq ($(ARCH),arm64)
    GOARCH ?= arm64
else
    $(error "This system's arch $(ARCH) isn't recognized/supported")
endif

OS := $(shell uname -s)
ifeq ($(OS),Darwin)
  GOOS ?= darwin
else ifeq ($(OS),Linux)
  GOOS ?= linux
else
    $(error "This system's os $(OS) isn't recognized/supported")
endif

.PHONY: build
build:
	@echo "\n🔧  Building $(IMAGE_TAG) Go binaries..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/admission-webhook .

//docker buildx build --push --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --platform linux/arm64,linux/amd64 -t $(IMAGE_REPO):$(IMAGE_TAG) .
.PHONY: docker-build
docker-build:
	@echo "\n📦 Building $(IMAGE_TAG) Docker image..."
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(IMAGE_REPO)
	docker buildx use craftbuilder
	docker buildx build --push --cache-to type=inline --build-arg buildkit.experimental.cache-from=golang --platform linux/arm64 -t $(IMAGE_REPO):$(IMAGE_TAG) .
