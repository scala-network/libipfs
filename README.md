# libIPFS

A C-style library library that wraps around go-ipfs and provides a simple API.

## Example

```cpp
#include "libipfs-linux.h"
#include <iostream>

int main() {

/* Starts the IPFS node */
std::cout << IPFSStartNode("./") << std::endl;

/* Add a custom bootstrap */
/* std::cout << BootstrapAdd("/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWJjbW3sQPpvsA5saTZfihc2xQ1dwN2xXog2xAZYifQFmR") << std::endl; */

/* Resolve an IPNS name */
std::cout << ResolveIPNS("12D3KooWJjbW3sQPpvsA5saTZfihc2xQ1dwN2xXog2xAZYifQFmR") << std::endl;

/* Add a directory to IPFS */
std::cout << AddDirectory("./test") << std::endl;

/* Stop IPFS Node */
std::cout << IPFSStopNode() << std::endl;
}
```

## Overview

libIPFS is used by the Scala Network Project to retrieve and publish critical information on to IPFS.

It runs a full blown IPFS instance and even exposes the HTTP API and Gateway of the underlying daemon.

### Why Go and not C or C++

Currently, no simple implementation or API exists for IPFS in C or C++. Instead of writing, or re-writing, large parts of IPFS in C or C++ we rather use Go and compile it to a C or C++ compatible library. IPFS is implemented in Go already.

###  Building

#### Requirements

* go >= 1.14.2
* make >= 4.2.1
* gcc and g++ >= 9.3.0

To build the library you can use the following commands, the outputs can be found in bin/

```bash
git clone https://github.com/scala-network/libipfs
cd libipfs/
make
```

###  LICENSE

View [LICENCE](https://github.com/scala-network/libipfs/blob/master/LICENSE)