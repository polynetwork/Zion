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

package rawdb

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	keyEpochPrefix = []byte("hs-ep")
)

func WriteEpoch(db ethdb.KeyValueWriter, id uint64, blob []byte) error {
	key := keyId(id)
	return db.Put(key, blob)
}

func ReadEpoch(db ethdb.Reader, id uint64) ([]byte, error) {
	key := keyId(id)
	return db.Get(key)
}

func keyId(id uint64) []byte {
	dat := uint64Bytes(id)
	return append(keyEpochPrefix, dat...)
}

func uint64Bytes(dat uint64) []byte {
	if dat == 0 {
		dat = math.MaxUint64
	}
	return new(big.Int).SetUint64(dat).Bytes()
}

func bytes2uint64(blob []byte) uint64 {
	dat := new(big.Int).SetBytes(blob).Uint64()
	if dat == math.MaxUint64 {
		dat = 0
	}
	return dat
}
