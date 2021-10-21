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

package backend

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
)

// API is a user facing RPC API to allow controlling the address and voting
// mechanisms of the HotStuff scheme.
type API struct {
	chain    consensus.ChainHeaderReader
	hotstuff *backend
}

// Proposals returns the current proposals the node tries to uphold and vote on.
func (api *API) Proposals() map[common.Address]bool {
	api.hotstuff.sigMu.RLock()
	defer api.hotstuff.sigMu.RUnlock()

	proposals := make(map[common.Address]bool)
	for address, auth := range api.hotstuff.proposals {
		proposals[address] = auth
	}
	return proposals
}

// todo: add/del candidate validators approach console or api
// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.hotstuff.sigMu.Lock()
	defer api.hotstuff.sigMu.Unlock()

	api.hotstuff.proposals[address] = auth
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.hotstuff.sigMu.Lock()
	defer api.hotstuff.sigMu.Unlock()

	delete(api.hotstuff.proposals, address)
}
