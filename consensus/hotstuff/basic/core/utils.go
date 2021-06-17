package core

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
)

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(MsgTypeCommit)})
	return buf.Bytes()
}
