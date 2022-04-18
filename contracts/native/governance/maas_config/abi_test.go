package maas_config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestABIShowJonString(t *testing.T) {
	t.Log(MaasConfigABI)
	for name, v := range ABI.Methods {
		t.Logf("method %s, id %s", name, hexutil.Encode(v.ID))
	}
	t.Log("\n")
}

var testAddresses = []common.Address{
	common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0125"),
	common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0126"),
	common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0127"),
	common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0128"),
	common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0129"),
}

func TestABIMethodContractName(t *testing.T) {

	enc, err := utils.PackOutputs(ABI, MethodName, contractName)
	assert.NoError(t, err)
	params := new(MethodContractNameOutput)
	assert.NoError(t, utils.UnpackOutputs(ABI, MethodName, params, enc))
	assert.Equal(t, contractName, params.Name)
}

func TestABIMethodChangeOwnerInput(t *testing.T) {
	expect := &MethodChangeOwnerInput{Addr: testAddresses[0]}
	enc, err := expect.Encode()
	assert.NoError(t, err)
	methodId := hexutil.Encode(crypto.Keccak256([]byte("changeOwner(address)"))[:4])
	t.Log("expected methodId of changeOwner ", methodId)
	t.Log("actual methodId of changedOwner ", hexutil.Encode(enc)[:10])
	assert.Equal(t, methodId, hexutil.Encode(enc)[:10])
	got := new(MethodChangeOwnerInput)
	assert.NoError(t, got.Decode(enc))
	assert.Equal(t, expect, got)
}

func TestABIMethodChangeOwnerOutput(t *testing.T) {
	var cases = []struct {
		Result bool
	}{
		{
			Result: true,
		},
		{
			Result: false,
		},
	}

	for _, testCase := range cases {
		output := &MethodChangeOwnerOutput{Success: testCase.Result}
		enc, err := output.Encode()
		assert.NoError(t, err)

		got := new(MethodChangeOwnerOutput)
		err = got.Decode(enc)
		assert.NoError(t, err)

		assert.Equal(t, output, got)
	}
}


func TestABIMethodGetOwnerOutput(t *testing.T) {
	var cases = []struct {
		Addr common.Address
	}{
		{
			Addr: testAddresses[0],
		},
		{
			Addr: testAddresses[1],
		},
	}

	for _, testCase := range cases {
		output := &MethodGetOwnerOutput{Addr: testCase.Addr}
		enc, err := output.Encode()
		assert.NoError(t, err)

		got := new(MethodGetOwnerOutput)
		err = got.Decode(enc)
		assert.NoError(t, err)

		assert.Equal(t, output, got)
	}
}

func TestMethodBlockAccountInput(t *testing.T) {
	var cases = []struct{
		Addr common.Address
		DoBlock bool
	} {
		{
			Addr: testAddresses[0],
			DoBlock: false,
		},
		{
			Addr: testAddresses[1],
			DoBlock: true,
		},
	}

	for _, testCase := range cases {
		output := &MethodBlockAccountInput{Addr: testCase.Addr, DoBlock: testCase.DoBlock}
		enc, err := output.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte("blockAccount(address,bool)"))[:4])
		t.Log("expected methodId of blockAccount ", methodId)
		t.Log("actual methodId of blockAccount ", hexutil.Encode(enc)[:10])
		t.Log("\n")

		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodBlockAccountInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodBlockAccountOutput(t *testing.T){
	var cases = []struct {
		Success bool
	}{
		{true},
		{false},
	}

	for _, testCase := range cases {
		output := &MethodBlockAccountOutput{Success: testCase.Success}
		enc, err := output.Encode()
		assert.NoError(t, err)

		got := new(MethodBlockAccountOutput)
		err = got.Decode(enc)
		assert.NoError(t, err)

		assert.Equal(t, got, output)
	}
}

func TestMethodIsBlockedInput(t *testing.T) {
	var cases = []struct{Addr common.Address} {
		{testAddresses[1]},
		{testAddresses[0]},
	}

	for _, testCase := range cases {
		output := &MethodIsBlockedInput{testCase.Addr}
		enc, err := output.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte("isBlocked(address)"))[:4])
		t.Log("expected methodId of isBlocked ", methodId)
		t.Log("actual methodId of isBlocked ", hexutil.Encode(enc)[:10])
		t.Log("\n")

		got := new(MethodIsBlockedInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodIsBlockedOutput(t *testing.T) {
	var cases = []struct{Success bool} {
		{true},
		{false},
	}

	for _, testCase := range cases {
		output := &MethodIsBlockedOutput{testCase.Success}
		enc, err := output.Encode()
		assert.NoError(t, err)

		got := new(MethodIsBlockedOutput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodGetBlacklistOutput(t *testing.T) {
	var cases = []struct{Result string} {
		{"Success"},
		{"Fail"},
	}

	for _, testCase := range cases {
		output := &MethodGetBlacklistOutput{testCase.Result}
		enc, err := output.Encode()
		assert.NoError(t, err)

		got := new(MethodGetBlacklistOutput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}
