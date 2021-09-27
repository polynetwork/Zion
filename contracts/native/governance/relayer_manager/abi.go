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
package relayer_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/relayer_manager_abi"
	polycomm "github.com/polynetwork/poly/common"
)

var (
	EventRegisterRelayer        = relayer_manager_abi.MethodRegisterRelayer
	EventApproveRegisterRelayer = relayer_manager_abi.MethodApproveRegisterRelayer
	EventRemoveRelayer          = relayer_manager_abi.MethodRemoveRelayer
	EventApproveRemoveRelayer   = relayer_manager_abi.MethodApproveRemoveRelayer
)

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(relayer_manager_abi.RelayerManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type RelayerListParam struct {
	AddressList []common.Address
	Address     common.Address
}

func (this *RelayerListParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.AddressList)))
	for _, v := range this.AddressList {
		sink.WriteVarBytes(v[:])
	}
	sink.WriteVarBytes(this.Address[:])
}

func (this *RelayerListParam) Deserialization(source *polycomm.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize AddressList length error")
	}
	addressList := make([]common.Address, 0)
	for i := 0; uint64(i) < n; i++ {
		address, eof := source.NextVarBytes()
		if eof {
			return fmt.Errorf("source.NextVarBytes, deserialize address error")
		}
		addr, err := common.AddressParseFromBytes(address)
		if err != nil {
			return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
		}
		addressList = append(addressList, addr)
	}

	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
	}
	this.AddressList = addressList
	this.Address = addr
	return nil
}

type ApproveRelayerParam struct {
	ID      uint64
	Address common.Address
}
