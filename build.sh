#!/bin/bash

export GO111MODULE="on"
go mod tidy
go mod vendor

PROJECT="autossh"
VERSION="v1.2.0"
BUILD=`date +%FT%T%z`
CURRENT_OS=`go env GOOS`
CURRENT_ARCH=`go env GOARCH`

function build() {
    os=$1
    arch=$2
    alias_name=$3
    package="${PROJECT}-${alias_name}-${arch}_${VERSION}"

    echo "build ${package} ..."
    mkdir -p "./releases/${package}"
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -mod=vendor -o "./releases/${package}/autossh" -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}" src/main/main.go

    # 如果是当前系统平台，则复制到根目录方便调试
    if [ "${os}" == "${CURRENT_OS}" ] && [ "${arch}" == "${CURRENT_ARCH}" ]; then
        echo "copy binary to root for debugging..."
        cp "./releases/${package}/autossh" .
    fi

    cp ./config.example.yml "./releases/${package}/config.yml"
    chmod +x ./install
    cp ./install "./releases/${package}/install"
    cd ./releases/
    zip -r "./${package}.zip" "./${package}"
    echo "clean ${package}"
    rm -rf "./${package}"
    cd ../
}

# OS X Mac
build darwin amd64 macOS
build darwin arm64 macOS

# Linux
build linux amd64 linux
build linux 386 linux
build linux arm linux
