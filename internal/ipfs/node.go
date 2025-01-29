package ipfs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/scala-network/libipfs/internal"
	"github.com/scala-network/libipfs/internal/utils"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/corerepo"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
)

func CreateRepo(repoPath string, bindPort int, bindIps []string) error {
	if !utils.IsValidPort(bindPort) {
		return fmt.Errorf("invalid port")
	}

	plugins, err := loader.NewPluginLoader(filepath.Join(repoPath, "plugins"))
	if err != nil {
		return err
	}

	if err := plugins.Initialize(); err != nil {
		return err
	}

	if err := plugins.Inject(); err != nil {
		return err
	}

	if !(utils.IsDir(repoPath)) {
		err := os.MkdirAll(repoPath, 0755)

		if err != nil {
			return err
		}
	} else {
		_, err := os.Stat(filepath.Join(repoPath, "config"))

		if err == nil {
			return nil
		}
	}

	cfg, err := config.Init(io.Discard, 2048)

	if err != nil {
		return err
	}

	var swarmAddresses []string

	if len(bindIps) == 0 {
		swarmAddresses = []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", bindPort),
			fmt.Sprintf("/ip6/::/tcp/%d", bindPort),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", bindPort),
			fmt.Sprintf("/ip6/::/udp/%d/quic", bindPort),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1/webtransport", bindPort),
			fmt.Sprintf("/ip6/::/udp/%d/quic-v1/webtransport", bindPort),
		}
	} else {
		for _, ip := range bindIps {
			if !utils.CheckBind(ip, bindPort) {
				return fmt.Errorf("failed to bind to %s:%d", ip, bindPort)
			}

			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip4/%s/tcp/%d", ip, bindPort))
			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip6/%s/tcp/%d", ip, bindPort))
			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip4/%s/udp/%d/quic", ip, bindPort))
			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip6/%s/udp/%d/quic", ip, bindPort))
			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip4/%s/udp/%d/quic-v1/webtransport", ip, bindPort))
			swarmAddresses = append(swarmAddresses, fmt.Sprintf("/ip6/%s/udp/%d/quic-v1/webtransport", ip, bindPort))
		}
	}

	cfg.Addresses.Swarm = swarmAddresses
	cfg.Addresses.NoAnnounce = internal.DefaultServerFilters
	cfg.Swarm.AddrFilters = internal.DefaultServerFilters
	cfg.Discovery.MDNS.Enabled = false
	cfg.Swarm.DisableNatPortMap = true

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return err
	}

	return nil
}

func CreateNode(ctx context.Context, repoPath string) (*core.IpfsNode, coreiface.CoreAPI, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    r,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, nil, err
	}

	api, err := coreapi.NewCoreAPI(node)

	if err != nil {
		return nil, nil, err
	}

	return node, api, nil
}

func GarbageCollect(node *core.IpfsNode) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return corerepo.GarbageCollect(node, ctx)
}

func CloseNode(node *core.IpfsNode) error {
	return node.Close()
}
