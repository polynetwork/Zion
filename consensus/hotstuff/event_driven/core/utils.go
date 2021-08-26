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
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
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

func (e *EventDrivenEngine) checkView(view *hotstuff.View) error {
	if e.curRound.Cmp(view.Round) != 0 || e.curHeight.Cmp(view.Height) != 0 {
		return errInvalidMessage
	}
	return nil
}

func (e *EventDrivenEngine) getMessageSeals(hash common.Hash, n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range e.messages.Votes(hash) {
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

func (e *EventDrivenEngine) aggregate(vote *Vote, size int) (*hotstuff.QuorumCert, error) {
	proposal := e.blkTree.GetBlockAndCheckHeight(vote.Hash, vote.View.Height)
	if proposal == nil {
		return nil, fmt.Errorf("last proposal %v not exist", vote.Hash)
	}

	seals := e.getMessageSeals(vote.Hash, size)
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

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func isTC(qc *hotstuff.QuorumCert) bool {
	if qc.Hash == utils.EmptyHash {
		return true
	}
	return false
}
