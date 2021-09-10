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

package backup


//// todo: 如果节点此时宕机怎么办？还是说允许所有的节点一起提交区块
//if existProposal := c.backend.GetProposal(committedBlock.Hash()); existProposal == nil {
//	//if c.isSelf(committedBlock.Coinbase()) {
//	//
//	//}
//
//} else {
//	c.logger.Trace("block already synced to chain reader", "")
//}


//func (c *core) forwardProposal() error {
//	logger := c.newSenderLogger("MSG_FORWARD_PROPOSAL")
//
//	proposal := c.smr.Proposal()
//	if proposal == nil || !bigEq(proposal.Number(), c.smr.Height()) {
//		return fmt.Errorf("no forward proposal")
//	}
//	justifyQC := c.smr.HighQC()
//	view := c.currentView()
//	msg := &MsgProposal{
//		Epoch:     c.smr.Epoch(),
//		View:      view,
//		Proposal:  proposal,
//		JustifyQC: justifyQC,
//	}
//	c.encodeAndBroadcast(MsgTypeProposal, msg)
//	logger.Trace("Forward proposal", "hash", msg.Proposal.Hash(), "justifyQC hash", msg.JustifyQC.Hash)
//	return nil
//}

//
//func (e *core) handleQC(src hotstuff.Validator, data *hotstuff.Message) error {
//	logger := e.newMsgLogger()
//
//	var (
//		qc     *hotstuff.QuorumCert
//		msgTyp = MsgTypeQC
//	)
//	if err := data.Decode(&qc); err != nil {
//		logger.Trace("Failed to decode", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	if err := e.signer.VerifyQC(qc, e.valset); err != nil {
//		logger.Trace("Failed to verify qc", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	if err := e.processQC(qc); err != nil {
//		logger.Trace("Failed to process qc", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	logger.Trace("Accept QC", "msg", msgTyp, "src", src.Address(), "qc", qc.Hash, "view", qc.View)
//	return nil
//}

// commit 3-chain
//
//committedBlock := c.blkPool.GetCommitBlock(lockQC.Hash)
//if committedBlock == nil {
//	c.logger.Trace("[Commit 3-Chain], failed to get commit block", "lockQC view", lockQC.View)
//	return
//}
//
//round := lockQC.Round()
//if existProposal := c.backend.GetProposal(committedBlock.Hash()); c.IsProposer() && existProposal == nil {
//	if err := c.backend.Commit(committedBlock); err != nil {
//		c.logger.Trace("[Commit 3-Chain], failed to commit", "err", err)
//	} else {
//		c.logger.Trace("[Commit 3-Chain], leader commit", "address", c.address, "hash", committedBlock.Hash(), "number", committedBlock.Number())
//	}
//}
//
//c.updateHighestCommittedRound(round)
//puredBlocks := c.blkPool.Pure(committedBlock.Hash())
//c.logger.Trace("[Commit 3-Chain], pured blocks", "hash lists", puredBlocks)

//lockQC := c.smr.LockQC()
//if lockQC == nil {
//	c.logger.Trace("[Commit 3-Chain]", "err", "lockQC is nil")
//	return
//}
//
//branch := c.blkPool.GetCommitBranch(lockQC.Hash)
//if branch == nil {
//	c.logger.Trace("[Commit 3-Chain], failed to get branch", "lockQC view", lockQC.View, "lockQC hash", lockQC.Hash)
//	return
//}
//
//for _, block := range branch {
//	if exist := c.chain.GetBlockByHash(block.Hash()); exist != nil {
//		continue
//	}
//	if !c.isSelf(block.Coinbase()) {
//		continue
//	}
//	if err := c.backend.Commit(block); err != nil {
//		c.logger.Trace("[Commit 3-Chain], failed to commit", "err", err)
//	} else {
//		c.logger.Trace("[Commit 3-Chain], commit success!", "hash", block.Hash(), "number", block.Number(), "proposer", block.Coinbase())
//	}
//}
//
//c.updateHighestCommittedRound(lockQC.Round())
//c.blkPool.Pure(lockQC.Hash)

//var parent *types.Block
//if w.currentBlock == nil {
//	parent = w.chain.CurrentBlock()
//} else {
//	parent = w.currentBlock
//}
//if parent.Time() >= uint64(timestamp) {
//	timestamp = int64(parent.Time() + 1)
//}