ifndef _include_buf_mk
_include_buf_mk = 1

include makefiles/shared.mk
include makefiles/go.mk

BUF := $(BIN)/buf
# This controls the version of buf to install and use.
BUF_VERSION := v1.50.0
# If true, Buf is installed from source instead of from releases
BUF_INSTALL_FROM_SOURCE := false

$(BUF): | $(BIN)
	$(info $(_bullet) Installing <buf>)
	GOBIN=$(BIN) $(GO) install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

generate: generate-buf

lint: buf-lint

.PHONY: generate-buf buf-lint buf-deps buf-local buf-https buf-ssh

generate-buf: $(BUF) buf-deps ## Generate stubs for protoc plugins
	$(info $(_bullet) Generating <buf>)
	
	PATH="$(BIN):$(PATH)" $(BUF) generate

buf-lint: $(BUF)
	$(info $(_bullet) Linting <buf>)
	$(BUF) lint

buf-deps: $(BUF)
	$(info $(_bullet) Installing buf dependencies)
	$(BUF) dep update

# local is what we run when testing locally.
# This does breaking change detection against our local git repository.

buf-local: buf-lint
	$(info $(_bullet) Checking for breaking changes <buf>)
	$(BUF) breaking --against '.git#branch=main'

# https is what we run when testing in most CI providers.
# This does breaking change detection against our remote HTTPS git repository.

buf-https: buf-lint
	$(info $(_bullet) Checking for breaking changes <buf>)
	$(BUF) breaking --against "$(HTTPS_GIT)#branch=main"

# ssh is what we run when testing in CI providers that provide ssh public key authentication.
# This does breaking change detection against our remote HTTPS ssh repository.
# This is especially useful for private repositories.

buf-ssh: buf-lint
	$(info $(_bullet) Checking for breaking changes <buf>)
	$(BUF) breaking --against "$(SSH_GIT)#branch=main"

endif