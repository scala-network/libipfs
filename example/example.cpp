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

    std::cout << "Resolving /ipns/ipfs.io to IPFS hash" << std::endl;
    std::cout << ResolveIPNSName((char*)"/ipns/ipfs.io") << std::endl;

    std::cout << "Publishing to IPNS" << std::endl;
    std::cout << PublishIPFSName((char*)"QmXL7PCYH8VtMxkVTPBXxpnpF893QrLHC5H5AKv2FaAExU") << std::endl;

    std::cout << "Stopping IPFS node" << std::endl;
    std::cout << IPFSStopNode() << std::endl;
}
