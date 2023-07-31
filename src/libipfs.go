package main

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/scala-network/libipfs/src/constants"
	"github.com/scala-network/libipfs/src/ipfs"
)

type Result struct {
	Status  string
	Message string
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C library
}

//export IPFSStartNode
func IPFSStartNode(dataPath *C.char, P2PPort int) *C.char {

	gDp := C.GoString(dataPath)

	if gDp == "" {
		gDp = constants.DefaultRepoPath
	}

	if P2PPort == 0 {
		P2PPort = constants.DefaultP2PPort
	}

	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node started on port %d", P2PPort),
	}

	err := ipfs.Start(gDp, P2PPort)

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to create IPFS node: %s\n", err)
		return toCJSONString(result)
	}

	return toCJSONString(result)
}

//export IPFSStopNode
func IPFSStopNode() *C.char {
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node stopped"),
	}

	err := ipfs.Stop()

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("IPFS node could not be stopped")
	}

	return toCJSONString(result)
}

//export ResolveIPNSName
func ResolveIPNSName(name *C.char) *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	resolvedPath, err := ipfs.ResolveName(C.GoString(name))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Name couldn't be resolved %s", err.Error())
	} else {
		result.Message = resolvedPath
	}

	return toCJSONString(result)
}

//export PublishIPFSName
func PublishIPFSName(ipfsHash *C.char) *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	name, err := ipfs.PublishName(C.GoString(ipfsHash))
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Couldn't publish to IPNS %s", err.Error())
	} else {
		result.Message = name
	}

	return toCJSONString(result)
}

//export GetPeerID
func GetPeerID() *C.char {
	peerId := ipfs.GetPeerID()
	result := Result{
		Status:  "ok",
		Message: peerId,
	}

	return toCJSONString(result)
}

//export IpfsAdd
func IpfsAdd(addPath *C.char) *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	ipfsHash, err := ipfs.Add(C.GoString(addPath))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Couldn't add to IPFS %s", err.Error())
	} else {
		result.Message = ipfsHash
	}

	return toCJSONString(result)
}

//export IpfsPin
func IpfsPin(ipfsHash *C.char) *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	err := ipfs.Pin(C.GoString(ipfsHash))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Couldn't pin to IPFS %s", err.Error())
	} else {
		result.Message = fmt.Sprintf("Pinned hash %s", C.GoString(ipfsHash))
	}

	return toCJSONString(result)
}

//export IpfsGetPinnedHashes
func IpfsGetPinnedHashes() *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	hashes, err := ipfs.GetPinnedHashes()

	hG := make(map[string][]string)
	hG["hashes"] = hashes
	hJ, err := json.Marshal(hG)

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Couldn't get pinned hashes from IPFS %s", err.Error())
	} else {
		result.Message = fmt.Sprintf("%s", hJ)
	}

	return toCJSONString(result)
}

//export IpfsGet
func IpfsGet(ipfsHash *C.char, downloadPath *C.char) *C.char {
	result := Result{
		Status:  "ok",
		Message: "",
	}

	err := ipfs.Get(C.GoString(ipfsHash), C.GoString(downloadPath))

	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Couldn't download from IPFS %s", err.Error())
	} else {
		result.Status = "ok"
		result.Message = fmt.Sprintf("Downloaded hash %s to %s", C.GoString(ipfsHash), C.GoString(downloadPath))
	}

	return toCJSONString(result)
}

func toCJSONString(result Result) *C.char {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("Fatal error converting result: %s", err.Error()))
	}
	return C.CString(string(resultJSON))
}
