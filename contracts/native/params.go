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
package native

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

// FailedTxGasUsage tx's gas usage should not be greater than an minimum fixed value if it execute failed.
const FailedTxGasUsage = uint64(100)

const (
	NativeNodeManager = "node_manager"
	NativeMaasConfig  = "maas_config"
)

var NativeContractAddrMap = map[string]common.Address{
	NativeNodeManager: utils.NodeManagerContractAddress,
	NativeMaasConfig:  utils.MaasConfigContractAddress,
}

func IsNativeContract(addr common.Address) bool {
	for _, v := range NativeContractAddrMap {
		if v == addr {
			return true
		}
	}
	return false
}
