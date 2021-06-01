package hotstuff

import "github.com/ethereum/go-ethereum/common"

type Validators interface {

	Add(val common.Address) error

	Del(val common.Address) error

	Find(val common.Address) bool

	GetLeader(view uint64) common.Address

	IsProposer(val common.Address) bool

	Get() []common.Address

	// Major retrieve the number of 2f + 1
	Major() bool
}
