package mock

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/rlp"
)

type QuorumCert struct {
	View          *core.View
	Code          core.MsgType
	Node          common.Hash // node hash but not block hash
	Proposer      common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

func (qc *QuorumCert) Height() *big.Int {
	if qc.View == nil {
		return common.Big0
	}
	return qc.View.Height
}

func (qc *QuorumCert) HeightU64() uint64 {
	return qc.Height().Uint64()
}

func (qc *QuorumCert) Round() *big.Int {
	if qc.View == nil {
		return common.Big0
	}
	return qc.View.Round
}

func (qc *QuorumCert) RoundU64() uint64 {
	return qc.Round().Uint64()
}

// Hash retrieve message hash but not proposal hash
func (qc *QuorumCert) SealHash() common.Hash {
	msg := core.NewCleanMessage(qc.View, qc.Code, qc.Node.Bytes())
	hash, _ := msg.Hash()
	return hash
}

func (qc *QuorumCert) Copy() *QuorumCert {
	enc, err := rlp.EncodeToBytes(qc)
	if err != nil {
		return nil
	}
	newQC := new(QuorumCert)
	if err := rlp.DecodeBytes(enc, &newQC); err != nil {
		return nil
	}
	return newQC
}
