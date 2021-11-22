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

package lock_proxy

import (
	"github.com/ethereum/go-ethereum/common"
)

// compareVals return true if `src` equals to `cmp`
func compareVals(v1, v2 []common.Address) bool {
	exist := func(addr common.Address, list []common.Address) bool {
		for _, v := range list {
			if addr == v {
				return true
			}
		}
		return false
	}

	contain := func(l1, l2 []common.Address) bool {
		for _, v := range l1 {
			if !exist(v, l2) {
				return false
			}
		}
		return true
	}

	if !contain(v1, v2) {
		return false
	}
	if !contain(v2, v1) {
		return false
	}

	return true
}
