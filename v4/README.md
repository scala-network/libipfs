# libznipfs

A C-style library implemented in Go to retrieve the seedlist from ZeroNet and IPFS

## Overview

libznipfs is used by the Stellite daemon to retrieve information from ZeroNet and
IPFS. It starts a full IPFS node and exposes some basic functionality to the
daemon.

When libznipfs is started from the daemon it runs a full IPFS node including the
HTTP API. This means standard IPFS commands can be used while the Stellite daemon
is running.

To use the API, you need to set the IPFS path to your Stellite data directory

Example:
`IPFS_PATH=~/.stellite/ipfs ipfs cat /ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv/readme`

This will print the default IPFS readme

### Why Go and not C or C++

Currently, no simple implementation or API exists for ZeroNet and IPFS in C or C++. Instead of writing, or re-writing, large parts of ZeroNet and IPFS in C or C++ we rather use Go and compile it to a C or C++ compatible library. IPFS is implemented in Go already and a Go library for ZeroNet already exist.

## Compiling

Note: It's easier to get your environment set up for building Stellite before starting. Building the library requires some of the same dependencies.

### Linux and MacOS

*Note on building: Current go-ipfs master seems broken, for now, use v0.4.17*

1. [Download](https://golang.org/dl/) and [install](https://golang.org/doc/install#tarball) Go 1.10 or higher
2. Set you GOPATH, or have it default to ~/go
3. Get and install the [Go IPFS source](https://github.com/ipfs/go-ipfs#build-from-source)
4. Grab our ZeroNet/IPFS library

```
go get -u -d github.com/stellitecoin/libznipfs/src
cd $GOPATH/src/github.com/stellitecoin/libznipfs
```

On Linux
`make build_linux`

On MacOS
`make build_macos`

You should now have libznipfs.a and libznipfs.h in the `bin` directory

If you want to build the Stellite daemon, copy `libznipfs.a` and `libznipfs.h` to the `external/libznipfs` directory in the daemon source.

### Windows

Install MSYS2 (x86_64)
https://www.msys2.org/

Open `MSYS2 MinGW 64-bit` as Administrator and from the terminal do:

1. [Download](https://golang.org/dl/) and [install](https://golang.org/doc/install#tarball) Go 1.10 or higher for Windows
	* You should use the 'archive' version and not MSI - ex. go1.10.3.windows-amd64.zip
	* You'll need to run `pacman -S unzip` to be able to extract the Go build
2. Set you GOPATH to /c/go by running `export GOPATH=/c/go`
3. Get and install the [Go IPFS source](https://github.com/ipfs/go-ipfs#build-from-source)
4. Grab our ZeroNet/IPFS library

```
go get -u -d github.com/stellitecoin/libznipfs/src
cd $GOPATH/src/github.com/stellitecoin/libznipfs
make build_windows
```

You should now have libznipfs.a and libznipfs.h in the `bin` directory

If you want to build the Stellite daemon, copy `libznipfs.a` and `libznipfs.h` to the `external/libznipfs` directory in the daemon source.
