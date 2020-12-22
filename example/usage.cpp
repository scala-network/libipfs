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
