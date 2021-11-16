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

package alloc_proxy

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/zion"
	"github.com/ethereum/go-ethereum/core/types"
)

func verifyHeaderAndCheckEpoch(header *types.Header, lastEpoch, curEpoch *nm.EpochInfo) error {
	nextEpochStartHeight, nextEpochVals, err := zion.VerifyHeader(header, lastEpoch.MemberList(), true)
	if err != nil {
		return err
	}
	if nextEpochStartHeight != curEpoch.StartHeight {
		return fmt.Errorf("epoch start height not match, expect %d got %d", nextEpochStartHeight, curEpoch.StartHeight)
	}
	if curEpochVals := curEpoch.MemberList(); !isSameVals(nextEpochVals, curEpochVals) {
		return fmt.Errorf("epoch validators not match, expect %v, got %v", nextEpochVals, curEpochVals)
	}
	return nil
}

func isSameVals(src, cmp []common.Address) bool {
	exist := func(a common.Address, l []common.Address) bool {
		for _, v := range l {
			if a == v {
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

	if !contain(src, cmp) {
		return false
	}
	if !contain(cmp, src) {
		return false
	}

	return true
}
