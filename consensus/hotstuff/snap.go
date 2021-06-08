package hotstuff

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
)

// todo: persist validators and candidates votes as snapshot, provider interface for query and verify
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

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		ValSet: s.ValSet.Copy(),
	}

	return cpy
}

func (s *Snapshot) major() int {
	return 2 * s.ValSet.F() + 1
}
