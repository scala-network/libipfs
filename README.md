# libIPFS
A C-style library library that wraps around go-ipfs as a library(not an executable that is embedded) and provides a very simple API.

## Example

```cpp
#include "libipfs-linux.h"
#include <iostream>

int main() {
    /* Starts the IPFS node */
    std::cout << IPFSStartNode("./") << std::endl;
    /* Resolve an IPNS name */
    std::cout << ResolveIPNSName("/ipns/ipfs.io") << std::endl;
    /* Stop IPFS Node */
    std::cout << IPFSStopNode() << std::endl;
}
```

## Overview

libIPFS is used by the Scala Network Project to retrieve and publish critical information on to IPFS.

It runs a barebones IPFS instance and provides functions to be called from C/C++.

### Why?

Currently, no simple implementation or API exists for IPFS in C or C++. Instead of writing, or re-writing, large parts of IPFS in C or C++ we rather use Go and compile it to a C or C++ compatible library. IPFS is implemented in Go already.

###  Building

#### Requirements

* go >= 1.16
* make >= 4.2.1
* gcc and g++ >= 9.3.0

To build the library you can use the following commands, the outputs can be found in bin/

```bash
go mod download
make build
```

You **need** an actual mac to build mac binaries, run:

```
make build build_macos_x64 or build_macos_arm64
```

###  LICENSE

View [LICENCE](https://github.com/scala-network/libipfs/blob/master/LICENSE)
