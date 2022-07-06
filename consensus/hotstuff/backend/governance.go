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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native"
	nmabi "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
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

// CheckPoint get `epochInfo` and `globalConfig` from governance contract and judge 2 things:
// 1.whether the block height is the right number to set new validators in block header while mining.
// 2.whether the block height is the right number to change epoch.
// return the flags and save epoch in lru cache.
func (s *backend) CheckPoint(state *state.StateDB, header *types.Header) (consensus.CheckPointStatus, uint64, error) {
	status := consensus.CheckPointStateUnknown
	if header.Number.Uint64() <= 1 {
		return status, 0, nil
	}

	// todo(fuk): 随时生效还是周期末生效
	globalConfig, err := s.getGlobalConfig(state, header.Number)
	if err != nil {
		return status, 0, err
	}

	start := s.snaps.start()
	height := header.Number.Uint64()
	epochLength := globalConfig.BlockPerEpoch.Uint64()
	if height == start && s.snaps.lastStart() != 0 {
		start = s.snaps.lastStart()
	}

	if height == start+epochLength-1 {
		status = consensus.CheckPointStatePrepare
	} else if height == start+epochLength { // apply new snapshot and start height is new.
		status = consensus.CheckPointStateStart
	} else if height == start+1 {
		status = consensus.CheckPointStateChange
	}
	return status, 0, nil
}

func (s *backend) Validators(height uint64) (vals hotstuff.ValidatorSet) {
	header := s.chain.GetHeaderByNumber(height)
	if header == nil {
		header = s.chain.GetHeaderByNumber(height - 1)
	}
	for vals == nil {
		extra, err := types.ExtractHotstuffExtraPayload(header.Extra)
		if err != nil {
			return nil
		}
		epochStartHeight := extra.Height
		s.chain.GetHeaderByNumber(epochStartHeight)
		vals = NewDefaultValSet(extra.Validators)
	}
	return
}

func (s *backend) ValidatorList(state *state.StateDB, height *big.Int) ([]common.Address, error) {
	epoch, err := s.getEpoch(state, height)
	if err != nil {
		return nil, err
	}
	return epoch.MemberList(), nil
}

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

// miner checkpoint，在epoch start添加新的validators
// 查询并存储snapshot
func (s *backend) applySnapshot(state *state.StateDB, header *types.Header, mining bool) error {
	height := header.Number
	govEpoch, err := s.getEpoch(state, height)
	if err != nil {
		return err
	}

	if govEpoch.StartHeight.Cmp(height) <= 0 {
		return fmt.Errorf("invalid height, expect %v, got %v", height.Uint64()+1, govEpoch.StartHeight)
	}

	// epoch id saved in chain data should be equal to epoch.maxId
	id := s.snaps.nextId()
	if govEpochID := govEpoch.ID.Uint64(); id != govEpochID {
		return fmt.Errorf("expect nextID to be %d, actual got %d", govEpochID, id)
	}

	snap := newSnapshot(id, govEpoch.StartHeight.Uint64(), govEpoch.MemberList())
	if !s.snaps.append(snap) {
		return fmt.Errorf("epoch already exist, %s", snap.String())
	}

	if err := snap.store(s.db); err != nil {
		return err
	}

	// todo(fuk): event should be sent by miner but not consensus.
	//if mining {
	//	s.SendValidatorsChange(snap.ValSet.AddressList())
	//}

	log.Info("[epoch]", "check point", snap.String())
	return nil
}

func (s *backend) getGlobalConfig(state *state.StateDB, height *big.Int) (*nm.GlobalConfig, error) {

	caller := s.signer.Address()
	ref := native.NewContractRef(state, caller, caller, height, common.EmptyHash, 0, nil)
	payload, err := new(nm.GetGlobalConfigParam).Encode()
	if err != nil {
		return nil, fmt.Errorf("encode GetGlobalConfig input failed: %v", err)
	}
	output, _, err := ref.NativeCall(caller, contractAddr, payload)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig native call failed: %v", err)
	}

	var (
		raw    []byte
		config = new(nm.GlobalConfig)
	)
	if err := utils.UnpackOutputs(nm.ABI, nmabi.MethodGetGlobalConfig, &raw, output); err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(raw, config); err != nil {
		return nil, err
	}

	log.Info("getGlobalConfig", "config.blockPerEpoch", config.BlockPerEpoch,
		"config.ConsensusValidatorNum", config.ConsensusValidatorNum)

	return config, nil
}

func (s *backend) getEpoch(state *state.StateDB, height *big.Int) (*nm.EpochInfo, error) {
	caller := s.signer.Address()
	ref := native.NewContractRef(state, caller, caller, height, common.EmptyHash, 0, nil)
	payload, err := new(nm.GetCurrentEpochInfoParam).Encode()
	if err != nil {
		return nil, fmt.Errorf("encode GetGlobalConfig input failed: %v", err)
	}
	output, _, err := ref.NativeCall(caller, contractAddr, payload)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig native call failed: %v", err)
	}

	var (
		raw   []byte
		epoch = new(nm.EpochInfo)
	)
	if err := utils.UnpackOutputs(nm.ABI, nmabi.MethodGetCurrentEpochInfo, &raw, output); err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(raw, epoch); err != nil {
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

func (s *backend) execEpochChange(ctx *systemTxContext) error {
	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}
	return s.executeTransaction(ctx, contractAddr, payload)
}
