package state

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type CacheDB StateDB

func (c *CacheDB) Put(key []byte, value []byte) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	c.Delete(key)

	s := (*StateDB)(c)
	so := s.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := c.key2Slot(key[common.AddressLength:])
		if len(value) <= common.HashLength-1 {
			c.putValue(so, slot, value, false)
			value = nil
		} else {
			c.putValue(so, slot, value[:common.HashLength-1], true)
			value = value[common.HashLength-1:]
		}

		for len(value) > 0 {
			slot := c.nextSlot(slot)
			if len(value) <= common.HashLength-1 {
				c.putValue(so, slot, value, false)
				break
			} else {
				c.putValue(so, slot, value[:common.HashLength-1], true)
				value = value[common.HashLength-1:]
			}
		}
	}

}

func (c *CacheDB) putValue(so *stateObject, slot common.Hash, value []byte, more bool) {
	if len(value) > common.HashLength-1 {
		panic("value should not exceed 31")
	}

	if more && len(value) != common.HashLength-1 {
		panic("value length should equal 31 when more is true")
	}

	if more {
		value = append([]byte{1}, value...)
	} else {
		padding := make([]byte, common.HashLength-len(value))
		padding[0] = byte(len(value) << 1)
		value = append(padding, value...)
	}

	s := (*StateDB)(c)
	hashValue := common.BytesToHash(value)
	so.SetState(s.db, slot, hashValue)
}

func (c *CacheDB) key2Slot(key []byte) common.Hash {
	key = crypto.Keccak256(key)
	return common.BytesToHash(key)
}

func (c *CacheDB) nextSlot(slot common.Hash) common.Hash {
	slotBytes := slot.Bytes()
	for offset := common.HashLength - 1; offset >= 0; offset-- {
		slotBytes[offset] = slotBytes[offset] + 1
		if slotBytes[offset] != 0 {
			break
		}
	}

	return c.key2Slot(slotBytes)
}

var (
	ErrValueNotExists = errors.New("value not exists")
)

func (c *CacheDB) Get(key []byte) ([]byte, error) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	stateObject := s.getStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if stateObject != nil {
		var result []byte
		slot := c.key2Slot(key[common.AddressLength:])
		value := stateObject.GetState(s.db, slot)
		meta := value[:][0]
		more := meta&1 == 1
		if more {
			result = append(result, value[1:]...)
		} else {
			if value == (common.Hash{}) {
				return nil, ErrValueNotExists
			}
			result = append(result, value[1:1+meta>>1]...)
		}

		for more {
			slot = c.nextSlot(slot)
			value = stateObject.GetState(s.db, slot)
			meta = value[:][0]
			more = meta&1 == 1
			if more {
				result = append(result, value[1:]...)
			} else {
				result = append(result, value[common.HashLength-meta>>1:]...)
			}
		}

		return result, nil
	}

	return nil, ErrValueNotExists
}

func (c *CacheDB) Delete(key []byte) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	s := (*StateDB)(c)
	so := s.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := c.key2Slot(key[common.AddressLength:])
		value := so.GetState(s.db, slot)
		so.SetState(s.db, slot, common.Hash{})
		more := value[:][0]&1 == 1
		for more {
			slot = c.nextSlot(slot)
			value = so.GetState(s.db, slot)
			so.SetState(s.db, slot, common.Hash{})
			more = value[:][0]&1 == 1
		}
	}
}
