#include "libznipfs-linux.h"
#include <stdio.h>
#include <iostream>

int main() {
  printf("This is a test C application to call libznipfs.\n");
  std::cout << IPFSStartNode("/tmp/torque") << std::endl;
}
