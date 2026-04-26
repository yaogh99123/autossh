PROJECT="autossh"
VERSION="v1.2.0"
BUILD=`date +%FT%T%z`
LDFLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all build vendor tidy clean

all: build

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

build: vendor
	mkdir -p ./bin
	go build -mod=vendor -o ./bin/${PROJECT} -ldflags ${LDFLAGS} src/main/main.go

clean:
	rm -rf ./bin
	rm -rf ./releases
	rm -rf ./vendor
