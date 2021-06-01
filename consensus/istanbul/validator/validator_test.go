package validator

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"testing"
)

func TestSubtraction(t *testing.T) {
	s1 := NewSet([]common.Address{
		common.BytesToAddress([]byte{3}),
		common.BytesToAddress([]byte{1}),
		common.BytesToAddress([]byte{4}),
		common.BytesToAddress([]byte{2}),
	}, istanbul.RoundRobin)

	s2 := NewSet([]common.Address{
		common.BytesToAddress([]byte{2}),
		common.BytesToAddress([]byte{5}),
		common.BytesToAddress([]byte{4}),
	}, istanbul.RoundRobin)

	addList, delList := Subtraction(s1, s2)
	for _, v := range addList{
		t.Logf("add %v", v.String())
	}
	for _, v := range delList{
		t.Logf("del %v", v.String())
	}
}
