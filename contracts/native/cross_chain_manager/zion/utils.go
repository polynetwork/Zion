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

package zion

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	internal "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
)

func verifyTx(proof, extra []byte, hdr *types.Header, contract common.Address) error {
	ethProof := new(ethapi.AccountResult)
	if err := json.Unmarshal(proof, ethProof); err != nil {
		return fmt.Errorf("VerifyFromEthProof, unmarshal proof error:%s", err)
	}

	proofResult, err := internal.VerifyAccountResult(ethProof, hdr, contract)
	if err != nil {
		return fmt.Errorf("VerifyFromEthProof, verifyMerkleProof error:%v", err)
	}
	if proofResult == nil {
		return fmt.Errorf("VerifyFromEthProof, verifyMerkleProof failed!")
	}
	if !internal.CheckProofResult(proofResult, extra) {
		return fmt.Errorf("VerifyFromEthProof, verify proof value hash failed, proof result:%x, extra:%x", proofResult, extra)
	}
	return nil
}
