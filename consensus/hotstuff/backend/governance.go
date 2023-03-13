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

package backend

import (
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

var (
	contractAddr = utils.NodeManagerContractAddress
	specMethod   = nm.GetSpecMethodID()
)

// FillHeader fulfill the header with validators for miner worker. there are 2 conditions:
// * governance epoch changed on chain, use the new validators for an new epoch start header.
// * governance epoch not changed, only set the old epoch start height in header.
func (s *backend) FillHeader(state *state.StateDB, header *types.Header) error {
	epoch, err := s.getGovernanceInfo(state)
	if err != nil {
		return err
	}

	start := epoch.StartHeight.Uint64()
	end := epoch.EndHeight.Uint64()
	height := header.Number.Uint64()
	if start == height {
		valset := NewDefaultValSet(epoch.MemberList())
		types.HotstuffHeaderFillWithValidators(header, valset.AddressList(), header.Number.Uint64(), end)
		log.Info("CheckPoint fill header", "start", start, "end", end, "current", height, "next validators", valset.String())
	} else {
		types.HotstuffHeaderFillWithValidators(header, nil, start, end)
	}
	return nil
}

// CheckPoint get `epochInfo` and `globalConfig` from governance contract and judge 2 things:
// 1.whether the block height is the right number to set new validators in block header while mining.
// 2.whether the block height is the right number to change epoch.
// return the flags and save epoch in lru cache.
func (s *backend) CheckPoint(height uint64) (uint64, bool) {
	if height <= 1 {
		return 0, false
	}

	state, err := s.chain.State()
	if err != nil {
		log.Warn("CheckPoint", "get state failed", err)
		return 0, false
	}
	epoch, err := s.getGovernanceInfo(state)
	if err != nil {
		log.Warn("CheckPoint", "get current epoch info, height", height, "err", err)
		return 0, false
	}

	// lock status to ensure that the action of restart engine wont recall CheckPoint twice.
	// and the action of lock should happen before restart.
	start := epoch.StartHeight.Uint64()
	if s.epochMu == 0 && height == start+1 {
		log.Trace("CheckPoint lock status", "height", height, "epoch start", start)
		atomic.StoreInt32(&s.epochMu, 1)
		return start, true
	}
	if s.epochMu == 1 && height > start+1 {
		log.Trace("CheckPoint unlock status", "height", height, "epoch start", start)
		atomic.StoreInt32(&s.epochMu, 0)
		return start, false
	}
	return start, false
}

// Validators get validators from backend by `consensus core`, param of `mining` is false denote need last epoch validators.
func (s *backend) Validators(height uint64, mining bool) hotstuff.ValidatorSet {
	if mining {
		return s.vals.Copy()
	}

	header := s.chain.GetHeaderByNumber(height)
	if header == nil {
		return nil
	}
	_, vals, err := s.getValidatorsByHeader(header, nil, s.chain)
	if err != nil {
		return nil
	}
	return vals
}

// IsSystemTransaction used by state processor while sync block.
func (s *backend) IsSystemTransaction(tx *types.Transaction, header *types.Header) (string, bool) {
	// consider that tx is deploy transaction, so the tx.to will be nil
	if tx == nil || len(tx.Data()) < 4 || tx.To() == nil {
		return "", false
	}
	if *tx.To() != contractAddr {
		return "", false
	}
	id := common.Bytes2Hex(tx.Data()[:4])
	if _, exist := specMethod[id]; !exist {
		return id, false
	}

	signer := types.MakeSigner(s.chain.Config(), header.Number)
	addr, err := signer.Sender(tx)
	if err != nil {
		return id, false
	}
	if header.Coinbase != addr {
		return id, false
	} else {
		return id, true
	}
}

// header height in front of state height
func (s *backend) execEpochChange(state *state.StateDB, header *types.Header, ctx *systemTxContext) error {

	epoch, err := s.getGovernanceInfo(state)
	if err != nil {
		return err
	}

	end := epoch.EndHeight.Uint64()
	height := header.Number.Uint64()
	if height != end-1 {
		return nil
	}

	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}
	if err := s.executeTransaction(ctx, contractAddr, payload); err != nil {
		return err
	}

	log.Info("Execute governance EpochChange", "end", end, "current", height)
	return nil
}

// getGovernanceInfo call governance contract method and retrieve related info.
func (s *backend) getGovernanceInfo(state *state.StateDB) (*nm.EpochInfo, error) {
	epoch, err := nm.GetCurrentEpochInfoFromDB(state)
	if err != nil {
		return nil, err
	}
	return epoch, nil
}

// execEndBlock execute governance contract method of `EndBlock`
func (s *backend) execEndBlock(ctx *systemTxContext) error {
	payload, err := new(nm.EndBlockParam).Encode()
	if err != nil {
		return err
	}
	return s.executeTransaction(ctx, contractAddr, payload)
}

// getValidatorsByHeader check if current header height is an new epoch start and retrieve the validators.
func (s *backend) getValidatorsByHeader(header, parent *types.Header, chain consensus.ChainHeaderReader) (
	bool, hotstuff.ValidatorSet, error) {

	// extract current header
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return false, nil, err
	}

	// the genesis block is an epoch start, and the validators stored in the field of `header.extra`
	if header.Number.Uint64() == 0 {
		return true, NewDefaultValSet(extra.Validators), nil
	}

	// if the block height equals to the `extra.height`, this block is an epoch start.
	// the the validators for this header is stored in last epoch start header.
	isEpoch := extra.StartHeight == header.Number.Uint64() && len(extra.Validators) > 0
	if isEpoch {
		if parent == nil {
			parent = chain.GetHeaderByHash(header.ParentHash)
		} else if parent.Hash() != header.ParentHash || parent.Number.Uint64()+1 != header.Number.Uint64() {
			return false, nil, consensus.ErrUnknownAncestor
		}
		if extra, err = types.ExtractHotstuffExtra(parent); err != nil {
			return isEpoch, nil, err
		}
	}

	epochHeader := s.getRecentHeader(extra.StartHeight, chain)
	if epochHeader == nil {
		return isEpoch, nil, fmt.Errorf("header %d neither on chain nor in lru cache", extra.StartHeight)
	}
	if extra, err = types.ExtractHotstuffExtra(epochHeader); err != nil {
		return isEpoch, nil, err
	}
	if extra.Validators == nil || len(extra.Validators) == 0 {
		return isEpoch, nil, fmt.Errorf("invalid epoch start header")
	}
	return isEpoch, NewDefaultValSet(extra.Validators), nil
}

// getRecentHeader in block sync module, the block headers are fetched in batches, and these headers will store in
// chain db until all block body and block receipts executing finished. but the action of `verifyHeader` is continuously.
// so save all epoch header in LRU is needed.
func (s *backend) getRecentHeader(height uint64, chain consensus.ChainHeaderReader) *types.Header {
	header := chain.GetHeaderByNumber(height)
	if header != nil {
		return header
	}

	data, ok := s.recents.Get(height)
	if !ok {
		return nil
	}

	header, ok = data.(*types.Header)
	if !ok {
		return nil
	}
	return header
}

func (s *backend) saveRecentHeader(header *types.Header) {
	s.recents.Add(header.Number.Uint64(), header)
}

// newEpochValidators prepare validators for next block.
func (s *backend) newEpochValidators() (vs hotstuff.ValidatorSet, err error) {
	var (
		header = s.chain.CurrentHeader()
		extra  *types.HotstuffExtra
	)

start:
	if extra, err = types.ExtractHotstuffExtra(header); err != nil {
		return
	}

	// the next block use parent extra.validators as valset
	if extra.StartHeight == header.Number.Uint64() {
		vs = NewDefaultValSet(extra.Validators)
		return
	}

	header = s.chain.GetHeaderByNumber(extra.StartHeight)
	goto start
}
