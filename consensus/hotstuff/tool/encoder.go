package tool

import (
	"bytes"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/event_driven/core"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Encode generate hotstuff genesis extra
func Encode(validators []common.Address) (string, error) {
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)

	ist := &types.HotstuffExtra{
		Validators:    validators,
		Seal:          make([]byte, types.HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...)), nil
}

func EncodeSalt(validators []common.Address, epoch, round uint64) (string, error) {
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)

	salt := &core.ExtraSalt{
		Epoch: epoch,
		Round: new(big.Int).SetUint64(round),
	}
	enc, err := core.Encode(salt)
	if err != nil {
		return "", err
	}

	ist := &types.HotstuffExtra{
		Validators:    validators,
		Seal:          make([]byte, types.HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
		Salt:          enc,
	}
	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...)), nil
}
