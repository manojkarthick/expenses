VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
LDFLAGS = -ldflags "-X expenses/cmd.Version=${VERSION}"
TIMESTAMP := $(shell date +%Y%m%d-%H%M%S)
DEVLDFLAGS = -ldflags "-X expenses/cmd.Version=dev-${TIMESTAMP}"

hello:
	echo "Hello"

dev-build:
	go build -v ${DEVLDFLAGS}

build:
	go build -v ${LDFLAGS}

run:
	go run main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: hello build