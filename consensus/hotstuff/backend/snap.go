package backend

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

// todo: use snap or reconfig validators group
func (s *backend) snap() hotstuff.ValidatorSet {
	return s.valset.Copy()
}
