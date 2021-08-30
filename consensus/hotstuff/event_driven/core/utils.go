/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func (e *EventDrivenEngine) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return e.signer.CheckSignature(e.valset, data, sig)
}

// todo
func (e *EventDrivenEngine) newLogger() log.Logger {
	logger := e.logger.New("state")
	return logger
}

func (e *EventDrivenEngine) address() common.Address {
	return e.addr
}

func (e *EventDrivenEngine) isSelf(addr common.Address) bool {
	return e.addr == addr
}

func (e *EventDrivenEngine) currentView() *hotstuff.View {
	return &hotstuff.View{
		Round:  new(big.Int).Set(e.curRound),
		Height: new(big.Int).Set(e.curHeight),
	}
}

func (e *EventDrivenEngine) checkProposer(proposer common.Address) error {
	if !e.valset.IsProposer(proposer) {
		return errNotFromProposer
	}
	return nil
}

func (e *EventDrivenEngine) checkEpoch(epoch uint64, height *big.Int) error {
	if e.epoch != epoch {
		return errInvalidHighQC
	}
	if height.Cmp(e.epochHeightStart) < 0 {
		return errInvalidEpoch
	}
	if height.Cmp(e.epochHeightEnd) > 0 {
		return errInvalidEpoch
	}
	return nil
}

func (e *EventDrivenEngine) checkView(view *hotstuff.View) error {
	if e.curRound.Cmp(view.Round) != 0 || e.curHeight.Cmp(view.Height) != 0 {
		return errInvalidMessage
	}
	return nil
}

func (e *EventDrivenEngine) checkJustifyQC(proposal hotstuff.Proposal, justifyQC *hotstuff.QuorumCert) error {
	if justifyQC == nil || justifyQC.View == nil || justifyQC.Hash == utils.EmptyHash || justifyQC.Proposer == utils.EmptyAddress {
		return fmt.Errorf("justifyQC fields may be empty or nil")
	}

	if justifyQC.View.Height.Cmp(new(big.Int).Sub(proposal.Number(), common.Big1)) != 0 {
		return fmt.Errorf("high qc height invalid")
	}

	if justifyQC.Hash != proposal.ParentHash() {
		return fmt.Errorf("justifyQC hash invalid")
	}

	vs := e.valset.Copy()
	vs.CalcProposerByIndex(justifyQC.View.Round.Uint64())
	proposer := vs.GetProposer().Address()
	if proposer != justifyQC.Proposer {
		return fmt.Errorf("invalid proposer")
	}

	highQC := e.blkPool.GetHighQC()
	return e.compareQC(highQC, justifyQC)
}

func (e *EventDrivenEngine) compareQC(expect, src *hotstuff.QuorumCert) error {
	if expect.Hash != src.Hash {
		return fmt.Errorf("qc hash expect %v, got %v", expect.Hash, src.Hash)
	}
	if expect.View.Cmp(src.View) != 0 {
		return fmt.Errorf("qc view expect %v, got %v", expect.View, src.View)
	}
	if expect.Proposer != src.Proposer {
		return fmt.Errorf("qc proposer expect %v, got %v", expect.Proposer, src.Proposer)
	}
	if !bytes.Equal(expect.Extra, src.Extra) {
		return fmt.Errorf("qc extra not same")
	}
	return nil
	// if !reflect.DeepEqual(expect, src) {
	//     return fmt.Errorf("qc not same")
	// }
}

// vote to highQC round + 1
func (e *EventDrivenEngine) checkVote(vote *Vote) error {
	if vote.View == nil || vote.Hash == utils.EmptyHash {
		return errInvalidVote
	}
	if vote.ParentHash == utils.EmptyHash || vote.ParentView == nil {
		return errInvalidVote
	}

	// vote view MUST be highQC view
	highQC := e.blkPool.GetHighQC()
	if new(big.Int).Sub(vote.View.Height, highQC.View.Height).Cmp(common.Big1) != 0 &&
		new(big.Int).Sub(vote.View.Round, highQC.View.Round).Cmp(common.Big1) != 0 {
		return errInvalidVote
	}
	return nil
}

func (e *EventDrivenEngine) getVoteSeals(hash common.Hash, n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range e.messages.Votes(hash) {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

func (e *EventDrivenEngine) getTimeoutSeals(round uint64, n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range e.messages.Timeouts(round) {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

func (e *EventDrivenEngine) Q() int {
	return e.valset.Q()
}

func (e *EventDrivenEngine) chain2Height() *big.Int {
	return new(big.Int).Add(e.epochHeightStart, common.Big2)
}

func (e *EventDrivenEngine) chain3Height() *big.Int {
	return new(big.Int).Add(e.epochHeightStart, common.Big3)
}

func (e *EventDrivenEngine) generateTimeoutEvent() *TimeoutEvent {
	return &TimeoutEvent{
		Epoch: e.epoch,
		View:  e.currentView(),
	}
}

func (e *EventDrivenEngine) aggregateQC(vote *Vote, size int) (*hotstuff.QuorumCert, error) {
	proposal := e.blkPool.GetBlockAndCheckHeight(vote.Hash, vote.View.Height)
	if proposal == nil {
		return nil, fmt.Errorf("last proposal %v not exist", vote.Hash)
	}

	seals := e.getVoteSeals(vote.Hash, size)
	sealedProposal, err := e.backend.PreCommit(proposal, seals)
	if err != nil {
		return nil, err
	}

	sealedBlock, ok := sealedProposal.(*types.Block)
	if !ok {
		return nil, errProposalConvert
	}

	extra := sealedBlock.Header().Extra
	qc := &hotstuff.QuorumCert{
		View:     vote.View,
		Hash:     sealedProposal.Hash(),
		Proposer: sealedProposal.Coinbase(),
		Extra:    extra,
	}

	return qc, nil
}

func (e *EventDrivenEngine) aggregateTC(event *TimeoutEvent, size int) *TimeoutCert {
	seals := e.getTimeoutSeals(event.View.Round.Uint64(), size)
	tc := &TimeoutCert{
		View:  event.View,
		Hash:  common.Hash{},
		Seals: seals,
	}
	return tc
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func isTC(qc *hotstuff.QuorumCert) bool {
	if qc.Hash == utils.EmptyHash {
		return true
	}
	return false
}

func sub1(num *big.Int) *big.Int {
	return new(big.Int).Sub(num, common.Big1)
}