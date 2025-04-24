cmake_minimum_required(VERSION 3.10)
project(libipfs_go)

if(APPLE)
    set(GOOS "darwin")
    find_library(SECURITY_FRAMEWORK Security REQUIRED)
    find_library(RESOLV_LIBRARY resolv REQUIRED)
elseif(UNIX)
    if(CMAKE_SYSTEM_NAME MATCHES "FreeBSD")
        set(GOOS "freebsd")
    else()
        set(GOOS "linux")
    endif()
elseif(WIN32)
    set(GOOS "windows")
else()
    message(FATAL_ERROR "Unsupported OS")
endif()

if(NOT DEFINED ARCH_ID)
  message(FATAL_ERROR "ARCH_ID must be passed from parent")
endif()

if(DEFINED ARCH_ID)
  message(STATUS "Using Scala ARCH_ID: ${ARCH_ID}")
  set(GOARCH_FROM_ARCH_ID ON)
  if(ARCH_ID STREQUAL "x86_64")
    set(GOARCH "amd64")
  elseif(ARCH_ID MATCHES "^(aarch64|arm64|armv8-a)$")
    set(GOARCH "arm64")
  elseif(ARCH_ID STREQUAL "riscv64")
    set(GOARCH "riscv64")
  else()
    message("Unsupported ARCH_ID: ${ARCH_ID} DEFAULTING TO amd64")
    set(GOARCH "amd64")
  endif()
else()
  message(STATUS "Deriving GOARCH from CMAKE_SYSTEM_PROCESSOR")
  string(TOLOWER "${CMAKE_SYSTEM_PROCESSOR}" SYSTEM_ARCH)
  message(STATUS "Detected system arch: ${SYSTEM_ARCH}")

  if(SYSTEM_ARCH STREQUAL "x86_64" OR SYSTEM_ARCH STREQUAL "amd64")
    set(GOARCH "amd64")
  elseif(SYSTEM_ARCH MATCHES "^(aarch64|arm64|armv8-a)$")
    set(GOARCH "arm64")
  elseif(SYSTEM_ARCH STREQUAL "riscv64")
    set(GOARCH "riscv64")
  else()
    message(FATAL_ERROR "Unsupported system processor architecture: ${SYSTEM_ARCH}")
  endif()
endif()

set(GO_OUT libipfs-${GOOS}-${GOARCH}.a)
set(GO_HEADER libipfs.h)

set(CROSS_CC "")
if(GOARCH STREQUAL "arm64")
  set(CROSS_CC "CC=aarch64-linux-gnu-gcc")
endif()

if(GOARCH STREQUAL "riscv64")
  set(CROSS_CC "CC=riscv64-linux-gnu-gcc")
endif()

if(GOOS STREQUAL "windows" AND CMAKE_C_PLATFORM_ID STREQUAL "MinGW")
  set(CROSS_CC "CC=x86_64-w64-mingw32-gcc-posix")
endif()

if(APPLE)
    	if(DEFINED scala_BINARY_DIR AND NOT "${scala_BINARY_DIR}" STREQUAL "")
		if(DEFINED CMAKE_OSX_SYSROOT)
			set(CROSS_CC "CC=clang --target=${CLANG_TARGET} --sysroot=${CMAKE_OSX_SYSROOT}")
		endif()
	endif()
endif()

if(FREEBSD)
    if(DEFINED scala_BINARY_DIR AND NOT "${scala_BINARY_DIR}" STREQUAL "")
        foreach(path IN LISTS CMAKE_CXX_IMPLICIT_LINK_DIRECTORIES)
            if("${path}" MATCHES "(.*/native)/usr/lib")
                set(FREEBSD_SYSROOT "${CMAKE_MATCH_1}")
                break()
            endif()
        endforeach()

        if(DEFINED FREEBSD_SYSROOT)
            set(CROSS_CC "CC=clang --target=x86_64-unknown-freebsd --sysroot=${FREEBSD_SYSROOT}")
        else()
            message(FATAL_ERROR "Could not detect FreeBSD sysroot from CMAKE_CXX_IMPLICIT_LINK_DIRECTORIES")
        endif()
    endif()
endif()

add_custom_command(
    OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/${GO_OUT} ${CMAKE_CURRENT_BINARY_DIR}/${GO_HEADER}
    COMMAND ${CMAKE_COMMAND} -E env
            ${CROSS_CC}
            CGO_ENABLED=1
            GOOS=${GOOS}
            GOARCH=${GOARCH}
            go build
            -buildmode=c-archive
            -trimpath
            "-ldflags=-s -w"
            -o ${CMAKE_CURRENT_BINARY_DIR}/${GO_OUT}
            ${CMAKE_CURRENT_SOURCE_DIR}/libipfs.go
    ${GO_EXTRA_COMMANDS}
    COMMAND ${CMAKE_COMMAND} -E copy_if_different
            ${CMAKE_CURRENT_BINARY_DIR}/libipfs-${GOOS}-${GOARCH}.h
            ${CMAKE_CURRENT_BINARY_DIR}/libipfs.h
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/libipfs.go
    COMMENT "Building Go IPFS library (${GOOS}/${GOARCH})"
)

add_custom_target(libipfs_build ALL
    DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/${GO_OUT}
)

add_library(libipfs STATIC IMPORTED GLOBAL)
set_target_properties(libipfs PROPERTIES
    IMPORTED_LOCATION ${CMAKE_CURRENT_BINARY_DIR}/${GO_OUT}
    INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR}
)

if(APPLE)
  set_target_properties(libipfs PROPERTIES
    INTERFACE_LINK_LIBRARIES "${SECURITY_FRAMEWORK};${RESOLV_LIBRARY}"
  )
elseif(UNIX AND NOT FREEBSD)
  set_target_properties(libipfs PROPERTIES
    INTERFACE_LINK_LIBRARIES "resolv"
  )
endif()

add_dependencies(libipfs libipfs_build)
