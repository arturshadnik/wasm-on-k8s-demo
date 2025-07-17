REGISTRY ?= quay.io/arturshadnik
IMAGE ?= $(REGISTRY)/counter
VERSION ?= 0.1.0
PLATFORM_NATIVE ?= linux/amd64,linux/arm64
PLATFORM_WASM ?= wasi/wasm32
DESTINATION ?= load # load, push

# Build
.PHONY: build
build:
	docker buildx build -t $(IMAGE):$(VERSION)-native -f build/Dockerfile --platform $(PLATFORM_NATIVE) --$(DESTINATION) .
	docker buildx build -t $(IMAGE):$(VERSION)-wasi -f build/Dockerfile.wasi --platform $(PLATFORM_WASM) --$(DESTINATION) .

build-bin:
	go build -o bin/counter cmd/main.go

build-bin-wasi:
	GOOS=wasip1 GOARCH=wasm go build -ldflags="-s -w" -o bin/counter.wasm cmd/main.go

run:
	INPUT_PATH=hack/test.csv OUTPUT_PATH=output/results.csv go run cmd/main.go

run-wasi: build-bin-wasi
	# map current directory to root inside the wasi module
	wasmedge --dir /:. --env INPUT_PATH=hack/test.csv --env OUTPUT_PATH=output/results.csv bin/counter.wasm

env-up:
	kind create cluster --name dev --config hack/kind-config.yaml --kubeconfig ~/Downloads/dev.kubeconfig

env-down:
	rm -rf ~/Downloads/dev.kubeconfig
	kind delete cluster --name dev
