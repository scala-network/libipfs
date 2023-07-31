# libIPFS

A C-style library library that wraps around go-ipfs as a library(not an executable that is embedded) and provides a very simple API.

## Example

```cpp
#include "libipfs-linux.h"
#include <iostream>

int main() {
    std::cout << "Starting IPFS node" << std::endl;
    std::cout << IPFSStartNode((char*)"", 0) << std::endl;

    std::cout << "Getting peer ID" << std::endl;
    std::cout << GetPeerID() << std::endl;

    std::cout << "Adding file to IPFS" << std::endl;
    std::cout << IpfsAdd((char*)"./test-file.jpg") << std::endl;

    std::cout << "Adding directory to IPFS" << std::endl;
    std::cout << IpfsAdd((char*)"./test-dir") <<std::endl;

    std::cout << "Downloading file from IPFS" << std::endl;
    std::cout << IpfsGet((char*)"QmXL7PCYH8VtMxkVTPBXxpnpF893QrLHC5H5AKv2FaAExU", (char*)"./test-download.jpg") << std::endl;

    std::cout << "Pinning file to IPFS" << std::endl;
    std::cout << IpfsPin((char*)"QmXL7PCYH8VtMxkVTPBXxpnpF893QrLHC5H5AKv2FaAExU") << std::endl;

    std::cout << "Locally Pinned Hashes" << std::endl;
    std::cout << IpfsGetPinnedHashes() << std::endl;

    std::cout << "Resolving /ipns/ipfs.io to IPFS hash" << std::endl;
    std::cout << ResolveIPNSName((char*)"/ipns/ipfs.io") << std::endl;

    std::cout << "Publishing to IPNS" << std::endl;
    std::cout << PublishIPFSName((char*)"QmXL7PCYH8VtMxkVTPBXxpnpF893QrLHC5H5AKv2FaAExU") << std::endl;

    std::cout << "Stopping IPFS node" << std::endl;
    std::cout << IPFSStopNode() << std::endl;
}
```

## Overview

libIPFS is used by the Scala Network Project to retrieve and publish critical information on to IPFS.

It runs a *barebones* IPFS instance and provides functions to be called from C/C++.

### Why?

Currently, no simple implementation or API exists for IPFS in C or C++. Instead of writing, or re-writing, large parts of IPFS in C or C++ we rather use Go and compile it to a C or C++ compatible library. IPFS is implemented in Go already.

## Building

### Requirements

* go >= 1.16
* make >= 4.2.1
* gcc and g++ >= 9.3.0

To build the library you can use the following commands, the outputs can be found in bin/

```
go mod download
make build
```
