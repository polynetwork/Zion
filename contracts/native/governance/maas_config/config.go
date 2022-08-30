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
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
)

var (
	gasTable = map[string]uint64{
		MethodName:         0,
		MethodChangeOwner:  30000,
		MethodGetOwner:     0,
		MethodBlockAccount: 30000,
		MethodIsBlocked:    0,
		MethodGetBlacklist: 0,

		MethodEnableGasManage:    30000,
		MethodSetGasManager:      30000,
		MethodIsGasManageEnabled: 0,
		MethodIsGasManager:       0,
		MethodGetGasManagerList:  0,

		MethodSetGasUsers:    30000,
		MethodIsGasUser:      0,
		MethodGetGasUserList: 0,

		MethodSetAdmins:    30000,
		MethodIsAdmin:      0,
		MethodGetAdminList: 0,
	}
)

func InitMaasConfig() {
	InitABI()
	native.Contracts[this] = RegisterMaasConfigContract
}

func RegisterMaasConfigContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodChangeOwner, ChangeOwner)
	s.Register(MethodGetOwner, GetOwner)

	s.Register(MethodBlockAccount, BlockAccount)
	s.Register(MethodIsBlocked, IsBlocked)
	s.Register(MethodGetBlacklist, GetBlacklist)

	s.Register(MethodEnableGasManage, EnableGasManage)
	s.Register(MethodSetGasManager, SetGasManager)
	s.Register(MethodIsGasManageEnabled, IsGasManageEnabled)
	s.Register(MethodIsGasManager, IsGasManager)
	s.Register(MethodGetGasManagerList, GetGasManagerList)

	s.Register(MethodSetGasUsers, SetGasUsers)
	s.Register(MethodIsGasUser, IsGasUser)
	s.Register(MethodGetGasUserList, GetGasUserList)

	s.Register(MethodSetAdmins, SetAdmins)
	s.Register(MethodIsAdmin, IsAdmin)
	s.Register(MethodGetAdminList, GetAdminList)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

// change owner
func ChangeOwner(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	// check caller == origin
	if err := contract.ValidateOwner(s, caller); err != nil {
		return utils.ByteFailed, errors.New("caller is not equal to origin")
	}

	// check owner
	currentOwner := getOwner(s)
	if currentOwner != common.EmptyAddress && caller != currentOwner {
		return utils.ByteFailed, errors.New("invalid authority for owner")
	}

	// decode input
	input := new(MethodChangeOwnerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("ChangeOwner", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// verify new owner address
	m := getAddressMap(s, blacklistKey)
	_, ok := m[input.Addr]
	if ok {
		err := errors.New("new owner address in blacklist")
		log.Trace("ChangeOwner", "invalid new owner", err)
		return utils.ByteFailed, err
	}

	// store owner
	set(s, ownerKey, input.Addr.Bytes())

	// emit event log
	if err := s.AddNotify(ABI, []string{EventChangeOwner}, common.BytesToHash(currentOwner.Bytes()), common.BytesToHash(input.Addr.Bytes())); err != nil {
		log.Trace("ChangeOwner", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventChangeOwner error")
	}

	return utils.ByteSuccess, nil
}

// get owner
func GetOwner(s *native.NativeContract) ([]byte, error) {
	output := &MethodAddressOutput{Addr: getOwner(s)}
	return output.Encode(MethodGetOwner)
}

func getOwner(s *native.NativeContract) common.Address {
	// get value
	value, _ := get(s, ownerKey)
	if len(value) == 0 {
		return common.EmptyAddress
	}
	return common.BytesToAddress(value)
}

func checkOwner(s *native.NativeContract) error {
	caller := s.ContractRef().CurrentContext().Caller
	origin := s.ContractRef().TxOrigin()
	if caller != origin {
		return errors.New("caller is not equal to origin")
	}

	if origin != getOwner(s) {
		return errors.New("invalid authority for owner")
	}
	return nil
}

func isAdmin(s *native.NativeContract) bool {
	origin := s.ContractRef().TxOrigin()
	m := getAddressMap(s, gasAdminListKey)
	_, ok := m[origin]
	return ok
}

func checkOwnerOrAdmin(s *native.NativeContract) error {
	caller := s.ContractRef().CurrentContext().Caller
	origin := s.ContractRef().TxOrigin()
	if caller != origin {
		return errors.New("caller is not equal to origin")
	}

	if origin != getOwner(s) && !isAdmin(s) {
		return errors.New("invalid authority for owner or admin")
	}
	return nil
}

// block account(add account to blacklist map) or unblock account
func BlockAccount(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check owner
	if err := checkOwner(s); err != nil {
		return utils.ByteFailed, err
	}

	// decode input
	input := new(MethodBlockAccountInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("blockAccount", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	currentOwner := getOwner(s)
	if input.Addr == currentOwner {
		err := errors.New("block owner is forbidden")
		log.Trace("blockAccount", "block owner is forbidden", err)
		return utils.ByteFailed, err
	}

	m := getAddressMap(s, blacklistKey)
	if input.DoBlock {
		m[input.Addr] = struct{}{}
	} else {
		delete(m, input.Addr)
	}

	value, err := json.Marshal(m)
	if err != nil {
		log.Trace("blockAccount", "encode value failed", err)
		return utils.ByteFailed, errors.New("encode value failed")
	}
	set(s, blacklistKey, value)

	// emit event log
	if err := s.AddNotify(ABI, []string{EventBlockAccount}, common.BytesToHash(input.Addr.Bytes()), input.DoBlock); err != nil {
		log.Trace("blockAccount", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventBlockAccount error")
	}

	return utils.ByteSuccess, nil
}

func getAddressMap(s *native.NativeContract, key []byte) map[common.Address]struct{} {
	value, _ := get(s, key)
	m := make(map[common.Address]struct{})
	if len(value) > 0 {
		if err := json.Unmarshal(value, &m); err != nil {
			log.Trace("getAddressMap", "decode value failed", err)
		}
	}
	return m
}

// check if account is blocked
func IsBlocked(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodIsBlockedInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("IsBlocked", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// get value
	m := getAddressMap(s, blacklistKey)
	_, ok := m[input.Addr]
	output := &MethodBoolOutput{Success: ok}

	return output.Encode(MethodIsBlocked)
}

// get blacklist json
func GetBlacklist(s *native.NativeContract) ([]byte, error) {
	// get value
	m := getAddressMap(s, blacklistKey)
	list := make([]common.Address, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	result, _ := json.Marshal(list)
	output := &MethodStringOutput{Result: string(result)}
	return output.Encode(MethodGetBlacklist)
}

// enable gas manage
func EnableGasManage(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check owner
	if err := checkOwner(s); err != nil {
		return utils.ByteFailed, err
	}

	// decode input
	input := new(MethodEnableGasManageInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("EnableGasManage", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// set enable status
	if input.DoEnable {
		set(s, gasManageEnableKey, utils.BYTE_TRUE)
	} else {
		del(s, gasManageEnableKey)
	}

	// emit event log
	if err := s.AddNotify(ABI, []string{EventEnableGasManage}, input.DoEnable); err != nil {
		log.Trace("EnableGasManage", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventEnableGasManage error")
	}

	return utils.ByteSuccess, nil
}

// check if gas manage is enabled
func IsGasManageEnabled(s *native.NativeContract) ([]byte, error) {
	// get value
	value, _ := get(s, gasManageEnableKey)
	output := &MethodBoolOutput{Success: len(value) > 0}
	return output.Encode(MethodIsGasManageEnabled)
}

// set gas manager address
func SetGasManager(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check owner
	if err := checkOwner(s); err != nil {
		return utils.ByteFailed, err
	}

	// decode input
	input := new(MethodSetGasManagerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("SetGasManager", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	m := getAddressMap(s, gasManagerListKey)
	if input.IsManager {
		m[input.Addr] = struct{}{}
	} else {
		delete(m, input.Addr)
	}

	value, err := json.Marshal(m)
	if err != nil {
		log.Trace("SetGasManager", "encode value failed", err)
		return utils.ByteFailed, errors.New("encode value failed")
	}
	set(s, gasManagerListKey, value)

	// emit event log
	if err := s.AddNotify(ABI, []string{EventSetGasManager}, common.BytesToHash(input.Addr.Bytes()), input.IsManager); err != nil {
		log.Trace("SetGasManager", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventSetGasManager error")
	}

	return utils.ByteSuccess, nil
}

// check if address is in gas manager list
func IsGasManager(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodIsGasManagerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("IsGasManager", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// get value
	m := getAddressMap(s, gasManagerListKey)
	_, ok := m[input.Addr]
	output := &MethodBoolOutput{Success: ok}

	return output.Encode(MethodIsGasManager)
}

// get gas manager list json
func GetGasManagerList(s *native.NativeContract) ([]byte, error) {
	// get value
	m := getAddressMap(s, gasManagerListKey)
	list := make([]common.Address, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	result, _ := json.Marshal(list)
	output := &MethodStringOutput{Result: string(result)}
	return output.Encode(MethodGetGasManagerList)
}

// set gas users
func SetGasUsers(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check owner
	if err := checkOwnerOrAdmin(s); err != nil {
		return utils.ByteFailed, err
	}

	// decode input
	input := new(MethodSetGasUsersInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("SetGasUsers", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	m := getAddressMap(s, gasUserListKey)
	for _, v := range input.Addrs {
		if input.AddOrRemove {
			m[v] = struct{}{}
		} else {
			delete(m, v)
		}
	}

	value, err := json.Marshal(m)
	if err != nil {
		log.Trace("SetGasUsers", "encode value failed", err)
		return utils.ByteFailed, errors.New("encode value failed")
	}
	set(s, gasUserListKey, value)

	// emit event log
	if err := s.AddNotify(ABI, []string{EventSetGasUsers}, input.Addrs, input.AddOrRemove); err != nil {
		log.Trace("SetGasUsers", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventSetGasUsers error")
	}

	return utils.ByteSuccess, nil
}

// check if address is in gas user list
func IsGasUser(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodIsGasUserInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("IsGasUser", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// get value
	m := getAddressMap(s, gasUserListKey)
	_, ok := m[input.Addr]
	output := &MethodBoolOutput{Success: ok}

	return output.Encode(MethodIsGasUser)
}

// get gas user list json
func GetGasUserList(s *native.NativeContract) ([]byte, error) {
	// get value
	m := getAddressMap(s, gasUserListKey)
	list := make([]common.Address, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	result, _ := json.Marshal(list)
	output := &MethodStringOutput{Result: string(result)}
	return output.Encode(MethodGetGasUserList)
}

// set admins
func SetAdmins(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check owner
	if err := checkOwner(s); err != nil {
		return utils.ByteFailed, err
	}

	// decode input
	input := new(MethodSetAdminsInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("SetAdmins", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	m := getAddressMap(s, gasAdminListKey)
	for _, v := range input.Addrs {
		if input.AddOrRemove {
			m[v] = struct{}{}
		} else {
			delete(m, v)
		}
	}

	value, err := json.Marshal(m)
	if err != nil {
		log.Trace("SetAdmins", "encode value failed", err)
		return utils.ByteFailed, errors.New("encode value failed")
	}
	set(s, gasAdminListKey, value)

	// emit event log
	if err := s.AddNotify(ABI, []string{EventSetAdmins}, input.Addrs, input.AddOrRemove); err != nil {
		log.Trace("SetAdmins", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventSetAdmins error")
	}

	return utils.ByteSuccess, nil
}

// check if address is in admin list
func IsAdmin(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodIsAdminInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("IsAdmin", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// get value
	m := getAddressMap(s, gasAdminListKey)
	_, ok := m[input.Addr]
	output := &MethodBoolOutput{Success: ok}

	return output.Encode(MethodIsAdmin)
}

// get admin list json
func GetAdminList(s *native.NativeContract) ([]byte, error) {
	m := getAddressMap(s, gasAdminListKey)
	list := make([]common.Address, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	result, _ := json.Marshal(list)
	output := &MethodStringOutput{Result: string(result)}
	return output.Encode(MethodGetAdminList)
}
