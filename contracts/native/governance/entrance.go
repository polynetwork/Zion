/*
 * Copyright (C) 2022 The zion network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */
package governance

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

func AssembleSystemTransactions(state *state.StateDB, height uint64) (types.Transactions, error) {
	// Genesis block has no system transaction?
	if height == 0 {
		return nil, nil
	}

	var txs types.Transactions
	systemSenderNonce := state.GetNonce(utils.SystemTxSender)
	// SystemTransaction: NodeManager.EndBlock
	{
		payload, err := new(nm.EndBlockParam).Encode()
		if err != nil {
			return nil, err
		}

		gas, err := core.IntrinsicGas(payload, nil, false, true, true)
		if err != nil {
			return nil, err
		}
		gas += nm.GasTable[node_manager_abi.MethodEndBlock]
		txs = append(txs, types.NewTransaction(systemSenderNonce, utils.NodeManagerContractAddress, common.Big0, gas, common.Big0, payload))
	}

	// SystemTransaction: NodeManager.ChangeEpoch
	{
		epoch, err := nm.GetCurrentEpochInfoFromDB(state)
		if err != nil {
			return nil, err
		}

		if epoch == nil || epoch.EndHeight == nil {
			return nil, fmt.Errorf("unexpected epoch or epoch end height missing")
		}

		if height + 1 == epoch.EndHeight.Uint64() {
			payload, err := new(nm.ChangeEpochParam).Encode()
			if err != nil {
				return nil, err
			}

			gas, err := core.IntrinsicGas(payload, nil, false, true, true)
			if err != nil {
				return nil, err
			}
			gas += nm.GasTable[node_manager_abi.MethodChangeEpoch]
			txs = append(txs, types.NewTransaction(systemSenderNonce + 1, utils.NodeManagerContractAddress, common.Big0, gas, common.Big0, payload))
		}
	}
	return txs, nil
}