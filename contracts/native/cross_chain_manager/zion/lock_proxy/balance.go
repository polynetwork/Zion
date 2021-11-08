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
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
)

var ErrUnSupported = errors.New("erc20 asset cross chain unsupported yet")

func getBalanceFor(s *native.NativeContract, fromAsset common.Address) (*big.Int, error) {
	if fromAsset == common.EmptyAddress {
		return s.StateDB().GetBalance(this), nil
	} else {
		return erc20Balance(s, fromAsset, this)
	}
}
/*
function _transferToContract(address fromAssetHash, uint256 amount) internal returns (bool) {
	if (fromAssetHash == address(0)) {
		// fromAssetHash === address(0) denotes user choose to lock ether
		// passively check if the received msg.value equals amount
		require(msg.value != 0, "transferred ether cannot be zero!");
		require(msg.value == amount, "transferred ether is not equal to amount!");
	} else {
		// make sure lockproxy contract will decline any received ether
		require(msg.value == 0, "there should be no ether transfer!");
		// actively transfer amount of asset from msg.sender to lock_proxy contract
		require(_transferERC20ToContract(fromAssetHash, _msgSender(), address(this), amount), "transfer erc20 asset to lock_proxy contract failed!");
	}
	return true;
}
*/
func transfer2Contract(s *native.NativeContract, fromAsset common.Address, amount *big.Int) error {
	return nil
}

/*
function _transferFromContract(address toAssetHash, address toAddress, uint256 amount) internal returns (bool) {
	if (toAssetHash == address(0x0000000000000000000000000000000000000000)) {
		// toAssetHash === address(0) denotes contract needs to unlock ether to toAddress
		// convert toAddress from 'address' type to 'address payable' type, then actively transfer ether
		address(uint160(toAddress)).transfer(amount);
	} else {
		// actively transfer amount of asset from msg.sender to lock_proxy contract
		require(_transferERC20FromContract(toAssetHash, toAddress, amount), "transfer erc20 asset to lock_proxy contract failed!");
	}
	return true;
}
*/
func transferFromContract(s *native.NativeContract, toAsset, toAddress common.Address, amount *big.Int) error {
	return nil
}

func onlySupportNativeToken(fromAsset common.Address) bool {
	if fromAsset == common.EmptyAddress {
		return true
	}
	return false
}

// todo: get erc20 balance
func erc20Balance(s *native.NativeContract, asset, user common.Address) (*big.Int, error) {
	return nil, ErrUnSupported
}
