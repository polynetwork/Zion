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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func (e *EventDrivenEngine) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return e.signer.CheckSignature(e.valset, data, sig)
}

// todo
func (e *EventDrivenEngine) newLogger() log.Logger {
	logger := e.logger.New("state")
	return logger
}

// todo:
func proposal2QC(proposal hotstuff.Proposal, round *big.Int) *hotstuff.QuorumCert {
	block := proposal.(*types.Block)
	h := block.Header()
	qc := new(hotstuff.QuorumCert)
	qc.View = &hotstuff.View{
		Height: block.Number(),
		Round:  round,
	}
	qc.Hash = h.Hash()
	qc.Proposer = h.Coinbase
	qc.Extra = h.Extra
	return qc
}
