package node_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const abijson = ``

func GetABI() abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return ab
}

type VBFTConfig struct {
	BlockMsgDelay        uint32          `json:"block_msg_delay"`
	HashMsgDelay         uint32          `json:"hash_msg_delay"`
	PeerHandshakeTimeout uint32          `json:"peer_handshake_timeout"`
	MaxBlockChangeView   uint32          `json:"max_block_change_view"`
	VrfValue             string          `json:"vrf_value"`
	VrfProof             string          `json:"vrf_proof"`
	Peers                []*VBFTPeerInfo `json:"peers"`
}

type VBFTPeerInfo struct {
	Index      uint32 `json:"index"`
	PeerPubkey string `json:"peerPubkey"`
	Address    string `json:"address"`
}

type RegisterPeerParam struct {
	PeerPubkey string
	Address    common.Address
}

type PeerParam struct {
	PeerPubkey string
	Address    common.Address
}

type PeerListParam struct {
	PeerPubkeyList []string
	Address        common.Address
}

type UpdateConfigParam struct {
	Configuration *Configuration
}

type Configuration struct {
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	MaxBlockChangeView   uint32
}
