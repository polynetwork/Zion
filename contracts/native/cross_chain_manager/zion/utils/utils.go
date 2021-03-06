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

package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	hlp "github.com/ethereum/go-ethereum/contracts/native/helper"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rlp"
)

func VerifyTx(proof []byte, hdr *types.Header, contract common.Address, extra []byte, checkResult bool) ([]byte, error) {
	ethProof := new(ethapi.AccountResult)
	if err := json.Unmarshal(proof, ethProof); err != nil {
		return nil, fmt.Errorf("VerifyFromEthProof, unmarshal proof failed, err:%v", err)
	}

	proofResult, err := hlp.VerifyAccountResult(ethProof, hdr, contract)
	if err != nil {
		return nil, fmt.Errorf("VerifyFromEthProof, verifyMerkleProof failed, err:%v", err)
	}
	if proofResult == nil {
		return nil, fmt.Errorf("VerifyFromEthProof, verifyMerkleProof failed, err:%s", "proof result is nil")
	}
	if checkResult && !hlp.CheckProofResult(proofResult, extra) {
		return nil, fmt.Errorf("VerifyFromEthProof, failed to check result, stored %s, got %s", hexutil.Encode(proofResult), hexutil.Encode(extra))
	}
	return proofResult, nil
}

func EncodeTxArgs(toAssetHash, toAddress []byte, amount *big.Int) ([]byte, error) {
	args := &scom.TxArgs{
		ToAssetHash: toAssetHash,
		ToAddress:   toAddress,
		Amount:      amount,
	}
	return rlp.EncodeToBytes(args)
}

func DecodeTxArgs(payload []byte) (*scom.TxArgs, error) {
	args := new(scom.TxArgs)
	if err := rlp.DecodeBytes(payload, args); err != nil {
		return nil, err
	}
	return args, nil
}

var (
	BytesTy, _   = abi.NewType("bytes", "", nil)
	AddrTy, _    = abi.NewType("address", "", nil)
	Uint64Ty, _  = abi.NewType("uint64", "", nil)
	Uint256Ty, _ = abi.NewType("uint256", "", nil)
)

type TunnelData struct {
	Caller     common.Address
	ToContract []byte
	Method     []byte
	TxData     []byte
}

func newTunnelArguments() abi.Arguments {
	return abi.Arguments{
		{Type: AddrTy, Name: "_caller"},
		{Type: BytesTy, Name: "_toContract"},
		{Type: BytesTy, Name: "_method"},
		{Type: BytesTy, Name: "_txData"},
	}
}

// bytes memory tunnelData = abi.encode(Utils.addressToBytes(msg.sender), _toContract, _method, _txData);
func (d *TunnelData) Encode() ([]byte, error) {
	args := newTunnelArguments()
	return args.Pack(d.Caller, d.ToContract, d.Method, d.TxData)
}

func (d *TunnelData) Decode(payload []byte) error {
	args := newTunnelArguments()
	list, err := args.Unpack(payload)
	if err != nil {
		return err
	}
	return args.Copy(d, list)
}

func EncodeMakeTxParams(paramTxHash []byte, crossChainId []byte, caller []byte,
	toChainID uint64, toContract []byte, method string, args []byte) (*scom.MakeTxParam, []byte, error) {

	tx := &scom.MakeTxParam{
		TxHash:              paramTxHash,
		CrossChainID:        crossChainId,
		FromContractAddress: caller,
		ToChainID:           toChainID,
		ToContractAddress:   toContract,
		Method:              method,
		Args:                args,
	}

	blob, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, nil, err
	}
	//txProof := crypto.Keccak256Hash(blob)
	return tx, blob, nil
}

func DecodeMakeTxParams(payload []byte) (*scom.MakeTxParam, error) {
	txParams := new(scom.MakeTxParam)
	if err := rlp.DecodeBytes(payload, txParams); err != nil {
		return nil, err
	}
	return txParams, nil
}

func GenerateCrossChainID(addr common.Address, paramTxHash []byte) []byte {
	blob := utils.EncodePacked(addr[:], paramTxHash[:])
	sum := sha256.Sum256(blob)
	return sum[:]
}
