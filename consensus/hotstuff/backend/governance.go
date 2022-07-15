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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	nmabi "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// 不会有api接口调用verifyHeader这种接口导致epoch顺序错乱。epoch可以使用slice而非map。
// 此外，可以针对map可以使用具体的结构。[]*epoch在具体使用调用过程中理论上应该是完全逆序(大的epoch在前)添加的。
// 轻节点的epoch不包含endHeight(无法从链上读取)。
//
// 轻节点理论上只会按asc顺序同步节点，verifyHeader更新一组validators，但是不会添加新epoch对应的endHeight。
// 那么，在根据epoch读取validators验证时，应该避开endHeight的查询。其他的接口可以继续使用endHeight。

var (
	contractAddr = utils.NodeManagerContractAddress
	specMethod   = nmabi.GetSpecMethodID()
)

func (s *backend) FillHeader(state *state.StateDB, header *types.Header) error {
	epoch, err := nm.GetCurrentEpochInfoFromDB(state)
	if err != nil {
		return err
	}

	start := epoch.StartHeight.Uint64()
	height := header.Number.Uint64()
	if start == height {
		valset := NewDefaultValSet(epoch.MemberList())
		types.HotstuffHeaderFillWithValidators(header, valset.AddressList(), header.Number.Uint64())
		log.Info("CheckPoint fill header", "start", start, "current", height, "state", s.chain.CurrentHeader().Number, "next validators", valset.String())
	} else {
		types.HotstuffHeaderFillWithValidators(header, nil, start)
	}
	return nil
}

// CheckPoint get `epochInfo` and `globalConfig` from governance contract and judge 2 things:
// 1.whether the block height is the right number to set new validators in block header while mining.
// 2.whether the block height is the right number to change epoch.
// return the flags and save epoch in lru cache.
// 矿工打包新区块时填充区块并判断是否需要重启共识, 此时state高度落后header高度1个区块。
// todo: globalConfig.epochLength should at least 3 block.
func (s *backend) CheckPoint(height uint64) {
	if height <= 1 {
		return
	}

	state, err := s.chain.State()
	if err != nil {
		log.Warn("CheckPoint", "get state failed", err)
		return
	}
	epoch, err := nm.GetCurrentEpochInfoFromDB(state)
	if err != nil {
		log.Warn("CheckPoint", "get current epoch info, height", height, "err", err)
		return
	}
	start := epoch.StartHeight.Uint64()
	if height == start+1 {
		s.restart()
	}
}

func (s *backend) Validators(hash common.Hash, mining bool) hotstuff.ValidatorSet {
	if mining {
		return s.vals.Copy()
	}

	header := s.chain.GetHeaderByHash(hash)
	_, vals, err := s.getValidatorsByHeader(header, nil, s.chain)
	if err != nil {
		return nil
	}
	return vals
}

func (s *backend) IsSystemTransaction(tx *types.Transaction, header *types.Header) bool {
	// consider that tx is deploy transaction, so the tx.to will be nil
	if tx == nil || len(tx.Data()) < 4 || tx.To() == nil {
		return false
	}
	if *tx.To() != contractAddr {
		return false
	}
	id := common.Bytes2Hex(tx.Data()[:4])
	if _, exist := specMethod[id]; !exist {
		return false
	}
	return true
}

// header height infront of state height
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

	log.Info("Execute governance EpochChange", "end", end, "current", height, "state", s.chain.CurrentHeader().Number.Uint64())
	return nil
}

func (s *backend) getGovernanceInfo(state *state.StateDB) (*nm.EpochInfo, error) {
	epoch, err := nm.GetCurrentEpochInfoFromDB(state)
	if err != nil {
		return nil, err
	}
	return epoch, nil
}

func (s *backend) execEndBlock(ctx *systemTxContext) error {
	payload, err := new(nm.EndBlockParam).Encode()
	if err != nil {
		return err
	}
	return s.executeTransaction(ctx, contractAddr, payload)
}

func (s *backend) getValidatorsByHeader(header, parent *types.Header, chain consensus.ChainHeaderReader) (
	bool, hotstuff.ValidatorSet, error) {

	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return false, nil, err
	}

	if header.Number.Uint64() == 0 {
		return true, NewDefaultValSet(extra.Validators), nil
	}

	isEpoch := extra.Height == header.Number.Uint64()
	if isEpoch {
		if parent == nil {
			parent = chain.GetHeaderByHash(header.ParentHash)
		}
		if extra, err = types.ExtractHotstuffExtra(parent); err != nil {
			return isEpoch, nil, err
		}
	}

	epoch := s.getRecentHeader(extra.Height, chain)
	if epoch == nil {
		return isEpoch, nil, fmt.Errorf("header %d neither on chain nor in lru cache", extra.Height)
	}
	if extra, err = types.ExtractHotstuffExtra(epoch); err != nil {
		return isEpoch, nil, err
	}

	return isEpoch, NewDefaultValSet(extra.Validators), nil
}

func (s *backend) saveRecentHeader(header *types.Header) {
	s.recents.Add(header.Number.Uint64(), header)
}

// todo(fuk): translate 同步时批量获取区块头，但是不会直接将区块头落账，等到所有的区块body执行完之后才会落账。
// 所以，这里有可能在chain上拿不到区块头，这时候就要从lru cache获取区块头
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

// initValidators prepare validators for next block.
func (s *backend) initValidators() (err error) {
	var (
		header = s.chain.CurrentHeader()
		extra  *types.HotstuffExtra
	)

start:
	if extra, err = types.ExtractHotstuffExtra(header); err != nil {
		return
	}

	// the next block use parent extra.validators as valset
	if extra.Height == header.Number.Uint64() {
		s.vals = NewDefaultValSet(extra.Validators)
		return
	}

	header = s.chain.GetHeaderByNumber(extra.Height)
	goto start
}
