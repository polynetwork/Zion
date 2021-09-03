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
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

// SMR state machine repo
type SMR struct {
	curEpoch      uint64
	curEpochStart *big.Int // [epochStart, epochEnd] is an closed interval, it represent height number but not round number
	curEpochEnd   *big.Int

	curRound  *big.Int
	curHeight *big.Int

	curHighestCommitRound *big.Int // used to calculate timeout duration
	curLatestVoteRound    *big.Int // latest vote round number

	curPendingRequest *types.Block
	curHighProposal   *types.Block
	curHighQC         *hotstuff.QuorumCert
	curLockQC         *hotstuff.QuorumCert
}

func newSMR() *SMR {
	return &SMR{}
}

func (s *SMR) Epoch() uint64 {
	return s.curEpoch
}

func (s *SMR) SetEpoch(e uint64) {
	s.curEpoch = e
}

func (s *SMR) EpochStart() *big.Int {
	return s.curEpochStart
}

func (s *SMR) EpochStartU64() uint64 {
	return s.EpochStart().Uint64()
}

func (s *SMR) SetEpochStart(height *big.Int) {
	s.curEpochStart = new(big.Int).Set(height)
}

func (s *SMR) EpochEnd() *big.Int {
	return s.curEpochEnd
}

func (s *SMR) EpochEndU64() uint64 {
	return s.EpochEnd().Uint64()
}

func (s *SMR) SetEpochEnd(height *big.Int) {
	s.curEpochEnd = new(big.Int).Set(height)
}

func (s *SMR) Round() *big.Int {
	return s.curRound
}

func (s *SMR) RoundU64() uint64 {
	return s.Round().Uint64()
}

func (s *SMR) SetRound(r *big.Int) {
	s.curRound = new(big.Int).Set(r)
}

func (s *SMR) Height() *big.Int {
	return s.curHeight
}

func (s *SMR) HeightU64() uint64 {
	return s.Height().Uint64()
}

func (s *SMR) SetHeight(h *big.Int) {
	s.curHeight = new(big.Int).Set(h)
}

func (s *SMR) HighCommitRound() *big.Int {
	return s.curHighestCommitRound
}

func (s *SMR) HighCommitRoundU64() uint64 {
	return s.HighCommitRound().Uint64()
}

func (s *SMR) SetHighCommitRound(r *big.Int) {
	s.curHighestCommitRound = new(big.Int).Set(r)
}

func (s *SMR) LatestVoteRound() *big.Int {
	return s.curLatestVoteRound
}

func (s *SMR) LatestVoteRoundU64() uint64 {
	return s.LatestVoteRound().Uint64()
}

func (s *SMR) SetLatestVoteRound(r *big.Int) {
	s.curLatestVoteRound = new(big.Int).Set(r)
}

func (s *SMR) GetRequest() *types.Block {
	if s.curPendingRequest != nil && s.curPendingRequest.Number().Cmp(s.curHeight) == 0 {
		return s.curPendingRequest
	}
	return nil
}

func (s *SMR) UpdateRequest(req *types.Block) {
	if req == nil {
		return
	}
	if s.curPendingRequest == nil || s.curPendingRequest.NumberU64() < req.NumberU64() {
		s.curPendingRequest = req
	}
}

func (s *SMR) LockQC() *hotstuff.QuorumCert {
	return s.curLockQC
}

func (s *SMR) SetLockQC(qc *hotstuff.QuorumCert) {
	if qc == nil || qc.View == nil {
		return
	}
	if s.curLockQC == nil || s.curLockQC.View == nil || s.curLockQC.View.Cmp(qc.View) < 0 {
		s.curLockQC = qc
	}
}

func (s *SMR) HighQC() *hotstuff.QuorumCert {
	return s.curHighQC
}

func (s *SMR) SetHighQC(qc *hotstuff.QuorumCert) {
	if qc == nil || qc.View == nil {
		return
	}
	if s.curHighQC == nil || s.curHighQC.View == nil || s.curHighQC.View.Cmp(qc.View) < 0 {
		s.curHighQC = qc
	}
}

func (s *SMR) Proposal() *types.Block {
	return s.curHighProposal
}

func (s *SMR) SetProposal(proposal *types.Block) {
	if proposal.NumberU64() < s.curHighProposal.NumberU64() {
		return
	}
	if proposal.Hash() == s.curHighProposal.Hash() {
		return
	}
	s.curHighProposal = proposal
}
