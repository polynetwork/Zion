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

package side_chain_manager

import (
	"io"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type SideChain struct {
	Owner       common.Address
	ChainID     uint64
	Router      uint64
	Name        string
	CCMCAddress []byte
	ExtraInfo   []byte
}

type Fee struct {
	View uint64
	Fee  *big.Int
}

type FeeInfo struct {
	StartHeight uint64
	FeeInfo     map[common.Address]*big.Int
}

type FeeVote struct {
	Address common.Address
	Fee     *big.Int
}

func (this *FeeInfo) EncodeRLP(w io.Writer) error {
	feeVote := make([]*FeeVote, 0, len(this.FeeInfo))
	for k, v := range this.FeeInfo {
		feeVote = append(feeVote, &FeeVote{k, v})
	}
	sort.SliceStable(feeVote, func(i, j int) bool {
		return feeVote[i].Fee.Cmp(feeVote[j].Fee) == 1
	})
	return rlp.Encode(w, []interface{}{this.StartHeight, feeVote})
}

func (this *FeeInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StartHeight uint64
		FeeVote     []*FeeVote
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.StartHeight = data.StartHeight

	feeInfo := make(map[common.Address]*big.Int, len(data.FeeVote))
	for _, v := range data.FeeVote {
		feeInfo[v.Address] = v.Fee
	}
	this.FeeInfo = feeInfo

	return nil
}

type RippleExtraInfo struct {
	Operator      common.Address
	Sequence      uint64
	Quorum        uint64
	SignerNum     uint64
	Pks           [][]byte
	ReserveAmount *big.Int
}

type AssetBind struct {
	AssetMap     map[uint64][]byte
	LockProxyMap map[uint64][]byte
}

type BindInfo struct {
	ChainId uint64
	Address []byte
}

func (this *AssetBind) EncodeRLP(w io.Writer) error {
	assetList := make([]*BindInfo, 0, len(this.AssetMap))
	for k, v := range this.AssetMap {
		assetList = append(assetList, &BindInfo{k, v})
	}
	sort.SliceStable(assetList, func(i, j int) bool {
		return assetList[i].ChainId > assetList[j].ChainId
	})

	lockProxyList := make([]*BindInfo, 0, len(this.LockProxyMap))
	for k, v := range this.LockProxyMap {
		lockProxyList = append(lockProxyList, &BindInfo{k, v})
	}
	sort.SliceStable(lockProxyList, func(i, j int) bool {
		return lockProxyList[i].ChainId > lockProxyList[j].ChainId
	})
	return rlp.Encode(w, []interface{}{assetList, lockProxyList})
}

func (this *AssetBind) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		AssetList     []*BindInfo
		LockProxyList []*BindInfo
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	assetMap := make(map[uint64][]byte, len(data.AssetList))
	for _, v := range data.AssetList {
		assetMap[v.ChainId] = v.Address
	}
	lockProxyMap := make(map[uint64][]byte, len(data.LockProxyList))
	for _, v := range data.LockProxyList {
		lockProxyMap[v.ChainId] = v.Address
	}
	this.AssetMap = assetMap
	this.LockProxyMap = lockProxyMap

	return nil
}
