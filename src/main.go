// package main implements libznipfs as an executable instead of a library
package main

import (
	"fmt"

	"github.com/scala-network/libznipfs/src/ipfs"
)

// main runs libznipfs and retrieves test data from ZeroNet and IPFS
func main() {

	fmt.Println("Test libznipfs as executable")

	baseDataPath := "/tmp"

	// The content retrieved from ZeroNet is the IPFS hash containing the nodelist
	ipfsHash := "Qmf2s5yVmbKkxJfNA4zzENYw11Mnq4HBz6RRmYNGRoYg57"
	fmt.Printf("Get IPFS hash from ZeroNet: %s\n", ipfsHash)

	ipfsNode, err := ipfs.New(baseDataPath)
	if err != nil {
		panic(err)
	}

	err = ipfsNode.Start()
	if err != nil {
		panic(err)
	}

	result, err := ipfsNode.Get(ipfsHash)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))

	err = ipfsNode.Stop()
	if err != nil {
		panic(err)
	}

}
