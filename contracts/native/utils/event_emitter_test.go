package utils

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestEventEmitter(t *testing.T) {
	name := "propose"
	abijson := `[
	{"type":"event","anonymous":false,"name":"` + name + `","inputs":[{"indexed":true,"name":"proposer","type":"address"},{"indexed":true,"name":"proposalId","type":"address"},{"indexed":false,"name":"value","type":"uint256"}]}
]`

	ab, err := abi.JSON(strings.NewReader(abijson))
	assert.NoError(t, err)

	contract := common.HexToAddress("0x05")
	blockNo := uint64(36)
	stateDB := NewTestStateDB()
	emmitter := NewEventEmitter(contract, blockNo, stateDB)

	proposer := common.HexToAddress("0x12")
	proposalID := common.HexToAddress("0x18")
	value := big.NewInt(120)

	topics := make([]common.Hash, 3)
	topics[0] = ab.Events[name].ID
	topics[1] = common.BytesToHash(proposer.Bytes())
	topics[2] = common.BytesToHash(proposalID.Bytes())

	emmitter.Event(topics, value.Bytes())

	hash := stateDB.BlockHash()
	data := stateDB.GetLogs(hash)
	assert.Equal(t, 1, len(data))

	event := data[0]
	assert.Equal(t, len(topics), len(event.Topics))
	assert.Equal(t, ab.Events[name].ID, event.Topics[0])
	assert.Equal(t, proposer, common.BytesToAddress(event.Topics[1].Bytes()))
	assert.Equal(t, proposalID, common.BytesToAddress(event.Topics[2].Bytes()))

	assert.Equal(t, value, new(big.Int).SetBytes(event.Data))
}
