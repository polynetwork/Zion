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

package signature_manager

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/signature_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddSignatureParam(t *testing.T) {
	inputTestCases := []AddSignatureParam{
		{common.EmptyAddress, big.NewInt(0), []byte{}, []byte{}},
		{common.HexToAddress("0x12"), big.NewInt(2), []byte{'a'}, []byte{'b'}},
	}

	for _, tc := range inputTestCases {
		enc, err := utils.PackMethod(ABI, MethodAddSignature, tc.Addr, tc.SideChainID, tc.Subject, tc.Signature)
		assert.NoError(t, err)

		var got AddSignatureParam
		assert.NoError(t, utils.UnpackMethod(ABI, MethodAddSignature, &got, enc))
		assert.ObjectsAreEqual(tc, got)
	}

	outputTestCases := []bool{true, false}
	for _, expect := range outputTestCases {
		raw, err := utils.PackOutputs(ABI, MethodAddSignature, expect)
		assert.NoError(t, err)

		var got bool
		assert.NoError(t, utils.UnpackOutputs(ABI, MethodAddSignature, &got, raw))
		assert.Equal(t, expect, got)
	}
}
