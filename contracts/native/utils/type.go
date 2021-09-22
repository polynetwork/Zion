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
package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	EmptyAddress = common.Address{}
	EmptyHash    = common.Hash{}
	EmptyBig     = big.NewInt(0)
	ByteFailed   = []byte{'0'}
	ByteSuccess  = []byte{'1'}
)

func Address2Hash(addr common.Address) common.Hash {
	if addr == EmptyAddress {
		return EmptyHash
	}
	return common.BytesToHash(addr.Bytes())
}

func Hash2Address(hash common.Hash) common.Address {
	if hash == EmptyHash {
		return EmptyAddress
	}
	return common.BytesToAddress(hash.Bytes())
}

func Address2Big(addr common.Address) *big.Int {
	if addr == EmptyAddress {
		return EmptyBig
	}
	return new(big.Int).SetBytes(addr.Bytes())
}

func Big2Address(num *big.Int) common.Address {
	if IsZero(num) {
		return EmptyAddress
	}
	return common.BytesToAddress(num.Bytes())
}

func Big2Hash(amount *big.Int) common.Hash {
	if IsZero(amount) {
		return EmptyHash
	}
	return common.BigToHash(amount)
}

func Hash2Big(hash common.Hash) *big.Int {
	if hash == EmptyHash {
		return EmptyBig
	}
	return new(big.Int).SetBytes(hash.Bytes())
}

func Bytes2Big(bz []byte) *big.Int {
	return new(big.Int).SetBytes(bz)
}

func Uint32Bytes(data uint32) []byte {
	num := new(big.Int).SetUint64(uint64(data))
	return num.Bytes()
}

func Uint64Bytes(data uint64) []byte {
	num := new(big.Int).SetUint64(data)
	return num.Bytes()
}

func Bool2Hash(data bool) common.Hash {
	buf := make([]byte, 2)
	value := 1
	if data == false {
		value = 0
	}
	binary.BigEndian.PutUint16(buf, uint16(value))
	return common.BytesToHash(buf)
}

func Hash2Bool(hash common.Hash) bool {
	if hash == EmptyHash {
		return false
	}
	data := binary.BigEndian.Uint16(hash[30:])
	if data == 0 {
		return false
	}
	return true
}

func Uint8Hash(data uint8) common.Hash {
	if data == 0 {
		return EmptyHash
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(data))
	return common.BytesToHash(buf)
}

func Hash2Uint8(hash common.Hash) uint8 {
	if hash == EmptyHash {
		return 0
	}
	data := binary.BigEndian.Uint16(hash[30:])
	if data >= math.MaxUint8 {
		return math.MaxUint8
	}
	return uint8(data)
}

func Uint32Hash(data uint32) common.Hash {
	if data == 0 {
		return EmptyHash
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, data)
	return common.BytesToHash(buf)
}

func Hash2Uint32(hash common.Hash) uint32 {
	if hash == EmptyHash {
		return 0
	}
	return binary.BigEndian.Uint32(hash[28:])
}

func Uint64Hash(data uint64) common.Hash {
	if data == 0 {
		return EmptyHash
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, data)
	return common.BytesToHash(buf)
}

func Hash2Uint64(hash common.Hash) uint64 {
	if hash == EmptyHash {
		return 0
	}
	return binary.BigEndian.Uint64(hash[24:])
}

func CombineAddress2Hash(addr1, addr2 common.Address) common.Hash {
	data := append(addr1.Bytes(), addr2.Bytes()...)
	return sha256.Sum256(data[:])
}

func Bool2Bytes(data bool) []byte {
	if data {
		return ByteSuccess
	}
	return ByteFailed
}

func IsZero(num *big.Int) bool {
	if num == EmptyBig || num.Cmp(EmptyBig) == 0 {
		return true
	}
	return false
}
