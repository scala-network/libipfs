#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := libipfs

default: build

# Embeds the ipfs binary for Linux
package_linux:
	esc -pkg ipfs -o src/ipfs/pack.go pack/linux

# Embeds the ipfs binary for Windows
package_windows:
	esc -pkg ipfs -o src/ipfs/pack.go pack/windows

# Embeds the ipfs binary for MacOS
package_macos:
	esc -pkg ipfs -o src/ipfs/pack.go pack/darwin

# Builds as executable for testing
build_test_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-linux-test ./src/*.go

build_linux: package_linux
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libipfs-linux.a ./src/libipfs.go

build_windows: package_windows
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build -buildmode=c-archive -o ./bin/libipfs-windows.a ./src/libipfs.go

build_macos: package_macos
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libipfs-macos.a ./src/libipfs.go

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
