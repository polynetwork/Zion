package utils

import "github.com/ethereum/go-ethereum/common"

type BtcNetType int

const (
	TyTestnet3 BtcNetType = iota
	TyRegtest
	TySimnet
	TyMainnet
)

var (
	HeaderSyncContractAddress        = common.HexToAddress("0xb2799bDE6831449d73C1F22CE815f773D0CafCc5")
	CrossChainManagerContractAddress = common.HexToAddress("0x5747C05FF236F8d18BB21Bc02ecc389deF853cae")
	SideChainManagerContractAddress  = common.HexToAddress("0x864Ff06eC5fFc75aB6eaf64263308ef5fa7d6637")
	NodeManagerContractAddress       = common.HexToAddress("0xA4Bf827047a08510722B2d62e668a72FCCFa232C")
	RelayerManagerContractAddress    = common.HexToAddress("0xA22f301D7Cb5b50dcA4a015b12EC0cc5f3971412")
	Neo3StateManagerContractAddress  = common.HexToAddress("0x5E839898821dB2A2F0eC9F8aAE7D7053744DB051")
)
