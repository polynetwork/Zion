package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	polycomm "github.com/polynetwork/poly/common"
)

const (
	KEY_PREFIX_BTC = "btc"

	KEY_PREFIX_BTC_VOTE = "btcVote"
	REQUEST             = "request"
	DONE_TX             = "doneTx"

	NOTIFY_MAKE_PROOF_EVENT = "makeProof"
)

type ChainHandler interface {
	MakeDepositProposal(service *native.NativeContract) (*MakeTxParam, error)
}

type InitRedeemScriptParam struct {
	RedeemScript string
}

func (this *InitRedeemScriptParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteString(this.RedeemScript)
}

func (this *InitRedeemScriptParam) Deserialization(source *polycomm.ZeroCopySource) error {
	redeemScript, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize redeemScript error")
	}

	this.RedeemScript = redeemScript
	return nil
}

type EntranceParam struct {
	SourceChainID         uint64 `json:"sourceChainId"`
	Height                uint32 `json:"height"`
	Proof                 []byte `json:"proof"`
	RelayerAddress        []byte `json:"relayerAddress"`
	Extra                 []byte `json:"extra"`
	HeaderOrCrossChainMsg []byte `json:"headerOrCrossChainMsg"`
}

func (this *EntranceParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint64(this.SourceChainID)
	sink.WriteUint32(this.Height)
	sink.WriteVarBytes(this.Proof)
	sink.WriteVarBytes(this.RelayerAddress)
	sink.WriteVarBytes(this.Extra)
	sink.WriteVarBytes(this.HeaderOrCrossChainMsg)
}

func (this *EntranceParam) Deserialization(source *polycomm.ZeroCopySource) error {
	sourceChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("EntranceParam deserialize sourcechainid error")
	}

	height, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("EntranceParam deserialize height error")
	}
	proof, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize proof error")
	}
	relayerAddr, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize relayerAddr error")
	}
	extra, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize txdata error")
	}
	headerOrCrossChainMsg, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize headerOrCrossChainMsg error")
	}
	this.SourceChainID = sourceChainID
	this.Height = height
	this.Proof = proof
	this.RelayerAddress = relayerAddr
	this.Extra = extra
	this.HeaderOrCrossChainMsg = headerOrCrossChainMsg
	return nil
}

type MakeTxParam struct {
	TxHash              []byte
	CrossChainID        []byte
	FromContractAddress []byte
	ToChainID           uint64
	ToContractAddress   []byte
	Method              string
	Args                []byte
}

func (this *MakeTxParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarBytes(this.TxHash)
	sink.WriteVarBytes(this.CrossChainID)
	sink.WriteVarBytes(this.FromContractAddress)
	sink.WriteUint64(this.ToChainID)
	sink.WriteVarBytes(this.ToContractAddress)
	sink.WriteVarBytes([]byte(this.Method))
	sink.WriteVarBytes(this.Args)
}

func (this *MakeTxParam) Deserialization(source *polycomm.ZeroCopySource) error {
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize txHash error")
	}
	crossChainID, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize crossChainID error")
	}
	fromContractAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize fromContractAddress error")
	}
	toChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize toChainID error")
	}
	toContractAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize toContractAddress error")
	}
	method, eof := source.NextString()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize method error")
	}
	args, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize args error")
	}

	this.TxHash = txHash
	this.CrossChainID = crossChainID
	this.FromContractAddress = fromContractAddress
	this.ToChainID = toChainID
	this.ToContractAddress = toContractAddress
	this.Method = method
	this.Args = args
	return nil
}

type ToMerkleValue struct {
	TxHash      []byte
	FromChainID uint64
	MakeTxParam *MakeTxParam
}

func (this *ToMerkleValue) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarBytes(this.TxHash)
	sink.WriteUint64(this.FromChainID)
	this.MakeTxParam.Serialization(sink)
}

func (this *ToMerkleValue) Deserialization(source *polycomm.ZeroCopySource) error {
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MerkleValue deserialize txHash error")
	}
	fromChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MerkleValue deserialize fromChainID error")
	}

	makeTxParam := new(MakeTxParam)
	err := makeTxParam.Deserialization(source)
	if err != nil {
		return fmt.Errorf("MerkleValue deserialize makeTxParam error:%s", err)
	}

	this.TxHash = txHash
	this.FromChainID = fromChainID
	this.MakeTxParam = makeTxParam
	return nil
}

type MultiSignParam struct {
	ChainID   uint64
	RedeemKey string
	TxHash    []byte
	Address   string
	Signs     [][]byte
}

func (this *MultiSignParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint64(this.ChainID)
	sink.WriteString(this.RedeemKey)
	sink.WriteVarBytes(this.TxHash)
	sink.WriteVarBytes([]byte(this.Address))
	sink.WriteUint64(uint64(len(this.Signs)))
	for _, v := range this.Signs {
		sink.WriteVarBytes(v)
	}
}

func (this *MultiSignParam) Deserialization(source *polycomm.ZeroCopySource) error {
	chainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize txHash error")
	}
	redeemKey, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize redeemKey error")
	}
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize txHash error")
	}
	address, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize address error")
	}
	n, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize signs length error")
	}
	signs := make([][]byte, 0)
	for i := 0; uint64(i) < n; i++ {
		v, eof := source.NextVarBytes()
		if eof {
			return fmt.Errorf("deserialize Signs error")
		}
		signs = append(signs, v)
	}

	this.ChainID = chainID
	this.RedeemKey = redeemKey
	this.TxHash = txHash
	this.Address = address
	this.Signs = signs
	return nil
}
