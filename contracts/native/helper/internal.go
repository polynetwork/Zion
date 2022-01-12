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

package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func VerifyTx(proof []byte, hdr *types.Header, contract common.Address, extra []byte, checkResult bool) ([]byte, error) {
	ethProof := new(ethapi.AccountResult)
	if err := json.Unmarshal(proof, ethProof); err != nil {
		return nil, fmt.Errorf("VerifyFromEthProof, unmarshal proof failed, err:%v", err)
	}

	proofResult, err := VerifyAccountResult(ethProof, hdr, contract)
	if err != nil {
		return nil, fmt.Errorf("VerifyFromEthProof, verifyMerkleProof failed, err:%v", err)
	}
	if proofResult == nil {
		return nil, fmt.Errorf("VerifyFromEthProof, verifyMerkleProof failed, err:%s", "proof result is nil")
	}
	if checkResult && !CheckProofResult(proofResult, extra) {
		return nil, fmt.Errorf("VerifyFromEthProof, failed to check result, stored %s, got %s", hexutil.Encode(proofResult), hexutil.Encode(extra))
	}
	return proofResult, nil
}

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
	storageKey := crypto.Keccak256(common.HexToHash(Replace0x(sp.Key)).Bytes())

	for _, prf := range sp.Proof {
		nodeList.Put(nil, common.Hex2Bytes(Replace0x(prf)))
	}
	ns := nodeList.NodeSet()

	val, err := trie.VerifyProof(proof.StorageHash, storageKey[:], ns)
	if err != nil {
		return nil, fmt.Errorf("verifyMerkleProof, verify storage proof error:%s\n", err)
	}
	return val, nil
}

func Replace0x(s string) string {
	return strings.Replace(strings.ToLower(s), "0x", "", 1)
}

func CheckProofResult(result, value []byte) bool {
	var s_temp []byte
	err := rlp.DecodeBytes(result, &s_temp)
	if err != nil {
		log.Errorf("checkProofResult, rlp.DecodeBytes error:%s\n", err)
		return false
	}
	//
	var s []byte
	for i := len(s_temp); i < 32; i++ {
		s = append(s, 0)
	}
	s = append(s, s_temp...)
	hash := crypto.Keccak256(value)

	return bytes.Equal(s, hash)
}
