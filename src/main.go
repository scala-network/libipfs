// package main implements libznipfs as an executable instead of a library
package main

import (
	"fmt"
	"os"
	"path/filepath"
	//"github.com/stellitecoin/libznipfs/src/ipfs"
)

// "github.com/stellitecoin/libznipfs/src/ipfs"
// "github.com/stellitecoin/libznipfs/src/zeronet"
//"github.com/stellitecoin/libznipfs/src/ipfs"
//"github.com/stellitecoin/libznipfs/src/zeronet"

// main runs libznipfs and retrieves test data from ZeroNet and IPFS
func main() {

	baseDataPath := "/tmp"

	ipfsHash := "QmP55FcTYsTJsP4a4aREgq5BrshyPpLHXhvPKXyPfzeHXg"

	// Testing new Go IPFS libraries
	fmt.Println("Testing new go-ipfs libraries")
	ipfsNode, err := NewIPFSNode(filepath.Join(baseDataPath, "ipfs-test"))
	if err != nil {
		fmt.Printf("Unable to create IPFS node: %s\n", err)
		os.Exit(1)
	}

	err = ipfsNode.Start(9001)
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
	fmt.Println(data)
	fmt.Println(string(data))
	//time.Sleep(time.Minute)

	os.Exit(0)

	//
	// OLD V4
	// fmt.Println("Test libznipfs as executable")
	//
	// baseDataPath := "/tmp"
	//
	// zn, err := zeronet.New(filepath.Join(baseDataPath, "zn-test"))
	// if err != nil {
	// 	fmt.Printf("Unable to create ZeroNet instance: %s\n", err)
	// 	os.Exit(1)
	// }
	//
	// // ZeroNet address 1KvWEyqhyHsU9y6UT8xYCFDC8Y1vKaNueX - Live JSON
	// // ZeroNet address 1CH9ApTd83RM8ggz35pZApnKoZqDf7wXyh - Testnet JSON
	// znAddress := "1KvWEyqhyHsU9y6UT8xYCFDC8Y1vKaNueX"
	//
	// fmt.Println("Fetching file from ZeroNet:", znAddress)
	// // This is a well-known ZeroNet address. We store the IPFS hash in ipfs.hash
	// content, err := zn.GetFile(znAddress, "ipfs.hash")
	// if err != nil {
	// 	fmt.Printf("Unable fetch from ZeroNet: %s\n", err)
	// 	os.Exit(1)
	// }
	//
	// // The content retrieved from ZeroNet is the IPFS hash containing the nodelist
	// ipfsHash := strings.TrimSpace(string(content))
	// fmt.Printf("Get IPFS hash from ZeroNet: %s\n", ipfsHash)
	//
	// ipfsNode, err := ipfs.New(filepath.Join(baseDataPath, "ipfs-test"))
	// if err != nil {
	// 	fmt.Printf("Unable to create IPFS node: %s\n", err)
	// 	os.Exit(1)
	// }
	//
	// err = ipfsNode.Start(5001)
	// if err != nil {
	// 	fmt.Printf("Unable to start IPFS node: %s\n", err)
	// 	os.Exit(1)
	// }
	//
	// data, err := ipfsNode.Get(ipfsHash)
	// if err != nil {
	// 	fmt.Printf("Unable fetch data from IPFS node: %s\n", err)
	// 	os.Exit(1)
	// }
	//
	// fmt.Println("Retrieved from IPFS")
	// fmt.Println(string(data))
	// ipfsNode.Stop()
}
