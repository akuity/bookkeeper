SHELL ?= /bin/bash

################################################################################
# Tests                                                                        #
#                                                                              #
# These targets are used by our continuous integration processes. Use these    #
# directly at your own risk -- they assume required tools (and correct         #
# versions thereof) to be present on your system.                              #
#                                                                              #
# If you prefer to executes these tasks in a container that is pre-loaded with #
# required tools, refer to the hacking section toward the bottom of this file. #
################################################################################

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test-unit
test-unit:
	go test \
		-v \
		-timeout=120s \
		-race \
		-coverprofile=coverage.txt \
		-covermode=atomic \
		./...

################################################################################
# Build: Targets to help build                                                 #
################################################################################

.PHONY: clean
clean:
	rm -rf build

.PHONY: build-base-image
build-base-image:
	mkdir -p build
	cp kargo-render-base.apko.yaml build
	docker run \
		--rm \
		-v $(dir $(realpath $(firstword $(MAKEFILE_LIST))))build:/build \
		-w /build \
		cgr.dev/chainguard/apko \
		build kargo-render-base.apko.yaml kargo-render-base kargo-render-base.tar.gz
	docker image load -i build/kargo-render-base.tar.gz

################################################################################
# Hack: Targets to help you hack                                               #
#                                                                              #
# These targets minimize required developer setup by executing in a container  #
# that is pre-loaded with required tools.                                      #
################################################################################

DOCKER_CMD := docker run \
	-it \
	--rm \
	-v gomodcache:/go/pkg/mod \
	-v $(dir $(realpath $(firstword $(MAKEFILE_LIST)))):/workspaces/kargo-render \
	-w /workspaces/kargo-render \
	kargo-render:dev-tools

.PHONY: hack-build-dev-tools
hack-build-dev-tools:
	docker build -f Dockerfile.dev -t kargo-render:dev-tools .

.PHONY: hack-lint
hack-lint: hack-build-dev-tools
	$(DOCKER_CMD) make lint

.PHONY: hack-test-unit
hack-test-unit: hack-build-dev-tools
	$(DOCKER_CMD) make test-unit

.PHONY: hack-build
hack-build: build-base-image
	docker build \
		--build-arg BASE_IMAGE=kargo-render-base \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		--build-arg GIT_TREE_STATE=$(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi) \
		--tag kargo-render:dev \
		.
