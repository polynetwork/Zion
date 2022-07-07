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
	vals, err := s.getValidatorsByHeader(header, nil, s.chain)
	if err != nil {
		return nil
	}
	return vals
}

//func (s *backend) CurrentEpoch() (uint64, []common.Address, error) {
//	statedb, err := s.chain.State()
//	if err != nil {
//		return 0, nil, err
//	}
//
//	current := s.chain.CurrentHeader()
//	epoch, err := nm.GetCurrentEpochInfoFromDB(statedb)
//	if err != nil {
//		return 0, nil, err
//	}
//	height := current.Number.Uint64()
//	start := epoch.StartHeight.Uint64()
//
//	// if consensus change epoch not finished, just read last epoch as validator
//	if epoch.ID.Uint64() > nm.StartEpochID.Uint64() && height > 1 && (start == height || start-height == 1) {
//		if epoch, err = nm.GetEpochInfoFromDB(statedb, new(big.Int).Sub(epoch.ID, big.NewInt(1))); err != nil {
//			return 0, nil, err
//		}
//	}
//
//	return epoch.StartHeight.Uint64(), epoch.MemberList(), nil
//}

func (s *backend) IsSystemTransaction(tx *types.Transaction, header *types.Header) bool {
	if tx == nil || len(tx.Data()) < 4 {
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

	config, epoch, err := s.getGovernanceInfo(state)
	if err != nil {
		return err
	}

	start := epoch.StartHeight.Uint64()
	height := header.Number.Uint64()
	epochLength := config.BlockPerEpoch.Uint64()
	if height != start+epochLength-1 {
		return nil
	}

	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}
	if err := s.executeTransaction(ctx, contractAddr, payload); err != nil {
		return err
	}

	log.Info("EpochChange", "start", start, "current", height, "state", s.chain.CurrentHeader().Number.Uint64())
	return nil
}

func (s *backend) getGovernanceInfo(state *state.StateDB) (*nm.GlobalConfig, *nm.EpochInfo, error) {
	config, err := nm.GetGlobalConfigFromDB(state)
	if err != nil {
		return nil, nil, err
	}
	epoch, err := nm.GetCurrentEpochInfoFromDB(state)
	if err != nil {
		return nil, nil, err
	}
	return config, epoch, nil
}

func (s *backend) execEndBlock(ctx *systemTxContext) error {
	payload, err := new(nm.EndBlockParam).Encode()
	if err != nil {
		return err
	}
	return s.executeTransaction(ctx, contractAddr, payload)
}

func (s *backend) getValidatorsByHeader(header, parent *types.Header, chain consensus.ChainHeaderReader) (hotstuff.ValidatorSet, error) {
	var epoch *types.Header

	// todo(fuk): add LRU
	// todo logic error
	extra, err := types.ExtractHotstuffExtraPayload(header.Extra)
	if err != nil {
		return nil, err
	}
	log.Info("---x1", "height", header.Number, "extra", extra.Validators)

	if header.Number.Uint64() == 0 {
		return NewDefaultValSet(extra.Validators), nil
	}

	if extra.Height != header.Number.Uint64() {
		epoch = chain.GetHeaderByNumber(extra.Height)
	} else {
		if parent == nil {
			parent = chain.GetHeaderByHash(header.ParentHash)
		}
		if extra, err = types.ExtractHotstuffExtra(parent); err != nil {
			return nil, err
		} else {
			epoch = chain.GetHeaderByNumber(extra.Height)
		}
	}

	if extra, err = types.ExtractHotstuffExtraPayload(epoch.Extra); err != nil {
		return nil, err
	} else {
		return NewDefaultValSet(extra.Validators), nil
	}
}
