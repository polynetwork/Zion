package hotstuff

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

var testValset = newDefaultSet([]common.Address{
	common.BigToAddress(big.NewInt(1)),
	common.BigToAddress(big.NewInt(2)),
	common.BigToAddress(big.NewInt(3)),
	common.BigToAddress(big.NewInt(4)),
})

func TestF(t *testing.T) {
	f := testValset.F()
	t.Log(f)
}
