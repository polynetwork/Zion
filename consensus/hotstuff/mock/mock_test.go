package mock

import (
	"testing"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestSimple
func TestSimple(t *testing.T) {
	sys := makeSystem(7)
	sys.Start()
	stop := make(chan struct{})
	<-stop
}
