package neo3

import (
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeserializeNeoMakeTxParam(t *testing.T) {
	value := helper.HexToBytes("20000000000000000000000000000000000000000000000000000000000000e269204df8f0fc252d2bd3f21c474510fed914cf5fb5ba98510ddfe83b3d6d5a3715ff14250e76987d838a75310c34bf422ea9f1ac4cc9060e0000000000000014cb569453781497dcb067b73d95b28802cb01553806756e6c6f636b4a149328aec1e84c93855e2fb4a01f5eb7ec15e1abd614e9cdc1efd22c74b5706f0068f79b69b46fa85a0d2035f50500000000000000000000000000000000000000000000000000000000")
	txParam, err := DeserializeNeoMakeTxParam(value)
	assert.Nil(t, err)
	assert.Equal(t, uint64(14), txParam.MakeTxParam.ToChainID)
}
