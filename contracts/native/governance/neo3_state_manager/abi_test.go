package neo3_state_manager

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateValidatorListParam(t *testing.T) {
	expect := &StateValidatorListParam{
		StateValidators: []string{
			"039b45040cc529966165ef5dff3d046a4960520ce616ae170e265d669e0e2de7f4",
			"0345e2bbda8d3d9e24d1e9ee61df15d4f435f69a44fe012d86e9cf9377baaa42cd",
			"023ccd59ec0fda27844984876ef2d440eca88e45c7401110210f7760cdcc73b5f7",
			"0392fbd1d809a3c62f7dcde8f25454a1570830a21e4b014b3f362a79baf413e115",
		},
		Address: common.HexToAddress("0x3"),
	}

	blob, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	got := new(StateValidatorListParam)
	err = rlp.DecodeBytes(blob, got)
	assert.NoError(t, err)

	assert.Equal(t, expect.StateValidators[0], got.StateValidators[0])
}
