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
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrKeyLen       = errors.New("cacheDB should only be used for native contract storage")
	ErrBigLen       = errors.New("big int length out of range")
	ErrNotExist     = errors.New("key not exist")
	ErrInvalidBytes = errors.New("invalid bytes")
)

type CacheDB StateDB

// support storage for type of `Address`
func (c *CacheDB) SetAddress(key []byte, value common.Address) error {
	hash := common.BytesToHash(value.Bytes())
	_, _, err := c.customSet(key, hash)
	return err
}

func (c *CacheDB) GetAddress(key []byte) (common.Address, error) {
	_, _, hash, err := c.customGet(key)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(hash.Bytes()), nil
}

func (c *CacheDB) DelAddress(key []byte) error {
	_, _, err := c.customDel(key)
	return err
}

// support storage for type of `Hash`
func (c *CacheDB) SetHash(key []byte, value common.Hash) error {
	_, _, err := c.customSet(key, value)
	return err
}

func (c *CacheDB) GetHash(key []byte) (common.Hash, error) {
	_, _, value, err := c.customGet(key)
	return value, err
}

func (c *CacheDB) DelHash(key []byte) error {
	_, _, err := c.customDel(key)
	return err
}

// support storage for type of `big`
func (c *CacheDB) SetBigInt(key []byte, value *big.Int) error {
	if len(value.Bytes()) > common.HashLength {
		return ErrBigLen
	}
	hash := common.BytesToHash(value.Bytes())
	_, _, err := c.customSet(key, hash)
	return err
}

func (c *CacheDB) GetBigInt(key []byte) (*big.Int, error) {
	_, _, raw, err := c.customGet(key)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(raw[:]), nil
}

func (c *CacheDB) DelBigInt(key []byte) error {
	_, _, err := c.customDel(key)
	return err
}

func (c *CacheDB) SetBytes(key []byte, value []byte) error {
	c.Put(key, value)
	return nil
}

func (c *CacheDB) GetBytes(key []byte) ([]byte, error) {
	return c.Get(key)
}

// custom functions
func (c *CacheDB) customSet(key []byte, value common.Hash) (addr common.Address, slot common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}

	s := (*StateDB)(c)
	s.SetState(addr, slot, value)
	return
}

func (c *CacheDB) customGet(key []byte) (addr common.Address, slot, value common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}
	s := (*StateDB)(c)
	value = s.GetState(addr, slot)
	return
}

func (c *CacheDB) customDel(key []byte) (addr common.Address, slot common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}
	s := (*StateDB)(c)
	s.SetState(addr, slot, common.Hash{})
	return
}

func parseKey(key []byte) (addr common.Address, slot common.Hash, err error) {
	if len(key) <= common.AddressLength {
		return common.Address{}, common.Hash{}, ErrKeyLen
	}
	addr = common.BytesToAddress(key[:common.AddressLength])
	slot = Key2Slot(key[common.AddressLength:])
	return
}

func (c *CacheDB) Put(key []byte, value []byte) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	c.Delete(key)

	s := (*StateDB)(c)
	so := s.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		if len(value) <= common.HashLength-1 {
			c.putValue(so, slot, value, false)
			value = nil
		} else {
			c.putValue(so, slot, value[:common.HashLength-1], true)
			value = value[common.HashLength-1:]
		}

		for len(value) > 0 {
			slot = c.nextSlot(slot)
			if len(value) <= common.HashLength-1 {
				c.putValue(so, slot, value, false)
				break
			} else {
				c.putValue(so, slot, value[:common.HashLength-1], true)
				value = value[common.HashLength-1:]
			}
		}
	}

}

func (c *CacheDB) putValue(so *stateObject, slot common.Hash, value []byte, more bool) {
	if len(value) > common.HashLength-1 {
		panic("value should not exceed 31")
	}

	if more && len(value) != common.HashLength-1 {
		panic("value length should equal 31 when more is true")
	}

	if more {
		value = append([]byte{1}, value...)
	} else {
		padding := make([]byte, common.HashLength-len(value))
		padding[0] = byte(len(value) << 1)
		value = append(padding, value...)
	}

	s := (*StateDB)(c)
	hashValue := common.BytesToHash(value)
	so.SetState(s.db, slot, hashValue)
}

func Key2Slot(key []byte) common.Hash {
	key = crypto.Keccak256(key)
	return common.BytesToHash(key)
}

func (c *CacheDB) nextSlot(slot common.Hash) common.Hash {
	slotBytes := slot.Bytes()
	for offset := common.HashLength - 1; offset >= 0; offset-- {
		slotBytes[offset] = slotBytes[offset] + 1
		if slotBytes[offset] != 0 {
			break
		}
	}

	return Key2Slot(slotBytes)
}

func (c *CacheDB) Get(key []byte) ([]byte, error) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	so := s.getStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		var result []byte
		slot := Key2Slot(key[common.AddressLength:])
		value := so.GetState(s.db, slot)
		meta := value[:][0]
		more := meta&1 == 1
		if more {
			result = append(result, value[1:]...)
		} else {
			if value == (common.Hash{}) {
				return nil, nil
			}
			result = append(result, value[common.HashLength-meta>>1:]...)
		}

		for more {
			slot = c.nextSlot(slot)
			value = so.GetState(s.db, slot)
			meta = value[:][0]
			more = meta&1 == 1
			if more {
				result = append(result, value[1:]...)
			} else {
				result = append(result, value[common.HashLength-meta>>1:]...)
			}
		}

		return result, nil
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
		value := so.GetState(s.db, slot)
		so.SetState(s.db, slot, common.Hash{})
		more := value[:][0]&1 == 1
		for more {
			slot = c.nextSlot(slot)
			value = so.GetState(s.db, slot)
			so.SetState(s.db, slot, common.Hash{})
			more = value[:][0]&1 == 1
		}
	}
}
