// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package native_client

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/maas_config"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

var ErrAccountBlocked = errors.New("account is in blacklist")
var ErrNotGasManager = errors.New("address is not in gas manager list")

func IsBlocked(state *state.StateDB, address *common.Address) bool {
	log.Debug("### isBlocked called")
	if address == nil {
		return false
	}
	caller := common.EmptyAddress
	ref := native.NewContractRef(state, caller, caller, big.NewInt(-1), common.EmptyHash, 0, nil)

	payload, err := (&maas_config.MethodIsBlockedInput{Addr: *address}).Encode()
	if err != nil {
		log.Error("[PackMethod]", "pack `isBlocked` input failed", err)
		return false
	}
	enc, _, err := ref.NativeCall(caller, utils.MaasConfigContractAddress, payload)
	if err != nil {
		return false
	}
	output := new(maas_config.MethodBoolOutput)
	if err := output.Decode(enc, maas_config.MethodIsBlocked); err != nil {
		log.Error("[native call]", "unpack `IsBlocked` output failed", err)
		return false
	}

	log.Debug("IsBlocked: " + address.String() + ", " + strconv.FormatBool(output.Success))
	return output.Success
}

func IsGasManageEnable(state *state.StateDB) bool {
	caller := common.EmptyAddress
	ref := native.NewContractRef(state, caller, caller, big.NewInt(-1), common.EmptyHash, 0, nil)

	payload, err := utils.PackMethod(maas_config.ABI, maas_config.MethodIsGasManageEnabled)
	if err != nil {
		log.Error("[PackMethod]", "pack `IsGasManageEnable` input failed", err)
		return false
	}
	enc, _, err := ref.NativeCall(caller, utils.MaasConfigContractAddress, payload)
	if err != nil {
		return false
	}
	output := new(maas_config.MethodBoolOutput)
	if err := output.Decode(enc, maas_config.MethodIsGasManageEnabled); err != nil {
		log.Error("[native call]", "unpack `IsGasManageEnable` output failed", err)
		return false
	}

	return output.Success
}

func IsGasManager(state *state.StateDB, address *common.Address) bool {
	if address == nil {
		return false
	}
	caller := common.EmptyAddress
	ref := native.NewContractRef(state, caller, caller, big.NewInt(-1), common.EmptyHash, 0, nil)

	payload, err := (&maas_config.MethodIsGasManagerInput{Addr: *address}).Encode()
	if err != nil {
		log.Error("[PackMethod]", "pack `IsGasManager` input failed", err)
		return false
	}
	enc, _, err := ref.NativeCall(caller, utils.MaasConfigContractAddress, payload)
	if err != nil {
		return false
	}
	output := new(maas_config.MethodBoolOutput)
	if err := output.Decode(enc, maas_config.MethodIsGasManager); err != nil {
		log.Error("[native call]", "unpack `IsGasManager` output failed", err)
		return false
	}

	return output.Success
}
