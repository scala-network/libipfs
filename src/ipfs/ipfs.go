package ipfs


import (
	"fmt"
	"os"
	"context"
	"path/filepath"
	"io/ioutil"
	"sync"

	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/scala-network/libipfs/src/constants"
	"github.com/scala-network/libipfs/src/utils"
	"github.com/libp2p/go-libp2p-core/peer"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	config "github.com/ipfs/go-ipfs-config"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	ma "github.com/multiformats/go-multiaddr"
)

var ipfsCoreAll *core.IpfsNode
var ipfsApiAll icore.CoreAPI
var ctxAll context.Context

func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) (error) {
	var wg sync.WaitGroup
	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
		}(peerInfo)
	}
	wg.Wait()
	return nil
}

func createRepo(ctx context.Context, dataPath string, P2PPort int) (string, error) {
	repoPath := dataPath
	
	if !(utils.IsDir(repoPath)) {
		err := os.MkdirAll(repoPath, 0755)

		if err != nil {
			return "", fmt.Errorf("Failed to get repo directory: %s", err)
		}
	} else {
		return repoPath, nil
	}

	cfg, err := config.Init(ioutil.Discard, 2048)

	swarmAddresses := []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", P2PPort),
		fmt.Sprintf("/ip6/::/tcp/%d", P2PPort),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", P2PPort),
		fmt.Sprintf("/ip6/::/udp/%d/quic", P2PPort),
	}

	cfg.Addresses.Swarm = swarmAddresses
	cfg.Addresses.NoAnnounce = constants.DefaultServerFilters
	cfg.Swarm.AddrFilters = constants.DefaultServerFilters
	cfg.Discovery.MDNS.Enabled = false
	cfg.Swarm.DisableNatPortMap = true

	if err != nil {
		return "", err
	}

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("Failed to init IPFS node: %s", err)
	}

	return repoPath, nil
}

func setupPlugins(externalPluginsPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}


func createNode(ctx context.Context, repoPath string) (*core.IpfsNode, icore.CoreAPI, error) {
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, nil, err
	}

	coreapi, _ := (coreapi.NewCoreAPI(node))

	return node, coreapi, nil
}


func spawnIpfsNode(ctx context.Context, dataPath string, P2PPort int) (*core.IpfsNode, icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, nil, err
	}

	repoPath, err := createRepo(ctx, dataPath, P2PPort)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	return createNode(ctx, repoPath)
}


func getUnixfsNode(path string) (files.Node, error) {
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


func Start(dataPath string, P2PPort int) (error) {
	var err error
	ctxAll = context.Background()
	ipfsCoreAll, ipfsApiAll, err = spawnIpfsNode(ctxAll, dataPath, P2PPort)
	err = connectToPeers(ctxAll, ipfsApiAll, constants.BootstrapNodes)

	if err != nil {
		return err
	}
	
	return nil
}

func Stop() (error) {
	err := ipfsCoreAll.Close()

	if err != nil {
		return err
	}
		
	return nil
}

func ResolveName(ipnsPath string) (string, error) {
	opts := []options.NameResolveOption{}
	resolvedName, err := ipfsApiAll.Name().Resolve(ctxAll, ipnsPath, opts...)

	if err != nil {
		return "", err
	} 

	return resolvedName.String(), nil
}

func PublishName(ipfsHash string) (string, error) {
	opts := []options.NamePublishOption{}

	pCid := icorepath.New(ipfsHash)
	ipnsEntry, err := ipfsApiAll.Name().Publish(ctxAll, pCid, opts...)
	
	if err != nil {
		return "", err
	}

	return ipnsEntry.Name(), nil
}

func GetPeerID() (string) {
	peerId := ipfsCoreAll.Identity
	return peerId.String()
}

func Add(AddPath string) (ipfsHash string, err error) {
	tempFile, err := getUnixfsNode(AddPath)

	if err != nil {
		return "", err
	}

	addedCid, err := ipfsApiAll.Unixfs().Add(ctxAll, tempFile)

	if err != nil {
		return "", nil
	}

	return addedCid.String(), nil
}

func Get(hash string, downloadPath string) (err error) {
	getCid := icorepath.New(hash)
	rootNode, err := ipfsApiAll.Unixfs().Get(ctxAll, getCid)

	if err != nil {
		return err
	}

	err = files.WriteTo(rootNode, downloadPath)

	if (err != nil) {
		return err
	}

	return nil
}