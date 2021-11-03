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
	storageKey := common.HexToHash(sp.Key)

	for _, prf := range sp.Proof {
		nodeList.Put(nil, common.FromHex(prf))
	}
	ns := nodeList.NodeSet()

	val, err := trie.VerifyProof(proof.StorageHash, storageKey[:], ns)
	if err != nil {
		return nil, fmt.Errorf("verifyMerkleProof, verify storage proof error:%s\n", err)
	}
	return val, nil
}
