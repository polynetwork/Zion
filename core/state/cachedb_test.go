package state

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/stretchr/testify/assert"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestNextSlot
func TestNextSlot(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	raw := []byte{byte(102), byte(102), byte(85)}
	hash := common.BytesToHash(raw)
	t.Log("hash is", hash.Hex())

	slot := hash
	for i := 0; i < 1; i++ {
		t.Log(slot.Hex())
		slot = c.nextSlot(slot)
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestCustomSet
func TestCustomSet(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	addr := common.HexToAddress("0x12345")
	key := append(addr.Bytes(), []byte("test")...)
	value := common.HexToHash("0x45678")

	_, _, err := c.customSet(key, value)
	assert.NoError(t, err)

	_, _, hash, err := c.customGet(key)
	assert.NoError(t, err)
	assert.Equal(t, value, hash)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestSetAddress
func TestSetAddress(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	addr := common.HexToAddress("0x12345")

	// normal address
	{
		key := append(addr.Bytes(), []byte("test1")...)
		value := common.HexToAddress("0xabc1234567")

		assert.NoError(t, c.SetAddress(key, value))
		got, err := c.GetAddress(key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)

		// delete address
		assert.NoError(t, c.DelAddress(key))
		got, err = c.GetAddress(key)
		assert.NoError(t, err)
		assert.Equal(t, common.EmptyAddress, got)
	}

	// empty address
	{
		key := append(addr.Bytes(), []byte("test2")...)
		value := common.EmptyAddress

		assert.NoError(t, c.SetAddress(key, value))
		got, err := c.GetAddress(key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestSetHash
func TestSetHash(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	addr := common.HexToAddress("0x12345")

	// normal hash
	{
		key := append(addr.Bytes(), []byte("test")...)
		value := common.HexToHash("0xabc12345679ab")
		assert.NoError(t, c.SetHash(key, value))
		got, err := c.GetHash(key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)

		// delete hash
		assert.NoError(t, c.DelHash(key))
		got, err = c.GetHash(key)
		assert.NoError(t, err)
		assert.Equal(t, common.EmptyHash, got)
	}

	// empty hash
	{
		key := append(addr.Bytes(), []byte("test")...)
		value := common.EmptyHash
		assert.NoError(t, c.SetHash(key, value))
		got, err := c.GetHash(key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestSetBigInt
func TestSetBigInt(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	addr := common.HexToAddress("0x12345")

	// normal big int
	{
		key := append(addr.Bytes(), []byte("test")...)
		value := big.NewInt(12345)
		assert.NoError(t, c.SetBigInt(key, value))
		got, err := c.GetBigInt(key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)

		// delete big int
		assert.NoError(t, c.DelBigInt(key))
		got, err = c.GetBigInt(key)
		assert.NoError(t, err)
		assert.Equal(t, 0, int(got.Uint64()))
	}

	// zero
	{
		key := append(addr.Bytes(), []byte("test")...)
		value := big.NewInt(0)
		assert.NoError(t, c.SetBigInt(key, value))
		got, err := c.GetBigInt(key)
		assert.NoError(t, err)
		assert.Equal(t, value.Sign(), got.Sign())
		assert.Zero(t, value.CmpAbs(got))
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/core/state -run TestSetBytes
func TestSetBytes(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)
	addr := common.HexToAddress("0x12345")

	// bytes length less than 32
	{
		key := append(addr.Bytes(), []byte("test1")...)
		expect := []byte{'a', 'b', 'c', 'd', 'e', '1'}
		assert.NoError(t, c.SetBytes(key, expect))
		t.Logf("end of setBytes")

		got, err := c.GetBytes(key)
		assert.NoError(t, err)
		assert.Equal(t, expect, got)
		t.Logf("end of getBytes")
	}

	// bytes length equals to 32
	{
		key := append(addr.Bytes(), []byte("test1")...)
		expect := common.HexToHash("0x123456789abcde").Bytes()
		assert.NoError(t, c.SetBytes(key, expect))

		got, err := c.GetBytes(key)
		assert.NoError(t, err)
		assert.Equal(t, expect, got)
	}

	// bytes with empty hash
	{
		key := append(addr.Bytes(), []byte("test1")...)
		expect := common.Hex2Bytes("12000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		assert.NoError(t, c.SetBytes(key, expect))

		got, err := c.GetBytes(key)
		assert.NoError(t, err)
		t.Log("got length", len(got), "expect length", len(expect))
		assert.Equal(t, expect, got)
	}
	{
		key := append(addr.Bytes(), []byte("test1")...)
		expect := common.Hex2Bytes("1230000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000022000")
		assert.NoError(t, c.SetBytes(key, expect))

		got, err := c.GetBytes(key)
		assert.NoError(t, err)
		t.Log("got length", len(got), "expect length", len(expect))
		assert.Equal(t, expect, got)
	}
}

func TestCacheDB(t *testing.T) {
	// Create an empty state database
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db), nil)
	c := (*CacheDB)(state)

	// Update it with some accounts
	for i := byte(0); i < 255; i++ {
		addr := common.BytesToAddress([]byte{i})

		value := []byte("0123456789ABCDEF0123456789ABCDEF1")
		key := append(addr[:], []byte("a")...)
		c.Put(key, value)
		valueBack, err := c.Get(key)
		if err != nil || !bytes.Equal(value, valueBack) {
			t.Fail()
		}

		c.Delete(key)
		v, err := c.Get(key)
		if v != nil || err != nil {
			t.Fail()
		}
	}

	testByteSize := 160
	testBytes := make([]byte, testByteSize)
	n, err := rand.Read(testBytes)
	if err != nil || n != testByteSize {
		t.Fail()
	}

	addr := common.BytesToAddress([]byte{0})
	key := append(addr[:], []byte("a")...)
	c.Put(key, testBytes)
	respBytes, _ := c.Get(key)
	if !bytes.Equal(testBytes, respBytes) {
		t.Fail()
	}

	{
		key, _ = hex.DecodeString("864ff06ec5ffc75ab6eaf64263308ef5fa7d663773696465436861696e0200000000000000")
		value, _ := hex.DecodeString("001b140000000000000000000000000000000000000000020000000000")
		c.Put(key, value)
		respBytes, _ = c.Get(key)
		if !bytes.Equal(value, respBytes) {
			t.Fail()
		}
	}

}
