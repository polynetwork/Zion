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
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "alloc proxy"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(ISideChainLockProxyABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.AllocProxyContractAddress
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

//function verifyHeaderAndExecuteTxInput(bytes calldata header, bytes calldata rawCrossTx, bytes calldata proof) external returns (bool);
type MethodVerifyHeaderAndExecuteTxInput struct {
	Header     []byte
	RawCrossTx []byte
	Proof      []byte
	Extra      []byte
}

func (i *MethodVerifyHeaderAndExecuteTxInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodVerifyHeaderAndExecuteTx, i.Header, i.RawCrossTx, i.Proof, i.Extra)
}
func (i *MethodVerifyHeaderAndExecuteTxInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodVerifyHeaderAndExecuteTx, i, payload)
}

//event BurnEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount, bytes crossTxId);
func emitBurnEvent(s *native.NativeContract, toChainId uint64, fromAddr, toAddr common.Address, amount *big.Int, crossTxId uint64) error {
	return s.AddNotify(ABI, []string{EventBurnEvent}, toChainId, fromAddr, toAddr, amount, utils.Uint64Bytes(crossTxId))
}

//event MintEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount);
func emitMintEvent(s *native.NativeContract, toChainId uint64, fromAddr, toAddr common.Address, amount *big.Int) error {
	return s.AddNotify(ABI, []string{EventMintEvent}, toChainId, fromAddr, toAddr, amount)
}
