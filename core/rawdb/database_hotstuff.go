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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	keyCurEpoch    = []byte("ht-cur-ep")
	keyEpochPrefix = []byte("ht-ep")
)

func WriteCurrentEpochHash(db ethdb.KeyValueWriter, hash common.Hash) error {
	return db.Put(keyCurEpoch, hash.Bytes())
}

func ReadCurrentEpochHash(db ethdb.Reader) (common.Hash, error) {
	blob, err := db.Get(keyCurEpoch)
	if err != nil {
		return common.EmptyHash, err
	}
	return common.BytesToHash(blob), nil
}

func WriteEpoch(db ethdb.KeyValueWriter, hash common.Hash, blob []byte) error {
	key := append(keyEpochPrefix, hash[:]...)
	return db.Put(key, blob)
}

func ReadEpoch(db ethdb.Reader, hash common.Hash) ([]byte, error) {
	key := append(keyEpochPrefix, hash[:]...)
	return db.Get(key)
}
