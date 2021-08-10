package side_chain_manager

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

type RegisterSideChainParam struct {
	Address      common.Address
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	CCMCAddress  []byte
	ExtraInfo    []byte
}

type ChainidParam struct {
	Chainid uint64
	Address common.Address
}

type RegisterRedeemParam struct {
	RedeemChainID   uint64
	ContractChainID uint64
	Redeem          []byte
	CVersion        uint64
	ContractAddress []byte
	Signs           [][]byte
}

type BtcTxParam struct {
	Redeem        []byte
	RedeemChainId uint64
	Sigs          [][]byte
	Detial        *BtcTxParamDetial
}

type BtcTxParamDetial struct {
	PVersion  uint64
	FeeRate   uint64
	MinChange uint64
}
