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

package maas_config

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "maas config"

const (

	// abi
	MaasConfigABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"BlockAccount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ChangeOwner\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"blockAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlacklist\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isBlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

	// method name
	MethodBlockAccount = "blockAccount"
	MethodChangeOwner  = "changeOwner"
	MethodGetOwner     = "getOwner"
	MethodGetBlacklist = "getBlacklist"
	MethodIsBlocked    = "isBlocked"
	MethodName         = "name"
	EventBlockAccount  = "BlockAccount"
	EventChangeOwner   = "ChangeOwner"

	// key prefix
	BLACKLIST = "blacklist"
	OWNER     = "owner"
)

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(MaasConfigABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.MaasConfigContractAddress
)

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
}

type MethodChangeOwnerInput struct {
	Addr common.Address
}

func (m *MethodChangeOwnerInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodChangeOwner, m.Addr)
}

func (m *MethodChangeOwnerInput) Decode(payload []byte) error {
	var data struct {
		Addr common.Address
	}
	if err := utils.UnpackMethod(ABI, MethodChangeOwner, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	return nil
}

type MethodChangeOwnerOutput struct {
	Success bool
}

func (m *MethodChangeOwnerOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodChangeOwner, m.Success)
}

func (m *MethodChangeOwnerOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodChangeOwner, m, payload)
}

type MethodGetOwnerOutput struct {
	Addr common.Address
}

func (m *MethodGetOwnerOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodGetOwner, m.Addr)
}

func (m *MethodGetOwnerOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodGetOwner, m, payload)
}

type MethodBlockAccountInput struct {
	Addr    common.Address
	DoBlock bool
}

func (m *MethodBlockAccountInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodBlockAccount, m.Addr, m.DoBlock)
}

func (m *MethodBlockAccountInput) Decode(payload []byte) error {
	var data struct {
		Addr    common.Address
		DoBlock bool
	}
	if err := utils.UnpackMethod(ABI, MethodBlockAccount, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	m.DoBlock = data.DoBlock
	return nil
}

type MethodBlockAccountOutput struct {
	Success bool
}

func (m *MethodBlockAccountOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodBlockAccount, m.Success)
}

func (m *MethodBlockAccountOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodBlockAccount, m, payload)
}

type MethodIsBlockedInput struct {
	Addr common.Address
}

func (m *MethodIsBlockedInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodIsBlocked, m.Addr)
}

func (m *MethodIsBlockedInput) Decode(payload []byte) error {
	var data struct {
		Addr common.Address
	}
	if err := utils.UnpackMethod(ABI, MethodIsBlocked, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	return nil
}

type MethodIsBlockedOutput struct {
	Success bool
}

func (m *MethodIsBlockedOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodIsBlocked, m.Success)
}

func (m *MethodIsBlockedOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodIsBlocked, m, payload)
}

type MethodGetBlacklistOutput struct {
	Result string
}

func (m *MethodGetBlacklistOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodGetBlacklist, m.Result)
}

func (m *MethodGetBlacklistOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodGetBlacklist, m, payload)
}
