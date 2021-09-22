package node_manager

import (
	"fmt"
	"io"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	polycomm "github.com/polynetwork/poly/common"
)

type Status uint8

func (this *Status) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint8(uint8(*this))
}

func (this *Status) Deserialization(source *polycomm.ZeroCopySource) error {
	status, eof := source.NextUint8()
	if eof {
		return fmt.Errorf("serialization.ReadUint8, deserialize status error: %v", io.ErrUnexpectedEOF)
	}
	*this = Status(status)
	return nil
}

type PeerPoolMap struct {
	PeerPoolMap map[string]*PeerPoolItem
}

func (this *PeerPoolMap) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.PeerPoolMap)))
	var peerPoolItemList []*PeerPoolItem
	for _, v := range this.PeerPoolMap {
		peerPoolItemList = append(peerPoolItemList, v)
	}
	sort.SliceStable(peerPoolItemList, func(i, j int) bool {
		return peerPoolItemList[i].PeerPubkey > peerPoolItemList[j].PeerPubkey
	})
	for _, v := range peerPoolItemList {
		v.Serialization(sink)
	}
}

func (this *PeerPoolMap) Deserialization(source *polycomm.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize PeerPoolMap length error")
	}
	peerPoolMap := make(map[string]*PeerPoolItem)
	for i := 0; uint64(i) < n; i++ {
		peerPoolItem := new(PeerPoolItem)
		if err := peerPoolItem.Deserialization(source); err != nil {
			return fmt.Errorf("deserialize peerPool error: %v", err)
		}
		peerPoolMap[peerPoolItem.PeerPubkey] = peerPoolItem
	}
	this.PeerPoolMap = peerPoolMap
	return nil
}

type PeerPoolItem struct {
	Index      uint32         //peer index
	PeerPubkey string         //peer pubkey
	Address    common.Address //peer owner
	Status     Status
}

func (this *PeerPoolItem) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint32(this.Index)
	sink.WriteString(this.PeerPubkey)
	sink.WriteVarBytes(this.Address[:])
	this.Status.Serialization(sink)
}

func (this *PeerPoolItem) Deserialization(source *polycomm.ZeroCopySource) error {
	index, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize index error")
	}
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("source.NextString, deserialize peerPubkey error")
	}
	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	status := new(Status)
	err := status.Deserialization(source)
	if err != nil {
		return fmt.Errorf("status.Deserialize. deserialize status error: %v", err)
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
	}

	this.Index = index
	this.PeerPubkey = peerPubkey
	this.Address = addr
	this.Status = *status
	return nil
}

type GovernanceView struct {
	View   uint32
	Height uint32
	TxHash common.Hash
}

func (this *GovernanceView) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint32(this.View)
	sink.WriteUint32(this.Height)
	sink.WriteHash(polycomm.Uint256(this.TxHash))
}

func (this *GovernanceView) Deserialization(source *polycomm.ZeroCopySource) error {
	view, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize view error")
	}
	height, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize height error")
	}
	txHash, eof := source.NextHash()
	if eof {
		return fmt.Errorf("source.NextHash, deserialize txHash error")
	}
	this.View = view
	this.Height = height
	this.TxHash = common.Hash(txHash)
	return nil
}

type ConsensusSigns struct {
	SignsMap map[common.Address]bool
}

func (this *ConsensusSigns) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.SignsMap)))
	var signsList []common.Address
	for k := range this.SignsMap {
		signsList = append(signsList, k)
	}
	sort.SliceStable(signsList, func(i, j int) bool {
		return signsList[i].Hex() > signsList[j].Hex()
	})
	for _, v := range signsList {
		sink.WriteVarBytes(v[:])
		sink.WriteBool(this.SignsMap[v])
	}
}

func (this *ConsensusSigns) Deserialization(source *polycomm.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize length of signsMap error")
	}
	signsMap := make(map[common.Address]bool)
	for i := 0; uint64(i) < n; i++ {
		address, eof := source.NextVarBytes()
		if eof {
			return fmt.Errorf("source.NextVarBytes, deserialize address error")
		}
		v, eof := source.NextBool()
		if eof {
			return fmt.Errorf("source.NextBool, deserialize v error")
		}
		addr, err := common.AddressParseFromBytes(address)
		if err != nil {
			return fmt.Errorf("common.AddressParseFromBytes, deserialize address error")
		}
		signsMap[addr] = v
	}
	this.SignsMap = signsMap
	return nil
}
