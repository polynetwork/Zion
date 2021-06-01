package pbft

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Validator struct {
	PublicKey        *ecdsa.PublicKey
	Address          common.Address
	VotingPower      *big.Int
	ProposerPriority int64
}

func NewValidator(pubkey *ecdsa.PublicKey) *Validator {
	v := &Validator{}
	v.PublicKey = pubkey
	v.Address = crypto.PubkeyToAddress(*pubkey)

	// todo: set voting power and proposer priority
	v.VotingPower = big.NewInt(0)
	v.ProposerPriority = -1

	return v
}

type ValSet struct {
	Validators       []*Validator
	ProposerIndex    int
	totalVotingPower *big.Int
}

func NewValSet(vals []*Validator) *ValSet {
	vs := new(ValSet)
	vs.Validators = vals
	vs.ProposerIndex = 0

	// todo: sort validators and set totalVoting power
	return vs
}
