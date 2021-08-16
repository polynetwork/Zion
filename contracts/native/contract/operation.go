package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
)

func ValidateOwner(n *native.NativeContract, address common.Address) error {
	if n.ContractRef().CheckWitness(address) == false {
		return fmt.Errorf("validateOwner, authentication failed!")
	}
	return nil
}
