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
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_lock_proxy_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "side chain lock proxy"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(ISideChainLockProxyABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.LockProxyContractAddress
)

// function name
type MethodContractNameInput struct{}

func (m *MethodContractNameInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodName)
}

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
}

//function burn(uint64 toChainId, address toAddress, uint256 amount) external returns (bool);
type MethodBurnInput struct {
	ToChainId uint64
	ToAddress common.Address
	Amount    *big.Int
}

func (i *MethodBurnInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodBurn, i.ToChainId, i.ToAddress, i.Amount)
}
func (i *MethodBurnInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodBurn, i, payload)
}

//function mint(bytes calldata argsBs, bytes calldata fromContractAddr, uint64 fromChainId) external returns (bool);
type MethodMintInput struct {
	ArgsBs           []byte
	FromContractAddr []byte
	FromChainId      uint64
}

func (i *MethodMintInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodMint, i.ArgsBs, i.FromContractAddr, i.FromChainId)
}
func (i *MethodMintInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodMint, i, payload)
}

//event BurnEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
func emitBurnEvent(s *native.NativeContract, fromAsset, fromAddr common.Address, toChainID uint64, toAsset, toAddr []byte, amount *big.Int) error {
	return s.AddNotify(ABI, []string{EventBurnEvent}, fromAsset, fromAddr, toChainID, toAsset, toAddr, amount.Bytes())
}

//event MintEvent(address toAssetHash, address toAddress, uint256 amount);
func emitMintEvent(s *native.NativeContract, toAsset, toAddr common.Address, amount *big.Int) error {
	return s.AddNotify(ABI, []string{EventMintEvent}, toAsset, toAddr, amount)
}
