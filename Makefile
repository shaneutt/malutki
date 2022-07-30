# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p build/
	go build -o build/malutki main.go

.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

.PHONY: clean
clean:
	rm -f build/*

# ------------------------------------------------------------------------------
# Images
# ------------------------------------------------------------------------------

IMAGE_NAME ?= ghcr.io/shaneutt/malutki
IMAGE_TAG ?= latest

.PHONY: docker.build
docker.build:
	docker buildx build \
		-f Dockerfile \
		--target distroless \
		-t $(IMAGE_NAME):$(IMAGE_TAG) .

# ------------------------------------------------------------------------------
# Test
# ------------------------------------------------------------------------------

.PHONY: test.e2e
test.e2e: docker.build
	MALUTKI_TEST_IMAGE="$(IMAGE_NAME):$(IMAGE_TAG)" GOFLAGS="-tags=e2e_tests" go test -v -race ./test/e2e/...
