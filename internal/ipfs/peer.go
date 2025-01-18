package ipfs

import (
	"fmt"

	"github.com/ipfs/kubo/core"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func PeerID(node *core.IpfsNode) string {
	return node.Identity.String()
}

func AddPeer(node *core.IpfsNode, api coreiface.CoreAPI, peerAddr string) error {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("invalid multiaddress: %v", err)
	}

	addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("failed to create AddrInfo from multiaddress: %v", err)
	}

	err = api.Swarm().Connect(node.Context(), *addrInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to peer: %v", err)
	}

	return nil
}

func RemovePeer(node *core.IpfsNode, api coreiface.CoreAPI, peerAddr string) error {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("invalid multiaddress: %v", err)
	}

	err = api.Swarm().Disconnect(node.Context(), maddr)
	if err != nil {
		return fmt.Errorf("failed to disconnect from peer: %v", err)
	}

	return nil
}
