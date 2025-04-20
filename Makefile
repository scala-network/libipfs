.PHONY: default build clean

VERSION := 3.0.1
LIB_NAME := libipfs
SRC := ./${LIB_NAME}.go
BIN_DIR := bin

PLATFORMS := \
	linux/amd64 \
    freebsd/amd64 \
    windows/amd64 \
    darwin/amd64 \
    darwin/arm64 \
    linux/riscv64 \
    linux/arm64

# Define cross-compilers for each architecture
CC_COMPILERS := \
    freebsd/amd64=clang --target=x86_64-unknown-freebsd \
    windows/amd64=x86_64-w64-mingw32-gcc \
    linux/amd64=gcc \
    linux/arm64=aarch64-linux-gnu-gcc \
    linux/riscv64=riscv64-linux-gnu-gcc \
    darwin/amd64=o64-clang \
    darwin/arm64=o64-clang

# Default target
default: build

# Clean target
clean:
	rm -rf $(BIN_DIR)/ && mkdir -p $(BIN_DIR) && cp -rf example.cpp $(BIN_DIR)/

# Function to get the cross-compiler based on GOOS/GOARCH
define GET_CC
$(shell echo $(CC_COMPILERS) | tr ' ' '\n' | grep -E '$(1)/$(2)=' | cut -d= -f2)
endef

# Build rule for each platform
define BUILD_RULE
build_$(1)_$(2):
	@echo "Building for $(1)/$(2)..."
	$(if $(filter $(1)/$(2),freebsd/amd64), \
		env SYSROOT=/usr/local/freebsd-sysroot CGO_CFLAGS="--sysroot=/usr/local/freebsd-sysroot" CGO_LDFLAGS="--sysroot=/usr/local/freebsd-sysroot",) \
	CGO_ENABLED=1 \
	GOOS=$(1) \
	GOARCH=$(2) \
	CC=$(call GET_CC,$(1),$(2)) \
	go build -buildmode=c-archive -ldflags="-s -w" -trimpath -o $(BIN_DIR)/${LIB_NAME}-$(1)-$(2).a $(SRC)
.PHONY: build_$(1)_$(2)
endef

# Generate build rules for all platforms
$(foreach plat, $(PLATFORMS), \
	$(eval $(call BUILD_RULE,$(word 1,$(subst /, ,$(plat))),$(word 2,$(subst /, ,$(plat))))))

# Build all platforms
build-all: $(foreach plat, $(PLATFORMS), build_$(subst /,_,$(plat)))

# Main build target
build: clean build-all