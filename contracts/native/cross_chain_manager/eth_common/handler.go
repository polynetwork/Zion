/*
 * Copyright (C) 2022 The Zion Authors
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

package eth_common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	icom "github.com/ethereum/go-ethereum/contracts/native/info_sync"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type Handler struct{}

func NewHandler() *Handler {
	return new(Handler)
}

func (h *Handler) MakeDepositProposal(service *native.NativeContract) (txParam *scom.MakeTxParam, err error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := side_chain_manager.GetSideChainObject(service, params.SourceChainID)
	if err != nil || sideChain == nil {
		err = fmt.Errorf("eth common handler  failed to get side chain instance, chain(%d) err: %v", params.SourceChainID, err)
		return
	}

	txParam, err = h.VerifyDepositProposal(service, sideChain, params)
	if err != nil {
		err = fmt.Errorf("eth common handler verify deposit proposal failure chain(%d):%s, err: %v", params.SourceChainID, sideChain.Name, err)
		return
	}

	err = scom.CheckDoneTx(service, txParam.CrossChainID, params.SourceChainID)
	if err != nil {
		err = fmt.Errorf("eth common handler check done transaction err: %v, chain(%d): %s", err, params.SourceChainID, sideChain.Name)
		return
	}

	err = scom.PutDoneTx(service, txParam.CrossChainID, params.SourceChainID)
	if err != nil {
		err = fmt.Errorf("eth common handler mark tx as done err: %v, chain(%d): %s", err, params.SourceChainID, sideChain.Name)
		return
	}
	return
}

func (h *Handler) VerifyDepositProposal(service *native.NativeContract,
	sideChain *side_chain_manager.SideChain, params *scom.EntranceParam) (txParam *scom.MakeTxParam, err error) {

	// Verify eth proof
	proof := new(Proof)
	err = json.Unmarshal(params.Proof, proof)
	if err != nil {
		err = fmt.Errorf("decode eth proof failed, err: %v", err)
		return
	}

	info, err := icom.GetRootInfo(service, sideChain.ChainID, params.Height)
	if err != nil {
		err = fmt.Errorf("get root info failure, err %v", err)
		return
	}
	if info == nil {
		err = fmt.Errorf("root info missing for height %d", params.Height)
		return
	}

	header, err := DecodeHeader(info)
	if err != nil {
		err = fmt.Errorf("decode root ifno failure, %v", err)
		return
	}

	err = VerifyCrossChainProof(crypto.Keccak256(params.Extra), proof, header.Root, sideChain.CCMCAddress)
	if err != nil {
		err = fmt.Errorf("VerifyCrossChainProof failed, err: %v", err)
		return
	}

	txParam, err = scom.DecodeTxParam(params.Extra)
	return
}

// Proof ...
type Proof struct {
	Address       string         `json:"address"`
	Balance       string         `json:"balance"`
	CodeHash      string         `json:"codeHash"`
	Nonce         string         `json:"nonce"`
	StorageHash   string         `json:"storageHash"`
	AccountProof  []string       `json:"accountProof"`
	StorageProofs []StorageProof `json:"storageProof"`
}

// StorageProof ...
type StorageProof struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Proof []string `json:"proof"`
}

// ProofAccount ...
type ProofAccount struct {
	Nonce    *big.Int
	Balance  *big.Int
	Storage  common.Hash
	Codehash common.Hash
}

// Verify account proof and contract storage proof
func VerifyCrossChainProof(value []byte, proof *Proof, root common.Hash, address []byte) (err error) {
	if root == (common.Hash{}) {
		return fmt.Errorf("empty root hash found in header")
	}
	nodeList := new(light.NodeList)
	for _, s := range proof.AccountProof {
		nodeList.Put(nil, common.Hex2Bytes(scom.Replace0x(s)))
	}
	ns := nodeList.NodeSet()
	addr := common.Hex2Bytes(scom.Replace0x(proof.Address))
	if !bytes.Equal(addr, address) {
		err = fmt.Errorf("contract address(%s) does not match with proof account(%s)", address, proof.Address)
		return
	}
	accountKey := crypto.Keccak256(addr)

	// Verify account proof
	accountValue, err := trie.VerifyProof(root, accountKey, ns)
	if err != nil {
		err = fmt.Errorf("account VerifyProof failure, err: %v", err)
		return
	}

	nonce, ok := new(big.Int).SetString(scom.Replace0x(proof.Nonce), 16)
	if !ok {
		err = fmt.Errorf("invalid account nonce: %s", proof.Nonce)
		return
	}
	balance, ok := new(big.Int).SetString(scom.Replace0x(proof.Balance), 16)
	if !ok {
		err = fmt.Errorf("invalid account balance: %s", proof.Balance)
		return
	}
	storageHash := common.HexToHash(proof.StorageHash)
	accountBytes, err := rlp.EncodeToBytes(&ProofAccount{
		Nonce:    nonce,
		Balance:  balance,
		Storage:  storageHash,
		Codehash: common.HexToHash(proof.CodeHash),
	})
	if err != nil {
		err = fmt.Errorf("rlp encode account value failed, err: %v", err)
		return
	}
	if !bytes.Equal(accountBytes, accountValue) {
		err = fmt.Errorf("account value does not match, wanted: %x, got: %x", accountBytes, accountValue)
		return
	}

	// Verify storage proof
	if len(proof.StorageProofs) != 1 {
		err = fmt.Errorf("invalid storage proof size, %v", proof.StorageProofs)
		return
	}
	sp := proof.StorageProofs[0]
	nodeList = new(light.NodeList)
	storageKey := crypto.Keccak256(common.HexToHash(sp.Key).Bytes())
	for _, p := range sp.Proof {
		nodeList.Put(nil, common.Hex2Bytes(scom.Replace0x(p)))
	}
	storageValue, err := trie.VerifyProof(storageHash, storageKey, nodeList.NodeSet())
	if err != nil {
		err = fmt.Errorf("account storage VerifyProof failure, err: %v", err)
		return
	}
	err = CheckProofResult(storageValue, value)
	if err != nil {
		err = fmt.Errorf("CheckProofResult failed, err: %v", err)
		return
	}
	return
}

// Check proof storage value hash
func CheckProofResult(result, value []byte) (err error) {
	var temp []byte
	err = rlp.DecodeBytes(result, &temp)
	if err != nil {
		err = fmt.Errorf("rlp decode proof result failed, err: %v", err)
		return
	}
	var hash []byte
	for i := len(temp); i < 32; i++ {
		hash = append(hash, 0)
	}
	hash = append(hash, temp...)
	if !bytes.Equal(hash, value) {
		err = fmt.Errorf("storage value does not match with proof result, wanted %x, got %x", result, value)
		return
	}
	return
}

type Header struct {
	Root common.Hash `json:"stateRoot" gencodec:"required"`
}

// Decode header
func DecodeHeader(data []byte) (h *Header, err error) {
	h = new(Header)
	err = json.Unmarshal(data, h)
	if err != nil {
		h = nil
	}
	return
}
