package hotstuff

import (
	"bytes"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsg(t *testing.T) {
	type Pig struct {
		Name    string
		Gender  uint8
		Address string
	}

	pig := &Pig{Name: "piggy", Gender: 2, Address: "England"}
	enc, err := rlp.EncodeToBytes(pig)
	assert.NoError(t, err)

	rec := new(Pig)
	err = rlp.DecodeBytes(enc, rec)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = rlp.Encode(&buf, pig)
	assert.NoError(t, err)

	//s := rlp.NewListStream(&buf, 0)
	//rec1 := new(Pig)
	//rlp.Decode()
	//err = s.Decode(rec1)
	//assert.NoError(t, err)

	//msg := MsgNewView{
	//	PrepareQC: QuorumCert{
	//		BlockHash: common.Hash{},
	//		ViewNum:   1,
	//		Type:      MsgTypeNewView,
	//		Signature: []byte("testsig"),
	//	},
	//	ViewNum:   1,
	//}
	//
	////var buf bytes.Buffer
	//buf, err := rlp.EncodeToBytes(msg)
	////err := msg.EncodeRLP(bufio.NewWriter(&buf))
	//assert.NoError(t, err)
	//length := uint64(len(buf))
	//t.Log(length)
	////rec := new(MsgNewView)
	////err = rlp.DecodeBytes(buf, rec)
	//////err = rec.DecodeRLP(rlp.NewListStream(&buf, length))
	////assert.NoError(t, err)
	////
	////t.Log(msg, rec)
}
