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
//
//import (
//	"testing"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/stretchr/testify/assert"
//)
//
//// notice: we need only 3 test case:
//// 1. `newView` send quorumCert, e.g: sendPreCommit, sendCommit
//// 2. `prepare` send msgNewProposal
//// 3. `prepareVote` send vote, e.g: sendPrepareVote, sendPreCommitVote, sendCommitVote
//// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandleMsg
//func TestHandleMsg(t *testing.T) {
//	N, H, R := 4, 5, 0
//
//	view := makeView(H, R)
//	sys := NewTestSystemWithBackend(N, H, R)
//
//	//closer := sys.Run(true)
//	//defer closer()
//
//	v0 := sys.backends[0]
//	r0 := v0.engine
//	vset := v0.Validators(common.EmptyHash, true)
//	sender := vset.GetByIndex(1)
//
//	// invalid message payload
//	{
//		payload := []byte{'1', 'a', 'b'}
//		assert.Equal(t, errFailedDecodeMessage, r0.handleMsg(payload))
//	}
//
//	// invalid sender
//	{
//		msg := &Message{
//			Code:    MsgTypeNewView,
//			View:    view,
//			Msg:     []byte{'1'},
//			Address: common.HexToAddress("0x12"),
//		}
//		payload, err := Encode(msg)
//		assert.NoError(t, err)
//		assert.Equal(t, errInvalidSigner, r0.handleMsg(payload))
//	}
//
//	// invalid msg type
//	{
//		msg := &Message{
//			Code:    100,
//			View:    view,
//			Msg:     []byte{'1'},
//			Address: sender.Address(),
//		}
//		assert.Equal(t, errInvalidMessage, r0.handleCheckedMsg(msg, sender))
//	}
//}
