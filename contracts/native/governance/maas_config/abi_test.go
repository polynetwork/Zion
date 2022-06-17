package maas_config

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
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
		output := &MethodBoolOutput{Success: testCase.Result}
		enc, err := output.Encode(MethodChangeOwner)
		assert.NoError(t, err)

		got := new(MethodBoolOutput)
		err = got.Decode(enc, MethodChangeOwner)
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
		output := &MethodAddressOutput{Addr: testCase.Addr}
		enc, err := output.Encode(MethodGetOwner)
		assert.NoError(t, err)

		got := new(MethodAddressOutput)
		err = got.Decode(enc, MethodGetOwner)
		assert.NoError(t, err)

		assert.Equal(t, output, got)
	}
}

func TestMethodBlockAccountInput(t *testing.T) {
	var cases = []struct {
		Addr    common.Address
		DoBlock bool
	}{
		{
			Addr:    testAddresses[0],
			DoBlock: false,
		},
		{
			Addr:    testAddresses[1],
			DoBlock: true,
		},
	}

	for _, testCase := range cases {
		output := &MethodBlockAccountInput{Addr: testCase.Addr, DoBlock: testCase.DoBlock}
		enc, err := output.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte("blockAccount(address,bool)"))[:4])
		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodBlockAccountInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodBlockAccountOutput(t *testing.T) {
	var cases = []struct {
		Success bool
	}{
		{true},
		{false},
	}

	for _, testCase := range cases {
		output := &MethodBoolOutput{Success: testCase.Success}
		enc, err := output.Encode(MethodBlockAccount)
		assert.NoError(t, err)

		got := new(MethodBoolOutput)
		err = got.Decode(enc, MethodBlockAccount)
		assert.NoError(t, err)

		assert.Equal(t, got, output)
	}
}

func TestMethodIsBlockedInput(t *testing.T) {
	var cases = []struct{ Addr common.Address }{
		{testAddresses[1]},
		{testAddresses[0]},
	}

	for _, testCase := range cases {
		input := &MethodIsBlockedInput{testCase.Addr}
		enc, err := input.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte("isBlocked(address)"))[:4])
		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodIsBlockedInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, input)
	}
}

func TestMethodIsBlockedOutput(t *testing.T) {
	var cases = []struct{ Success bool }{
		{true},
		{false},
	}

	for _, testCase := range cases {
		output := &MethodBoolOutput{testCase.Success}
		enc, err := output.Encode(MethodIsBlocked)
		assert.NoError(t, err)

		got := new(MethodBoolOutput)
		err = got.Decode(enc, MethodIsBlocked)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodGetBlacklistOutput(t *testing.T) {
	var cases = []struct{ Result string }{
		{"Success"},
		{"Fail"},
	}

	for _, testCase := range cases {
		output := &MethodStringOutput{testCase.Result}
		enc, err := output.Encode(MethodGetBlacklist)
		assert.NoError(t, err)

		got := new(MethodStringOutput)
		err = got.Decode(enc, MethodGetBlacklist)
		assert.NoError(t, err)
		assert.Equal(t, got, output)
	}
}

func TestMethodEnableGasManageInput(t *testing.T) {
	var cases = []struct{ DoEnable bool }{
		{true},
		{false},
	}

	for _, testCase := range cases {
		input := &MethodEnableGasManageInput{testCase.DoEnable}
		enc, err := input.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte(MethodEnableGasManage + "(bool)"))[:4])
		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodEnableGasManageInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, input)
	}
}

func TestMethodSetGasManagerInput(t *testing.T) {
	var cases = []struct {
		Addr    common.Address
		IsWhite bool
	}{
		{
			Addr:    testAddresses[0],
			IsWhite: false,
		},
		{
			Addr:    testAddresses[1],
			IsWhite: true,
		},
	}

	for _, testCase := range cases {
		input := &MethodSetGasManagerInput{testCase.Addr, testCase.IsWhite}
		enc, err := input.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte(MethodSetGasManager + "(address,bool)"))[:4])
		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodSetGasManagerInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, input)
	}
}

func TestMethodIsGasManagerInput(t *testing.T) {
	var cases = []struct{ Addr common.Address }{
		{testAddresses[1]},
		{testAddresses[0]},
	}

	for _, testCase := range cases {
		input := &MethodIsGasManagerInput{testCase.Addr}
		enc, err := input.Encode()
		assert.NoError(t, err)

		methodId := hexutil.Encode(crypto.Keccak256([]byte(MethodIsGasManager + "(address)"))[:4])
		assert.Equal(t, methodId, hexutil.Encode(enc)[:10])

		got := new(MethodIsGasManagerInput)
		err = got.Decode(enc)
		assert.NoError(t, err)
		assert.Equal(t, got, input)
	}
}
