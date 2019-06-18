#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := libznipfs

default: build

package:
	esc -pkg ipfs -o src/ipfs/pack.go pack/linux

# Builds as executable for testing
build_test_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-linux-test ./src/*.go

build_linux:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libznipfs-linux.a ./src/libznipfs.go

build_windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build -buildmode=c-archive -o ./bin/libznipfs-windows.a ./src/libznipfs.go

build_macos:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libznipfs-mac.a ./src/libznipfs.go

build: build_linux \
	build_windows \
	build_macos

run: build_test_linux
	LOG_FORMAT=Text \
	LOG_LEVEL=Debug \
	./bin/${APP_NAME}-linux-test

run_race:
	GOOS=linux \
	GOARCH=amd64 \
	go run -race ./src/*.go

clean:
	rm ./bin/*
