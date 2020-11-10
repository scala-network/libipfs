// package name: libipfs
package main

import "C"
import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"github.com/scala-network/libipfs/src/ipfs"
)

var ipfsNode *ipfs.IPFS

// Result holds the seedlist and any error that occurred in the process
// for the daemon to use
type Result struct {
	// Status for the result
	Status string
	// Message to be displayed
	Message string
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C library
}

/**
 * libipfs implements the C-style library for fetching information
 * from ZeroNet and IPFS.
 * Here we only have 3 exported functions that can be called from C
 */

//export IPFSStartNode
// IPFSStartNode starts the IPFS node and initializes ZeroNet
func IPFSStartNode(dataPath *C.char) *C.char {
	// result is marshalled to JSON before being returned to the daemon
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node started on port 5001"),
	}
	var err error
	basePath := C.GoString(dataPath)

	ipfsNode, err = ipfs.New(filepath.Join(basePath, "ipfs"))
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to create IPFS node: %s\n", err)
		return toCJSONString(result)
	}

	err = ipfsNode.Start()
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to start IPFS node: %s\n", err)
	}

	return toCJSONString(result)
}

//export IPFSStopNode
func IPFSStopNode() *C.char{
	result:= Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node stopped"),
	}

    err := ipfsNode.Stop()
    
    if err != nil {
        result.Status = "err"
        result.Message = fmt.Sprintf("IPFS node could not be stopped")
    }
    
    return toCJSONString(result)
}

//export ResolveIPNS
func ResolveIPNS(peerID *C.char) *C.char{

	var err error
	var hash string
	var result Result

	hash, err = ipfsNode.Resolve(C.GoString(peerID))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Could not resolve peer ID: %s\n",err)
	}else{
		result.Status = "ok"
		result.Message = fmt.Sprintf("%s", hash)
	}

	return toCJSONString(result)
}

//export AddDirectory
func AddDirectory(directory *C.char) *C.char{

	var err error
	var hash string
	var result Result

	hash, err = ipfsNode.AddDirectory(C.GoString(directory))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("%s",err)
	}else{
		result.Status = "ok"
		result.Message = fmt.Sprintf("%s", hash)
	}

	return toCJSONString(result)
}

//export BootstrapAdd
func BootstrapAdd(peer *C.char) *C.char{

	var err error
	var response string
	var result Result

	peerArray := []string{C.GoString(peer)}
	response, err = ipfsNode.BootstrapAdd(peerArray)

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("%s",err)
	}else{
		result.Status = "ok"
		result.Message = fmt.Sprintf("%s", response)
	}

	return toCJSONString(result)
}

//export GetPeerID
func GetPeerID() *C.char{

	var err error
	var response string
	var result Result

	response, err = ipfsNode.GetPeerID()

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("%s",err)
	}else{
		result.Status = "ok"
		result.Message = fmt.Sprintf("%s", response)
	}

	return toCJSONString(result)
}

//export PublishToIPNS
func PublishToIPNS(contentHash *C.char) *C.char{

	var err error
	var response string
	var result Result

	response, err = ipfsNode.PublishName(C.GoString(contentHash))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("%s",err)
	}else{
		result.Status = "ok"
		result.Message = fmt.Sprintf("%s", response)
	}

	return toCJSONString(result)
}


// toCJSONString marshals the error result into JSON for the daemon to
// understand and returns it in the required C format
func toCJSONString(result Result) *C.char {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("Fatal error converting result: %s", err))
	}
	return C.CString(string(resultJSON))
}
