package ipfs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/ipfs/go-ipfs/assets"
	oldcmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/path/resolver"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmV6FjemM1K8oXjrvuq3wuVWWoU2TLDPmNnKrxHzY3v6Ai/go-multiaddr-net"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"

	ipfspath "github.com/ipfs/go-ipfs/path"
	uio "github.com/ipfs/go-ipfs/unixfs/io"

	log "github.com/Sirupsen/logrus"
)

// IPFS implements communication over IPFS. Most of the code for the
// node has been taken from the official core IPFS daemon command
// https://github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/daemon.go
type IPFS struct {
	keypairBits int
	configPath  string
	context     *oldcmds.Context
	cancelFunc  context.CancelFunc
}

// New constructs a new IPFS node
func New(dataPath string) (*IPFS, error) {

	log.SetLevel(log.ErrorLevel)
	log.SetOutput(ioutil.Discard)
	os.MkdirAll(dataPath, 0744)

	// TODO: Check if the descriptor limit need to be increases
	// if err := utilmain.ManageFdLimit(); err != nil {
	// 	fmt.Printf("Unable to set file limits: %s", err)
	// }

	instance := IPFS{
		keypairBits: 2048,
	}

	instance.context = &oldcmds.Context{
		ConfigRoot: dataPath,
		ReqLog:     &oldcmds.ReqLog{},
	}
	instance.configPath = instance.context.ConfigRoot

	// Initialize the IPFS node's defaults
	if !fsrepo.IsInitialized(instance.configPath) {
		err := instance.initNode(ioutil.Discard, instance.configPath)
		if err != nil {
			panic(err)
		}
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(instance.context.ConfigRoot)
	if err != nil {
		if err == fsrepo.ErrNeedMigration {
			return nil, fmt.Errorf(
				"IPFS repo needs to be migrated. Either delete '%s' or run IPFS migrations",
				instance.context.ConfigRoot)
		}
	}

	// Configure
	nodeConfig := &core.BuildCfg{
		Repo:      repo,
		Permanent: true,
		Online:    true,
		DisableEncryptedConnections: false,
		ExtraOpts: map[string]bool{
			"pubsub": false,
			"ipnsps": false,
			"mplex":  false,
		},
		// DHT is the default, non-experimental routing option
		Routing: core.DHTOption,
	}

	cancelContext, cancelFunc := context.WithCancel(context.Background())
	node, err := core.NewNode(cancelContext, nodeConfig)
	if err != nil {
		return nil, err
	}
	node.SetLocal(false)

	instance.cancelFunc = cancelFunc
	instance.context.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	return &instance, nil
}

// Start the IPFS node and API
func (ipfs *IPFS) Start(apiPort int) error {
	apiAddr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", apiPort)
	apiMaddr, err := ma.NewMultiaddr(apiAddr)
	if err != nil {
		return fmt.Errorf("IPFS API: invalid API address: %q (err: %s)", apiAddr, err)
	}

	apiLis, err := manet.Listen(apiMaddr)
	if err != nil {
		return fmt.Errorf("IPFS API: manet.Listen(%s) failed: %s", apiMaddr, err)
	}
	// we might have listened to /tcp/0 - lets see what we are listing on
	apiMaddr = apiLis.Multiaddr()

	var opts = []corehttp.ServeOption{
		corehttp.CheckVersionOption(),
		corehttp.CommandsOption(*ipfs.context),
		corehttp.WebUIOption,
		corehttp.VersionOption(),
		corehttp.LogOption(),
	}

	node, err := ipfs.context.ConstructNode()
	if err != nil {
		return fmt.Errorf("IPFS API: ConstructNode() failed: %s", err)
	}

	if err := node.Repo.SetAPIAddr(apiMaddr); err != nil {
		return fmt.Errorf("IPFS API: SetAPIAddr() failed: %s", err)
	}

	go func() {
		err := corehttp.Serve(node, manet.NetListener(apiLis), opts...)
		// TODO: Find a better way to pass errors back
		if err != nil {
			fmt.Printf("An API error occurred: %s\n", err)
			panic(err)
		}
	}()
	return nil
}

// Get an object from IPFS and return it as bytes
func (ipfs *IPFS) Get(hash string) ([]byte, error) {

	node, err := ipfs.context.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("IPFS API: ConstructNode() failed: %s", err)
	}

	dagResolver := &resolver.Resolver{
		DAG:         node.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	dagNode, err := core.Resolve(
		ipfs.context.Context(),
		node.Namesys,
		dagResolver,
		ipfspath.Path(fmt.Sprintf("/ipfs/%s", hash)))
	if err != nil {
		return nil, err
	}

	reader, err := uio.NewDagReader(ipfs.context.Context(), dagNode, node.DAG)
	if err != nil {
		panic(err)
	}

	return ioutil.ReadAll(reader)
}

// Stop the IPFS node
func (ipfs *IPFS) Stop() {
	ipfs.cancelFunc()
}

// initNode sets up the IPFS node with default settings
func (ipfs *IPFS) initNode(out io.Writer, repoRoot string) error {

	if err := ipfs.isWritable(repoRoot); err != nil {
		return err
	}

	conf, err := config.Init(out, ipfs.keypairBits)
	if err != nil {
		return err
	}
	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}
	if err := addDefaultAssets(out, repoRoot); err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoRoot)
}

func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	err = nd.SetupOfflineRouting()
	if err != nil {
		return err
	}

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}

// isWritable checks if the given path is writable
func (ipfs *IPFS) isWritable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		var fi *os.File
		fi, err = os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesn't exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}
	return err
}

// addDefaultAssets adds the default IPFS data to the node for hosting
func addDefaultAssets(out io.Writer, repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	_, err = assets.SeedInitDocs(nd)
	if err != nil {
		return fmt.Errorf("init: seeding init docs failed: %s", err)
	}

	return err
}
