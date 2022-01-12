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
package okex

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/okex/ethsecp256k1"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	tbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

// Handler ...
type Handler struct {
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

// NewCDC ...
func NewCDC() *codec.Codec {
	cdc := codec.New()

	ethsecp256k1.RegisterCodec(cdc)
	return cdc
}

type CosmosHeader struct {
	Header  types.Header
	Commit  *types.Commit
	Valsets []*types.Validator
}

// SyncGenesisHeader ...
func (h *Handler) SyncGenesisHeader(native *native.NativeContract) (err error) {
	ctx := native.ContractRef().CurrentContext()
	param := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncGenesisHeader, param, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, hscommon.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}

	// get genesis header from input parameters
	cdc := NewCDC()
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
		Height:             header.Header.Height,
		NextValidatorsHash: header.Header.NextValidatorsHash,
		ChainID:            header.Header.ChainID,
		BlockHash:          header.Header.Hash(),
	}); err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, failed to PutEpochSwitchInfo: %v", err)
	}
	return nil
}

func (h *Handler) SyncBlockHeader(native *native.NativeContract) error {
	params := &hscommon.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}

	cdc := NewCDC()
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
		if info.Height >= myHeader.Header.Height {
			log.Debugf("SyncBlockHeader, height %d is lower or equal than epoch switching height %d",
				myHeader.Header.Height, info.Height)
			continue
		}
		if err = VerifyCosmosHeader(&myHeader, info); err != nil {
			return fmt.Errorf("SyncBlockHeader, failed to verify header: %v", err)
		}
		info.NextValidatorsHash = myHeader.Header.NextValidatorsHash
		info.Height = myHeader.Header.Height
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

// SyncCrossChainMsg ...
func (h *Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}

func GetEpochSwitchInfo(service *native.NativeContract, chainId uint64) (*CosmosEpochSwitchInfo, error) {
	raw, err := service.GetCacheDB().Get(
		utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(hscommon.EPOCH_SWITCH), utils.GetUint64Bytes(chainId)))
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch switching height: %v", err)
	}
	info := new(CosmosEpochSwitchInfo)
	if err = rlp.DecodeBytes(raw, info); err != nil {
		return nil, fmt.Errorf("failed to deserialize CosmosEpochSwitchInfo: %v", err)
	}
	return info, nil
}

func PutEpochSwitchInfo(service *native.NativeContract, chainId uint64, info *CosmosEpochSwitchInfo) error {
	blob, err := rlp.EncodeToBytes(info)
	if err != nil {
		return err
	}
	service.GetCacheDB().Put(
		utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(hscommon.EPOCH_SWITCH), utils.GetUint64Bytes(chainId)), blob)
	return notifyEpochSwitchInfo(service, chainId, info)
}

func notifyEpochSwitchInfo(native *native.NativeContract, chainID uint64, info *CosmosEpochSwitchInfo) error {
	return native.AddNotify(hscommon.ABI, []string{"OKEpochSwitchInfoEvent"}, chainID, info.BlockHash.String(), info.Height,
		info.NextValidatorsHash.String(), info.ChainID, native.ContractRef().BlockHeight())
}

type CosmosEpochSwitchInfo struct {
	// The height where validators set changed last time. Poly only accept
	// header and proof signed by new validators. That means the header
	// can not be lower than this height.
	Height int64

	// Hash of the block at `Height`. Poly don't save the whole header.
	// So we can identify the content of this block by `BlockHash`.
	BlockHash tbytes.HexBytes

	// The hash of new validators set which used to verify validators set
	// committed with proof.
	NextValidatorsHash tbytes.HexBytes

	// The cosmos chain-id of this chain basing Cosmos-sdk.
	ChainID string
}

func (m *CosmosEpochSwitchInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID})
}

func (m *CosmosEpochSwitchInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height             int64
		BlockHash          tbytes.HexBytes
		NextValidatorsHash tbytes.HexBytes
		ChainID            string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID = data.Height, data.BlockHash, data.NextValidatorsHash, data.ChainID
	return nil
}

func VerifyCosmosHeader(myHeader *CosmosHeader, info *CosmosEpochSwitchInfo) error {
	// now verify this header
	valset := types.NewValidatorSet(myHeader.Valsets)
	if !bytes.Equal(info.NextValidatorsHash, valset.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, block validator is not right, next validator hash: %s, "+
			"validator set hash: %s", info.NextValidatorsHash.String(), hex.EncodeToString(valset.Hash()))
	}
	if !bytes.Equal(myHeader.Header.ValidatorsHash, valset.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, block validator is not right!, header validator hash: %s, "+
			"validator set hash: %s", myHeader.Header.ValidatorsHash.String(), hex.EncodeToString(valset.Hash()))
	}
	if myHeader.Commit.GetHeight() != myHeader.Header.Height {
		return fmt.Errorf("VerifyCosmosHeader, commit height is not right! commit height: %d, "+
			"header height: %d", myHeader.Commit.GetHeight(), myHeader.Header.Height)
	}
	if !bytes.Equal(myHeader.Commit.BlockID.Hash, myHeader.Header.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, commit hash is not right!, commit block hash: %s,"+
			" header hash: %s", myHeader.Commit.BlockID.Hash.String(), hex.EncodeToString(valset.Hash()))
	}
	if err := myHeader.Commit.ValidateBasic(); err != nil {
		return fmt.Errorf("VerifyCosmosHeader, commit is not right! err: %s", err.Error())
	}
	if valset.Size() != len(myHeader.Commit.Signatures) {
		return fmt.Errorf("VerifyCosmosHeader, the size of precommits is not right!")
	}
	talliedVotingPower := int64(0)
	for idx, commitSig := range myHeader.Commit.Signatures {
		if commitSig.Absent() {
			continue // OK, some precommits can be missing.
		}
		_, val := valset.GetByIndex(idx)
		// Validate signature.
		precommitSignBytes := myHeader.Commit.VoteSignBytes(info.ChainID, idx)
		if !val.PubKey.VerifyBytes(precommitSignBytes, commitSig.Signature) {
			return fmt.Errorf("VerifyCosmosHeader, Invalid commit -- invalid signature: %v", commitSig)
		}
		// Good precommit!
		if myHeader.Commit.BlockID.Equals(commitSig.BlockID(myHeader.Commit.BlockID)) {
			talliedVotingPower += val.VotingPower
		}
	}
	if talliedVotingPower <= valset.TotalVotingPower()*2/3 {
		return fmt.Errorf("VerifyCosmosHeader, voteing power is not enough!")
	}

	return nil
}
