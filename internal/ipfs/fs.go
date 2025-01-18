package ipfs

import (
	"os"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/core"
	coreiface "github.com/ipfs/kubo/core/coreiface"
)

func GetUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func Add(node *core.IpfsNode, api coreiface.CoreAPI, path string) (string, error) {
	file, err := GetUnixfsNode(path)

	if err != nil {
		return "", err
	}

	cid, err := api.Unixfs().Add(node.Context(), file)

	if err != nil {
		return "", err
	}

	return cid.String(), nil
}

func Get(node *core.IpfsNode, api coreiface.CoreAPI, fullPath string, destination string, pin bool) (string, error) {
	p, err := path.NewPath(fullPath)

	if err != nil {
		return "", err
	}

	rootNodeFile, err := api.Unixfs().Get(node.Context(), p)
	if err != nil {
		return "", err
	}

	if files.WriteTo(rootNodeFile, destination) != nil {
		return "", err
	}

	if pin {
		return destination, Pin(node, api, fullPath)
	}

	return destination, nil
}
