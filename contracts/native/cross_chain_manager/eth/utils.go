/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package eth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ecom "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth/types"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	common2 "github.com/ethereum/go-ethereum/contracts/native/info_sync/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"math/big"
)

func verifyFromEthTx(native *native.NativeContract, proof, extra []byte, fromChainID uint64, height uint32,
	sideChain *side_chain_manager.SideChain) (*scom.MakeTxParam, error) {
	value, err := common2.GetRootInfo(native, fromChainID, height)
	if err != nil {
		return nil, fmt.Errorf("verifyFromEthTx, GetCrossChainInfo error:%s", err)
	}
	header := &Header{}
	err = json.Unmarshal(value, header)
	if err != nil {
		return nil, fmt.Errorf("verifyFromEthTx, json unmarshal header error:%s", err)
	}

	ethProof := new(ETHProof)
	err = json.Unmarshal(proof, ethProof)
	if err != nil {
		return nil, fmt.Errorf("verifyFromEthTx, unmarshal proof error:%s", err)
	}

	if len(ethProof.StorageProofs) != 1 {
		return nil, fmt.Errorf("verifyFromEthTx, incorrect proof format")
	}

	//todo 1. verify the proof with header
	//determine where the k and v from
	proofResult, err := VerifyMerkleProof(ethProof, header, sideChain.CCMCAddress)
	if err != nil {
		return nil, fmt.Errorf("verifyFromEthTx, verifyMerkleProof error:%v", err)
	}
	if proofResult == nil {
		return nil, fmt.Errorf("verifyFromEthTx, verifyMerkleProof failed!")
	}

	if !CheckProofResult(proofResult, extra) {
		return nil, fmt.Errorf("verifyFromEthTx, verify proof value hash failed, proof result:%x, extra:%x", proofResult, extra)
	}

	txParam, err := scom.DecodeTxParam(extra)
	if err != nil {
		return nil, fmt.Errorf("verifyFromEthTx, deserialize merkleValue error:%s", err)
	}
	return txParam, nil
}

// used by quorum
func VerifyMerkleProofLegacy(ethProof *ETHProof, blockData *types.Header, contractAddr []byte) ([]byte, error) {
	return VerifyMerkleProof(ethProof, To1559(blockData), contractAddr)
}

func VerifyMerkleProof(ethProof *ETHProof, blockData *Header, contractAddr []byte) ([]byte, error) {
	//1. prepare verify account
	nodeList := new(light.NodeList)

	for _, s := range ethProof.AccountProof {
		p := scom.Replace0x(s)
		nodeList.Put(nil, ecom.Hex2Bytes(p))
	}
	ns := nodeList.NodeSet()

	addr := ecom.Hex2Bytes(scom.Replace0x(ethProof.Address))
	if !bytes.Equal(addr, contractAddr) {
		return nil, fmt.Errorf("verifyMerkleProof, contract address is error, proof address: %s, side chain address: %s", ethProof.Address, hex.EncodeToString(contractAddr))
	}
	acctKey := crypto.Keccak256(addr)

	// 2. verify account proof
	acctVal, err := trie.VerifyProof(blockData.Root, acctKey, ns)
	if err != nil {
		return nil, fmt.Errorf("verifyMerkleProof, verify account proof error:%s\n", err)
	}

	nounce := new(big.Int)
	_, ok := nounce.SetString(scom.Replace0x(ethProof.Nonce), 16)
	if !ok {
		return nil, fmt.Errorf("verifyMerkleProof, invalid format of nounce:%s\n", ethProof.Nonce)
	}

	balance := new(big.Int)
	_, ok = balance.SetString(scom.Replace0x(ethProof.Balance), 16)
	if !ok {
		return nil, fmt.Errorf("verifyMerkleProof, invalid format of balance:%s\n", ethProof.Balance)
	}

	storageHash := ecom.HexToHash(scom.Replace0x(ethProof.StorageHash))
	codeHash := ecom.HexToHash(scom.Replace0x(ethProof.CodeHash))

	acct := &ProofAccount{
		Nounce:   nounce,
		Balance:  balance,
		Storage:  storageHash,
		Codehash: codeHash,
	}

	acctrlp, err := rlp.EncodeToBytes(acct)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(acctrlp, acctVal) {
		return nil, fmt.Errorf("verifyMerkleProof, verify account proof failed, wanted:%v, get:%v", acctrlp, acctVal)
	}

	//3.verify storage proof
	nodeList = new(light.NodeList)
	if len(ethProof.StorageProofs) != 1 {
		return nil, fmt.Errorf("verifyMerkleProof, invalid storage proof format")
	}

	sp := ethProof.StorageProofs[0]
	storageKey := crypto.Keccak256(ecom.HexToHash(scom.Replace0x(sp.Key)).Bytes())

	for _, prf := range sp.Proof {
		nodeList.Put(nil, ecom.Hex2Bytes(scom.Replace0x(prf)))
	}

	ns = nodeList.NodeSet()
	val, err := trie.VerifyProof(storageHash, storageKey, ns)
	if err != nil {
		return nil, fmt.Errorf("verifyMerkleProof, verify storage proof error:%s\n", err)
	}

	return val, nil
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
