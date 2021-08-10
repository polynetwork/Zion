package cross_chain_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const abijson = ``

func GetABI() abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return ab
}

type EntranceParam struct {
	SourceChainID         uint64 `json:"sourceChainId"`
	Height                uint32 `json:"height"`
	Proof                 []byte `json:"proof"`
	RelayerAddress        []byte `json:"relayerAddress"`
	Extra                 []byte `json:"extra"`
	HeaderOrCrossChainMsg []byte `json:"headerOrCrossChainMsg"`
}

type MultiSignParam struct {
	ChainID   uint64
	RedeemKey string
	TxHash    []byte
	Address   string
	Signs     [][]byte
}

type BlackChainParam struct {
	ChainID uint64
}
