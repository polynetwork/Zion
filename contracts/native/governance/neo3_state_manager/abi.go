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
package neo3_state_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	EventRegisterStateValidator = "evtRegisterStateValidator"
	EventApproveRegisterStateValidator = "evtApproveRegisterStateValidator"
	EventRemoveStateValidator = "evtRemoveStateValidator"
	EventApproveRemoveStateValidator   = "evtApproveRemoveStateValidator"

)

const abiJson = `[
   {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventRegisterStateValidator + `","type":"event"},
   {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventApproveRegisterStateValidator + `","type":"event"},
   {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventRemoveStateValidator + `","type":"event"},
   {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventApproveRemoveStateValidator + `","type":"event"},
   {"inputs":[],"name":"` + MethodContractName + `","outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable","type":"function"},
   {"inputs":[],"name":"` + MethodGetCurrentStateValidator + `","outputs":[{"internalType":"bytes","name":"Validator","type":"bytes"}],"stateMutability":"nonpayable","type":"function"},   
   {"inputs":[{"internalType":"uint64","name":"ID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveRegisterStateValidator + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
   {"inputs":[{"internalType":"uint64","name":"ID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveRemoveStateValidator + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
   {"inputs":[{"internalType":"string[]","name":"StateValidators","type":"string[]"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodRegisterStateValidator + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
   {"inputs":[{"internalType":"string[]","name":"StateValidators","type":"string[]"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodRemoveStateValidator + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

func GetABI() *abi.ABI {
	//ab, err := abi.JSON(strings.NewReader(neo3_state_manager_abi.Neo3StateManagerABI))
	ab, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type StateValidatorListParam struct {
	StateValidators []string // public key strings in encoded format, each is 33 bytes in []byte
	Address         common.Address
}

func (this *StateValidatorListParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{this.StateValidators, this.Address})
}
func (this *StateValidatorListParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StateValidators []string
		Address         common.Address
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.StateValidators, this.Address = data.StateValidators, data.Address
	return nil
}

type ApproveStateValidatorParam struct {
	ID      uint64 // StateValidatorApproveID
	Address common.Address
}
