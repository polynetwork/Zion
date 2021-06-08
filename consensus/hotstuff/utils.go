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

package hotstuff

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func CheckValidatorSignature(valSet ValidatorSet, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	if _, val := valSet.GetByAddress(signer); val != nil {
		return val.Address(), nil
	}

	return common.Address{}, errUnauthorizedAddress
}

func PrepareVoteSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(MsgTypeVote)})
	return buf.Bytes()
}

// PrepareCommittedSeal returns a committed seal for the given hash
// todo: change `MsgTypeCommit` to `MsgTypeDecide`
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(MsgTypeCommit)})
	return buf.Bytes()
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func Signers(header *types.Header) ([]common.Address, error) {
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return []common.Address{}, err
	}

	var addrs []common.Address
	proposalSeal := PrepareCommittedSeal(header.Hash())

	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			return nil, errInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (s *roundState) decode(msg p2p.Msg) ([]byte, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, errDecodeFailed
	}
	return data, nil
}

func (c *roundState) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	valSet := c.snap.ValSet
	return CheckValidatorSignature(valSet, data, sig)
}

func (s *roundState) broadcast(payload []byte) error {
	s.msgCh <- &InnerMsg{Payload:payload}
	return s.gossip(payload)
}

// Broadcast implements istanbul.Backend.Gossip
// todo:
//  1. record message
//  3. peer msg lru
func (s *roundState) gossip(payload []byte) error {
	targets := make(map[common.Address]bool)
	for _, v := range s.snap.ValSet.List() {
		if v.Address() != s.address {
			targets[v.Address()] = true
		}
	}
	if s.broadcaster != nil && len(targets) > 0 {
		ps := s.broadcaster.FindPeers(targets)
		for _, p := range ps {
			go p.Send(P2PHotstuffMsg, payload)
		}
	}
	return nil
}

func (s *roundState) unicast(to common.Address, payload []byte) error {
	if to == s.address {
		s.msgCh <- &InnerMsg{Payload: payload}
		return nil
	}
	peer := s.broadcaster.FindPeer(to)
	if peer == nil {
		return fmt.Errorf("can't find p2p peer of %s", to.Hex())
	}
	go peer.Send(P2PHotstuffMsg, payload)
	return nil
}

func (s *roundState) extend(ancestor, block *types.Header) bool {
	b := block
	for b = s.fetchParentHeader(b); b != nil; {
		if b.Hash() == ancestor.Hash() {
			return true
		}
	}
	return false
}

func (s *roundState) fetchParentHeader(header *types.Header) *types.Header {
	return s.fetchHeader(header.ParentHash)
}

func (s *roundState) fetchHeader(hash common.Hash) *types.Header {
	return s.store.GetHeader(hash)
}

func (s *roundState) fetchParentBlockWithHeader(header *types.Header) *types.Block {
	return s.fetchBlock(header.ParentHash, header.Number.Uint64() - 1)
}

func (s *roundState) fetchParentBlock(block *types.Block) *types.Block {
	return s.fetchBlock(block.ParentHash(), block.NumberU64() - 1)
}

// todo: 从chainReader中拿到block，阻塞式等待
func (s *roundState) fetchBlock(hash common.Hash, view uint64) *types.Block {
	return s.store.GetBlock(hash, view)
	//if block != nil {
	//	return block, nil
	//}
	//
	//s.waitProposal.Wait()
	//block = s.chain.GetBlock(hash, view)
	//if block == nil {
	//	return nil, fmt.Errorf("block (%s %d) not arrived", hash.Hex(), view)
	//}
	//return block, nil
}

func (s *roundState) finalizeMessage(encoder rlp.Encoder, typ MsgType) ([]byte, error) {
	var buf bytes.Buffer
	if err := encoder.EncodeRLP(&buf); err != nil {
		return nil, err
	}

	msg := &Message{
		Code:      typ,
		Msg:       buf.Bytes(),
	}

	// Add sender address
	msg.Address = s.address

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = s.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// GetLeader get leader with view number approach round robin
func (s *roundState) getLeader() common.Address {
	s.snap.ValSet.CalcProposer(common.Address{}, s.qcHigh.Number.Uint64())
	return s.snap.ValSet.GetProposer().Address()
}

func (s *roundState) isLeader() bool {
	leader := s.getLeader()
	return leader == s.address
}

func (s *roundState) isValidator(val common.Address) bool {
	index, _ := s.snap.ValSet.GetByAddress(val)
	return index >= 0
}
