ifndef _include_k6_mk
_include_k6_mk :=1

include makefiles/shared.mk
include makefiles/git.mk

K6 := $(BIN)/k6
K6_VERSION ?= v0.34.0
K6_TESTS_PATH ?= test/k6
K6_TEST ?= $(K6_TESTS_PATH)/load.js
K6_HTTP_BASE_URL ?= http://localhost:8080
K6_GRPC_BASE_URL ?= localhost:8080
K6_GRPC_PROTO_DIR ?= ./proto
K6OS ?= linux
K6_ARCHIVE_TYPE = tar.gz
K6_UNARCHIVE = @tar -xf k6-$(K6_VERSION)-$(K6OS)-$(ARCH).$(K6_ARCHIVE_TYPE)

ifeq ($(OS),darwin)
	K6OS = macos
	K6_ARCHIVE_TYPE = zip
	K6_UNARCHIVE = @unzip k6-$(K6_VERSION)-$(K6OS)-$(ARCH).$(K6_ARCHIVE_TYPE)
endif

$(K6): | $(BIN)
	$(info $(_bullet) Installing <k6>)
	curl -sSfL https://github.com/grafana/k6/releases/download/$(K6_VERSION)/k6-$(K6_VERSION)-$(K6OS)-$(ARCH).$(K6_ARCHIVE_TYPE) -o k6-$(K6_VERSION)-$(K6OS)-$(ARCH).$(K6_ARCHIVE_TYPE)
	@$(K6_UNARCHIVE)
	@mv k6-$(K6_VERSION)-$(K6OS)-$(ARCH)/k6 $(K6)
	@rm -rf k6-$(K6_VERSION)-$(K6OS)-$(ARCH).$(K6_ARCHIVE_TYPE)
	chmod u+x $(K6)

smoke-test: smoke-test-k6

load-test: load-test-k6

stress-test: stress-test-k6

soak-test: soak-test-k6

.PHONY: smoke-test-k6 load-test-k6 stress-test-k6 soak-test-k6

smoke-test-k6: K6_TEST_PROFILE := smoke
smoke-test-k6: $(K6) git-submodules
	$(info $(_bullet) Running $(K6_TEST_PROFILE) tests <k6>)
	K6_TEST_PROFILE=$(K6_TEST_PROFILE) \
	K6_HTTP_BASE_URL="$(K6_HTTP_BASE_URL)" \
	K6_GRPC_BASE_URL="$(K6_GRPC_BASE_URL)" \
	K6_GRPC_PROTO_DIR="$(K6_GRPC_PROTO_DIR)" \
	K6_GRPC_PROTO_FILES="$(K6_GRPC_PROTO_FILES)" \
	$(K6) run $(K6_TEST)

load-test-k6: K6_TEST_PROFILE := load
load-test-k6: $(K6) git-submodules
	$(info $(_bullet) Running $(K6_TEST_PROFILE) tests <k6>)
	K6_TEST_PROFILE=$(K6_TEST_PROFILE) \
	K6_HTTP_BASE_URL="$(K6_HTTP_BASE_URL)" \
	K6_GRPC_BASE_URL="$(K6_GRPC_BASE_URL)" \
	K6_GRPC_PROTO_DIR="$(K6_GRPC_PROTO_DIR)" \
	K6_GRPC_PROTO_FILES="$(K6_GRPC_PROTO_FILES)" \
	$(K6) run $(K6_TEST)

stress-test-k6: K6_TEST_PROFILE := stress
stress-test-k6: $(K6) git-submodules
	$(info $(_bullet) Running $(K6_TEST_PROFILE) tests <k6>)
	K6_TEST_PROFILE=$(K6_TEST_PROFILE) \
	K6_HTTP_BASE_URL="$(K6_HTTP_BASE_URL)" \
	K6_GRPC_BASE_URL="$(K6_GRPC_BASE_URL)" \
	K6_GRPC_PROTO_DIR="$(K6_GRPC_PROTO_DIR)" \
	K6_GRPC_PROTO_FILES="$(K6_GRPC_PROTO_FILES)" \
	$(K6) run $(K6_TEST)

soak-test-k6: K6_TEST_PROFILE := soak
soak-test-k6: $(K6) git-submodules
	$(info $(_bullet) Running $(K6_TEST_PROFILE) tests <k6>)
	K6_TEST_PROFILE=$(K6_TEST_PROFILE) \
	K6_HTTP_BASE_URL="$(K6_HTTP_BASE_URL)" \
	K6_GRPC_BASE_URL="$(K6_GRPC_BASE_URL)" \
	K6_GRPC_PROTO_DIR="$(K6_GRPC_PROTO_DIR)" \
	K6_GRPC_PROTO_FILES="$(K6_GRPC_PROTO_FILES)" \
	$(K6) run $(K6_TEST)

endif