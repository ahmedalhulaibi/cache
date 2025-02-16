DOCKER_REST_PORT_MAPPING ?= 8080:8080
DOCKER_GRPC_PORT_MAPPING ?= 8090:8090

K3D_REGISTRIES_YAML ?= ./infra/k3d-registries.yaml
K3D_CLUSTER_NAME ?= cache-api
K3D_HOST_PORT ?= 8888

GO_SRC_DIR ?= ./go

SSH_GIT := ssh://git@github.com/ahmedalhulaibi/cache-api.git
HTTPS_GIT := https://github.com/ahmedalhulaibi/cache-api.git
PKG := github.com/ahmedalhulaibi/cache-api
GOPRIVATE := github.com/ahmedalhulaibi

K6_HTTP_BASE_URL ?= http://localhost:8888/cache-api
K6_GRPC_BASE_URL ?= localhost:8888
K6_GRPC_PROTO_DIR ?= ../..
K6_GRPC_PROTO_FILES ?= $(shell find ./proto | grep "\.proto" | tr '\n' ';')

include makefiles/shared.mk
include makefiles/git.mk
include makefiles/go.mk
include makefiles/buf.mk
include makefiles/docker.mk
include makefiles/kubectl.mk
include makefiles/k3d.mk
include makefiles/skaffold.mk
include makefiles/k6.mk
include makefiles/linkerd.mk

build: build-api 

bootstrap: bootstrap-deployment
bootstrap: .ssh/id_rsa
run-skaffold: .ssh/id_rsa

vendor: vendor-go
	cd $(GO_SRC_DIR) && GOBIN=$(BIN) $(GO) install  \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	istio.io/tools/cmd/protoc-gen-docs

.ssh/id_rsa:
	@mkdir -p .ssh
	@ln -sf $(HOME)/.ssh/id_rsa .ssh/id_rsa

.PHONY: build-api build-docker pr-ready run-local

build-api:
	$(info $(_bullet) Building <api>)
	cd $(GO_SRC_DIR) && $(GO) build -o bin/api ./cmd/api && cd -

run-local: build-api
	$(info $(_bullet) Running <api>)
	$(GO_SRC_DIR)/bin/api

build-docker: # .ssh/id_rsa
	$(info $(_bullet) Building docker <api>)
	docker build --no-cache -f ./go/Dockerfile --build-arg BUILDPKG=cmd/api --build-arg PKG=$(PKG) --build-arg GOPRIVATE=$(GOPRIVATE) $(GO_SRC_DIR) -t ghcr.io/ahmedalhulaibi/cache-api:latest

run-docker: build-docker
	$(info $(_bullet) Running docker <api>)
	docker run -p $(DOCKER_REST_PORT_MAPPING) -p $(DOCKER_GRPC_PORT_MAPPING) ghcr.io/ahmedalhulaibi/cache-api:latest

bootstrap-deployment: $(KUBECTL) ## Bootstrap deployment
	$(info $(_bullet) Bootstraping <deployment>)
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -k infra/k8s/bootstrap/overlays/local

pr-ready: format vendor generate lint integration-test git-dirty
