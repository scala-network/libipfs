#include "libznipfs-linux.h"
#include <stdio.h>
#include <iostream>

int main() {
  printf("This is a test C application to call libznipfs.\n");

  IPFSStartNode("/tmp/torque");

  std::cout << ZNIPFSGetSeedList("1KvWEyqhyHsU9y6UT8xYCFDC8Y1vKaNueX");

  IPFSStopNode();
}
