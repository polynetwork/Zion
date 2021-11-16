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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestCmpAddressList(t *testing.T) {
	addr1 := common.HexToAddress("0x12")
	addr2 := common.HexToAddress("0x13")
	addr3 := common.HexToAddress("0x9a")
	addr4 := common.HexToAddress("0xbf")
	addr5 := common.HexToAddress("0x175")
	addr6 := common.HexToAddress("0x666")
	addr7 := common.EmptyAddress

	var testcases = []struct {
		L1     []common.Address
		L2     []common.Address
		Expect bool
	}{
		{
			[]common.Address{addr1, addr2, addr3, addr4},
			[]common.Address{addr1, addr3, addr4, addr2},
			true,
		},
		{
			[]common.Address{addr5, addr6},
			[]common.Address{addr1, addr5, addr6},
			false,
		},
		{
			[]common.Address{addr1, addr6},
			[]common.Address{addr1, addr6, addr1},
			true,
		},
		{
			[]common.Address{addr7},
			[]common.Address{addr7},
			true,
		},
		{
			[]common.Address{},
			[]common.Address{addr7},
			false,
		},
	}

	for _, v := range testcases {
		assert.Equal(t, v.Expect, compareVals(v.L1, v.L2))
	}
}
