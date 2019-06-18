// package main implements libznipfs as an executable instead of a library
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/contribute-torque/libznipfs/src/zeronet"
	//	shell "github.com/ipfs/go-ipfs-api"
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

	// ZeroNet address 1KvWEyqhyHsU9y6UT8xYCFDC8Y1vKaNueX - Live JSON
	// ZeroNet address 1CH9ApTd83RM8ggz35pZApnKoZqDf7wXyh - Testnet JSON
	znAddress := "1KvWEyqhyHsU9y6UT8xYCFDC8Y1vKaNueX"

	fmt.Println("Fetching file from ZeroNet:", znAddress)
	// This is a well-known ZeroNet address. We store the IPFS hash in ipfs.hash
	content, err := zn.GetFile(znAddress, "ipfs.hash")
	if err != nil {
		fmt.Printf("Unable fetch from ZeroNet: %s\n", err)
		os.Exit(1)
	}

	// The content retrieved from ZeroNet is the IPFS hash containing the nodelist
	ipfsHash := strings.TrimSpace(string(content))
	fmt.Printf("Get IPFS hash from ZeroNet: %s\n", ipfsHash)

	//
	// NOTE Start IPFS extract and run daemon
	//
	// fmt.Println("Testing new embedded ipfs")
	//
	// //baseDataPath := "/tmp"
	//
	// fileBytes, err := FSByte(false, "/pack/linux/ipfs")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// f, err := os.OpenFile("/tmp/ipfs", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// n, err := f.Write(fileBytes)
	// if err != nil {
	// 	panic(err)
	// }
	// f.Close()
	//
	// fmt.Println("Wrote", n)
	// fmt.Println(len(fileBytes))
	//
	// cmd := exec.Command("/tmp/ipfs", "daemon")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// op, err := cmd.CombinedOutput()
	// fmt.Println(string(op))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(string(op))
	//
	// NOTE: End IPFS extract and run daemon
	//

	//
	// NOTE: Fetch via IPFS API
	//
	// sh := shell.NewShell("localhost:5001")
	// err := sh.Get("QmbQVPLwUSbLQL3tkQGFRFWDTCjRs6LSJDPEjG24BnNuhD", "/tmp/testfile.txt")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// NOTE: End Fetch via IPFS API
	//

}
