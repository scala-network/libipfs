package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/assets"
	oldcmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	log "github.com/sirupsen/logrus"
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

// NewIPFSNode constructs a new IPFS node
func NewIPFSNode(dataPath string) (*IPFS, error) {

	log.SetLevel(log.ErrorLevel)
	log.SetOutput(ioutil.Discard)
	os.MkdirAll(dataPath, 0744)

	instance := IPFS{
		keypairBits: 2048,
	}

	plugins, err := instance.loadPlugins(instance.configPath)
	if err != nil {
		panic(err)
	}

	instance.context = &oldcmds.Context{
		ConfigRoot: dataPath,
		ReqLog:     &oldcmds.ReqLog{},
		Plugins:    plugins,
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

		return nil, fmt.Errorf("Unable to create repo: %s", err)
	}
	err = repo.SetConfigKey("Discovery.MDNS.Enabled", false)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to disable MDNS: %s", err)
	}

	// Configure
	nodeConfig := &core.BuildCfg{
		Repo:                        repo,
		Permanent:                   true,
		Online:                      true,
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
	node.IsDaemon = true

	instance.cancelFunc = cancelFunc
	instance.context.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	return &instance, nil
}

// Start the IPFS node and API
func (ipfs *IPFS) Start(apiPort int) error {
	fmt.Println("Start IPFS Node")
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
		/*corehttp.CheckVersionOption(),*/
		corehttp.CommandsOption(*ipfs.context),
		/*corehttp.WebUIOption,*/
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
		defer fmt.Sprintf("\n\nThere goes the IPFS node!\n\n")
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
	fmt.Println("Get from IPFS", hash)

	// 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	node, err := ipfs.context.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("IPFS API: ConstructNode() failed: %s", err)
	}
	api, err := coreapi.NewCoreAPI(node, options.Api.FetchBlocks(false))
	if err != nil {
		return nil, err
	}

	// We only take an IPFS hash as input
	parsedPath, err := iface.ParsePath(hash)
	if err != nil {
		return nil, err
	}

	ipfsNode, err := api.Unixfs().Get(ctx, parsedPath)
	if err != nil {
		return nil, err
	}

	// We only support getting a file directly
	reader := files.ToFile(ipfsNode)
	if reader == nil {
		return nil, fmt.Errorf("IPFS Get: Hash '%s' is not a regular file", hash)
	}

	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("IPFS Get: Unable to read contents: %s", err)
	}
	defer reader.Close()

	return contents, nil
}

// Stop the IPFS node
func (ipfs *IPFS) Stop() {
	fmt.Println("Stop IPFS Node")
	//ipfs.cancelFunc()
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

	// err = nd.SetupOfflineRouting()
	// if err != nil {
	// 	return err
	// }

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

func (ipfs *IPFS) loadPlugins(repoPath string) (*loader.PluginLoader, error) {
	pluginpath := filepath.Join(repoPath, "plugins")

	// check if repo is accessible before loading plugins
	var plugins *loader.PluginLoader
	err := ipfs.isWritable(repoPath)
	if err != nil {
		pluginpath = ""
	}
	plugins, err = loader.NewPluginLoader(pluginpath)
	if err != nil {
		return nil, fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return nil, fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return nil, fmt.Errorf("error running plugins: %s", err)
	}
	return plugins, nil
}
