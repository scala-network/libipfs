#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

LIB_NAME := libipfs

default: build

clean:
		rm -rf bin/ && mkdir bin/

build_linux:
		CGO_ENABLED=1 \
		GOOS=linux \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-linux.a ./src/${LIB_NAME}.go

build_windows:
		CGO_ENABLED=1 \
		GOOS=windows \
		GOARCH=amd64 \
		CC=x86_64-w64-mingw32-gcc \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-windows.a ./src/${LIB_NAME}.go

build_macos_x64:
		CGO_ENABLED=1 \
		GOOS=darwin \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-macos-x64.a ./src/${LIB_NAME}.go

build_macos_arm64:
		CGO_ENABLED=1 \
		GOOS=darwin \
		GOARCH=arm64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-macos-arm64.a ./src/${LIB_NAME}.go

build: clean \
	build_linux \
	build_windows \
	build_macos_x64 \
	build_macos_arm64