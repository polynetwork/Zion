#!/bin/bash

clear
go test -v -count 1 github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core -run TestNewRound
