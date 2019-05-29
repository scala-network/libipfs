// package main implements libznipfs as an executable instead of a library
package main

import (
	shell "github.com/ipfs/go-ipfs-api"
)

// main runs libznipfs and retrieves test data from ZeroNet and IPFS
func main() {

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
	sh := shell.NewShell("localhost:5001")
	err := sh.Get("QmbQVPLwUSbLQL3tkQGFRFWDTCjRs6LSJDPEjG24BnNuhD", "/tmp/testfile.txt")
	if err != nil {
		panic(err)
	}
	//
	// NOTE: End Fetch via IPFS API
	//

}
