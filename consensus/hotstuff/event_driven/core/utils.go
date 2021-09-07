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
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func (c *core) newMsgLogger(msgtyp interface{}) log.Logger {
	return c.logger.New("view", c.currentView(), "msg", msgtyp)
}

func (c *core) newSenderLogger(msgtyp string) log.Logger {
	return c.logger.New("view", c.currentView(), "msg", msgtyp)
}

func (c *core) newLogger() log.Logger {
	return c.logger.New("view", c.currentView())
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func copyNum(src *big.Int) *big.Int {
	return new(big.Int).Set(src)
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

func bigEq(src, dst *big.Int) bool {
	return src.Cmp(dst) == 0
}

func bigEq0(num *big.Int) bool {
	return num.Cmp(common.Big0) == 0
}
