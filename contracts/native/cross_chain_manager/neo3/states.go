package neo3

import (
	"fmt"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/io"
	"github.com/joeqian10/neo3-gogogo/mpt"
)

type CrossChainMsg struct {
	*mpt.StateRoot
}

func (this *CrossChainMsg) GetScriptHash() (*helper.UInt160, error) {
	if len(this.Witnesses) == 0 {
		return nil, fmt.Errorf("NeoCrossChainMsg.Witness incorrect length")
	}
	verificationScriptBs, err := crypto.Base64Decode(this.Witnesses[0].Verification) // base64
	if err != nil {
		return nil, fmt.Errorf("NeoCrossChainMsg.Witness.Verification decode error: %s", err)
	}
	if len(verificationScriptBs) == 0 {
		return nil, fmt.Errorf("NeoCrossChainMsg.Witness.VerificationScript is empty")
	}
	scriptHash := helper.UInt160FromBytes(crypto.Hash160(verificationScriptBs))
	return scriptHash, nil
}

func (this *CrossChainMsg) GetMessage(magic uint32) ([]byte, error) {
	buff2 := io.NewBufBinaryWriter()
	this.SerializeUnsigned(buff2.BinaryWriter)
	if buff2.Err != nil {
		return nil, fmt.Errorf("neo3-gogogo mpt.StateRoot SerializeUnsigned error: %s", buff2.Err)
	}
	hash := helper.UInt256FromBytes(crypto.Sha256(buff2.Bytes()))

	buf := io.NewBufBinaryWriter()
	buf.BinaryWriter.WriteLE(magic)
	buf.BinaryWriter.WriteLE(hash)
	if buf.Err != nil {
		return nil, fmt.Errorf("NeoCrossChainMsg.GetMessage write hash error: %s", buf.Err)
	}
	return buf.Bytes(), nil
}

func DeserializeCrossChainMsg(source []byte) (*CrossChainMsg, error) {
	ccm := &CrossChainMsg{}
	br := io.NewBinaryReaderFromBuf(source)
	ccm.Deserialize(br)
	if br.Err != nil {
		return nil, fmt.Errorf("neo3 StateRoot.Deserialzie error: %v", br.Err)
	}
	return ccm, nil
}

func SerializeCrossChainMsg(ccm *CrossChainMsg) ([]byte, error) {
	bw := io.NewBufBinaryWriter()
	ccm.Serialize(bw.BinaryWriter)
	if bw.Err != nil {
		return nil, fmt.Errorf("neo3 StateRoot.Serialize error: %v", bw.Err)
	}
	return bw.Bytes(), nil
}

type NeoMakeTxParam struct {
	*scom.MakeTxParam
}

func (this *NeoMakeTxParam) Deserialize(br *io.BinaryReader) {
	this.TxHash = br.ReadVarBytes()
	this.CrossChainID = br.ReadVarBytes()
	this.FromContractAddress = br.ReadVarBytes()
	br.ReadLE(&this.ToChainID)
	this.ToContractAddress = br.ReadVarBytes()
	this.Method = string(br.ReadVarBytes())
	this.Args = br.ReadVarBytes()
}

func (this *NeoMakeTxParam) Serialize(bw *io.BinaryWriter) {
	bw.WriteVarBytes(this.TxHash)
	bw.WriteVarBytes(this.CrossChainID)
	bw.WriteVarBytes(this.FromContractAddress)
	bw.WriteLE(this.ToChainID)
	bw.WriteVarBytes(this.ToContractAddress)
	bw.WriteVarBytes([]byte(this.Method))
	bw.WriteVarBytes(this.Args)
}

func DeserializeNeoMakeTxParam(source []byte) (*NeoMakeTxParam, error) {
	param := new(NeoMakeTxParam)
	param.MakeTxParam = new(scom.MakeTxParam)
	br := io.NewBinaryReaderFromBuf(source)
	param.Deserialize(br)
	if br.Err != nil {
		return nil, br.Err
	}
	return param, nil
}

func SerializeNeoMakeTxParam(param *NeoMakeTxParam) ([]byte, error) {
	bbw := io.NewBufBinaryWriter()
	param.Serialize(bbw.BinaryWriter)
	if bbw.Err != nil {
		return nil, bbw.Err
	}
	return bbw.Bytes(), nil
}
