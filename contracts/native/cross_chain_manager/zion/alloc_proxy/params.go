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
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

type CrossTx struct {
	ToChainId   uint64
	FromAddress common.Address
	ToAddress   common.Address
	Amount      *big.Int
	Index       uint64

	hash atomic.Value
}

func (tx *CrossTx) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{tx.ToChainId, tx.FromAddress, tx.ToAddress, tx.Amount, tx.Index})
}

func (tx *CrossTx) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ToChainId   uint64
		FromAddress common.Address
		ToAddress   common.Address
		Amount      *big.Int
		Index       uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	tx.ToChainId, tx.FromAddress, tx.ToAddress, tx.Amount, tx.Index =
		data.ToChainId, data.FromAddress, data.ToAddress, data.Amount, data.Index
	return nil
}

func (tx *CrossTx) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		ToChainId   uint64
		FromAddress common.Address
		ToAddress   common.Address
		Amount      *big.Int
		Index       uint64
	}{
		ToChainId:   tx.ToChainId,
		FromAddress: tx.FromAddress,
		ToAddress:   tx.ToAddress,
		Amount:      tx.Amount,
		Index:       tx.Index,
	}
	v := utils.RLPHash(inf)
	tx.hash.Store(v)
	return v
}

func EncodeHeader(h *types.Header) ([]byte, error) {
	return h.MarshalJSON()
}

func DecodeHeader(payload []byte) (*types.Header, error) {
	h := new(types.Header)
	if err := h.UnmarshalJSON(payload); err != nil {
		return nil, err
	}
	return h, nil
}

func EncodeEpoch(epoch *nm.EpochInfo) ([]byte, error) {
	return rlp.EncodeToBytes(epoch)
}

func DecodeEpoch(payload []byte) (*nm.EpochInfo, error) {
	epoch := new(nm.EpochInfo)
	if err := rlp.DecodeBytes(payload, epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}
