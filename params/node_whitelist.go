package params

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var instance *nodeWhiteConfig
var mu sync.RWMutex
var once sync.Once

var (
	// ErrNodeWhitelist is returned if the node address is not in whitelist
	ErrNodeWhitelist = errors.New("### Node address is not in whitelist")
)

// node whitelist config file
type nodeWhiteConfigFile struct {
	NodeWhiteEnable bool     `json:"nodeWhiteEnable"` // only the node address in NodeWhitelist can connect to this node with p2p network when the IpWhitelistEnable is true, default false
	NodeWhitelist   []string `json:"nodeWhitelist"`   // node address whitelist
}

type nodeWhiteConfig struct {
	nodeWhiteEnable bool
	nodeWhitelist   map[string]struct{}
}

func StartNodeWhiteLoadTask(filepath string) {
	go func() {
		var initialStat fs.FileInfo

		for {
			stat, err := os.Stat(filepath)
			if err != nil {
				log.Warn("StartNodeWhiteLoadTask", "os.Stat error, filepath: "+filepath, err)
				time.Sleep(time.Second)
				continue
			}

			if initialStat == nil || stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				err = loadNodeWhiteConfig(filepath)
				initialStat = stat
				if err != nil {
					log.Warn("StartNodeWhiteLoadTask", "loadNodeWhiteConfig error, filepath: "+filepath, err)
				}
			}

			time.Sleep(time.Second)
		}
	}()
}

func loadNodeWhiteConfig(path string) error {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var res *nodeWhiteConfigFile
	if err := json.Unmarshal(config, &res); err != nil {
		return err
	}

	instance = GetNodeWhiteConfig()
	mu.Lock()
	defer mu.Unlock()
	instance.nodeWhiteEnable = res.NodeWhiteEnable
	instance.nodeWhitelist = convertToAddressMap(res.NodeWhitelist)
	return nil
}

func convertToAddressMap(list []string) map[string]struct{} {
	var result = make(map[string]struct{})
	for _, s := range list {
		// address format check
		if common.IsHexAddress(s) {
			// store address as lower case
			result[strings.ToLower(s)] = struct{}{}
		}
	}
	return result
}

func GetNodeWhiteConfig() *nodeWhiteConfig {
	once.Do(func() {
		instance = &nodeWhiteConfig{}
	})

	return instance
}

func (config *nodeWhiteConfig) CheckNodeWhitelist(addr *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if !config.nodeWhiteEnable {
		return true
	}

	return isAddressInMap(addr, config.nodeWhitelist)
}

func isAddressInMap(address *common.Address, m map[string]struct{}) bool {
	if address == nil {
		return false
	}

	return isInMap(strings.ToLower(address.String()), m)
}

func isInMap(s string, m map[string]struct{}) bool {
	if len(m) == 0 {
		return false
	}

	_, ok := m[s]
	return ok
}
