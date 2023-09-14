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

 package community

 import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

// storage key prefix
const (
	SKP_COMMUNITY_INFO                = "st_community_info"
)

var (
	PercentDecimal = new(big.Int).Exp(big.NewInt(10), big.NewInt(4), nil)
)

type CommunityInfo struct {
	CommunityRate    *big.Int
	CommunityAddress common.Address
}

func StoreCommunityInfo(s *state.StateDB, communityRate *big.Int, communityAddress common.Address) (*CommunityInfo, error) {
	cache := (*state.CacheDB)(s)
	communityInfo := &CommunityInfo{
		CommunityRate:    communityRate,
		CommunityAddress: communityAddress,
	}
	if err := setGenesisCommunityInfo(cache, communityInfo); err != nil {
		return nil, err
	}
	return communityInfo, nil
}
 
 func setGenesisCommunityInfo(s *state.CacheDB, communityInfo *CommunityInfo) error {
	 if communityInfo.CommunityRate.Cmp(PercentDecimal) > 0 {
		 return fmt.Errorf("setGenesisCommunityInfo, CommunityRate over size")
	 }
	 key := communityInfoKey()
	 store, err := rlp.EncodeToBytes(communityInfo)
	 if err != nil {
		 return fmt.Errorf("setCommunityInfo, serialize community info error: %v", err)
	 }
	 customSet(s, key, store)
	 return nil
 }
 
 func SetCommunityInfo(s *native.NativeContract, communityInfo *CommunityInfo) error {
	 if communityInfo.CommunityRate.Cmp(PercentDecimal) > 0 {
		 return fmt.Errorf("setCommunityInfo, CommunityRate over size")
	 }
	 key := communityInfoKey()
	 store, err := rlp.EncodeToBytes(communityInfo)
	 if err != nil {
		 return fmt.Errorf("setCommunityInfo, serialize community info error: %v", err)
	 }
	 set(s, key, store)
	 return nil
 }
 
 func GetCommunityInfoImpl(s *native.NativeContract) (*CommunityInfo, error) {
	 communityInfo := new(CommunityInfo)
	 key := communityInfoKey()
	 store, err := get(s, key)
	 if err != nil {
		 return nil, fmt.Errorf("GetCommunityInfoImpl, get store error: %v", err)
	 }
	 if err := rlp.DecodeBytes(store, communityInfo); err != nil {
		 return nil, fmt.Errorf("GetCommunityInfoImpl, deserialize community info error: %v", err)
	 }
	 return communityInfo, nil
 }
 
 func GetCommunityInfoFromDB(s *state.StateDB) (*CommunityInfo, error) {
	 cache := (*state.CacheDB)(s)
	 communityInfo := new(CommunityInfo)
	 key := communityInfoKey()
	 store, err := customGet(cache, key)
	 if err != nil {
		 return nil, fmt.Errorf("GetCommunityInfoFromDB, get store error: %v", err)
	 }
	 if err := rlp.DecodeBytes(store, communityInfo); err != nil {
		 return nil, fmt.Errorf("GetCommunityInfoFromDB, deserialize community info error: %v", err)
	 }
	 return communityInfo, nil
 }
 
 // ====================================================================
 //
 // storage basic operations
 //
 // ====================================================================
 
 func get(s *native.NativeContract, key []byte) ([]byte, error) {
	 return customGet(s.GetCacheDB(), key)
 }
 
 func set(s *native.NativeContract, key, value []byte) {
	 customSet(s.GetCacheDB(), key, value)
 }
 
 func del(s *native.NativeContract, key []byte) {
	 customDel(s.GetCacheDB(), key)
 }
 
 func customGet(db *state.CacheDB, key []byte) ([]byte, error) {
	 value, err := db.Get(key)
	 if err != nil {
		 return nil, err
	 } else if value == nil || len(value) == 0 {
		 return nil, errors.New("EOF")
	 } else {
		 return value, nil
	 }
 }
 
 func customSet(db *state.CacheDB, key, value []byte) {
	 db.Put(key, value)
 }
 
 func customDel(db *state.CacheDB, key []byte) {
	 db.Delete(key)
 }
 
 func communityInfoKey() []byte {
	 return utils.ConcatKey(utils.NodeManagerContractAddress, []byte(SKP_COMMUNITY_INFO))
 }
 