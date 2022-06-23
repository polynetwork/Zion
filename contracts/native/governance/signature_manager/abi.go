package signature_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/signature_manager_abi"
)

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(signature_manager_abi.SignatureManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}
