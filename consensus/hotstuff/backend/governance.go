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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
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
	contractAddr       = utils.NodeManagerContractAddress
	startEpochID       = nm.StartEpochID.Uint64()
	genesisEpochStart  = uint64(0)
	genesisEpochLength = nm.GenesisBlockPerEpoch.Uint64()
)

func init() {
	core.StoreGenesis = func(db ethdb.Database, header *types.Header) error {
		extra, err := types.ExtractHotstuffExtra(header)
		if err != nil {
			return err
		}
		epoch := &snapshot{
			ID:     startEpochID,
			Start:  genesisEpochStart,
			ValSet: NewDefaultValSet(extra.Validators),
		}
		return epoch.store(db)
	}
}

// CheckPoint get `epochInfo` and `globalConfig` from governance contract and judge 2 things:
// 1.whether the block height is the right number to set new validators in block header while mining.
// 2.whether the block height is the right number to change epoch.
// return the flags and save epoch in lru cache.
func (s *backend) CheckPoint(height uint64) (beforeChange, changing, resetValidators bool) {
	if height <= 1 {
		return
	} else if height == s.maxEnd() {
		beforeChange = true
	} else if height == s.nextStart() {
		changing = true
	} else if height == s.maxStart() + 1 {
		resetValidators = true
	}
	return
}

// miner checkpoint，在epoch start添加新的validators

func (s *backend) SavePoint(state *state.StateDB, height *big.Int, mining bool) error {
	govEpoch, err := s.getEpoch(state, height)
	if err != nil {
		return err
	}
	if start := govEpoch.StartHeight.Uint64(); start != s.nextStart() {
		return fmt.Errorf("expect start height %d, got %d", s.nextStart(), start)
	}
	// epoch id saved in chain data should be equal to epoch.maxId
	id := s.nextId()
	if govEpochID := govEpoch.ID.Uint64(); id != govEpochID {
		return fmt.Errorf("expect nextID to be %d, actual got %d", govEpochID, id)
	}

	config, err := s.getGlobalConfig(state, height)
	if err != nil {
		return err
	}

	epoch := &snapshot{
		ID:     id,
		Start:  govEpoch.StartHeight.Uint64(),
		End:    govEpoch.StartHeight.Uint64() + config.BlockPerEpoch.Uint64(),
		ValSet: NewDefaultValSet(govEpoch.MemberList()),
	}
	if !s.appendEpoch(epoch) {
		return fmt.Errorf("epoch already exist, %s", epoch.String())
	}

	if err := epoch.store(s.db); err != nil {
		return err
	}

	// todo(fuk): event should be sent by miner but not consensus.
	if mining {
		s.SendValidatorsChange(epoch.ValSet.AddressList())
	}

	log.Info("[epoch]", "check point", epoch.String())
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

// LoadEpochs read epoch from database and append in slice cache while consensus engine started.
func (s *backend) LoadEpochs() {
	id := startEpochID

	for {
		epoch := new(snapshot)
		if err := epoch.load(s.db, id); err != nil {
			return
		}
		s.appendEpoch(epoch)
		id = s.nextId()
		log.Info("[epoch]", "load epoch", epoch.String())
	}
}

func (s *backend) GetEpoch(height uint64) *snapshot {
	for _, epoch := range s.epochs {
		if height >= epoch.Start {
			return epoch
		}
	}
	// `height` is an uint number, and the min epoch height is 0, so the function will return in above loop
	return nil
}

func (s *backend) Validators(height uint64) hotstuff.ValidatorSet {
	epoch := s.GetEpoch(height)
	return epoch.ValSet.Copy()
}

func (s *backend) ValidatorList(height uint64) []common.Address {
	return s.Validators(height).AddressList()
}

// SyncEpoch light mode
func (s *backend) SyncEpoch(parent, header *types.Header) error {
	if parent.Number.Uint64() == 0 {
		return nil
	}

	parentExt, err := types.ExtractHotstuffExtra(parent)
	if err != nil {
		return err
	}
	if parentExt.Validators == nil || len(parentExt.Validators) == 0 {
		return nil
	}

	epoch := &snapshot{
		ID:     s.nextId(),
		Start:  header.Number.Uint64(),
		ValSet: NewDefaultValSet(parentExt.Validators),
	}

	if s.appendEpoch(epoch) {
		if err := epoch.store(s.db); err != nil {
			return err
		}
		log.Info("[epoch]", "sync epoch", epoch.String())
	}

	return nil
}

func (s *backend) DumpEpochs() string {
	str := ""
	for _, v := range s.epochs {
		str += v.String() + "\r\n"
	}
	return str
}

// -----------------------------------------------------------------
// epoch store and read functions
// -----------------------------------------------------------------

// s.epochs is an desc list
func (s *backend) appendEpoch(epoch *snapshot) bool {
	s.epochMu.Lock()
	defer s.epochMu.Unlock()

	if len(s.epochs) == 0 {
		s.epochs = append(s.epochs, epoch)
		return true
	}

	// already exist
	if epoch.ID <= s.maxID() || epoch.Start <= s.maxStart() {
		return false
	}

	s.epochs = append(s.epochs, epoch)
	s.epochs[0], s.epochs[len(s.epochs)-1] = s.epochs[len(s.epochs)-1], s.epochs[0]

	return true
}

func (s *backend) maxID() uint64 {
	if s.epochs == nil {
		return startEpochID
	}
	return s.epochs[0].ID
}

func (s *backend) maxStart() uint64 {
	if s.epochs == nil {
		return genesisEpochStart
	}
	return s.epochs[0].Start
}

func (s *backend) maxEnd() uint64 {
	if s.epochs == nil {
		return genesisEpochLength
	}
	return s.epochs[0].End
}

func (s *backend) nextStart() uint64 {
	if s.epochs == nil {
		return genesisEpochStart
	}
	return s.epochs[0].End + 1
}

func (s *backend) nextId() uint64 {
	if s.epochs == nil {
		return startEpochID
	}
	return s.maxID() + 1
}
