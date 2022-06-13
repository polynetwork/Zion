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

package node_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

var (
	gasTable = map[string]uint64{
		MethodName:             0,
		MethodPropose:          30000,
		MethodVote:             30000,
		MethodEpoch:            0,
		MethodGetEpochByID:     0,
		MethodProof:            0,
		MethodGetChangingEpoch: 0,
	}
)

const (
	// The minimum distance between two adjacent epochs is 60 blocks
	MinEpochValidPeriod uint64 = 60
	// The default value of distance for two adjacent epochs
	DefaultEpochValidPeriod uint64 = 86400
	// The max distance between two adjacent epochs
	MaxEpochValidPeriod uint64 = 86400 * 10
	// Consensus engine allows at least 4 validators, this denote F >= 1
	MinProposalPeersLen int = 4
	// Consensus engine allows at most 100 validators, this denote F <= 33
	MaxProposalPeersLen int = 100
	// Every validator can propose at most 6 proposals in an epoch
	MaxProposalNumPerEpoch int = 6
	// Proposal should be voted and passed in period
	MinVoteEffectivePeriod uint64 = 10
)

func InitNodeManager() {
	InitABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodCreateValidator, CreateValidator)
	s.Register(MethodUpdateValidator, UpdateValidator)
	s.Register(MethodStake, Stake)
	s.Register(MethodUnStake, UnStake)
}

func CreateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &CreateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodCreateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CreateValidator, unpack params error: %v", err)
	}

	// check pub key
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, decode pubkey error: %v", err)
	}
	pubkey, err := crypto.DecompressPubkey(dec)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, decompress pubkey error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pubkey)
	if addr == common.EmptyAddress {
		return nil, fmt.Errorf("invalid pubkey")
	}
	if addr == caller || addr == params.ProposalAddress {
		return nil, fmt.Errorf("stake, consensus and proposal address can not be duplicate")
	}

	// check commission
	globalConfig, err := GetGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, GetGlobalConfig error: %v", err)
	}
	if params.Commission.Sign() == -1 {
		return nil, fmt.Errorf("CreateValidator, commission must be positive")
	}
	if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
		return nil, fmt.Errorf("CreateValidator, commission can not greater than globalConfig.MaxCommission: %s",
			globalConfig.MaxCommission.String())
	}

	// check desc
	if uint64(len(params.Desc)) > globalConfig.MaxDescLength {
		return nil, fmt.Errorf("CreateValidator, desc length more than limit %d", globalConfig.MaxDescLength)
	}

	// check to see if the pubkey has been registered before
	_, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, GetValidator error: %v", err)
	}
	if found {
		return nil, fmt.Errorf("CreateValidator, validator already exist")
	}

	// check initial stake
	if globalConfig.MinInitialStake.Cmp(params.InitStake) == 1 {
		return nil, fmt.Errorf("CreateValidator, initial stake %s is less than min initial stake %s",
			params.InitStake.String(), globalConfig.MinInitialStake.String())
	}

	// store validator
	validator := &Validator{
		StakeAddress:    caller,
		ConsensusPubkey: params.ConsensusPubkey,
		ProposalAddress: params.ProposalAddress,
		Commission:      &Commission{Rate: params.Commission, UpdateHeight: height},
		Status:          Unlocked,
		Jailed:          false,
		UnlockTime:      common.Big0,
		TotalStake:      params.InitStake,
		SelfStake:       params.InitStake,
		Desc:            params.Desc,
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, setValidator error: %v", err)
	}
	// add validator to all validators pool
	err = addToAllValidators(s, validator.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, addToAllValidators error: %v", err)
	}

	// deposit native token
	err = deposit(s, caller, params.InitStake, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, deposit error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodCreateValidator}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UpdateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &UpdateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateValidator, unpack params error: %v", err)
	}

	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, decode pubkey error: %v", err)
	}
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, get validator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UpdateValidator, can not found record")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("UpdateValidator, stake address is not caller")
	}
	globalConfig, err := GetGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, GetGlobalConfig error: %v", err)
	}

	if params.ProposalAddress != common.EmptyAddress {
		validator.ProposalAddress = params.ProposalAddress
	}

	if params.Commission != common.Big0 {
		// check commission
		if params.Commission.Sign() == -1 {
			return nil, fmt.Errorf("UpdateValidator, commission must be positive")
		}
		if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
			return nil, fmt.Errorf("UpdateValidator, commission can not greater than globalConfig.MaxCommission: %s",
				globalConfig.MaxCommission.String())
		}
		validator.Commission = &Commission{Rate: params.Commission, UpdateHeight: height}
	}

	if params.Desc != "" {
		// check desc
		if uint64(len(params.Desc)) > globalConfig.MaxDescLength {
			return nil, fmt.Errorf("UpdateValidator, desc length more than limit %d", globalConfig.MaxDescLength)
		}
		validator.Desc = params.Desc
	}

	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodUpdateValidator}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Stake(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &StakeParam{}
	if err := utils.UnpackMethod(ABI, MethodStake, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Stake, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("Stake, decode pubkey error: %v", err)
	}

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("Stake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("Stake, validator is not exist")
	}

	// deposit native token
	err = deposit(s, caller, params.Amount, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, deposit error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UnStake(s *native.NativeContract) ([]byte, error) {

}

///////////////////
// Propose participant propose new `epoch change` schema
func Propose(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight().Uint64()
	proposer := s.ContractRef().TxOrigin()
	caller := ctx.Caller

	// check authority
	curEpoch, err := getCurrentEpoch(s)
	if err != nil {
		log.Trace("checkConsensusSign", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	if err := CheckAuthority(proposer, caller, curEpoch); err != nil {
		log.Trace("propose", "check authority failed", err, "tx origin", proposer.Hex())
		return utils.ByteFailed, ErrInvalidAuthority
	}

	// decode input
	input := new(MethodProposeInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("propose", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}

	peers := input.Peers
	startHeight := input.StartHeight
	// check peers, try to match all peer's public key and address
	if peers == nil || peers.List == nil || len(peers.List) == 0 {
		log.Trace("propose", "check peers", "peer list is nil")
		return utils.ByteFailed, ErrInvalidPeers
	}
	if len(peers.List) < MinProposalPeersLen || len(peers.List) > MaxProposalPeersLen {
		log.Trace("propose", "check peers number",
			fmt.Errorf("propose, peers length should be in range of [%d, %d]",
				MinProposalPeersLen, MaxProposalPeersLen))
		return utils.ByteFailed, ErrPeersNum
	}
	for _, peer := range peers.List {
		if err := checkPeer(peer); err != nil {
			log.Trace("propose", "check peer public key", "public key not match address")
			return utils.ByteFailed, ErrInvalidPubKey
		}
	}

	// check peers, number for proposal's peers should be at least 2/3 of old members
	if curEpoch.OldMemberNum(peers) < curEpoch.QuorumSize() {
		log.Trace("propose", "check old members", "proposal peers should be at least 2/3 old members")
		return utils.ByteFailed, ErrOldParticipantsNumber
	}

	// proposal start height should be in range of [height + minEpochValidPeriod, height + maxEpochValidPeriod]
	if startHeight > 0 {
		latestStartHeight := height + MinEpochValidPeriod
		farawayStartHeight := height + MaxEpochValidPeriod
		if startHeight < latestStartHeight || startHeight > farawayStartHeight {
			log.Trace("propose", "check start height", fmt.Errorf("propose, proposal start height should be in range of [%d,  %d]",
				latestStartHeight, farawayStartHeight))
			return utils.ByteFailed, ErrProposalStartHeight
		}
	} else {
		startHeight = height + DefaultEpochValidPeriod
	}

	// generate new epoch as proposal
	epochID := curEpoch.ID + 1
	sort.Sort(peers)
	epoch := &EpochInfo{
		ID:          epochID,
		Peers:       peers,
		StartHeight: startHeight,
		Proposer:    proposer,
		Status:      ProposalStatusPropose,
	}
	proposal := epoch.Hash()

	// check duplicate proposal and validator's proposals number
	if checkProposal(s, epochID, proposal) {
		log.Trace("propose", "check proposal hash, dump proposal", proposal.Hex())
		return utils.ByteFailed, ErrDuplicateProposal
	}
	if num := proposalsNum(s, epochID, proposer); num >= MaxProposalNumPerEpoch {
		log.Trace("propose", "check validator proposal number, expect < ", MaxProposalNumPerEpoch, "got", num)
		return utils.ByteFailed, ErrProposalsNum
	}

	if err := storeEpoch(s, epoch); err != nil {
		log.Trace("propose", "store epoch failed", err)
		return utils.ByteFailed, ErrStorage
	}
	if err := storeProposal(s, epoch.ID, epoch.Hash()); err != nil {
		log.Trace("propose", "store proposal hash failed", err)
		return utils.ByteFailed, ErrStorage
	}

	// vote to self proposal
	if err := storeVote(s, proposal, proposer); err != nil {
		log.Trace("propose", "store vote failed", err)
		return utils.ByteFailed, ErrStorage
	}
	storeVoteTo(s, epochID, proposer, proposal)

	// emit event log
	if err := emitEventProposed(s, epoch); err != nil {
		log.Trace("propose", "emit event log failed", err)
		return utils.ByteFailed, ErrEmitLog
	}

	log.Debug("propose", "validator send an proposal", proposer.Hex(), "epoch", epoch.String())
	return utils.ByteSuccess, nil
}

// Vote participants vote to proposal
func Vote(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	voter := s.ContractRef().TxOrigin()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight().Uint64()

	// check authority
	curEpoch, err := getCurrentEpoch(s)
	if err != nil {
		log.Trace("vote", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	if err := CheckAuthority(voter, caller, curEpoch); err != nil {
		log.Trace("vote", "check authority failed", err, "voter", voter.Hex())
		return utils.ByteFailed, ErrInvalidAuthority
	}

	// decode and check epoch info
	input := new(MethodVoteInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("vote", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}
	epochID := input.EpochID
	proposal := input.EpochHash

	if expectEpochID := curEpoch.ID + 1; epochID != expectEpochID {
		log.Trace("vote", "check epoch ID failed, expect", expectEpochID, "got", curEpoch.ID)
		return utils.ByteFailed, ErrInvalidInput
	}
	if !findProposal(s, epochID, proposal) {
		log.Trace("vote", "find proposal failed", proposal.Hex())
		return utils.ByteFailed, ErrProposalNotExist
	}
	epoch, err := getEpoch(s, proposal)
	if err != nil {
		log.Trace("vote", "get epoch failed", proposal.Hex())
		return utils.ByteFailed, ErrEpochNotExist
	}
	if epoch.Status == ProposalStatusPassed {
		log.Trace("vote", "epoch status err", "proposal already passed", "epoch", epoch.Hash().Hex(), "epoch ID", epoch.ID)
		return utils.ByteFailed, ErrProposalPassed
	}
	if epochID != epoch.ID {
		log.Trace("vote", "check epoch id failed, expect", epoch.ID, "got", epochID)
		return utils.ByteFailed, ErrInvalidEpoch
	}
	if proposal != epoch.Hash() {
		log.Trace("vote", "check epoch hash failed, expect", proposal.Hex(), "got", epoch.Hash().Hex())
		return utils.ByteFailed, ErrInvalidEpoch
	}

	// vote should be finished before start height
	if height+MinVoteEffectivePeriod >= epoch.StartHeight {
		log.Trace("vote", "too late to change epoch", "consensus need some time to restart")
		return utils.ByteFailed, ErrVoteHeight
	}

	// already reach quorum size
	sizeBeforeVote := voteSize(s, proposal)
	if sizeBeforeVote >= curEpoch.QuorumSize() {
		log.Trace("vote", "check size", "already reach quorum size", "num", sizeBeforeVote, "quorum size", curEpoch.QuorumSize())
		return utils.ByteSuccess, nil
	}

	// filter duplicate vote or delete old vote
	lastVote2 := findVoteTo(s, epochID, voter)
	if lastVote2 != common.EmptyHash {
		if lastVote2 == proposal {
			log.Trace("vote", "check vote", "duplicate vote", "proposal", proposal.Hex(), "vote", voter.Hex())
			return utils.ByteSuccess, nil
		}
		delVoteTo(s, epochID, voter)
		if err := deleteVote(s, proposal, voter); err != nil {
			log.Trace("vote", "delete last voted proposal failed", err, "proposal", proposal.Hex(), "vote", voter.Hex())
			return utils.ByteFailed, ErrStorage
		}
	}

	log.Debug("vote", "validator vote to proposal", epoch.Hash(), "voter", voter.Hex(), "epoch ID", epochID)
	// store vote
	storeVoteTo(s, input.EpochID, voter, proposal)
	if err := storeVote(s, proposal, voter); err != nil {
		log.Trace("vote", "store vote failed", err)
		return utils.ByteFailed, ErrStorage
	}

	sizeAfterVote := voteSize(s, proposal)
	groupSize := len(curEpoch.Members())
	if err := emitEventVoted(s, input.EpochID, proposal, sizeAfterVote, groupSize); err != nil {
		log.Trace("vote", "emit voted log failed", err)
		return utils.ByteFailed, ErrEmitLog
	}

	// change epoch point:
	// 1. update status and store current epoch
	// 2. store current epoch proof
	// 3. emit event log
	// 4. dirty job which used to clear all useless storage
	// 5. pub epoch change event to miner worker
	if sizeAfterVote == curEpoch.QuorumSize() {
		epoch.Status = ProposalStatusPassed
		if err := storeEpoch(s, epoch); err != nil {
			log.Trace("vote", "store passed epoch failed", err)
			return utils.ByteFailed, ErrStorage
		}

		storeCurrentEpochHash(s, epoch.Hash())
		storeEpochProof(s, epoch.ID, epoch.Hash())
		if err := emitEpochChange(s, curEpoch, epoch); err != nil {
			log.Trace("vote", "emit epoch change log failed", err)
			return utils.ByteFailed, ErrEmitLog
		}

		dirtyJob(s, curEpoch, epoch)
		log.Debug("vote", "proposal passed", epoch.Hash())
	}

	return utils.ByteSuccess, nil
}

// dirtyJob filter current epoch and clear storage of `epoch`, `proposal`, `vote`, `voteTo`
func dirtyJob(s *native.NativeContract, last, cur *EpochInfo) {
	proposals, _ := getProposals(s, cur.ID)
	for _, v := range proposals {
		if v == cur.Hash() {
			continue
		}

		delEpoch(s, v)
		if err := delProposal(s, cur.ID, v); err != nil {
			log.Error("dirtyJob", "dirty job failed", err)
		}

		clearVotes(s, v)
		if last != nil && last.Peers != nil && last.Peers.List != nil {
			for _, v := range last.Peers.List {
				delVoteTo(s, cur.ID, v.Address)
			}
		}
	}

	list, err := getProposals(s, cur.ID)
	if err != nil {
		log.Warn("dirtyJob", "check proposal number after dirty job, can't get proposals", err)
	}
	if len(list) < 0 || len(list) > 1 {
		log.Warn("dirtyJob", "check proposal number after dirty job, expect", 1, "got", len(list))
	}
}

// GetCurrentEpoch retrieve current effective epoch info
func GetCurrentEpoch(s *native.NativeContract) ([]byte, error) {
	epoch, err := getCurrentEpoch(s)
	if err != nil {
		log.Trace("epoch", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	output := &MethodEpochOutput{Epoch: epoch}
	return output.Encode()
}

func GetEpochWithStateDB(db *state.StateDB) (*EpochInfo, error) {
	ctx := generateEmptyContext(db)
	return getCurrentEpoch(ctx)
}

func GetEpochByHeight(db *state.StateDB, height uint64) (*EpochInfo, error) {
	ctx := generateEmptyContext(db)
	epoch, err := getCurrentEpoch(ctx)
	if err != nil {
		return nil, err
	}

	for height < epoch.StartHeight {
		if epoch, err = getEffectiveEpochByID(ctx, epoch.ID-1); err != nil {
			return nil, err
		}
	}

	return epoch, nil
}

func GetChangingEpoch(s *native.NativeContract) ([]byte, error) {
	curEpochHash, err := getCurrentEpochHash(s)
	if err != nil {
		return utils.ByteFailed, err
	}
	epoch, err := getEpoch(s, curEpochHash)
	if err != nil {
		return utils.ByteFailed, err
	}

	height := s.ContractRef().BlockHeight().Uint64()
	if height > epoch.StartHeight {
		return utils.ByteFailed, fmt.Errorf("epoch invalid")
	}
	output := &MethodEpochOutput{Epoch: epoch}
	return output.Encode()
}

// GetEpochByID retrieve history effective epoch with epochID
func GetEpochByID(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodGetEpochByIDInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("getEpochByID", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}

	epoch, err := getEffectiveEpochByID(s, input.EpochID)
	if err != nil {
		log.Trace("getEpochByID", "get history epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	output := MethodEpochOutput{Epoch: epoch}
	return output.Encode()
}

func GetEpochProof(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodProofInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("proof", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}

	proof, err := getEpochProof(s, input.EpochID)
	if err != nil {
		log.Trace("proof", "get current epoch proof failed", err)
		return utils.ByteFailed, ErrEpochProofNotExist
	}
	output := &MethodProofOutput{Hash: proof}
	return output.Encode()
}

func CheckConsensusSigns(s *native.NativeContract, method string, input []byte, signer common.Address) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	log.Trace("checkConsensusSign", "method", method, "input", hexutil.Encode(input), "signer", signer.Hex())

	// get epoch info
	epoch, err := getCurrentEpoch(s)
	if err != nil {
		log.Trace("checkConsensusSign", "get current epoch failed", err)
		return false, ErrEpochNotExist
	}

	// check authority
	if err := CheckAuthority(signer, caller, epoch); err != nil {
		log.Trace("checkConsensusSign", "check authority failed", err)
		return false, ErrInvalidAuthority
	}

	// get or set consensus sign info
	sign := &ConsensusSign{Method: method, Input: input}
	if exist, err := getSign(s, sign.Hash()); err != nil {
		if err.Error() == "EOF" {
			if err := storeSign(s, sign); err != nil {
				log.Trace("checkConsensusSign", "store sign failed", err, "hash", sign.Hash().Hex())
				return false, ErrStorage
			} else {
				log.Trace("checkConsensusSign", "store sign, hash", sign.Hash().Hex())
			}
		} else {
			log.Trace("checkConsensusSign", "get sign failed", err, "hash", sign.Hash().Hex())
			return false, ErrConsensusSignNotExist
		}
	} else if exist.Hash() != sign.Hash() {
		log.Trace("checkConsensusSign", "check sign hash failed, expect", exist.Hash().Hex(), "got", sign.Hash().Hex())
		return false, ErrInvalidSign
	}

	// check duplicate signature
	if findSigner(s, sign.Hash(), signer) {
		log.Trace("checkConsensusSign", "signer already exist", signer.Hex(), "hash", sign.Hash().Hex())
		return false, ErrDuplicateSigner
	}

	// do not store redundancy sign
	sizeBeforeSign := getSignerSize(s, sign.Hash())
	log.Trace("checkConsensusSign", "sign hash", sign.Hash().Hex(), "size before sign", sizeBeforeSign)
	if sizeBeforeSign >= epoch.QuorumSize() {
		return false, nil
	}

	// store signer address and emit event log
	if err := storeSigner(s, sign.Hash(), signer); err != nil {
		log.Trace("checkConsensusSign", "store signer failed", err, "hash", sign.Hash().Hex())
		return false, ErrStorage
	}
	sizeAfterSign := getSignerSize(s, sign.Hash())
	if err := emitConsensusSign(s, sign, signer, sizeAfterSign); err != nil {
		log.Trace("checkConsensusSign", "emit consensus sign log failed", err, "hash", sign.Hash().Hex())
		return false, ErrEmitLog
	}
	log.Trace("checkConsensusSign", "sign hash", sign.Hash().Hex(), "size after sign", sizeAfterSign)

	return sizeAfterSign >= epoch.QuorumSize(), nil
}

func EpochChangeAtNextBlock(curHeight, epochStartHeight uint64) bool {
	if curHeight+1 == epochStartHeight {
		return true
	}
	return false
}
