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

package backup
//
//import (
//	"math/big"
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//)
//
//// PaceMaker paceMaker is an standalone module which used to keep consensus liveness.
//// it driven the proposal round in sequence by an `QC` or `TC` which means that the consensus need at least
//// `2F + 1` valid vote message or timeout message to agree that the consensus engine can be driven into
//// the next round.
//type PaceMaker struct {
//	epoch uint64
//
//	currentHeight,
//	currentRound *big.Int
//
//	timtoutMsgs           map[uint64]common.Address // map round to message sender collection
//	highestCommittedRound *big.Int                  // record last committed block which used to calculate timeout duration
//	timer                 *time.Timer               // internal timer
//	sender                *EventSender
//}
//
//func NewPaceMaker(epoch uint64, initRound, initHighestCommittedRound *big.Int, sender *EventSender) *PaceMaker {
//	pm := new(PaceMaker)
//	pm.currentRound = initRound
//	pm.highestCommittedRound = initHighestCommittedRound
//	pm.timtoutMsgs = make(map[uint64]common.Address)
//	pm.sender = sender
//	return pm
//}
//
//func (p *PaceMaker) CurrentEpoch() uint64 {
//	return p.epoch
//}
//
//func (p *PaceMaker) CurrentHeight() *big.Int {
//	return p.currentHeight
//}
//
//func (p *PaceMaker) CurrentRound() *big.Int {
//	return p.currentRound
//}
//
//// ProcessLocalTimeout broadcast timeout message to all and record the message sender and round
//func (p *PaceMaker) ProcessLocalTimeout(sender common.Address, round *big.Int) {
//
//}
//
//// ProcessRemoteTimeout record timeout info and drive consensus into next round
//// if timeout message count arrived an quorum size.
//func (p *PaceMaker) ProcessRemoteTimeout(sender common.Address, round *big.Int) {
//
//}
//
//// AdvanceRound drive the consensus to the next round by qc/tc
//func (p *PaceMaker) AdvanceRound(qc *hotstuff.QuorumCert) error {
//	return nil
//}
//
//func (p *PaceMaker) stopTimer() {
//	if p.timer != nil {
//		p.timer.Stop()
//	}
//}
//
//// todo: add in config
//const standardTmoDuration = time.Second * 4
//
//func (p *PaceMaker) newRoundChangeTimer() {
//	p.stopTimer()
//
//	index := p.currentRound.Uint64() - 1
//	if p.highestCommittedRound.Uint64() != 0 {
//		if p.currentRound.Uint64()-p.highestCommittedRound.Uint64() < 3 {
//			index = 0
//		} else {
//			index = p.currentRound.Uint64() - p.highestCommittedRound.Uint64() - 3
//		}
//	}
//
//	tmoDuration := time.Duration(index) * standardTmoDuration
//	p.timer = time.AfterFunc(tmoDuration, func() {
//		p.sender.sendEvent(TimeoutEvent{
//			Epoch: p.epoch,
//			Round: nil,
//		})
//	})
//}
