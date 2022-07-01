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
func (s *backend) CheckPoint(height uint64) (beforeChange, changing, resetValidators bool) {
	if height <= 1 {
		return
	} else if height == s.snaps.end() {
		beforeChange = true
	} else if height == s.snaps.start() {
		changing = true
	} else if height == s.snaps.start()+1 {
		resetValidators = true
	}
	return
}

func (s *backend) Validators(height uint64) hotstuff.ValidatorSet {
	epoch := s.snaps.get(height)
	return epoch.ValSet.Copy()
}

func (s *backend) ValidatorList(height uint64) []common.Address {
	return s.Validators(height).AddressList()
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
func (s *backend) applySnapshot(state *state.StateDB, height *big.Int, mining bool) error {
	govEpoch, err := s.getEpoch(state, height)
	if err != nil {
		return err
	}
	if start := govEpoch.StartHeight.Uint64(); start != s.snaps.nextStart() {
		return fmt.Errorf("expect start height %d, got %d", s.snaps.nextStart(), start)
	}
	// epoch id saved in chain data should be equal to epoch.maxId
	id := s.snaps.nextId()
	if govEpochID := govEpoch.ID.Uint64(); id != govEpochID {
		return fmt.Errorf("expect nextID to be %d, actual got %d", govEpochID, id)
	}

	config, err := s.getGlobalConfig(state, height)
	if err != nil {
		return err
	}

	snap := newSnapshot(id, govEpoch.StartHeight.Uint64(), config.BlockPerEpoch.Uint64(), govEpoch.MemberList())
	if !s.snaps.append(snap) {
		return fmt.Errorf("epoch already exist, %s", snap.String())
	}

	if err := snap.store(s.db); err != nil {
		return err
	}

	// todo(fuk): event should be sent by miner but not consensus.
	if mining {
		s.SendValidatorsChange(snap.ValSet.AddressList())
	}

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
	output, _, err := ref.NativeCall(caller, utils.NodeManagerContractAddress, payload)
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

func (s *backend) execEpochChange(ctx *systemTxContext) error {
	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}
	return s.executeSystemTx(ctx, contractAddr, payload)
}

func (s *backend) execEndBlock(ctx *systemTxContext) error {
	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}

	return s.executeSystemTx(ctx, contractAddr, payload)
}
