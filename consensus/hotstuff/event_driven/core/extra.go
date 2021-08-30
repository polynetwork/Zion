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

package core

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// 这里存在一个问题，就是是否需要将epochid,round以及justifyQC放到区块头内，如果放到区块头内，这些内容只在共识层面有用，其他的地方都没有用。
// 加入系统刚启动，这时候round==height, 不同的节点等待时间不一样，但是最终会进入到同一个round
// 如果4个节点有3个在同一个round，而另一个节点刚进入，round与当前节点不一致，怎么驱动其进入到比较高的round的问题，接收到新的qc时会进入。
// 如果4个节点有两个在一个round1，另外两个在另一个round2, 因为形不成qc/tc，这样一来，必然要将当前epoch以及round等信息持久化。因为即便将共识信息
// 写入到区块也会存在这种问题
// 所以应该引入3个机制:
// 1. 超时消息重发
// 2. 记录snapshot，snap应该包含当前所在epoch,currentRound,currentHeight,lastVoteRound,highestCommitRound,
// block tree pure之后剩下的区块(最多3个，qc不需要记录，可以根据这些区块形成新的qc)
func generateExtra(header *types.Header, valSet hotstuff.ValidatorSet, epoch uint64, round *big.Int) ([]byte, error) {
	var (
		buf  bytes.Buffer
		vals = valSet.AddressList()
	)

	fmt.Println("-----------2")
	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.HotstuffExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity-len(header.Extra))...)
	}
	fmt.Println("-----------3")
	buf.Write(header.Extra[:types.HotstuffExtraVanity])
	fmt.Println("-----------4")
	salt := &ExtraSalt{
		Epoch: epoch,
		Round: round,
	}
	fmt.Println("-----------5")
	saltEnc, err := Encode(salt)
	if err != nil {
		return nil, err
	}
	fmt.Println("-----------6")
	ist := &types.HotstuffExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
		Salt:          saltEnc,
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

func extraProposal(proposal hotstuff.Proposal) (salt *ExtraSalt, qc *hotstuff.QuorumCert, err error) {
	block, ok := proposal.(*types.Block)
	if !ok {
		return nil, nil, errProposalConvert
	}
	return extraHeader(block.Header())
}

func extraHeader(h *types.Header) (salt *ExtraSalt, qc *hotstuff.QuorumCert, err error) {
	var extra *types.HotstuffExtra

	qc = new(hotstuff.QuorumCert)
	qc.Hash = h.Hash()
	qc.Proposer = h.Coinbase
	qc.Extra = h.Extra

	if extra, err = types.ExtractHotstuffExtra(h); err != nil {
		return
	}
	if err = rlp.DecodeBytes(extra.Salt, &salt); err != nil {
		return
	}

	qc.View = &hotstuff.View{
		Height: h.Number,
		Round:  salt.Round,
	}
	return
}
