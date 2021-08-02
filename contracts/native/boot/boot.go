package boot

import (
	"github.com/ethereum/go-ethereum/contracts/native/governance"
)

func InitialNativeContracts() {
	governance.InitGovernance()
}
