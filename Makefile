VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
LDFLAGS = -ldflags "-X github.com/manojkarthick/expenses/cmd.Version=${VERSION}"
TIMESTAMP := $(shell date +%Y%m%d-%H%M%S)
DEVLDFLAGS = -ldflags "-X github.com/manojkarthick/expenses/cmd.Version=dev-${TIMESTAMP}"

hello:
	echo "Hello"

dev-build:
	go build -v ${DEVLDFLAGS}

build:
	go build -v ${LDFLAGS}

run:
	go run main.go

all: hello build