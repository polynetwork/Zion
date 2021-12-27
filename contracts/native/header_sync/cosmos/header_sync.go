/*
 * Copyright (C) 2020 The poly network Authors
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

package cosmos

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type CosmosHandler struct{}

func NewCosmosHandler() *CosmosHandler {
	return &CosmosHandler{}
}

func newCDC() *codec.Codec {
	cdc := codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{}, multisig.PubKeyMultisigThresholdAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.PrivKeyAminoName, nil)
	return cdc
}

func (this *CosmosHandler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	param := &scom.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncGenesisHeader, param, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, scom.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}

	// get genesis header from input parameters
	cdc := newCDC()
	var header CosmosHeader
	err = cdc.UnmarshalBinaryBare(param.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader: %s", err)
	}
	// check if has genesis header
	info, err := GetEpochSwitchInfo(native, param.ChainID)
	if err == nil && info != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, genesis header had been initialized")
	}
	if err := PutEpochSwitchInfo(native, param.ChainID, &CosmosEpochSwitchInfo{
		Height:             uint64(header.Header.Height),
		NextValidatorsHash: header.Header.NextValidatorsHash,
		ChainID:            header.Header.ChainID,
		BlockHash:          header.Header.Hash(),
	}); err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, failed to PutEpochSwitchInfo: %v", err)
	}
	return nil
}

func (this *CosmosHandler) SyncBlockHeader(native *native.NativeContract) error {
	params := &scom.SyncBlockHeaderParam{}

	ctx := native.ContractRef().CurrentContext()
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
		return err
	}

	cdc := newCDC()
	cnt := 0
	info, err := GetEpochSwitchInfo(native, params.ChainID)
	if err != nil {
		return fmt.Errorf("SyncBlockHeader, get epoch switching height failed: %v", err)
	}
	for _, v := range params.Headers {
		var myHeader CosmosHeader
		err := cdc.UnmarshalBinaryBare(v, &myHeader)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader failed to unmarshal header: %v", err)
		}
		if bytes.Equal(myHeader.Header.NextValidatorsHash, myHeader.Header.ValidatorsHash) {
			continue
		}
		if info.Height >= uint64(myHeader.Header.Height) {
			log.Debugf("SyncBlockHeader, height %d is lower or equal than epoch switching height %d",
				myHeader.Header.Height, info.Height)
			continue
		}
		if err = VerifyCosmosHeader(&myHeader, info); err != nil {
			return fmt.Errorf("SyncBlockHeader, failed to verify header: %v", err)
		}
		info.NextValidatorsHash = myHeader.Header.NextValidatorsHash
		info.Height = uint64(myHeader.Header.Height)
		info.BlockHash = myHeader.Header.Hash()
		cnt++
	}
	if cnt == 0 {
		return fmt.Errorf("no header you commited is useful")
	}
	if err := PutEpochSwitchInfo(native, params.ChainID, info); err != nil {
		return fmt.Errorf("SyncBlockHeader, failed to PutEpochSwitchInfo: %v", err)
	}
	return nil
}

func (this *CosmosHandler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}
