// package name: libznipfs
package main

import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stellitecoin/libznipfs/src/ipfs"
	"github.com/stellitecoin/libznipfs/src/zeronet"
)

var zn *zeronet.ZeroNet
var ipfsNode *ipfs.IPFS

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C library
}

/**
 * libznipfs implements the C-style library for fetching information
 * from ZeroNet and IPFS.
 * Here we only have 3 exported functions that can be called from C
 * At this time, any errors that occur in fetching data would should the
 * daemon down
 */

//export StartNode
func StartNode(dataPath *C.char) {
	// Start the ZN/IPFS node

	var err error
	basePath := C.GoString(dataPath)

	zn, err = zeronet.New(filepath.Join(basePath, "zeronet"))
	if err != nil {
		fmt.Printf("Unable to create ZeroNet instance: %s\n", err)
		os.Exit(1)
	}

	ipfsNode, err = ipfs.New(filepath.Join(basePath, "ipfs"))
	if err != nil {
		fmt.Printf("Unable to create IPFS node: %s\n", err)
		os.Exit(1)
	}

	err = ipfsNode.Start(5001)
	if err != nil {
		fmt.Printf("Unable to start IPFS node: %s\n", err)
		os.Exit(1)
	}
}

//export GetSeedList
func GetSeedList(zeroNetAddress *C.char) *C.char {
	// GetSeedList Retrieve the seedlist using ZeroNet and IPFS
	// Returns the address list from the given ZeroNet address

	address := C.GoString(zeroNetAddress)

	// This is a well-known ZeroNet address. We store the IPFS hash in ipfs.hash
	content, err := zn.GetFile(address, "ipfs.hash")
	if err != nil {
		fmt.Printf("Unable fetch from ZeroNet: %s\n", err)
		os.Exit(1)
	}
	ipfsHash := strings.TrimSpace(string(content))

	data, err := ipfsNode.Get(ipfsHash)
	if err != nil {
		fmt.Printf("Unable fetch data from IPFS node: %s\n", err)
		os.Exit(1)
	}

	return C.CString(string(data))
}

//export StopNode
func StopNode() {
	// Stop the ZN/IPFS node
	ipfsNode.Stop()
}
