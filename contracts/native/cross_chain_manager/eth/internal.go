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

package eth

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func VerifyAccountResult(proof *ethapi.AccountResult, header *types.Header, contractAddr common.Address) ([]byte, error) {
	if err := VerifyAccountProof(proof, header, contractAddr); err != nil {
		return nil, err
	}
	return VerifyStorageProof(proof)
}

func VerifyAccountProof(proof *ethapi.AccountResult, header *types.Header, contractAddr common.Address) error {
	nodeList := new(light.NodeList)
	for _, s := range proof.AccountProof {
		nodeList.Put(nil, common.FromHex(s))
	}
	ns := nodeList.NodeSet()

	if proof.Address != contractAddr {
		return fmt.Errorf("verifyMerkleProof, contract address is error, proof address: %s, side chain address: %s", proof.Address.Hex(), contractAddr.Hex())
	}
	acctKey := crypto.Keccak256(proof.Address[:])

	// 2. verify account proof
	acctVal, err := trie.VerifyProof(header.Root, acctKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify account proof error:%s\n", err)
	}

	acct := &state.Account{
		Nonce:    uint64(proof.Nonce),
		Balance:  proof.Balance.ToInt(),
		Root:     proof.StorageHash,
		CodeHash: proof.CodeHash[:],
	}

	acctrlp, err := rlp.EncodeToBytes(acct)
	if err != nil {
		return err
	}

	if !bytes.Equal(acctrlp, acctVal) {
		return fmt.Errorf("verifyMerkleProof, verify account proof failed, wanted:%v, get:%v", acctrlp, acctVal)
	}

	return nil
}

func VerifyStorageProof(proof *ethapi.AccountResult) ([]byte, error) {
	nodeList := new(light.NodeList)
	if len(proof.StorageProof) != 1 {
		return nil, fmt.Errorf("verifyMerkleProof, invalid storage proof format")
	}
	sp := proof.StorageProof[0]
	storageKey := crypto.Keccak256(common.HexToHash(scom.Replace0x(sp.Key)).Bytes())

	for _, prf := range sp.Proof {
		nodeList.Put(nil, common.Hex2Bytes(scom.Replace0x(prf)))
	}
	ns := nodeList.NodeSet()

	val, err := trie.VerifyProof(proof.StorageHash, storageKey[:], ns)
	if err != nil {
		return nil, fmt.Errorf("verifyMerkleProof, verify storage proof error:%s\n", err)
	}
	return val, nil
}

func CheckStorageResult(result, value []byte) error {
	if value == nil || result == nil {
		return fmt.Errorf("invalid result or value")
	}

	// length of `result` and `value` should be 33 && 32
	if len(result) != common.HashLength+1 {
		return fmt.Errorf("proof result is an rlp string and it's length should be 33")
	}
	if len(value) != common.HashLength {
		return fmt.Errorf("value is an full hash and it's length should be 32")
	}

	var s_temp []byte
	err := rlp.DecodeBytes(result, &s_temp)
	if err != nil {
		return err
	}
	var s []byte
	for i := len(s_temp); i < 32; i++ {
		s = append(s, 0)
	}
	s = append(s, s_temp...)

	value = cacheDBRecover(value)
	valueHash := common.BytesToHash(value)
	proofHash := common.BytesToHash(s)
	if proofHash != valueHash {
		return fmt.Errorf("proof result expect %s, got value %s", proofHash.Hex(), valueHash.Hex())
	}
	return nil
}

// todo(fuk): 1.cachedb将value的存储记录了一个标志位，后续state_object变更后需要改回来, 2. value 无需再次keccak256
func cacheDBRecover(value []byte) []byte {
	return append([]byte{1}, value[:common.HashLength-1]...)
}
