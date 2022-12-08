package mock

import (
	"testing"

	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestSimple
func TestSimple(t *testing.T) {
	sys := makeSystem(7)
	sys.Start()
	sys.Close(10)
}

// net scale is 7, 2 of them send fake message of newView with wrong height.
func TestMockCase1(t *testing.T) {
	H, R, fR := uint64(10), uint64(0), uint64(1)

	sys := makeSystem(7)
	sys.RepoHook(2, func(data []byte) {
		if h, r := sys.Leader().api.CurrentSequence(); h == H && r == R {
			var msg core.Message
			if err := rlp.DecodeBytes(data, &msg); err != nil {
				t.Errorf("failed to decode message, err: %v", err)
			}
			if msg.Code == core.MsgTypeNewView {
				t.Log("-----------xxxxxxxxxxxxxxx", fR)
			}
		}
	})
	sys.Start()
}
