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

package fetcher

import "github.com/ethereum/go-ethereum/common"

// txAnnounce is the notification of the availability of a batch
// of new transactions in the network.
type nodeAnnounce struct {
	origin string        	// Identifier of the peer originating the notification
	exist []common.Address // Batch of transaction hashes being announced
}


type NodeFetcher struct {
	quit chan struct{}
}

func (f *NodeFetcher) Notify(peer string, filter []common.Address) error {
	return nil
}

func (f *NodeFetcher) Enqueue(peer string, filter []common.Address, direct bool) error {
	return nil
}

// hash notifications and block fetches until termination requested.
func (f *NodeFetcher) Start() {
	go f.loop()
}

// Stop terminates the announcement based synchroniser, canceling all pending
// operations.
func (f *NodeFetcher) Stop() {
	close(f.quit)
}

func (f *NodeFetcher) loop() {

}