#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := liznipfs

default: build

# Builds as executable for testing
build_test_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-linux-test ./src/main.go

build_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libznipfs.a libznipfs.go

build_windows:
	GOOS=windows \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libznipfs.a libznipfs.go

build_macos:
	GOOS=darwin \
	GOARCH=amd64 \
	go build -buildmode=c-archive -o ./bin/libznipfs.a libznipfs.go

build: build_linux \
	build_windows \
	build_macos

run: build_test_linux
	LOG_FORMAT=Text \
	LOG_LEVEL=Debug \
	./bin/${APP_NAME}-linux-test

clean:
	rm ./bin/*
