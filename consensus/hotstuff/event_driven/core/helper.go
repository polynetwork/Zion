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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) newMsgLogger(msgtyp interface{}) log.Logger {
	return c.logger.New("view", c.currentView(), "msg", msgtyp)
}

func (c *core) newSenderLogger(msgtyp string) log.Logger {
	return c.logger.New("view", c.currentView(), "msg", msgtyp)
}

func (c *core) newLogger() log.Logger {
	return c.logger.New("view", c.currentView())
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return c.signer.CheckSignature(c.valset, data, sig)
}

func (c *core) isSelf(addr common.Address) bool {
	return c.address == addr
}

func (c *core) currentView() *hotstuff.View {
	return &hotstuff.View{
		Round:  new(big.Int).Set(c.smr.Round()),
		Height: new(big.Int).Set(c.smr.Height()),
	}
}

func (c *core) getProposer() common.Address {
	proposer := c.valset.GetProposer()
	return proposer.Address()
}

func (c *core) checkProposer(proposer common.Address) error {
	if !c.valset.IsProposer(proposer) {
		return errNotFromProposer
	}
	return nil
}

func (c *core) checkEpoch(epoch uint64, height *big.Int) error {
	if c.smr.Epoch() != epoch {
		return errInvalidHighQC
	}
	if height.Cmp(c.smr.EpochStart()) < 0 {
		return errInvalidEpoch
	}
	if height.Cmp(c.smr.EpochEnd()) > 0 {
		return errInvalidEpoch
	}
	return nil
}

func (c *core) checkView(code hotstuff.MsgType, view *hotstuff.View) error {
	if cmp := view.Cmp(c.currentView()); cmp > 0 {
		return errFutureMessage
	} else if cmp < 0 {
		return errOldMessage
	}

	// validator receive vote before proposal, allow to receive msg `timeout` and `tc` before proposal
	if c.smr.State() == StateNewRound && code == MsgTypeVote {
		return errFutureMessage
	}

	// validator should not receive duplicate proposal msg, allow to receive msg `vote`,`timeout` and `tc`
	if c.smr.State() == StateProposed && code == MsgTypeProposal {
		return errOldMessage
	}

	// validator should not receive another proposal after vote, allow to receive msg `timeout` and `tc`
	if c.smr.State() == StateVoted && code == MsgTypeProposal {
		return errOldMessage
	}

	return nil
}

func (c *core) checkProposalView(proposal *types.Block, view *hotstuff.View) (*hotstuff.View, error) {
	salt, _, err := extraProposal(proposal)
	if err != nil {
		return nil, err
	}

	proposalView := newView(salt.Round, proposal.Number())
	if view.Cmp(proposalView) != 0 {
		return nil, fmt.Errorf("proposal view expect %v, got %v", view, proposalView)
	}
	return proposalView, nil
}

func (c *core) checkJustifyQC(proposal hotstuff.Proposal, justifyQC *hotstuff.QuorumCert) error {
	if justifyQC == nil {
		return fmt.Errorf("justifyQC is nil")
	}
	if !bigEq0(justifyQC.Height()) && justifyQC.View == nil {
		return fmt.Errorf("justifyQC view is nil")
	}
	if justifyQC.Hash == common.EmptyHash {
		return fmt.Errorf("justifyQC hash is empty")
	}
	if !bigEq0(justifyQC.Height()) && justifyQC.Proposer == common.EmptyAddress {
		return fmt.Errorf("justifyQC proposer is empty")
	}
	if justifyQC.Hash != proposal.ParentHash() {
		return fmt.Errorf("justifyQC hash extendship invalid")
	}
	if _, eq := bigSub1Eq(proposal.Number(), justifyQC.Height()); !eq {
		return fmt.Errorf("justifyQC height invalid")
	}

	if !bigEq0(justifyQC.Height()) {
		vs := c.valset.Copy()
		vs.CalcProposerByIndex(justifyQC.View.Round.Uint64())
		proposer := vs.GetProposer().Address()
		if proposer != justifyQC.Proposer {
			return fmt.Errorf("justifyQC proposer expect %v got %v", proposer, justifyQC.Proposer)
		}
	}

	return nil
}

func (c *core) compareQC(expect, src *hotstuff.QuorumCert) error {
	if expect == nil || expect.View == nil {
		return fmt.Errorf("invalid expect qc")
	}
	if src == nil || src.View == nil {
		return fmt.Errorf("invalid src qc")
	}
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
	// todo(fuk): or implement this with `reflect.DeepEqual(expect, src)`
	return nil
}

func (c *core) getVoteSeals(hash common.Hash, n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range c.messages.Votes(hash) {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

func (c *core) getTimeoutSeals(round uint64, n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range c.messages.Timeouts(round) {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

func (c *core) Q() int {
	return c.valset.Q()
}

func (c *core) chain2Height() *big.Int {
	return new(big.Int).Add(c.smr.EpochStart(), common.Big2)
}

func (c *core) isChain2() bool {
	return c.smr.Height().Cmp(c.chain2Height()) >= 0
}

func (c *core) chain3Height() *big.Int {
	return new(big.Int).Add(c.smr.EpochStart(), common.Big3)
}

func (c *core) isChain3() bool {
	return c.smr.Height().Cmp(c.chain3Height()) >= 0
}

func (c *core) generateTimeoutEvent() *TimeoutEvent {
	tm := &TimeoutEvent{
		Epoch: c.smr.Epoch(),
		View:  c.currentView(),
	}
	tm.Digest = tm.Hash()
	return tm
}

func (c *core) aggregateQC(vote *Vote, size int) (*hotstuff.QuorumCert, *types.Block, error) {
	if vote == nil || vote.View == nil {
		return nil, nil, fmt.Errorf("invalid vote")
	}
	proposal := c.blkPool.GetBlockAndCheckHeight(vote.Hash, vote.View.Height)
	if proposal == nil {
		return nil, nil, fmt.Errorf("last proposal %v not exist", vote.Hash)
	}

	seals := c.getVoteSeals(vote.Hash, size)
	sealedProposal, err := c.backend.PreCommit(proposal, seals)
	if err != nil {
		return nil, nil, err
	}

	sealedBlock, ok := sealedProposal.(*types.Block)
	if !ok {
		return nil, nil, errProposalConvert
	}

	extra := sealedBlock.Header().Extra
	qc := &hotstuff.QuorumCert{
		View:     vote.View,
		Hash:     sealedProposal.Hash(),
		Proposer: sealedProposal.Coinbase(),
		Extra:    extra,
	}

	return qc, proposal, nil
}

func (c *core) aggregateTC(event *TimeoutEvent, size int) *TimeoutCert {
	seals := c.getTimeoutSeals(event.View.Round.Uint64(), size)
	tc := &TimeoutCert{
		View:  event.View,
		Hash:  event.Digest,
		Seals: seals,
	}
	return tc
}

func (c *core) updateHighQCAndProposal(qc *hotstuff.QuorumCert, proposal *types.Block) error {
	c.smr.SetHighQC(qc)
	c.smr.SetProposal(proposal)

	qcBlock := c.blkPool.GetBlockByHash(qc.Hash)
	if qcBlock == nil {
		return fmt.Errorf("qc block %v not exist", qc.Hash)
	}

	sealedProposal, err := c.backend.ForwardCommit(qcBlock, qc.Extra)
	if err != nil {
		return fmt.Errorf("forward commit err %v", err)
	}
	sealedBlock, ok := sealedProposal.(*types.Block)
	if !ok {
		return errProposalConvert
	}
	if err := c.blkPool.AddBlock(sealedBlock, qc.View.Round); err != nil {
		return fmt.Errorf("add block err %v", err)
	}
	return nil
}

func (c *core) nextValSet() hotstuff.ValidatorSet {
	vs := c.valset.Copy()
	vs.CalcProposerByIndex(c.smr.Round().Uint64() + 1)
	return vs
}

func (c *core) nextProposer() common.Address {
	vs := c.valset.Copy()
	vs.CalcProposerByIndex(c.smr.Round().Uint64() + 1)
	proposer := vs.GetProposer()
	return proposer.Address()
}
