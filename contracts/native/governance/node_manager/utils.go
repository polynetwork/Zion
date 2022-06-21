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

package node_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func nativeTransfer(s *native.NativeContract, from, to common.Address, amount *big.Int) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if !core.CanTransfer(s.StateDB(), from, amount) {
		return fmt.Errorf("%s insufficient balance", from.Hex())
	}
	core.Transfer(s.StateDB(), from, to, amount)
	return nil
}

func CheckConsensusSigns(s *native.NativeContract, method string, input []byte, signer common.Address) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	log.Trace("checkConsensusSign", "method", method, "input", hexutil.Encode(input), "signer", signer.Hex())

	// get epoch info
	epoch, err := GetCurrentEpochInfo(s)
	if err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, GetCurrentEpochInfo error: %v", err)
	}

	// check authority
	if err := CheckValidatorAuthority(signer, caller, epoch); err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, CheckValidatorAuthority error: %v", err)
	}

	// get or set consensus sign info
	sign := &ConsensusSign{Method: method, Input: input}
	if exist, err := getSign(s, sign.Hash()); err != nil {
		if err.Error() == "EOF" {
			if err := storeSign(s, sign); err != nil {
				return false, fmt.Errorf("CheckConsensusSigns, storeSign error: %v, hash %s", err, sign.Hash().Hex())
			} else {
				log.Trace("checkConsensusSign", "store sign, hash", sign.Hash().Hex())
			}
		} else {
			return false, fmt.Errorf("CheckConsensusSigns, get sign error: %v, hash %s", err, sign.Hash().Hex())
		}
	} else if exist.Hash() != sign.Hash() {
		return false, fmt.Errorf("CheckConsensusSigns, check sign hash failed, expect: %s, got %s", exist.Hash().Hex(), sign.Hash().Hex())
	}

	// check duplicate signature
	if findSigner(s, sign.Hash(), signer) {
		return false, fmt.Errorf("CheckConsensusSigns, signer already exist: %s, hash %s", signer.Hex(), sign.Hash().Hex())
	}

	// do not store redundancy sign
	sizeBeforeSign := getSignerSize(s, sign.Hash())
	log.Trace("checkConsensusSign", "sign hash", sign.Hash().Hex(), "size before sign", sizeBeforeSign)
	if sizeBeforeSign >= epoch.ValidatorQuorumSize() {
		return false, nil
	}

	// store signer address and emit event log
	if err := storeSigner(s, sign.Hash(), signer); err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, store signer failed: %s, hash %s", err, sign.Hash().Hex())
	}
	sizeAfterSign := getSignerSize(s, sign.Hash())
	if err := s.AddNotify(ABI, []string{"CheckConsensusSigns"}, sign.Method, sign.Input, signer, uint64(sizeAfterSign)); err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, emit consensus sign log failed: %s, hash %s", err, sign.Hash().Hex())
	}
	log.Trace("checkConsensusSign", "sign hash", sign.Hash().Hex(), "size after sign", sizeAfterSign)

	return sizeAfterSign >= epoch.ValidatorQuorumSize(), nil
}

func CheckVoterSigns(s *native.NativeContract, method string, input []byte, signer common.Address) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	log.Trace("CheckVoterSigns", "method", method, "input", hexutil.Encode(input), "signer", signer.Hex())

	// get epoch info
	epoch, err := GetCurrentEpochInfo(s)
	if err != nil {
		return false, fmt.Errorf("CheckVoterSigns, GetCurrentEpochInfo error: %v", err)
	}

	// check authority
	if err := CheckVoterAuthority(signer, caller, epoch); err != nil {
		return false, fmt.Errorf("CheckVoterSigns, CheckValidatorAuthority error: %v", err)
	}

	// get or set consensus sign info
	sign := &ConsensusSign{Method: method, Input: input}
	if exist, err := getSign(s, sign.Hash()); err != nil {
		if err.Error() == "EOF" {
			if err := storeSign(s, sign); err != nil {
				return false, fmt.Errorf("CheckVoterSigns, storeSign error: %v, hash %s", err, sign.Hash().Hex())
			} else {
				log.Trace("CheckVoterSigns", "store sign, hash", sign.Hash().Hex())
			}
		} else {
			return false, fmt.Errorf("CheckVoterSigns, get sign error: %v, hash %s", err, sign.Hash().Hex())
		}
	} else if exist.Hash() != sign.Hash() {
		return false, fmt.Errorf("CheckVoterSigns, check sign hash failed, expect: %s, got %s", exist.Hash().Hex(), sign.Hash().Hex())
	}

	// check duplicate signature
	if findSigner(s, sign.Hash(), signer) {
		return false, fmt.Errorf("CheckVoterSigns, signer already exist: %s, hash %s", signer.Hex(), sign.Hash().Hex())
	}

	// do not store redundancy sign
	sizeBeforeSign := getSignerSize(s, sign.Hash())
	log.Trace("CheckVoterSigns", "sign hash", sign.Hash().Hex(), "size before sign", sizeBeforeSign)
	if sizeBeforeSign >= epoch.VoterQuorumSize() {
		return false, nil
	}

	// store signer address and emit event log
	if err := storeSigner(s, sign.Hash(), signer); err != nil {
		return false, fmt.Errorf("CheckVoterSigns, store signer failed: %s, hash %s", err, sign.Hash().Hex())
	}
	sizeAfterSign := getSignerSize(s, sign.Hash())
	if err := s.AddNotify(ABI, []string{"CheckVoterSigns"}, sign.Method, sign.Input, signer, uint64(sizeAfterSign)); err != nil {
		return false, fmt.Errorf("CheckVoterSigns, emit consensus sign log failed: %s, hash %s", err, sign.Hash().Hex())
	}
	log.Trace("CheckVoterSigns", "sign hash", sign.Hash().Hex(), "size after sign", sizeAfterSign)

	return sizeAfterSign >= epoch.VoterQuorumSize(), nil
}

func CheckValidatorAuthority(origin, caller common.Address, epoch *EpochInfo) error {
	if epoch == nil || epoch.Validators == nil {
		return fmt.Errorf("invalid epoch")
	}
	if origin == common.EmptyAddress || caller == common.EmptyAddress {
		return fmt.Errorf("origin/caller is empty address")
	}
	if origin != caller {
		return fmt.Errorf("origin must be caller")
	}
	for _, v := range epoch.Validators {
		if v.Address == origin {
			return nil
		}
	}
	return fmt.Errorf("tx origin %s is not valid validator", origin.Hex())
}

func CheckVoterAuthority(origin, caller common.Address, epoch *EpochInfo) error {
	if epoch == nil || epoch.Voters == nil {
		return fmt.Errorf("invalid epoch")
	}
	if origin == common.EmptyAddress || caller == common.EmptyAddress {
		return fmt.Errorf("origin/caller is empty address")
	}
	if origin != caller {
		return fmt.Errorf("origin must be caller")
	}
	for _, v := range epoch.Voters {
		if v.Address == origin {
			return nil
		}
	}
	return fmt.Errorf("tx origin %s is not valid validator", origin.Hex())
}

func EpochChangeAtNextBlock(curHeight, epochStartHeight uint64) bool {
	if curHeight+1 == epochStartHeight {
		return true
	}
	return false
}
