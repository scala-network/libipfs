package ipfs

import (
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/core"
	coreiface "github.com/ipfs/kubo/core/coreiface"
)

func Pin(node *core.IpfsNode, api coreiface.CoreAPI, toPin string) error {
	p, err := path.NewPath(toPin)

	if err != nil {
		return err
	}

	return api.Pin().Add(node.Context(), p)
}

func ListPinned(node *core.IpfsNode, api coreiface.CoreAPI) ([]string, error) {
	pins, err := api.Pin().Ls(node.Context())

	if err != nil {
		return nil, err
	}

	var pinned []string

	for pin := range pins {
		pinned = append(pinned, pin.Path().String())
	}

	return pinned, nil
}

func Unpin(node *core.IpfsNode, api coreiface.CoreAPI, toRemove string) error {
	p, err := path.NewPath(toRemove)

	if err != nil {
		return err
	}

	return api.Pin().Rm(node.Context(), p)
}
