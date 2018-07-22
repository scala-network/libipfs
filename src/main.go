// package main implements libznipfs as an executable instead of a library
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stellitecoin/libznipfs/src/ipfs"
	"github.com/stellitecoin/libznipfs/src/zeronet"
)

// main runs libznipfs and retrieves test data from ZeroNet and IPFS
func main() {

	fmt.Println("Test libznipfs as executable")

	baseDataPath := "/tmp"

	zn, err := zeronet.New(filepath.Join(baseDataPath, "zn-test"))
	if err != nil {
		fmt.Printf("Unable to create ZeroNet instance: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Fetching file from ZeroNet")
	// This is a well-known ZeroNet address. We store the IPFS hash in ipfs.hash
	content, err := zn.GetFile("1FAiQ7MddvavaRF6b47fPEY4nSBVJUbCXf", "ipfs.hash")
	if err != nil {
		fmt.Printf("Unable fetch from ZeroNet: %s\n", err)
		os.Exit(1)
	}

	// The content retrieved from ZeroNet is the IPFS hash containing the nodelist
	ipfsHash := strings.TrimSpace(string(content))
	fmt.Printf("Get IPFS hash from ZeroNet: %s\n", ipfsHash)

	ipfsNode, err := ipfs.New(filepath.Join(baseDataPath, "ipfs-test"))
	if err != nil {
		fmt.Printf("Unable to create IPFS node: %s\n", err)
		os.Exit(1)
	}

	err = ipfsNode.Start(5001)
	if err != nil {
		fmt.Printf("Unable to start IPFS node: %s\n", err)
		os.Exit(1)
	}

	data, err := ipfsNode.Get(ipfsHash)
	if err != nil {
		fmt.Printf("Unable fetch data from IPFS node: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Retrieved from IPFS")
	fmt.Println(string(data))
	ipfsNode.Stop()
}
