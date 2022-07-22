package signature_manager

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/signature_manager_abi"
	"github.com/ethereum/go-ethereum/rlp"
)

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(ISignatureManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type AddSignatureParam struct {
	Addr        common.Address
	SideChainID *big.Int
	Subject     []byte
	Signature   []byte
}

func (p *AddSignatureParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.Addr, p.SideChainID, p.Subject, p.Signature})
}

func (p *AddSignatureParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Address     common.Address
		SideChainID *big.Int
		Subject     []byte
		Signature   []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	p.Addr, p.SideChainID, p.Subject, p.Signature = data.Address, data.SideChainID, data.Subject, data.Signature
	return nil
}
