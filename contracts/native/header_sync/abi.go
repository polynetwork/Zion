package header_sync

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

type SyncGenesisHeaderParam struct {
	ChainID       uint64
	GenesisHeader []byte
}

type SyncBlockHeaderParam struct {
	ChainID uint64
	Address common.Address
	Headers [][]byte
}

type SyncCrossChainMsgParam struct {
	ChainID        uint64
	Address        common.Address
	CrossChainMsgs [][]byte
}
