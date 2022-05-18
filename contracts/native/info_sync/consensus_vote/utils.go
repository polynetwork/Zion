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

package consensus_vote

import (
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/log"
)

func CheckConsensusSigns(s *native.NativeContract, input []byte) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	// get epoch info
	epochBytes, err := node_manager.GetCurrentEpoch(s)
	if err != nil {
		log.Trace("checkConsensusSign", "get current epoch bytes failed", err)
		return false, node_manager.ErrEpochNotExist
	}
	output := new(node_manager.MethodEpochOutput)
	output.Decode(epochBytes)
	epoch := output.Epoch

	// check authority
	if err := node_manager.CheckAuthority(caller, caller, epoch); err != nil {
		log.Trace("checkConsensusSign", "check authority failed", err)
		return false, node_manager.ErrInvalidAuthority
	}

	// get or set consensus sign info
	msg := &VoteMessage{Input: input}
	if exist, err := getVoteMessage(s, msg.Hash()); err != nil {
		if err.Error() == "EOF" {
			if err := storeVoteMessage(s, msg); err != nil {
				log.Trace("checkConsensusSign", "store sign failed", err, "hash", msg.Hash().Hex())
				return false, node_manager.ErrStorage
			}
		} else {
			log.Trace("checkConsensusSign", "get sign failed", err, "hash", msg.Hash().Hex())
			return false, node_manager.ErrConsensusSignNotExist
		}
	} else if exist.Hash() != msg.Hash() {
		log.Trace("checkConsensusSign", "check sign hash failed, expect", exist.Hash().Hex(), "got", msg.Hash().Hex())
		return false, node_manager.ErrInvalidSign
	}

	// check duplicate signature
	if findSigner(s, msg.Hash(), caller) {
		log.Trace("checkConsensusSign", "signer already exist", caller.Hex(), "hash", msg.Hash().Hex())
		return false, nil
	}

	// store signer address and check quorum
	ok, err := storeSignerAndCheckQuorum(s, msg.Hash(), caller, epoch.QuorumSize())
	if err != nil {
		log.Trace("checkConsensusSign", "store signer failed", err, "hash", msg.Hash().Hex())
		return false, node_manager.ErrStorage
	}

	return ok, nil
}
