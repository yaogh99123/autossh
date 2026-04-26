PROJECT="autossh"
VERSION="v1.2.1"
BUILD=`date +%FT%T%z`
LDFLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all build vendor tidy clean

all: build

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

build: vendor
	go build -mod=vendor -o ${PROJECT} -ldflags ${LDFLAGS} src/main/main.go

clean:
	rm -f ${PROJECT}
	rm -rf ./bin
	rm -rf ./releases
	rm -rf ./vendor
