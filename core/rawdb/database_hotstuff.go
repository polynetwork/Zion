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
	keyCurEpoch    = []byte("hs-cur-ep-ht")
	keyEpochPrefix = []byte("hs-ep")
)

func WriteCurrentEpochHeight(db ethdb.KeyValueWriter, height uint64) error {
	return db.Put(keyCurEpoch, uint64Bytes(height))
}

func ReadCurrentEpochHeight(db ethdb.Reader) (uint64, error) {
	blob, err := db.Get(keyCurEpoch)
	if err != nil {
		return 0, err
	}
	return bytes2uint64(blob), nil
}

func WriteEpoch(db ethdb.KeyValueWriter, height uint64, blob []byte) error {
	key := keyHeight(height)
	return db.Put(key, blob)
}

func ReadEpoch(db ethdb.Reader, height uint64) ([]byte, error) {
	key := keyHeight(height)
	return db.Get(key)
}

func keyHeight(height uint64) []byte {
	dat := uint64Bytes(height)
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
