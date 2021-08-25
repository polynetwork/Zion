package state

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

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
}