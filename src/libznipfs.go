// package name: libznipfs
package main

import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stellitecoin/libznipfs/src/ipfs"
	"github.com/stellitecoin/libznipfs/src/zeronet"
)

var zn *zeronet.ZeroNet
var ipfsNode *ipfs.IPFS
var checkpoints map[int]string
var checkpointMutex sync.RWMutex

// Result holds the seedlist and any error that occurred in the process
// for the daemon to use
type Result struct {
	// Status for the result
	Status string
	// Message to be displayed
	Message string
	// The seedlist
	Seedlist []string
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C library
}

/**
 * libznipfs implements the C-style library for fetching information
 * from ZeroNet and IPFS.
 * Here we only have 3 exported functions that can be called from C
 */

//export IPFSStartNode
func IPFSStartNode(dataPath *C.char) *C.char {
	// Start the ZN/IPFS node
	ipfsPort := 5001
	// result is marshalled to JSON before being returned to the daemon
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node started on port %d", ipfsPort),
	}
	var err error
	basePath := C.GoString(dataPath)

	zn, err = zeronet.New(filepath.Join(basePath, "zeronet"))
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to create ZeroNet instance: %s\n", err)
		return toCJSONString(result)
	}

	ipfsNode, err = ipfs.New(filepath.Join(basePath, "ipfs"))
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to create IPFS node: %s\n", err)
		return toCJSONString(result)
	}

	err = ipfsNode.Start(ipfsPort)
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to start IPFS node: %s\n", err)
	}

	return toCJSONString(result)
}

//export ZNIPFSGetSeedList
// ZNIPFSGetSeedList retrieves the seedlist using ZeroNet and IPFS and returns
// it as JSON to the daemon. We use a named return here to ensure any
// lower level panic's (from 3rd party libs) are captured back to the daemon
func ZNIPFSGetSeedList(zeroNetAddress *C.char) (resultJSON *C.char) {
	// This defer/recover block captures any lower lever panics that might
	// occur in the 3rd party IPFS and ZeroNet libraries. It prevents the
	// daemon from crashing should such an error occur.
	defer func() {
		if r := recover(); r != nil {
			resultJSON = toCJSONString(Result{
				Status:  "err",
				Message: fmt.Sprintf("Unable to fetch seedlist from IPFS and ZeroNet: %s", r),
			})
			return
		}
	}()

	// Returns the address list from the given ZeroNet address
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("Seedlist retrieved from ZeroNet and IPFS"),
	}

	address := C.GoString(zeroNetAddress)

	// This is a well-known ZeroNet address. We store the IPFS hash in ipfs.hash
	content, err := zn.GetFile(address, "ipfs.hash")
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable fetch from ZeroNet: %s\n", err)
		resultJSON = toCJSONString(result)
		return
	}
	ipfsHash := strings.TrimSpace(string(content))

	data, err := ipfsNode.Get(ipfsHash)
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable fetch data from IPFS node: %s\n", err)
		resultJSON = toCJSONString(result)
		return
	}

	// data contains a JSON array with the seed list
	err = json.Unmarshal(data, &result.Seedlist)
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Invalid seedlist format: %s\n", err)
	}

	// If the seedlist was in the correct format is has been stored in
	// result.Seedlist and can be returned without reassigning
	resultJSON = toCJSONString(result)
	return
}

//export IPFSStopNode
func IPFSStopNode() {
	// Stop the ZN/IPFS node
	ipfsNode.Stop()
}

//export ZNStartCheckpointCollection
// ZNStartCheckpointCollection starts watching for new checkpoints
func ZNStartCheckpointCollection(
	dataPathC *C.char,
	checkpointZeroNetAddressC *C.char,
) {
	dataPath := C.GoString(dataPathC)
	checkpointZeroNetAddress := C.GoString(checkpointZeroNetAddressC)
	fmt.Printf("Start checking for checkpoints at %s (path:%s)\n", checkpointZeroNetAddress, dataPath)
	go checkpointFetchLoop(checkpointZeroNetAddress,
		dataPath,
		10,             // Hold at most 10 checkpoints
		time.Second*20, // Check for a new checkpoint every 20 seconds
	)
}

//export ZNGetCheckpointAt
// GetCheckpoint returns the checkpoint for the requested height. We only keep
// the last 10 checkpoints since start and return a blank if the checkpoint was
// not found
func ZNGetCheckpointAt(heightC C.int) *C.char {

	height := int(heightC)
	fmt.Printf("ZeroNet: Getting checkpoint at height: %d\n", height)

	checkpointMutex.RLock()
	defer checkpointMutex.RUnlock()
	if hash, ok := checkpoints[height]; ok {
		return C.CString(hash)
	}
	return C.CString("")
}

// checkpointFetchLoop will fetch the latest checkpoint from ZeroNet at a
// specified interval and store it in checkpoints for easy daemon retrieval
//
// We do this in an async manner because fetching from ZeroNet can take 2+
// seconds. Doing synchronous fetches would halt the daemon for a few seconds.
func checkpointFetchLoop(
	checkpointsZeroNetAddress string,
	baseDataPath string,
	maxCheckpoints int,
	interval time.Duration) {
	zn, err := zeronet.New(filepath.Join(baseDataPath, "zn-checkpoints"))
	if err != nil {
		fmt.Printf("ZeroNet: Unable to create ZeroNet instance: %s\n", err)
		os.Exit(1)
	}

	var height int
	var checkpointInfo []string
	var checkpointCount int
	checkpoints = make(map[int]string)

	for {
		content, err := zn.GetFile(checkpointsZeroNetAddress, "index.html")
		if err != nil {
			// If we could not fetch from ZeroNet, try again later
			// fmt.Printf("ZeroNet: Unable fetch from ZeroNet: %s\n", err)
			goto sleep
		}
		// Checkpoints are in the format `height:hash`
		checkpointInfo = strings.Split(string(content), ":")
		if len(checkpointInfo) != 2 {
			fmt.Println("ZeroNet: Checkpoint info is in an incorrect format")
			goto sleep
		}

		height, err = strconv.Atoi(checkpointInfo[0])
		if err != nil {
			fmt.Println("ZeroNet: Height from checkpoint is not a number")
			goto sleep
		}

		checkpointMutex.Lock()
		if _, ok := checkpoints[height]; ok {
			checkpointMutex.Unlock()
			// We already have this checkpoint, wait for the next one
			goto sleep
		}

		// Ass the new checkpoint
		checkpoints[height] = checkpointInfo[1]
		checkpointCount = len(checkpoints)
		// If we have too many checkpoints, remove the oldest one
		if checkpointCount > maxCheckpoints {
			// Go maps are randomised at runtime, we need to sort it and remove
			// the smaller one
			var keys []int
			for k := range checkpoints {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			delete(checkpoints, keys[0])
		}
		fmt.Printf("ZeroNet: Cached new checkpoint at height %d with hash %s. Total cached checkpoints: %d\n",
			height,
			checkpointInfo[1],
			len(checkpoints))
		checkpointMutex.Unlock()

	sleep:
		time.Sleep(interval)
	}
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
