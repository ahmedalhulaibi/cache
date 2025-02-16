ifndef _include_linkerd_mk
_include_linkerd_mk := 1

include makefiles/kubectl.mk
include makefiles/shared.mk

LINKERD := $(BIN)/linkerd
LINKERD_VERSION ?= stable-2.11.1

$(LINKERD): | $(BIN)
	$(info $(_bullet) Installing <linkerd>)
	curl -sSfL https://github.com/linkerd/linkerd2/releases/download/$(LINKERD_VERSION)/linkerd2-cli-$(LINKERD_VERSION)-$(OS)-$(ARCH) -o $(LINKERD)
	chmod u+x $(LINKERD)

bootstrap-linkerd: $(KUBECTL) $(LINKERD)
	$(info $(_bullet) Bootstrapping <linkerd>)
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -k infra/k8s/linkerd/base
	$(LINKERD) check
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -f infra/k8s/linkerd/base/linkerd2-viz.yaml
	$(LINKERD) check
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -f infra/k8s/linkerd/base/linkerd2-viz-ingress.yaml
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -f infra/k8s/linkerd/base/linkerd2-jaeger.yaml
	$(LINKERD) check

endif