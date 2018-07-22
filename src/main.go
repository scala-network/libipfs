// package main implements libznipfs as an executable instead of a library
package main

import (
	"fmt"
	"os"

	"github.com/stellitecoin/libznipfs/src/zeronet"
)

// main runs libznipfs and retrieves test data from ZeroNet and IPFS
func main() {

	fmt.Println("Test libznipfs as executable")

	zn, err := zeronet.New("/tmp/zn-test")
	if err != nil {
		fmt.Printf("Unable to create ZeroNet instance: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Fetching file from ZeroNet")
	// This is a well-known ZeroNet address. We store the IPFS hash in ipfs.has
	content, err := zn.GetFile("1FAiQ7MddvavaRF6b47fPEY4nSBVJUbCXf", "ipfs.hash")
	if err != nil {
		fmt.Printf("Unable fetch from ZeroNet: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(content))
}
