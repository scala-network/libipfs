package main

import (
	"C"
	"context"
	"encoding/json"
	"fmt"
)

import (
	"github.com/scala-network/libipfs/internal/ipfs"

	"github.com/ipfs/kubo/core"
	iface "github.com/ipfs/kubo/core/coreiface"
)

type Result struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

var (
	ctx        context.Context
	cancelFunc context.CancelFunc
	node       *core.IpfsNode
	api        iface.CoreAPI
)

func init() {
	ctx, cancelFunc = context.WithCancel(context.Background())
}

func main() {
	// The main function is required to make the CGO compiler
	// compile the package as a C library.
}

//export Start
func Start(dataPath *C.char, P2PPort int) *C.char {
	err := ipfs.CreateRepo(C.GoString(dataPath), P2PPort, []string{})
	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to create repo: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	node, api, err = ipfs.CreateNode(ctx, C.GoString(dataPath))
	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to create node: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	peerId := ipfs.PeerID(node)

	result := Result{
		Status: "success",
		Data:   map[string]string{"peerId": peerId},
	}

	return toCJSONString(result)
}

//export Stop
func Stop() *C.char {
	err := ipfs.CloseNode(node)

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to close node: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export GarbageCollect
func GarbageCollect() *C.char {
	err := ipfs.GarbageCollect(node)

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to garbage collect: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export Add
func Add(path *C.char) *C.char {
	cid, err := ipfs.Add(node, api, C.GoString(path))

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to add file: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   map[string]string{"cid": cid},
	}

	return toCJSONString(result)
}

//export Pin
func Pin(cid *C.char) *C.char {
	err := ipfs.Pin(node, api, C.GoString(cid))

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to pin file: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export Unpin
func Unpin(cid *C.char) *C.char {
	err := ipfs.Unpin(node, api, C.GoString(cid))

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to unpin file: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export Get
func Get(fullPath *C.char, destination *C.char, pin bool) *C.char {
	_, err := ipfs.Get(node, api, C.GoString(fullPath), C.GoString(destination), pin)

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to get file/directory: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export PeerID
func PeerID() *C.char {
	peerId := ipfs.PeerID(node)

	result := Result{
		Status: "success",
		Data:   map[string]string{"peerId": peerId},
	}

	return toCJSONString(result)
}

//export AddPeer
func AddPeer(peerAddr *C.char) *C.char {
	err := ipfs.AddPeer(node, api, C.GoString(peerAddr))

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to add peer: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

//export RemovePeer
func RemovePeer(peerAddr *C.char) *C.char {
	err := ipfs.RemovePeer(node, api, C.GoString(peerAddr))

	if err != nil {
		result := Result{
			Status: "error",
			Data:   map[string]string{"error": fmt.Sprintf("Failed to remove peer: %s", err.Error())},
		}
		return toCJSONString(result)
	}

	result := Result{
		Status: "success",
		Data:   nil,
	}

	return toCJSONString(result)
}

func toCJSONString(result Result) *C.char {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("fatal error converting result: %s", err.Error()))
	}
	return C.CString(string(resultJSON))
}
