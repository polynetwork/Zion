package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
)

const (
	OP_RETURN_SCRIPT_FLAG   = byte(0xcc)
	BTC_TX_PREFIX           = "btctx"
	BTC_FROM_TX_PREFIX      = "btcfromtx"
	UTXOS                   = "utxos"
	STXOS                   = "stxos"
	MULTI_SIGN_INFO         = "multiSignInfo"
	MAX_FEE_COST_PERCENTS   = 1.0
	MAX_SELECTING_TRY_LIMIT = 1000000
	SELECTING_K             = 4.0
)

func getNetParam(service *native.NativeContract, chainId uint64) (*chaincfg.Params, error) {
	side, err := side_chain_manager.GetSideChain(service, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get bitcoin net parameter: %v", err)
	}
	if side == nil {
		return nil, fmt.Errorf("side chain info for chainId: %d is not registered", chainId)
	}
	if side.CCMCAddress == nil || len(side.CCMCAddress) != 8 {
		return nil, fmt.Errorf("CCMCAddress is nil or its length is not 8")
	}
	switch utils.BtcNetType(binary.LittleEndian.Uint64(side.CCMCAddress)) {
	case utils.TyTestnet3:
		return &chaincfg.TestNet3Params, nil
	case utils.TyRegtest:
		return &chaincfg.RegressionNetParams, nil
	case utils.TySimnet:
		return &chaincfg.SimNetParams, nil
	default:
		return &chaincfg.MainNetParams, nil
	}
}

func chooseUtxos(native *native.NativeContract, chainID uint64, amount int64, outs []*wire.TxOut, rk []byte, m, n int) ([]*Utxo, int64, int64, error) {
	utxoKey := hex.EncodeToString(rk)
	utxos, err := getUtxos(native, chainID, utxoKey)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("chooseUtxos, getUtxos error: %v", err)
	}
	sort.Sort(sort.Reverse(utxos))
	detail, err := side_chain_manager.GetBtcTxParam(native, rk, chainID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("chooseUtxos, failed to get btcTxParam: %v", err)
	}
	if detail == nil {
		return nil, 0, 0, fmt.Errorf("chooseUtxos, no btcTxParam is set for redeem key %s", hex.EncodeToString(rk))
	}
	cs := &CoinSelector{
		sortedUtxos: utxos,
		target:      uint64(amount),
		maxP:        MAX_FEE_COST_PERCENTS,
		tries:       MAX_SELECTING_TRY_LIMIT,
		mc:          detail.MinChange,
		k:           SELECTING_K,
		txOuts:      outs,
		feeRate:     detail.FeeRate,
		m:           m,
		n:           n,
	}
	result, sum, fee := cs.Select()
	if result == nil || len(result) == 0 {
		return nil, 0, 0, fmt.Errorf("chooseUtxos, current utxo is not enough")
	}
	stxos, err := getStxos(native, chainID, utxoKey)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("chooseUtxos, failed to get stxos: %v", err)
	}
	stxos.Utxos = append(stxos.Utxos, result...)
	putStxos(native, chainID, utxoKey, stxos)

	toSort := new(Utxos)
	toSort.Utxos = result
	sort.Sort(sort.Reverse(toSort))
	idx := 0
	for _, v := range toSort.Utxos {
		for utxos.Utxos[idx].Op.String() != v.Op.String() {
			idx++
		}
		utxos.Utxos = append(utxos.Utxos[:idx], utxos.Utxos[idx+1:]...)
	}
	putUtxos(native, chainID, utxoKey, utxos)
	return result, int64(sum), int64(fee), nil
}

func putTxos(k string, native *native.NativeContract, chainID uint64, txoKey string, txos *Utxos) {
	chainIDBytes := utils.GetUint64Bytes(chainID)
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(k), chainIDBytes, []byte(txoKey))
	sink := polycomm.NewZeroCopySink(nil)
	txos.Serialization(sink)
	native.GetCacheDB().Put(key, cstates.GenRawStorageItem(sink.Bytes()))
}

func getTxos(k string, native *native.NativeContract, chainID uint64, txoKey string) (*Utxos, error) {
	chainIDBytes := utils.GetUint64Bytes(chainID)
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(k), chainIDBytes, []byte(txoKey))
	store, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("get%s, get btcTxStore error: %v", k, err)
	}
	txos := &Utxos{
		Utxos: make([]*Utxo, 0),
	}
	if store != nil {
		utxosBytes, err := cstates.GetValueFromRawStorageItem(store)
		if err != nil {
			return nil, fmt.Errorf("get%s, deserialize from raw storage item err:%v", k, err)
		}
		err = txos.Deserialization(polycomm.NewZeroCopySource(utxosBytes))
		if err != nil {
			return nil, fmt.Errorf("get%s, utxos.Deserialization err:%v", k, err)
		}
	}
	return txos, nil
}

func getStxos(native *native.NativeContract, chainID uint64, stxoKey string) (*Utxos, error) {
	stxos, err := getTxos(STXOS, native, chainID, stxoKey)
	return stxos, err
}

func getUtxos(native *native.NativeContract, chainID uint64, utxoKey string) (*Utxos, error) {
	utxos, err := getTxos(UTXOS, native, chainID, utxoKey)
	return utxos, err
}

func putStxos(native *native.NativeContract, chainID uint64, stxoKey string, stxos *Utxos) {
	putTxos(STXOS, native, chainID, stxoKey, stxos)
}

func putUtxos(native *native.NativeContract, chainID uint64, utxoKey string, utxos *Utxos) {
	putTxos(UTXOS, native, chainID, utxoKey, utxos)
}

func getStxoAmts(service *native.NativeContract, chainID uint64, txIns []*wire.TxIn, redeemKey string) ([]uint64, *Utxos, error) {
	stxos, err := getStxos(service, chainID, redeemKey)
	if err != nil {
		return nil, nil, fmt.Errorf("getStxoAmts, failed to get stxos: %v", err)
	}
	amts := make([]uint64, len(txIns))
	for i, in := range txIns {
		toDel := -1
		for j, v := range stxos.Utxos {
			if bytes.Equal(in.PreviousOutPoint.Hash[:], v.Op.Hash) && in.PreviousOutPoint.Index == v.Op.Index {
				amts[i] = v.Value
				toDel = j
				break
			}
		}
		if toDel < 0 {
			return nil, nil, fmt.Errorf("getStxoAmts, %d txIn not found in stxos", i)
		}
		stxos.Utxos = append(stxos.Utxos[:toDel], stxos.Utxos[toDel+1:]...)
	}

	return amts, stxos, nil
}

func verifySigs(sigs [][]byte, addr string, addrs []btcutil.Address, redeem []byte, tx *wire.MsgTx,
	pkScripts [][]byte, amts []uint64) error {
	if len(sigs) != len(tx.TxIn) {
		return fmt.Errorf("not enough sig, only %d sigs but %d required", len(sigs), len(tx.TxIn))
	}
	var signerAddr btcutil.Address = nil
	for _, a := range addrs {
		if a.EncodeAddress() == addr {
			signerAddr = a
		}
	}

	if signerAddr == nil {
		return fmt.Errorf("address %s not found in redeem script", addr)
	}

	for i, sig := range sigs {
		if len(sig) < 1 {
			return fmt.Errorf("length of no.%d sig is less than 1", i)
		}
		tSig := sig[:len(sig)-1]
		pSig, err := btcec.ParseDERSignature(tSig, btcec.S256())
		if err != nil {
			return fmt.Errorf("failed to parse no.%d sig: %v", i, err)
		}
		var hash []byte
		switch c := txscript.GetScriptClass(pkScripts[i]); c {
		case txscript.MultiSigTy, txscript.ScriptHashTy:
			hash, err = txscript.CalcSignatureHash(redeem, txscript.SigHashType(sig[len(sig)-1]), tx, i)
			if err != nil {
				return fmt.Errorf("failed to calculate sig hash: %v", err)
			}
		case txscript.WitnessV0ScriptHashTy:
			sh := txscript.NewTxSigHashes(tx)
			hash, err = txscript.CalcWitnessSigHash(redeem, sh, txscript.SigHashType(sig[len(sig)-1]), tx, i, int64(amts[i]))
			if err != nil {
				return fmt.Errorf("failed to calculate sig hash: %v", err)
			}
		default:
			return fmt.Errorf("script %s not supported", c)
		}
		if !pSig.Verify(hash, signerAddr.(*btcutil.AddressPubKey).PubKey()) {
			return fmt.Errorf("verify no.%d sig and not pass", i+1)
		}
	}

	return nil
}

func putBtcMultiSignInfo(native *native.NativeContract, txid []byte, multiSignInfo *MultiSignInfo) error {
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(MULTI_SIGN_INFO), txid)
	sink := polycomm.NewZeroCopySink(nil)
	multiSignInfo.Serialization(sink)
	native.GetCacheDB().Put(key, cstates.GenRawStorageItem(sink.Bytes()))
	return nil
}

func getBtcMultiSignInfo(native *native.NativeContract, txid []byte) (*MultiSignInfo, error) {
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(MULTI_SIGN_INFO), txid)
	multiSignInfoStore, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("getBtcMultiSignInfo, get multiSignInfoStore error: %v", err)
	}

	multiSignInfo := &MultiSignInfo{
		MultiSignInfo: make(map[string][][]byte),
	}
	if multiSignInfoStore != nil {
		multiSignInfoBytes, err := cstates.GetValueFromRawStorageItem(multiSignInfoStore)
		if err != nil {
			return nil, fmt.Errorf("getBtcMultiSignInfo, deserialize from raw storage item err:%v", err)
		}
		err = multiSignInfo.Deserialization(polycomm.NewZeroCopySource(multiSignInfoBytes))
		if err != nil {
			return nil, fmt.Errorf("getBtcMultiSignInfo, deserialize multiSignInfo err:%v", err)
		}
	}
	return multiSignInfo, nil
}

func addSigToTx(sigMap *MultiSignInfo, addrs []btcutil.Address, redeem []byte, tx *wire.MsgTx, pkScripts [][]byte) error {
	for i := 0; i < len(tx.TxIn); i++ {
		var (
			script []byte
			err    error
		)
		builder := txscript.NewScriptBuilder()
		switch c := txscript.GetScriptClass(pkScripts[i]); c {
		case txscript.MultiSigTy, txscript.ScriptHashTy:
			builder.AddOp(txscript.OP_FALSE)
			for _, addr := range addrs {
				signs, ok := sigMap.MultiSignInfo[addr.EncodeAddress()]
				if !ok {
					continue
				}
				val := signs[i]
				builder.AddData(val)
			}
			if c == txscript.ScriptHashTy {
				builder.AddData(redeem)
			}
			script, err = builder.Script()
			if err != nil {
				return fmt.Errorf("failed to build sigscript for input %d: %v", i, err)
			}
			tx.TxIn[i].SignatureScript = script
		case txscript.WitnessV0ScriptHashTy:
			data := make([][]byte, len(sigMap.MultiSignInfo)+2)
			idx := 1
			for _, addr := range addrs {
				signs, ok := sigMap.MultiSignInfo[addr.EncodeAddress()]
				if !ok {
					continue
				}
				data[idx] = signs[i]
				idx++
			}
			data[idx] = redeem
			tx.TxIn[i].Witness = wire.TxWitness(data)
		default:
			return fmt.Errorf("addSigToTx, type of no.%d utxo is %s which is not supported", i, c)
		}
	}
	return nil
}

// This function needs to input the input and output information of the transaction
// and the lock time. Function build a raw transaction without signature and return it.
// This function uses the partial logic and code of btcd to finally return the
// reference of the transaction object.
func getUnsignedTx(txIns []*wire.TxIn, outs []*wire.TxOut, changeOut *wire.TxOut, locktime *int64) (*wire.MsgTx, error) {
	if locktime != nil && (*locktime < 0 || *locktime > int64(wire.MaxTxInSequenceNum)) {
		return nil, fmt.Errorf("getUnsignedTx, locktime %d out of range", *locktime)
	}

	// Add all transaction inputs to a new transaction after performing
	// some validity checks.
	mtx := wire.NewMsgTx(wire.TxVersion)
	for _, in := range txIns {
		if locktime != nil && *locktime != 0 {
			in.Sequence = wire.MaxTxInSequenceNum - 1
		}
		mtx.AddTxIn(in)
	}
	for _, out := range outs {
		mtx.AddTxOut(out)
	}
	if changeOut.Value > 0 {
		mtx.AddTxOut(changeOut)
	}
	// Set the Locktime, if given.
	if locktime != nil {
		mtx.LockTime = uint32(*locktime)
	}

	return mtx, nil
}

func getTxOuts(amounts map[string]int64, netParam *chaincfg.Params) ([]*wire.TxOut, error) {
	outs := make([]*wire.TxOut, 0)
	for encodedAddr, amount := range amounts {
		// Decode the provided address.
		addr, err := btcutil.DecodeAddress(encodedAddr, netParam)
		if err != nil {
			return nil, fmt.Errorf("getTxOuts, decode addr fail: %v", err)
		}

		if !addr.IsForNet(netParam) {
			return nil, fmt.Errorf("getTxOuts, addr is not for %s", netParam.Name)
		}

		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, fmt.Errorf("getTxOuts, failed to generate pay-to-address script: %v", err)
		}

		txOut := wire.NewTxOut(amount, pkScript)
		outs = append(outs, txOut)
	}

	return outs, nil
}

func getLockScript(redeem []byte, netParam *chaincfg.Params) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write(redeem)
	witAddr, err := btcutil.NewAddressWitnessScriptHash(hasher.Sum(nil), netParam)
	if err != nil {
		return nil, fmt.Errorf("getChangeTxOut, failed to get witness address: %v", err)
	}
	script, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return nil, fmt.Errorf("getChangeTxOut, failed to get p2sh script: %v", err)
	}
	return script, nil
}

func putBtcFromInfo(native *native.NativeContract, txid []byte, btcFromInfo *BtcFromInfo) error {
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(BTC_FROM_TX_PREFIX), txid)
	sink := polycomm.NewZeroCopySink(nil)
	btcFromInfo.Serialization(sink)
	native.GetCacheDB().Put(key, cstates.GenRawStorageItem(sink.Bytes()))
	return nil
}

func getBtcFromInfo(native *native.NativeContract, txid []byte) (*BtcFromInfo, error) {
	key := utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(BTC_FROM_TX_PREFIX), txid)
	btcFromInfoStore, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("getBtcFromInfo, get multiSignInfoStore error: %v", err)
	}
	btcFromInfo := new(BtcFromInfo)
	if btcFromInfoStore == nil {
		return nil, fmt.Errorf("getBtcFromInfo, can not find any record")
	}
	multiSignInfoBytes, err := cstates.GetValueFromRawStorageItem(btcFromInfoStore)
	if err != nil {
		return nil, fmt.Errorf("getBtcFromInfo, deserialize from raw storage item err:%v", err)
	}
	err = btcFromInfo.Deserialization(polycomm.NewZeroCopySource(multiSignInfoBytes))
	if err != nil {
		return nil, fmt.Errorf("getBtcFromInfo, deserialize multiSignInfo err:%v", err)
	}
	return btcFromInfo, nil
}
