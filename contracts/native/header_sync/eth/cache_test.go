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

package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	docache := func(epoch uint64, data []uint32) {
		caches.addCache(epoch, data)
		caches.deleteCaches()
	}

	docache(12, []uint32{1, 2, 3})
	docache(13, []uint32{1, 2, 3})
	assert.Equal(t, 2, len(caches.items))

	docache(14, []uint32{1, 2, 3})
	docache(15, []uint32{1, 2, 3})
	docache(16, []uint32{1, 2, 3})
	assert.Equal(t, 3, len(caches.items))

	var ok bool
	_, ok = caches.items[14]
	assert.True(t, ok)
	_, ok = caches.items[15]
	assert.True(t, ok)
	_, ok = caches.items[16]
	assert.True(t, ok)
}
