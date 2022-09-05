/*
 * Copyright (C) 2022 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type MaasGenesis struct {
	Config struct {
		ChainId             uint64 `json:"chainId"`
		HomesteadBlock      uint64 `json:"homesteadBlock"`
		Eip150Block         uint64 `json:"eip150Block"`
		Eip155Block         uint64 `json:"eip155Block"`
		Eip158Block         uint64 `json:"eip158Block"`
		ByzantiumBlock      uint64 `json:"byzantiumBlock"`
		ConstantinopleBlock uint64 `json:"constantinopleBlock"`
		PetersburgBlock     uint64 `json:"petersburgBlock"`
		IstanbulBlock       uint64 `json:"istanbulBlock"`
		HotStuff            struct {
			Protocol string `json:"protocol"`
		} `json:"hotstuff"`
	} `json:"config"`
	Alloc map[string]struct {
		PublicKey string `json:"publicKey"`
		Balance   string `json:"balance"`
	} `json:"alloc"`
	Coinbase   string `json:"coinbase"`
	Difficulty string `json:"difficulty"`
	ExtraData  string `json:"extraData"`
	GasLimit   string `json:"gasLimit"`
	Nonce      string `json:"nonce"`
	Mixhash    string `json:"mixhash"`
	ParentHash string `json:"parentHash"`
	Timestamp  string `json:"timestamp"`
}

func (m *MaasGenesis) Encode() (string, error) {
	genesisJson, err := json.MarshalIndent(m, "", "\t")
	return string(genesisJson), err
}

func (m *MaasGenesis) Decode(data string) error {
	dataBytes := []byte(data)
	err := json.Unmarshal(dataBytes, m)
	return err
}

// Default genesis is not a valid block
// it's used only for initialization purpose
func (m *MaasGenesis) Default() {
	m.Config = struct {
		ChainId             uint64 `json:"chainId"`
		HomesteadBlock      uint64 `json:"homesteadBlock"`
		Eip150Block         uint64 `json:"eip150Block"`
		Eip155Block         uint64 `json:"eip155Block"`
		Eip158Block         uint64 `json:"eip158Block"`
		ByzantiumBlock      uint64 `json:"byzantiumBlock"`
		ConstantinopleBlock uint64 `json:"constantinopleBlock"`
		PetersburgBlock     uint64 `json:"petersburgBlock"`
		IstanbulBlock       uint64 `json:"istanbulBlock"`
		HotStuff            struct {
			Protocol string `json:"protocol"`
		} `json:"hotstuff"`
	}{
		10898, 0, 0, 0, 0, 0, 0, 0, 0,
		struct {
			Protocol string `json:"protocol"`
		}{
			"basic",
		},
	}

	m.Coinbase = "0x0000000000000000000000000000000000000000"
	m.Difficulty = "0x1"
	//"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f8daf893940f45dafcd39e1c59202e91306a1671cb7f8884be944ba9732e2358f41e682b7a7aab71e614e08383df946519d82d761c275f01de6f197542de924296928e946cfd8d31a55cdf03303e280792c7d6ce855601f394782833979973d83cd48332e09938b95f9ba32b5094990f7aafa09fea4583c0c72063b306cde54e1e8f949bcd01b46b98254ee3611eb1501a12780343f7d2b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
	m.GasLimit = "0xffffffff"
	m.Nonce = "0x4510809143055965"
	m.Mixhash = "0x0000000000000000000000000000000000000000000000000000000000000000"
	m.ParentHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
	m.Timestamp = "0x00"
}

func DumpGenesis(filePaths [3]string, contents [3]string) error {
	for k, file_ := range filePaths {
		err := ensureBaseDir(file_)
		if err != nil {
			return err
		}
		f, err := os.Create(file_)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write([]byte(contents[k]))
		if err != nil {
			return err
		}
	}
	return nil
}

// If there is no such directory, it will creates this dir
func ensureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

func DefaultBasePath() (basePath string) {
	exePath, _ := os.Executable()
	dir := filepath.Dir(exePath) + "/"
	return dir
}

type Node struct {
	Address  string      `json:"address"`
	NodeKey  string      `json:"nodeKey"`
	KeyStore interface{} `json:"keystore"`
	PubKey   string      `json:"pubKey"`
	Static   string      `json:"static"`
}

func (m *Node) Encode() (string, error) {
	genesisJson, err := json.MarshalIndent(m, "", "\t")
	return string(genesisJson), err
}

func (m *Node) Decode(data string) error {
	dataBytes := []byte(data)
	err := json.Unmarshal(dataBytes, m)
	return err
}

// This is only used in unit test
func DeleteBasePath(basePath string) {
	dir, _ := ioutil.ReadDir(basePath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{basePath, d.Name()}...))
	}
	os.Remove(basePath)
}
