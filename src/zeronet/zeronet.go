package zeronet

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/donovansolms/ZeroGo/site_manager"
	"github.com/donovansolms/ZeroGo/utils"

	log "github.com/Sirupsen/logrus"
)

// ZeroNet contains the functionality to communicate with ZeroNet
type ZeroNet struct {
	mutex       sync.Mutex
	dataPath    string
	siteManager *site_manager.SiteManager
}

// New constructs and initializes ZeroNet and returns the instance
func New(dataPath string) (*ZeroNet, error) {

	// We discard log output to avoid the underlying libraries to print
	// stuff out in the daemon
	log.SetLevel(log.ErrorLevel)
	log.SetOutput(ioutil.Discard)

	// Set up ZeroNet
	utils.SetDataPath(dataPath)
	siteManager := site_manager.NewSiteManager()

	// Create instance
	zn := ZeroNet{
		dataPath:    dataPath,
		siteManager: siteManager,
	}

	return &zn, nil
}

// GetFile retrieves the given filename from the ZeroNet address
// and returns the contents as bytes
func (zn *ZeroNet) GetFile(address string, filename string) ([]byte, error) {
	// Create the ZeroNet data path and add certificates
	os.MkdirAll(utils.GetDataPath(), 0744)
	utils.CreateCerts()

	// NOTE: We remove the site for the site to sync correctly. The ZeroGo lib
	// doesn't yet implement the ZeroNet update protocol, once implemented, the
	// site will update itselt
	zn.siteManager.Remove(address)

	// Fetch the site from ZeroNet and wait for the operation to complete
	site := zn.siteManager.Get(address)
	done := site.Wait()
	if !done {
		return []byte{}, errors.New("Timeout!!!")
	}

	// Get the specified file from the ZeroNet address
	data, err := site.GetFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
