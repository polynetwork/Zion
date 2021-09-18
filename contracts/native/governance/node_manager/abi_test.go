package node_manager

import (
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

func TestABIMethod(t *testing.T) {
	name := "CheckConsensusSignsEvent"
	abijson := `[
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"signs","type":"uint64"}],"name":"CheckConsensusSignsEvent","type":"event"}
]`
	ab, err := abi.JSON(strings.NewReader(abijson))
	assert.NoError(t, err)

	type Input struct {
		signs         uint64
	}

	expectInput := &Input{
		signs:        uint64(123),
	}

	payload, err := utils.PackEvents(&ab, name, expectInput.signs)
	assert.NoError(t, err)


	inputData := &Input{}
	err = utils.UnpackMethod(&ab, name, inputData, payload)
	////assert.NoError(t, err)
	////
	//assert.True(t, reflect.DeepEqual(expectInput, inputData))

}


