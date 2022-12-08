package mock

import (
	"testing"
	"time"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestSimple
func TestSimple(t *testing.T) {
	sys := makeSystem(4)
	sys.Start()
	timer := time.AfterFunc(30*time.Second, func() {
		sys.Stop()
	})
	<-timer.C
}
