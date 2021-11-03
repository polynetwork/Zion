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

package eccd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/eccd_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	ErrABINotExist    = errors.New("abi method not exist")
	ErrABIInputLength = errors.New("abi input length error")
	ErrABIInputType   = errors.New("abi input type error")
)

const contractName = "node manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IEthCrossChainDataABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.RelayChainECCDContractAddress
)

type MethodPutCurEpochStartHeightInput struct {
	CurEpochStartHeight uint32
}

func (m *MethodPutCurEpochStartHeightInput) Encode() ([]byte, error) {
	mth, ok := ABI.Methods[MethodPutCurEpochStartHeight]
	if !ok {
		return nil, ErrABINotExist
	}
	return mth.Inputs.Pack(m.CurEpochStartHeight)
}
func (m *MethodPutCurEpochStartHeightInput) Decode(payload []byte) error {
	mth, ok := ABI.Methods[MethodPutCurEpochStartHeight]
	if !ok {
		return ErrABINotExist
	}

	list, err := mth.Inputs.Unpack(payload)
	if err != nil {
		return err
	}

	if len(list) != 1 {
		return ErrABIInputLength
	}
	iter := list[0]
	value, ok := iter.(uint32)
	if !ok {
		return ErrABIInputType
	}

	m.CurEpochStartHeight = value
	return nil
}

type MethodPutCurEpochStartHeightOutput struct {
	Succeed bool
}
