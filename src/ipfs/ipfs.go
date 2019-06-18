package ipfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// IPFS implements communication over IPFS.
//
// We package the official IPFS daemon release for each platform into
// libznipfs. This ensures the daemon operates correctly and has the
// added benefit of being easy to maintain
type IPFS struct {
	daemonPath string
}

// New constructs a new IPFS node
func New(dataPath string) (*IPFS, error) {

	log.SetLevel(log.ErrorLevel)
	log.SetOutput(ioutil.Discard)
	os.MkdirAll(dataPath, 0744)

	binaryName := "ipfs"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fileBytes, err := FSByte(false,
		fmt.Sprintf("/pack/%s/%s", runtime.GOOS, binaryName))
	if err != nil {
		panic(err)
	}

	daemonPath := filepath.Join(dataPath, binaryName)

	outFile, err := os.OpenFile(
		daemonPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0755)
	if err != nil {
		return nil, err
	}

	_, err = outFile.Write(fileBytes)
	if err != nil {
		return nil, err
	}
	outFile.Close()

	instance := IPFS{
		daemonPath: daemonPath,
	}

	return &instance, nil
}

// Start the IPFS node and API
func (ipfs *IPFS) Start(apiPort int) error {

	// TODO: Need to run ipfs init at least once before daemon will work

	cmd := exec.Command(ipfs.daemonPath, "daemon")
	op, err := cmd.CombinedOutput()
	fmt.Println(string(op))
	if err != nil {
		return err
	}

	fmt.Println(string(op))

	// apiAddr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", apiPort)
	// apiMaddr, err := ma.NewMultiaddr(apiAddr)
	// if err != nil {
	// 	return fmt.Errorf("IPFS API: invalid API address: %q (err: %s)", apiAddr, err)
	// }
	//
	// apiLis, err := manet.Listen(apiMaddr)
	// if err != nil {
	// 	return fmt.Errorf("IPFS API: manet.Listen(%s) failed: %s", apiMaddr, err)
	// }
	// // we might have listened to /tcp/0 - lets see what we are listing on
	// apiMaddr = apiLis.Multiaddr()
	//
	// var opts = []corehttp.ServeOption{
	// 	corehttp.CheckVersionOption(),
	// 	corehttp.CommandsOption(*ipfs.context),
	// 	corehttp.WebUIOption,
	// 	corehttp.VersionOption(),
	// 	corehttp.LogOption(),
	// }
	//
	// node, err := ipfs.context.ConstructNode()
	// if err != nil {
	// 	return fmt.Errorf("IPFS API: ConstructNode() failed: %s", err)
	// }
	//
	// if err := node.Repo.SetAPIAddr(apiMaddr); err != nil {
	// 	return fmt.Errorf("IPFS API: SetAPIAddr() failed: %s", err)
	// }
	//
	// go func() {
	// 	defer fmt.Sprintf("\n\nThere goes the IPFS node!\n\n")
	// 	err := corehttp.Serve(node, manet.NetListener(apiLis), opts...)
	// 	// TODO: Find a better way to pass errors back
	// 	if err != nil {
	// 		fmt.Printf("An API error occurred: %s\n", err)
	// 		panic(err)
	// 	}
	// }()
	return nil
}

// Get an object from IPFS and return it as bytes
func (ipfs *IPFS) Get(hash string) ([]byte, error) {

	// node, err := ipfs.context.ConstructNode()
	// if err != nil {
	// 	return nil, fmt.Errorf("IPFS API: ConstructNode() failed: %s", err)
	// }
	//
	// dagResolver := &resolver.Resolver{
	// 	DAG:         node.DAG,
	// 	ResolveOnce: uio.ResolveUnixfsOnce,
	// }
	//
	// dagNode, err := core.Resolve(
	// 	ipfs.context.Context(),
	// 	node.Namesys,
	// 	dagResolver,
	// 	ipfspath.Path(fmt.Sprintf("/ipfs/%s", hash)))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// reader, err := uio.NewDagReader(ipfs.context.Context(), dagNode, node.DAG)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// return ioutil.ReadAll(reader)
	return nil, nil
}

// Stop the IPFS node
func (ipfs *IPFS) Stop() {
	//ipfs.cancelFunc()
}
