package pbft

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

type VoteType uint8

const (
	VoteTypePreVote   = 1
	VoteTypePreCommit = 2
)

type Vote struct {
	Hrs       *HRS
	Validator common.Address
	PublicKey *ecdsa.PublicKey
	MsgType   VoteType
}

type VoteSet struct {
	height uint64

	prevotes map[int64][]*Vote // mapping round to pre-vote vote slice
	precommits map[int64][]*Vote // mapping round to pre-commit vote slice

	mtx sync.RWMutex
}

// todo: do not allow dumplicate votes
func (vs *VoteSet) Add(v *Vote) {

}

// todo:
func (vs *VoteSet) HasMajor(round int64, typ VoteType) bool {
	return false
}
