package neo3_state_manager

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

type StateValidatorListParam struct {
	StateValidators []string       // public key strings in encoded format, each is 33 bytes in []byte
	Address         common.Address // for check witness?
}

type ApproveStateValidatorParam struct {
	ID      uint64         // StateValidatorApproveID
	Address common.Address // for check witness?
}
