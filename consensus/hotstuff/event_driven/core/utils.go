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

package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func copyNum(src *big.Int) *big.Int {
	return new(big.Int).Set(src)
}

func newView(round, height *big.Int) *hotstuff.View {
	return &hotstuff.View{
		Round:  new(big.Int).Set(round),
		Height: new(big.Int).Set(height),
	}
}

func isTC(qc *hotstuff.QuorumCert) bool {
	if qc.Hash == utils.EmptyHash {
		return true
	}
	return false
}

func bigAdd1(num *big.Int) *big.Int {
	return new(big.Int).Add(num, common.Big1)
}

func bigAdd1Eq(src, dst *big.Int) (*big.Int, bool) {
	sum := new(big.Int).Add(src, common.Big1)
	return sum, sum.Cmp(dst) == 0
}

func bigSub1(num *big.Int) *big.Int {
	return new(big.Int).Sub(num, common.Big1)
}

func bigSub1Eq(src, dst *big.Int) (*big.Int, bool) {
	res := bigSub1(src)
	return res, res.Cmp(dst) == 0
}

func bigCmp(src, dst *big.Int) int {
	return src.Cmp(dst)
}

func bigEq(src, dst *big.Int) bool {
	return src.Cmp(dst) == 0
}

func bigLt(src, dst *big.Int) bool {
	return src.Cmp(dst) < 0
}

func bigGt(src, dst *big.Int) bool {
	return src.Cmp(dst) > 0
}

func bigEq0(num *big.Int) bool {
	return num.Cmp(common.Big0) == 0
}

// convert inner MsgType to hotstuff MsgType
func convertUpMsgType(data interface{}) hotstuff.MsgType {
	code := data.(uint64)
	return MsgType(code)
}

// convert hotstuff MsgType to inner MsgType
func convertDownMsgType(typ hotstuff.MsgType) MsgType {
	code := typ.Value()
	return MsgType(code)
}
