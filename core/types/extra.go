// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// HotstuffDigest represents a hash of "Hotstuff practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine
	HotstuffDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	HotstuffExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	HotstuffExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidHotstuffHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidHotstuffHeaderExtra = errors.New("invalid extra data format")
)

type HotstuffExtra struct {
	StartHeight   uint64           // denote the epoch start height
	EndHeight     uint64           // the epoch end height
	Validators    []common.Address // consensus participants address for next epoch, and in the first block, it contains all genesis validators. keep empty if no epoch change.
	Seal          []byte           // proposer signature
	CommittedSeal [][]byte         // consensus participants signatures and it's size should be greater than 2/3 of validators
	Salt          []byte           // omit empty
}

// Dump only used for debug or test
func (ist *HotstuffExtra) Dump() string {
	seals := []string{}
	for _, v := range ist.CommittedSeal {
		seals = append(seals, hexutil.Encode(v))
	}
	return fmt.Sprintf("{StartHeight: %v, EndHeight: %v, Validators: %v}", ist.StartHeight, ist.EndHeight, ist.Validators)
}

// ExtractHotstuffExtra extracts all values of the HotstuffExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractHotstuffExtra(h *Header) (*HotstuffExtra, error) {
	if h == nil || h.Extra == nil {
		return nil, fmt.Errorf("invalid header")
	}
	return ExtractHotstuffExtraPayload(h.Extra)
}

func ExtractHotstuffExtraPayload(extra []byte) (*HotstuffExtra, error) {
	if len(extra) < HotstuffExtraVanity {
		return nil, ErrInvalidHotstuffHeaderExtra
	}

	var hotstuffExtra *HotstuffExtra
	err := rlp.DecodeBytes(extra[HotstuffExtraVanity:], &hotstuffExtra)
	if err != nil {
		return nil, ErrInvalidHotstuffHeaderExtra
	}
	return hotstuffExtra, nil
}

// HotstuffFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Istanbul hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func HotstuffFilteredHeader(h *Header) *Header {
	newHeader := CopyHeader(h)
	extra, err := ExtractHotstuffExtra(newHeader)
	if err != nil {
		return nil
	}

	// fields of `seal` and `committedSeal` related with consensus voting and DON't participants in hash generating
	extra.Seal = []byte{}
	extra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:HotstuffExtraVanity], payload...)

	return newHeader
}

func HotstuffHeaderFillWithValidators(header *Header, vals []common.Address, epochStartHeight uint64, epochEndHeight uint64) error {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < HotstuffExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, HotstuffExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:HotstuffExtraVanity])

	if vals == nil {
		vals = []common.Address{}
	}
	ist := &HotstuffExtra{
		StartHeight:   epochStartHeight,
		EndHeight:     epochEndHeight,
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
		Salt:          []byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return err
	}
	header.Extra = append(buf.Bytes(), payload...)
	return nil
}

func (h *Header) SetSeal(seal []byte) error {
	extra, err := ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}
	extra.Seal = seal
	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}
	h.Extra = append(h.Extra[:HotstuffExtraVanity], payload...)
	return nil
}

func (h *Header) SetCommittedSeal(committedSeals [][]byte) error {
	extra, err := ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}

	extra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(extra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:HotstuffExtraVanity], payload...)
	return nil
}

func GenerateExtraWithSignature(epochStartHeight, epochEndHeight uint64, vals []common.Address, seal []byte, committedSeal [][]byte) ([]byte, error) {
	var (
		buf   bytes.Buffer
		extra []byte
	)

	extra = append(extra, bytes.Repeat([]byte{0x00}, HotstuffExtraVanity-len(extra))...)
	buf.Write(extra[:HotstuffExtraVanity])

	if vals == nil {
		vals = []common.Address{}
	}
	ist := &HotstuffExtra{
		StartHeight:   epochStartHeight,
		EndHeight:     epochEndHeight,
		Validators:    vals,
		Seal:          seal,
		CommittedSeal: committedSeal,
		Salt:          []byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}
	extra = append(buf.Bytes(), payload...)
	return extra, nil
}
