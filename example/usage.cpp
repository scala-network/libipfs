#include "libipfs-linux.h"
#include <iostream>
#include <cstring>
#include <cstdint>
#include <unistd.h>

uint8_t strContains(char* string, char* toFind)
{
    uint8_t slen = strlen(string);
    uint8_t tFlen = strlen(toFind);
    uint8_t found = 0;

    if( slen >= tFlen )
    {
        for(uint8_t s=0, t=0; s<slen; s++)
        {
            do{

                if( string[s] == toFind[t] )
                {
                    if( ++found == tFlen ) return 1;
                    s++;
                    t++;
                }
                else { s -= found; found=0; t=0; }

              }while(found);
        }
        return 0;
    }
    else return -1;
}

int main() {
  /* Starts the IPFS node */
  char* ipfsStart = IPFSStartNode("./");

  if(strContains(ipfsStart, "IPFS node started on port 5001") == 1){
	std::cout << ipfsStart << std::endl;
	sleep(20);
	std::cout << BootstrapAdd("/ip4/144.91.88.34/tcp/4001/p2p/12D3KooWCmwuJ4qf6hGBjeGodApT3MjQ4mpB6ZxeUh74hviHCh87") << std::endl;
	std::cout << ResolveIPNS("/ipns/12D3KooWCmwuJ4qf6hGBjeGodApT3MjQ4mpB6ZxeUh74hviHCh87") << std::endl;
	std::cout << AddDirectory("./test") << std::endl;
	std::cout << PublishToIPNS("QmcwV6oxW4XZisKdEK4V4j8ouA75qfdT8dBtm61mbruoKF");
	std::cout << IPFSStopNode() << std::endl;
  }
}
