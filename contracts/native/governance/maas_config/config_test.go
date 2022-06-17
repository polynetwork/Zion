package maas_config

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitMaasConfig()
	os.Exit(m.Run())
}

var (
	testStateDB  *state.StateDB
	testEmptyCtx *native.NativeContract

	testSupplyGas uint64 = 100000000000000000
	testCaller    common.Address
)

func generateNativeContractRef(origin common.Address, blockNum int) *native.ContractRef {
	token := make([]byte, common.HashLength)
	rand.Read(token)
	hash := common.BytesToHash(token)
	return native.NewContractRef(testStateDB, origin, origin, big.NewInt(int64(blockNum)), hash, testSupplyGas, nil)
}

func generateNativeContract(origin common.Address, blockNum int) *native.NativeContract {
	ref := generateNativeContractRef(origin, blockNum)
	return native.NewNativeContract(testStateDB, ref)
}

func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	testCaller = testAddresses[0]
}

func TestChangeAndGetOwner(t *testing.T) {
	type TestCase struct {
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Expect        error
		ReturnData    bool
	}

	cases := []*TestCase{
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodChangeOwnerInput{Addr: testAddresses[0]}
				c.Payload, _ = input.Encode()
			},
			Expect:     nil,
			ReturnData: true,
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				payload, err := utils.PackMethod(ABI, MethodGetOwner)
				assert.NoError(t, err)
				c.Payload = payload
			},
			Expect: nil,
		},
	}

	resetTestContext()
	ctx := generateNativeContract(testCaller, 3)

	for k, v := range cases {
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		if k == 0 {
			res, err := strconv.ParseBool(string(result))
			assert.NoError(t, err)
			t.Log("changeOwner result: ", res)
			assert.Equal(t, v.ReturnData, res)
		}
		if k == 1 {
			t.Log("getOwner result: ", hexutil.Encode(result))
			assert.Equal(t, common.HexToAddress(common.Bytes2Hex(result)), testAddresses[0])
		}
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func setDefaultOwner(ctx *native.NativeContract) {
	input := &MethodChangeOwnerInput{Addr: testAddresses[0]}
	payload, _ := input.Encode()
	result, _, err := ctx.ContractRef().NativeCall(testCaller, this, payload)
	if err != nil {
		panic(err)
	}
	res, _ := strconv.ParseBool(string(result))
	if !res {
		panic("setDefaultOwner error")
	}
}

func TestMethodBlockAccount(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		ReturnData    []byte
		Expect        error
	}
	cases := []*TestCase{
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodBlockAccountInput{Addr: testAddresses[3], DoBlock: true}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'0'},
			Expect:     errors.New("invalid authority for owner"),
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodBlockAccountInput{Addr: testAddresses[3], DoBlock: true}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'1'},
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodBlockAccountInput{Addr: testAddresses[3], DoBlock: false}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'1'},
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodBlockAccountInput{Addr: testAddresses[4], DoBlock: false}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'1'},
			Expect:     nil,
		},
	}

	resetTestContext()
	for _, v := range cases {
		ctx := generateNativeContract(testCaller, v.BlockNum)

		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log("blockAccount result: ", string(result))
		assert.Equal(t, v.ReturnData, result)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func blockTestAccount(ctx *native.NativeContract) {
	input := &MethodBlockAccountInput{Addr: testAddresses[3], DoBlock: true}
	payload, _ := input.Encode()
	result, _, err := ctx.ContractRef().NativeCall(testCaller, this, payload)
	if err != nil || result[0] != utils.ByteSuccess[0] {
		panic("blockTestAccount err: " + err.Error())
	}
}

func TestMethodGetBlacklist(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		BlackList     []common.Address
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				c.Payload, _ = utils.PackMethod(ABI, MethodGetBlacklist)
			},
			BlackList: []common.Address{},
			Expect:    nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				blockTestAccount(ctx)
				c.Payload, _ = utils.PackMethod(ABI, MethodGetBlacklist)
			},
			BlackList: []common.Address{testAddresses[3]},
			Expect:    nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)

		got := new(MethodStringOutput)
		err = got.Decode(result, MethodGetBlacklist)
		assert.NoError(t, err)
		list := make([]common.Address, 1)
		json.Unmarshal([]byte(got.Result), &list)
		t.Log("blackList result: ", list)
		assert.Equal(t, list, v.BlackList)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodIsBlocked(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		ReturnData    bool
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodIsBlockedInput{Addr: testAddresses[3]}
				c.Payload, _ = input.Encode()
			},
			ReturnData: false,
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				blockTestAccount(ctx)
				input := &MethodIsBlockedInput{Addr: testAddresses[3]}
				c.Payload, _ = input.Encode()
			},
			ReturnData: true,
			Expect:     nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		got := new(MethodBoolOutput)
		got.Decode(result, MethodIsBlocked)
		t.Log("isBlocked result: ", got.Success)
		assert.Equal(t, got.Success, v.ReturnData)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodName(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		ReturnData    string
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				c.Payload, _ = utils.PackMethod(ABI, MethodName)
			},
			ReturnData: contractName,
			Expect:     nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		got := new(MethodContractNameOutput)
		got.Decode(result)
		t.Log("name result: ", got.Name)
		assert.Equal(t, got.Name, v.ReturnData)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func encodeMethodBoolOutput(result bool, methodName string) []byte {
	enc, _ := (&MethodBoolOutput{result}).Encode(methodName)
	return enc
}

func encodeMethodStringOutput(result string, methodName string) []byte {
	enc, _ := (&MethodStringOutput{result}).Encode(methodName)
	return enc
}

func TestMethodSetGasManager(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		ReturnData    []byte
		Expect        error
	}
	cases := []*TestCase{
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodSetGasManagerInput{Addr: testAddresses[3], IsManager: true}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'0'},
			Expect:     errors.New("invalid authority for owner"),
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodSetGasManagerInput{Addr: testAddresses[3], IsManager: true}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'1'},
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodIsGasManagerInput{Addr: testAddresses[3]}
				c.Payload, _ = input.Encode()
			},
			ReturnData: encodeMethodBoolOutput(true, MethodIsGasManager),
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				c.Payload, _ = utils.PackMethod(ABI, MethodGetGasManagerList)
			},
			ReturnData: encodeMethodStringOutput("[\""+strings.ToLower(testAddresses[3].String())+"\"]", MethodGetGasManagerList),
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodSetGasManagerInput{Addr: testAddresses[3], IsManager: false}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'1'},
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodIsGasManagerInput{Addr: testAddresses[3]}
				c.Payload, _ = input.Encode()
			},
			ReturnData: encodeMethodBoolOutput(false, MethodIsGasManager),
			Expect:     nil,
		},
		{
			BlockNum: 3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				c.Payload, _ = utils.PackMethod(ABI, MethodGetGasManagerList)
			},
			ReturnData: encodeMethodStringOutput("[]", MethodGetGasManagerList),
			Expect:     nil,
		},
	}

	resetTestContext()
	for _, v := range cases {
		ctx := generateNativeContract(testCaller, v.BlockNum)

		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		assert.Equal(t, v.ReturnData, result)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodEnableGasManage(t *testing.T) {
	type TestCase struct {
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Expect        error
		ReturnData    []byte
	}

	cases := []*TestCase{
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				payload, err := utils.PackMethod(ABI, MethodIsGasManageEnabled)
				assert.NoError(t, err)
				c.Payload = payload
			},
			ReturnData: encodeMethodBoolOutput(false, MethodIsGasManageEnabled),
			Expect:     nil,
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodEnableGasManageInput{DoEnable: true}
				c.Payload, _ = input.Encode()
			},
			ReturnData: []byte{'0'},
			Expect:     errors.New("invalid authority for owner"),
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodEnableGasManageInput{DoEnable: true}
				c.Payload, _ = input.Encode()
			},
			Expect:     nil,
			ReturnData: []byte{'1'},
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				payload, err := utils.PackMethod(ABI, MethodIsGasManageEnabled)
				assert.NoError(t, err)
				c.Payload = payload
			},
			ReturnData: encodeMethodBoolOutput(true, MethodIsGasManageEnabled),
			Expect:     nil,
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				setDefaultOwner(ctx)
				input := &MethodEnableGasManageInput{DoEnable: false}
				c.Payload, _ = input.Encode()
			},
			Expect:     nil,
			ReturnData: []byte{'1'},
		},
		{
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				payload, err := utils.PackMethod(ABI, MethodIsGasManageEnabled)
				assert.NoError(t, err)
				c.Payload = payload
			},
			ReturnData: encodeMethodBoolOutput(false, MethodIsGasManageEnabled),
			Expect:     nil,
		},
	}

	resetTestContext()
	ctx := generateNativeContract(testCaller, 3)

	for _, v := range cases {
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		assert.Equal(t, v.ReturnData, result)
	}
}
