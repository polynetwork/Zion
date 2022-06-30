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
package common

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var ErrTxAlreadyImported = errors.New("tx already imported")

func Replace0x(s string) string {
	return strings.Replace(strings.ToLower(s), "0x", "", 1)
}

func PutDoneTx(native *native.NativeContract, crossChainID []byte, chainID uint64) error {
	native.GetCacheDB().Put(doneTxKey(chainID, crossChainID), crossChainID)
	return nil
}

func CheckDoneTx(native *native.NativeContract, crossChainID []byte, chainID uint64) error {
	value, err := native.GetCacheDB().Get(doneTxKey(chainID, crossChainID))
	if err != nil {
		return fmt.Errorf("checkDoneTx, native.GetCacheDB().Get error: %v", err)
	}
	if value != nil {
		return ErrTxAlreadyImported
	}
	return nil
}

func NotifyMakeProof(native *native.NativeContract, merkleValueHex string, key string) error {
	return native.AddNotify(ABI, []string{NOTIFY_MAKE_PROOF_EVENT}, merkleValueHex, native.ContractRef().BlockHeight().Uint64(), key)
}

func NotifyReplenish(native *native.NativeContract, txHashes []string, chainId uint64) error {
	err := native.AddNotify(ABI, []string{REPLENISH_EVENT}, txHashes, chainId)
	if err != nil {
		return fmt.Errorf("NotifyReplenish failed: %v", err)
	}
	return nil
}

func Uint256ToBytes(num *big.Int) []byte {
	if num == nil {
		return common.EmptyHash[:]
	}
	return common.LeftPadBytes(num.Bytes(), 32)
}

func BytesToUint256(data []byte) *big.Int {
	if data == nil || len(data) == 0 {
		return common.Big0
	}
	return new(big.Int).SetBytes(common.TrimLeftZeroes(data))
}

func doneTxKey(chainID uint64, crossChainID []byte) []byte {
	contract := utils.CrossChainManagerContractAddress
	return utils.ConcatKey(contract, []byte(DONE_TX), utils.GetUint64Bytes(chainID), crossChainID)
}
