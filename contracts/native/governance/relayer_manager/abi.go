package relayer_manager

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

type RelayerListParam struct {
	AddressList []common.Address
	Address     common.Address
}

type ApproveRelayerParam struct {
	ID      uint64
	Address common.Address
}
