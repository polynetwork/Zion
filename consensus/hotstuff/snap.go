package hotstuff

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
)

type Snapshot struct {
	ValSet ValidatorSet
}

func newSnapshot(vals ValidatorSet) *Snapshot {
	return &Snapshot{ValSet: vals}
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}