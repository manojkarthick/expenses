PACKAGE_NAME := github.com/manojkarthick/expenses
GOLANG_VERSION := 1.17.5
GOLANG_CROSS_VERSION := v$(GOLANG_VERSION)

VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
LDFLAGS = -ldflags "-X github.com/manojkarthick/expenses/cmd.Version=${VERSION}"
TIMESTAMP := $(shell date +%Y%m%d-%H%M%S)
DEVLDFLAGS = -ldflags "-X github.com/manojkarthick/expenses/cmd.Version=dev-${TIMESTAMP}"

.PHONY: start
start:
	echo "Building expenses!"

.PHONY: dev-build
dev-build:
	go build -v ${DEVLDFLAGS}

.PHONY: build
build:
	go build -v ${LDFLAGS}

.PHONY: run
run:
	go run main.go

.PHONY: all
all: start build

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e PACKAGE_VERSION=$(VERSION) \
		-e GOLANG_VERSION=$(GOLANG_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		troian/golang-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e GOLANG_VERSION=$(GOLANG_VERSION) \
		-e PACKAGE_VERSION=$(VERSION) \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		troian/golang-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist
