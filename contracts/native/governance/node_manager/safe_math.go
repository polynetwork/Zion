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

package node_manager

import (
	"errors"
	"math/big"
)

type Dec struct {
	I *big.Int
}

var (
	DecNegError = errors.New("decimal is negative")
	DecNilError = errors.New("decimal is nil")
)

func NewDecFromBigInt(i *big.Int) Dec {
	return Dec{I: i}
}

func (t Dec) IsNil() bool       { return t.I == nil }           // is token nil
func (t Dec) IsZero() bool      { return (t.I).Sign() == 0 }    // is equal to zero
func (t Dec) IsNegative() bool  { return (t.I).Sign() == -1 }   // is negative
func (t Dec) IsPositive() bool  { return (t.I).Sign() == 1 }    // is positive
func (t Dec) Equal(t2 Dec) bool { return (t.I).Cmp(t2.I) == 0 } // equal decimals
func (t Dec) GT(t2 Dec) bool    { return (t.I).Cmp(t2.I) > 0 }  // greater than
func (t Dec) GTE(t2 Dec) bool   { return (t.I).Cmp(t2.I) >= 0 } // greater than or equal
func (t Dec) LT(t2 Dec) bool    { return (t.I).Cmp(t2.I) < 0 }  // less than
func (t Dec) LTE(t2 Dec) bool   { return (t.I).Cmp(t2.I) <= 0 } // less than or equal

// BigInt returns a copy of the underlying big.Int.
func (t Dec) BigInt() *big.Int {
	cp := new(big.Int)
	return cp.Set(t.I)
}

// addition
func (t Dec) Add(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	res := new(big.Int).Add(t.I, t2.I)

	return Dec{res}, nil
}

// subtraction
func (t Dec) Sub(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	res := new(big.Int).Sub(t.I, t2.I)
	r := Dec{res}
	if r.IsNegative() {
		return Dec{nil}, DecNegError
	}
	return r, nil
}

// multiplication with token decimal
func (t Dec) Mul(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	mul := new(big.Int).Mul(t.I, t2.I)

	return Dec{mul}, nil
}

// multiplication with token decimal
func (t Dec) MulWithTokenDecimal(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	mul := new(big.Int).Mul(t.I, t2.I)
	chopped := new(big.Int).Div(mul, TokenDecimal)

	return Dec{chopped}, nil
}

// multiplication with percent decimal
func (t Dec) MulWithPercentDecimal(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	mul := new(big.Int).Mul(t.I, t2.I)
	chopped := new(big.Int).Div(mul, PercentDecimal)

	return Dec{chopped}, nil
}

// div with token decimal
func (t Dec) DivWithTokenDecimal(t2 Dec) (Dec, error) {
	if t.IsNil() || t2.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() || t2.IsNegative() {
		return Dec{nil}, DecNegError
	}
	// multiply precision
	mul := new(big.Int).Mul(t.I, TokenDecimal)
	div := new(big.Int).Div(mul, t2.I)

	return Dec{div}, nil
}

// div uint64 without ratio
func (t Dec) DivUint64(i uint64) (Dec, error) {
	if t.IsNil() {
		return Dec{nil}, DecNilError
	}
	if t.IsNegative() {
		return Dec{nil}, DecNegError
	}

	div := new(big.Int).Div(t.I, new(big.Int).SetUint64(i))

	return Dec{div}, nil
}
