ifndef _include_k3d_mk
_include_k3d_mk :=1

include makefiles/shared.mk

K3D := $(BIN)/k3d
K3D_VERSION ?= v4.4.8
K3D_CLUSTER_NAME ?= local
K3D_HOST_PORT ?= 80
K3D_REGISTRIES_YAML ?= ./k3d-registries.yaml

BOOTSTRAP_CONTEXT := k3d-$(K3D_CLUSTER_NAME)

$(K3D): | $(BIN)
	$(info $(_bullet) Installing <k3d>)
	curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | PATH="$(BIN):$(PATH)" K3D_INSTALL_DIR=$(BIN) USE_SUDO="false" TAG=$(K3D_VERSION) bash
	chmod u+x $(K3D)

bootstrap: bootstrap-k3d

.PHONY: bootstrap-k3d clean-k3d

bootstrap-k3d: .bootstrap-registry .bootstrap-cluster

clean-k3d: $(K3D)
	$(K3D) registry delete $(K3D_CLUSTER_NAME).registry.localhost
	$(K3D) cluster delete $(K3D_CLUSTER_NAME)
	rm .bootstrap-registry .bootstrap-cluster


.bootstrap-registry: $(K3D)
	$(K3D) registry create $(K3D_CLUSTER_NAME).registry.localhost --port 12345 && touch .bootstrap-registry

.bootstrap-cluster: $(K3D)
	$(K3D) cluster create $(K3D_CLUSTER_NAME) \
	--api-port 6550 \
	-p "$(K3D_HOST_PORT):80@loadbalancer" \
	--registry-use $(K3D_CLUSTER_NAME).registry.localhost:12345 \
	--registry-config $(K3D_REGISTRIES_YAML) \
	&& touch .bootstrap-cluster

endif