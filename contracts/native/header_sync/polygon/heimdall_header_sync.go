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
package polygon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	polygonTypes "github.com/ethereum/go-ethereum/contracts/native/header_sync/polygon/types"
	polygonCmn "github.com/ethereum/go-ethereum/contracts/native/header_sync/polygon/types/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type HeimdallHandler struct {
}

// NewHeimdallHandler ...
func NewHeimdallHandler() *HeimdallHandler {
	return &HeimdallHandler{}
}

type CosmosHeader struct {
	Header  polygonTypes.Header
	Commit  *polygonTypes.Commit
	Valsets []*polygonTypes.Validator
}

// SyncGenesisHeader ...
func (h *HeimdallHandler) SyncGenesisHeader(native *native.NativeContract) (err error) {
	ctx := native.ContractRef().CurrentContext()
	param := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncGenesisHeader, param, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract param deserialize error: %v", err)
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
	cdc := polygonTypes.NewCDC()
	var header CosmosHeader
	err = cdc.UnmarshalBinaryBare(param.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("HeimdallHandler SyncGenesisHeader: %s", err)
	}
	// check if has genesis header
	info, err := GetEpochSwitchInfo(native, param.ChainID)
	if err == nil && info != nil {
		return fmt.Errorf("HeimdallHandler SyncGenesisHeader, genesis header had been initialized")
	}
	if header.Header.Height < 0 {
		return fmt.Errorf("HeidallHandler SyncGenesisHeader, header height invalid")
	}
	if err := PutEpochSwitchInfo(native, param.ChainID, &CosmosEpochSwitchInfo{
		Height:             uint64(header.Header.Height),
		NextValidatorsHash: header.Header.NextValidatorsHash,
		ChainID:            header.Header.ChainID,
		BlockHash:          header.Header.Hash(),
	}); err != nil {
		return fmt.Errorf("HeimdallHandler SyncGenesisHeader, PutEpochSwitchInfo: %s", err)
	}
	return nil
}

func (h *HeimdallHandler) SyncBlockHeader(native *native.NativeContract) error {
	params := &hscommon.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}

	cdc := polygonTypes.NewCDC()
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
		if myHeader.Header.Height < 0 {
			return fmt.Errorf("SyncBlockHeader, header height invalid")
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

// SyncCrossChainMsg ...
func (h *HeimdallHandler) SyncCrossChainMsg(native *native.NativeContract) error {
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
		utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(hscommon.EPOCH_SWITCH), utils.GetUint64Bytes(chainId)),
		blob)
	return nil
}

type CosmosEpochSwitchInfo struct {
	// The height where validators set changed last time. Poly only accept
	// header and proof signed by new validators. That means the header
	// can not be lower than this height.
	Height uint64

	// Hash of the block at `Height`. Poly don't save the whole header.
	// So we can identify the content of this block by `BlockHash`.
	BlockHash polygonCmn.HexBytes

	// The hash of new validators set which used to verify validators set
	// committed with proof.
	NextValidatorsHash polygonCmn.HexBytes

	// The cosmos chain-id of this chain basing Cosmos-sdk.
	ChainID string
}

func (m *CosmosEpochSwitchInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID})
}

func (m *CosmosEpochSwitchInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height             uint64
		BlockHash          polygonCmn.HexBytes
		NextValidatorsHash polygonCmn.HexBytes
		ChainID            string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID = data.Height, data.BlockHash, data.NextValidatorsHash, data.ChainID
	return nil
}

func VerifySpan(native *native.NativeContract, heimdallPolyChainID uint64, proof *CosmosProof) (span *Span, err error) {
	info, err := GetEpochSwitchInfo(native, heimdallPolyChainID)
	if err != nil {
		err = fmt.Errorf("HeimdallHandler failed to get epoch switching height: %v", err)
		return
	}

	if err = VerifyCosmosHeader(&proof.Header, info); err != nil {
		return nil, fmt.Errorf("HeimdallHandler failed to verify cosmos header: %v", err)
	}

	if len(proof.Proof.Ops) != 2 {
		err = fmt.Errorf("proof size wrong")
		return
	}
	if !bytes.Equal(proof.Proof.Ops[1].Key, []byte("bor")) {
		err = fmt.Errorf("wrong module for proof")
		return
	}

	prt := rootmulti.DefaultProofRuntime()

	err = prt.VerifyValue(&proof.Proof, proof.Header.Header.AppHash, proof.Value.Kp, proof.Value.Value)
	if err != nil {
		err = fmt.Errorf("validateHeaderExtraField VerifyValue error: %s", err)
		return
	}

	heimdallSpan := &polygonTypes.HeimdallSpan{}
	err = polygonTypes.NewCDC().UnmarshalBinaryBare(proof.Value.Value, heimdallSpan)
	if err != nil {
		err = fmt.Errorf("validateHeaderExtraField heimdallSpan UnmarshalBinaryBare error: %s", err)
		return
	}

	span, err = SpanFromHeimdall(heimdallSpan)
	return
}

func VerifyCosmosHeader(myHeader *CosmosHeader, info *CosmosEpochSwitchInfo) error {
	// now verify this header
	valset := polygonTypes.NewValidatorSet(myHeader.Valsets)
	if !bytes.Equal(info.NextValidatorsHash, valset.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, block validator is not right, next validator hash: %s, "+
			"validator set hash: %s", info.NextValidatorsHash.String(), hex.EncodeToString(valset.Hash()))
	}
	if !bytes.Equal(myHeader.Header.ValidatorsHash, valset.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, block validator is not right!, header validator hash: %s, "+
			"validator set hash: %s", myHeader.Header.ValidatorsHash.String(), hex.EncodeToString(valset.Hash()))
	}
	if myHeader.Commit.Height() != myHeader.Header.Height {
		return fmt.Errorf("VerifyCosmosHeader, commit height is not right! commit height: %d, "+
			"header height: %d", myHeader.Commit.Height(), myHeader.Header.Height)
	}
	if !bytes.Equal(myHeader.Commit.BlockID.Hash, myHeader.Header.Hash()) {
		return fmt.Errorf("VerifyCosmosHeader, commit hash is not right!, commit block hash: %s,"+
			" header hash: %s", myHeader.Commit.BlockID.Hash.String(), hex.EncodeToString(valset.Hash()))
	}
	if err := myHeader.Commit.ValidateBasic(); err != nil {
		return fmt.Errorf("VerifyCosmosHeader, commit is not right! err: %s", err.Error())
	}
	if valset.Size() != myHeader.Commit.Size() {
		return fmt.Errorf("VerifyCosmosHeader, the size of precommits is not right!")
	}
	talliedVotingPower := int64(0)
	for _, commitSig := range myHeader.Commit.Precommits {
		if commitSig == nil {
			continue
		}
		idx := commitSig.ValidatorIndex
		_, val := valset.GetByIndex(idx)
		if val == nil {
			return fmt.Errorf("VerifyCosmosHeader, validator %d doesn't exist!", idx)
		}
		if commitSig.Type != polygonTypes.PrecommitType {
			return fmt.Errorf("VerifyCosmosHeader, commitSig.Type(%d) wrong", commitSig.Type)
		}
		// Validate signature.
		precommitSignBytes := myHeader.Commit.VoteSignBytes(info.ChainID, idx)
		if !val.PubKey.VerifyBytes(precommitSignBytes, commitSig.Signature) {
			return fmt.Errorf("VerifyCosmosHeader, Invalid commit -- invalid signature: %v", commitSig)
		}
		// Good precommit!
		if myHeader.Commit.BlockID.Equals(commitSig.BlockID) {
			talliedVotingPower += val.VotingPower
		}
	}
	if talliedVotingPower <= valset.TotalVotingPower()*2/3 {
		return fmt.Errorf("VerifyCosmosHeader, voteing power is not enough!")
	}

	return nil
}
