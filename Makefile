.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

LIB_NAME := libipfs
default: build

clean:
		rm -rf bin/ && mkdir bin/

build_linux_x64:
		CGO_ENABLED=1 \
		GOOS=linux \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-linux.a ./src/${LIB_NAME}.go

build_linux_arm64:
		CGO_ENABLED=1 \
		GOOS=linux \
		GOARCH=arm64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-linux-arm64.a ./src/${LIB_NAME}.go

build_windows_x64:
		CGO_ENABLED=1 \
		GOOS=windows \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-windows.a ./src/${LIB_NAME}.go

build_linux_riscv:
		CGO_ENABLED=1 \
		GOOS=linux \
		GOARCH=riscv64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-linux-riscv64.a ./src/${LIB_NAME}.go

build_darwin_x64:
		CGO_ENABLED=1 \
		GOOS=darwin \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-darwin.a ./src/${LIB_NAME}.go

build_darwin_arm64:
		CGO_ENABLED=1 \
		GOOS=darwin \
		GOARCH=arm64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-darwin-arm64.a ./src/${LIB_NAME}.go

build_freebsd_x64:
		CGO_ENABLED=1 \
		GOOS=freebsd \
		GOARCH=amd64 \
		go build -buildmode=c-archive -o ./bin/${LIB_NAME}-freebsd.a ./src/${LIB_NAME}.go

build-all: build_linux_x64 build_linux_arm64 build_windows_x64 build_linux_riscv build_darwin_x64 build_darwin_arm64 build_freebsd_x64

build: clean \
	build-all