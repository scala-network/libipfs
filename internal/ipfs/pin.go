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
	pinCh := make(chan coreiface.Pin, 1)
	var pinned []string

	errCh := make(chan error, 1)
	go func() {
		errCh <- api.Pin().Ls(node.Context(), pinCh)
	}()

	for pin := range pinCh {
		pinned = append(pinned, pin.Path().String())
	}

	if err := <-errCh; err != nil {
		return nil, err
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
