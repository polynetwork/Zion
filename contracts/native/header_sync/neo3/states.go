package neo3

import (
	"fmt"
	"github.com/joeqian10/neo3-gogogo/block"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/io"
)

type BlockHeader struct {
	*block.Header
}

func (this *BlockHeader) GetMessage(magic uint32) ([]byte, error) {
	buff2 := io.NewBufBinaryWriter()
	this.SerializeUnsigned(buff2.BinaryWriter)
	if buff2.Err != nil {
		return nil, fmt.Errorf("neo3 Header SerializeUnsigned error: %v", buff2.Err)
	}
	hash := helper.UInt256FromBytes(crypto.Sha256(buff2.Bytes()))

	buf := io.NewBufBinaryWriter()
	buf.BinaryWriter.WriteLE(magic)
	buf.BinaryWriter.WriteLE(hash)
	if buf.Err != nil {
		return nil, fmt.Errorf("neo3 BlockHeader.GetMessage write hash error: %v", buf.Err)
	}
	return buf.Bytes(), nil
}

func DeserializeNeo3Header(source []byte) (*BlockHeader, error) {
	h := &BlockHeader{}
	br := io.NewBinaryReaderFromBuf(source)
	h.Deserialize(br)
	if br.Err != nil {
		return nil, fmt.Errorf("neo3 Header.Deserialize error: %v", br.Err)
	}
	return h, nil
}

func SerializeNeo3Header(h *BlockHeader) ([]byte, error) {
	bw := io.NewBufBinaryWriter()
	h.Serialize(bw.BinaryWriter)
	if bw.Err != nil {
		return nil, fmt.Errorf("neo3 Header.Serialize error: %v", bw.Err)
	}
	return bw.Bytes(), nil
}
