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
package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type CacheDB StateDB

func (c *CacheDB) Put(key []byte, value []byte) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	so := s.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		so.SetState(c.db, slot, value)
	}

}

func Key2Slot(key []byte) common.Hash {
	key = crypto.Keccak256(key)
	return common.BytesToHash(key)
}

func (c *CacheDB) Get(key []byte) ([]byte, error) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	so := s.getStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		value := so.GetState(s.db, slot)
		return value, nil
	}

	return nil, nil
}

func (c *CacheDB) Delete(key []byte) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	so := s.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		so.SetState(s.db, slot, nil)
	}
}
