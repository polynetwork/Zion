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

package hotstuff

type SelectProposerPolicy uint64

const (
	RoundRobin SelectProposerPolicy = iota
	Sticky
	VRF
)

type Config struct {
	RequestTimeout uint64               `toml:",omitempty"` // The timeout for each Istanbul round in milliseconds.
	BlockPeriod    uint64               `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second for basic hotstuff and mill-seconds for event-driven
	LeaderPolicy   SelectProposerPolicy `toml:",omitempty"` // The policy for speaker selection
	Test           bool                 `toml:",omitempty"`
	Epoch          uint64               `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
}

// todo: modify request timeout, and miner recommit default value is 3s. recommit time should be > blockPeriod
var DefaultBasicConfig = &Config{
	RequestTimeout: 6000,
	BlockPeriod:    3,
	LeaderPolicy:   RoundRobin,
	Epoch:          30000,
	Test:           false,
}

var DefaultEventDrivenConfig = &Config{
	RequestTimeout: 4000,
	BlockPeriod:    2000,
	LeaderPolicy:   RoundRobin,
	Epoch:          0,
	Test:           false,
}
