# MIT License
#
# Copyright (c) 2020 Zbigniew Mandziejewicz
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

ifndef _include_go_mk
_include_go_mk = 1

include makefiles/shared.mk

GO ?= go
FORMAT_FILES ?= .
GO_SRC_DIR ?= .

GOFUMPT := $(BIN)/gofumpt
GOFUMPT_VERSION ?= v0.7.0

GOLANGCILINT := $(BIN)/golangci-lint
GOLANGCILINT_VERSION ?= v1.64.5
GOLANGCILINT_CONCURRENCY ?= 16

GOIMPORTS := $(BIN)/goimports
GOIMPORTS_VERSION ?= v0.30.0

$(GOFUMPT): | $(BIN)
	$(info $(_bullet) Installing <gofumpt>)
	GOBIN=$(BIN) $(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(GOLANGCILINT): | $(BIN)
	$(info $(_bullet) Installing <golangci-lint>)
	GOBIN=$(BIN) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

$(GOIMPORTS): | $(BIN)
	$(info $(_bullet) Installing <goimports>)
	GOBIN=$(BIN) $(GO) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

clean: clean-go

deps: deps-go

vendor: vendor-go

format: format-go

lint: lint-go

test: test-go

test-coverage: test-coverage-go

integration-test: integration-test-go

.PHONY: deps-go format-go lint-go test-go test-coverage-go integration-test-go

clean-go: ## Clean Go
	$(info $(_bullet) Cleaning <go>)
	cd $(GO_SRC_DIR) && rm -rf vendor/

deps-go: ## Download Go dependencies
	$(info $(_bullet) Downloading dependencies <go>)
	cd $(GO_SRC_DIR) && $(GO) mod download

vendor-go: ## Vendor Go dependencies
	$(info $(_bullet) Vendoring dependencies <go>)
	cd $(GO_SRC_DIR) && $(GO) mod vendor

format-go: $(GOFUMPT) $(GOIMPORTS) ## Format Go code
	$(info $(_bullet) Formatting code)
	cd $(GO_SRC_DIR) && $(GOIMPORTS) -w $(FORMAT_FILES)
	cd $(GO_SRC_DIR) && $(GOFUMPT) -w $(FORMAT_FILES)

lint-go: $(GOLANGCILINT)
	$(info $(_bullet) Linting <go>) 
	cd $(GO_SRC_DIR) && $(GOLANGCILINT) run --concurrency $(GOLANGCILINT_CONCURRENCY) ./...

test-go: ## Run Go tests
	$(info $(_bullet) Running tests <go>)
	cd $(GO_SRC_DIR) && $(GO) test ./...
	
test-coverage-go: ## Run Go tests with coverage
	$(info $(_bullet) Running tests with coverage <go>) 
	cd $(GO_SRC_DIR) && $(GO) test -cover ./...

integration-test-go: ## Run Go integration tests
	$(info $(_bullet) Running integration tests <go>) 
	cd $(GO_SRC_DIR) && $(GO) test -tags integration -count 1 ./...

endif