package utils

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ontio/ontology-crypto/vrf"
	vbftconfig "github.com/polynetwork/poly/consensus/vbft/config"
)

func ConcatKey(contract common.Address, args ...[]byte) []byte {
	temp := contract[:]
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func GetUint32Bytes(num uint32) []byte {
	var p [4]byte
	binary.LittleEndian.PutUint32(p[:], num)
	return p[:]
}

func GetBytesUint32(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(b[:])
}

func GetBytesUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(b[:])
}

func GetUint64Bytes(num uint64) []byte {
	var p [8]byte
	binary.LittleEndian.PutUint64(p[:], num)
	return p[:]
}

func ValidatePeerPubKeyFormat(pubkey string) error {
	pk, err := vbftconfig.Pubkey(pubkey)
	if err != nil {
		return fmt.Errorf("failed to parse pubkey")
	}
	if !vrf.ValidatePublicKey(pk) {
		return fmt.Errorf("invalid for VRF")
	}
	return nil
}
