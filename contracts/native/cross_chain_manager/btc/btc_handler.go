package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	ABI *abi.ABI
)

type BTCHandler struct {
}

func NewBTCHandler() *BTCHandler {
	return &BTCHandler{}
}

func (this *BTCHandler) MultiSign(service *native.NativeContract, params *scom.MultiSignParam) error {

	multiSignInfo, err := getBtcMultiSignInfo(service, params.TxHash)
	if err != nil {
		return fmt.Errorf("MultiSign, getBtcMultiSignInfo error: %v", err)
	}

	_, ok := multiSignInfo.MultiSignInfo[params.Address]
	if ok {
		return fmt.Errorf("MultiSign, address %s already sign", params.Address)
	}

	redeemScript, err := side_chain_manager.GetBtcRedeemScriptBytes(service, params.RedeemKey, params.ChainID)
	if err != nil {
		return fmt.Errorf("MultiSign, get btc redeem script with redeem key %v from db error: %v", params.RedeemKey, err)
	}

	netParam, err := getNetParam(service, params.ChainID)
	if err != nil {
		return fmt.Errorf("MultiSign, %v", err)
	}

	_, addrs, n, err := txscript.ExtractPkScriptAddrs(redeemScript, netParam)
	if err != nil {
		return fmt.Errorf("MultiSign, failed to extract pkscript addrs: %v", err)
	}
	if len(multiSignInfo.MultiSignInfo) == n {
		return fmt.Errorf("MultiSign, already enough signature: %d", n)
	}

	txb, err := service.GetCacheDB().Get(utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(BTC_TX_PREFIX),
		params.TxHash))
	if err != nil {
		return fmt.Errorf("MultiSign, failed to get tx %s from cacheDB: %v", hex.EncodeToString(params.TxHash), err)
	}

	mtx := wire.NewMsgTx(wire.TxVersion)
	err = mtx.BtcDecode(bytes.NewBuffer(txb), wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		return fmt.Errorf("MultiSign, failed to decode tx: %v", err)
	}

	pkScripts := make([][]byte, len(mtx.TxIn))
	for i, in := range mtx.TxIn {
		pkScripts[i] = in.SignatureScript
		in.SignatureScript = nil
	}
	amts, stxos, err := getStxoAmts(service, params.ChainID, mtx.TxIn, params.RedeemKey)
	if err != nil {
		return fmt.Errorf("MultiSign, failed to get stxos: %v", err)
	}
	err = verifySigs(params.Signs, params.Address, addrs, redeemScript, mtx, pkScripts, amts)
	if err != nil {
		return fmt.Errorf("MultiSign, failed to verify: %v", err)
	}
	multiSignInfo.MultiSignInfo[params.Address] = params.Signs
	err = putBtcMultiSignInfo(service, params.TxHash, multiSignInfo)
	if err != nil {
		return fmt.Errorf("MultiSign, putBtcMultiSignInfo error: %v", err)
	}

	if len(multiSignInfo.MultiSignInfo) != n {
		service.AddNotify(ABI, []string{"btcTxMultiSign"}, params.TxHash, multiSignInfo.MultiSignInfo)
	} else {
		err = addSigToTx(multiSignInfo, addrs, redeemScript, mtx, pkScripts)
		if err != nil {
			return fmt.Errorf("MultiSign, failed to add sig to tx: %v", err)
		}
		var buf bytes.Buffer
		err = mtx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			return fmt.Errorf("MultiSign, failed to encode msgtx to bytes: %v", err)
		}

		witScript, err := getLockScript(redeemScript, netParam)
		if err != nil {
			return fmt.Errorf("MultiSign, failed to get lock script: %v", err)
		}
		utxos, err := getUtxos(service, params.ChainID, params.RedeemKey)
		if err != nil {
			return fmt.Errorf("MultiSign, getUtxos error: %v", err)
		}
		txid := mtx.TxHash()
		for i, v := range mtx.TxOut {
			if bytes.Equal(witScript, v.PkScript) {
				newUtxo := &Utxo{
					Op: &OutPoint{
						Hash:  txid[:],
						Index: uint32(i),
					},
					Value:        uint64(v.Value),
					ScriptPubkey: v.PkScript,
				}
				utxos.Utxos = append(utxos.Utxos, newUtxo)
			}
		}
		putUtxos(service, params.ChainID, params.RedeemKey, utxos)
		btcFromTxInfo, err := getBtcFromInfo(service, params.TxHash)
		if err != nil {
			return fmt.Errorf("MultiSign, failed to get from tx hash %s from cacheDB: %v",
				hex.EncodeToString(params.TxHash), err)
		}
		putStxos(service, params.ChainID, params.RedeemKey, stxos)
		service.AddNotify(ABI, []string{"btcTxToRelay"}, btcFromTxInfo.FromChainID, params.ChainID,
			hex.EncodeToString(buf.Bytes()), hex.EncodeToString(btcFromTxInfo.FromTxHash), params.RedeemKey)
	}
	return nil
}
