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

ifndef include_skaffold_mk
_include_skaffold_mk := 1

include makefiles/shared.mk
include makefiles/kubectl.mk

SKAFFOLD := $(BIN)/skaffold
SKAFFOLD_VERSION ?= 1.39.9

$(SKAFFOLD): | $(BIN)
	$(info $(_bullet) Installing <skaffold>)
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-$(ARCH) -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

debug: debug-skaffold

deploy: deploy-skaffold

dev: dev-skaffold

run: run-skaffold

.PHONY: clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold

clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold: $(SKAFFOLD) $(KUBECTL)
clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold: export PATH := $(BIN):$(shell echo $$PATH)

clean-skaffold: ## Clean Skaffold
	$(info $(_bullet) Cleaning <skaffold>)
	! kubectl config current-context &>/dev/null || \
	$(SKAFFOLD) delete

build-skaffold: ## Build artifacts with Skaffold
	$(info $(_bullet) Building artifacts with <skaffold>)
	$(SKAFFOLD) build

deploy-skaffold: build-skaffold ## Deploy artifacts with Skaffold
	$(info $(_bullet) Deploying with <skaffold>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy --force --build-artifacts -

run-skaffold: ## Run with Skaffold
	$(info $(_bullet) Running stack with <skaffold>)
	$(SKAFFOLD) run --force

dev-skaffold: ## Run in development mode with Skaffold
	$(info $(_bullet) Running stack in development mode with <skaffold>)
	$(SKAFFOLD) dev --force --port-forward 

debug-skaffold: ## Run in debugging mode with Skaffold
	$(info $(_bullet) Running stack in debugging mode with <skaffold>)
	$(SKAFFOLD) debug --force --port-forward

endif