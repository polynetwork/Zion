package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MethodContractName        = "name"
	MethodImportOuterTransfer = "importOuterTransfer"
	MethodMultiSign           = "MultiSign"
	MethodBlackChain          = "BlackChain"
	MethodWhiteChain          = "WhiteChain"
)

var ABI *abi.ABI

func init() {
	ABI = GetABI()
}

const abijson = `[
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"merkleValueHex","type":"string"},{"indexed":false,"internalType":"uint64","name":"BlockHeight","type":"uint64"},{"indexed":false,"internalType":"string","name":"key","type":"string"}],"name":"` + NOTIFY_MAKE_PROOF_EVENT + `","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"TxHash","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"sink","type":"bytes"}],"name":"btcTxMultiSignEvent","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"FromChainID","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"ChainID","type":"uint64"},{"indexed":false,"internalType":"string","name":"buf","type":"string"},{"indexed":false,"internalType":"string","name":"FromTxHash","type":"string"},{"indexed":false,"internalType":"string","name":"RedeemKey","type":"string"}],"name":"btcTxToRelayEvent","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"rk","type":"string"},{"indexed":false,"internalType":"string","name":"buf","type":"string"},{"indexed":false,"internalType":"uint64[]","name":"amts","type":"uint64[]"}],"name":"makeBtcTxEvent","type":"event"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"}],"name":"` + MethodBlackChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[],"name":"` + MethodContractName + `","outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"SourceChainID","type":"uint64"},{"internalType":"uint32","name":"Height","type":"uint32"},{"internalType":"bytes","name":"Proof","type":"bytes"},{"internalType":"bytes","name":"RelayerAddress","type":"bytes"},{"internalType":"bytes","name":"Extra","type":"bytes"},{"internalType":"bytes","name":"HeaderOrCrossChainMsg","type":"bytes"}],"name":"` + MethodImportOuterTransfer + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"},{"internalType":"string","name":"RedeemKey","type":"string"},{"internalType":"bytes","name":"TxHash","type":"bytes"},{"internalType":"string","name":"Address","type":"string"},{"internalType":"bytes[]","name":"Signs","type":"bytes[]"}],"name":"` + MethodMultiSign + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"}],"name":"` + MethodWhiteChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type BlackChainParam struct {
	ChainID uint64
}
