/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */
package side_chain_manager

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

var netParam = &chaincfg.TestNet3Params

func GetSideChainApply(native *native.NativeContract, chanid uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chanid)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getRegisterSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getRegisterSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putSideChainApply(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainId)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putRegisterSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY), chainidByte), blob)
	return nil
}

func GetSideChain(native *native.NativeContract, chainID uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainIDByte := utils.GetUint64Bytes(chainID)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(SIDE_CHAIN),
		chainIDByte))
	if err != nil {
		return nil, fmt.Errorf("getSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	}
	return nil, nil

}

func PutSideChain(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainId)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(SIDE_CHAIN), chainidByte), blob)
	return nil
}

func getUpdateSideChain(native *native.NativeContract, chanid uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chanid)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getUpdateSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getUpdateSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putUpdateSideChain(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainId)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putUpdateSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte), blob)
	return nil
}

func getQuitSideChain(native *native.NativeContract, chainid uint64) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chainid)

	chainIDStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return fmt.Errorf("getQuitSideChain, get registerSideChainRequestStore error: %v", err)
	}
	if chainIDStore != nil {
		return nil
	}
	return fmt.Errorf("getQuitSideChain, no record")
}

func putQuitSideChain(native *native.NativeContract, chainid uint64) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chainid)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST), chainidByte),chainidByte)
	return nil
}

func GetContractBind(native *native.NativeContract, redeemChainID, contractChainID uint64,
	redeemKey []byte) (*ContractBinded, error) {
	contract := utils.SideChainManagerContractAddress
	redeemChainIDByte := utils.GetUint64Bytes(redeemChainID)
	contractChainIDByte := utils.GetUint64Bytes(contractChainID)
	contractBindStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(REDEEM_BIND),
		redeemChainIDByte, contractChainIDByte, redeemKey))
	if err != nil {
		return nil, fmt.Errorf("GetContractBind, get contractBindStore error: %v", err)
	}
	if contractBindStore != nil {
		cb := new(ContractBinded)
		err = rlp.DecodeBytes(contractBindStore, cb)
		if err != nil {
			return nil, fmt.Errorf("GetContractBind, deserialize BindContract err:%v", err)
		}
		return cb, nil
	} else {
		return nil, nil
	}

}

func putContractBind(native *native.NativeContract, redeemChainID, contractChainID uint64,
	redeemKey, contractAddress []byte, cver uint64) error {
	contract := utils.SideChainManagerContractAddress
	redeemChainIDByte := utils.GetUint64Bytes(redeemChainID)
	contractChainIDByte := utils.GetUint64Bytes(contractChainID)
	bc := &ContractBinded{
		Contract: contractAddress,
		Ver:      cver,
	}
	blob, err := rlp.EncodeToBytes(bc)
	if err != nil {
		return fmt.Errorf("putContractBind, contractBind.Serialization error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(REDEEM_BIND),
		redeemChainIDByte, contractChainIDByte, redeemKey), blob)
	return nil
}

func putBindSignInfo(native *native.NativeContract, message []byte, multiSignInfo *BindSignInfo) error {
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(BIND_SIGN_INFO), message)
	blob, err := rlp.EncodeToBytes(multiSignInfo)
	if err != nil {
		return fmt.Errorf("putBindSignInfo, BindSignInfo.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(key, blob)
	return nil
}

func getBindSignInfo(native *native.NativeContract, message []byte) (*BindSignInfo, error) {
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(BIND_SIGN_INFO), message)
	bindSignInfoStore, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("getBtcMultiSignInfo, get multiSignInfoStore error: %v", err)
	}

	bindSignInfo := &BindSignInfo{
		BindSignInfo: make(map[string][]byte),
	}
	if bindSignInfoStore != nil {
		err = rlp.DecodeBytes(bindSignInfoStore, bindSignInfo)
		if err != nil {
			return nil, fmt.Errorf("getBtcMultiSignInfo, deserialize multiSignInfo err:%v", err)
		}
	}
	return bindSignInfo, nil
}

func putBtcTxParam(native *native.NativeContract, redeemKey []byte, redeemChainId uint64, detail *BtcTxParamDetial) error {
	redeemChainIdBytes := utils.GetUint64Bytes(redeemChainId)
	blob, err := rlp.EncodeToBytes(detail)
	if err != nil {
		return err
	}

	native.GetCacheDB().Put(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(BTC_TX_PARAM), redeemKey,
		redeemChainIdBytes), blob)
	return nil
}

func GetBtcTxParam(native *native.NativeContract, redeemKey []byte, redeemChainId uint64) (*BtcTxParamDetial, error) {
	redeemChainIdBytes := utils.GetUint64Bytes(redeemChainId)
	store, err := native.GetCacheDB().Get(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(BTC_TX_PARAM), redeemKey,
		redeemChainIdBytes))
	if err != nil {
		return nil, fmt.Errorf("GetBtcTxParam, get btcTxParam error: %v", err)
	}
	if store != nil {
		detial := new(BtcTxParamDetial)
		if err := rlp.DecodeBytes(store, detial); err != nil {
			return nil, fmt.Errorf("GetBtcTxParam, DecodeBytes error: %v", err)
		}
		return detial, nil
	}
	return nil, nil
}

func verifyRedeemRegister(param *RegisterRedeemParam, addrs []btcutil.Address) (map[string][]byte, error) {
	r := make([]byte, len(param.Redeem))
	copy(r, param.Redeem)
	cverBytes := utils.GetUint64Bytes(param.CVersion)
	fromChainId := utils.GetUint64Bytes(param.RedeemChainID)
	toChainId := utils.GetUint64Bytes(param.ContractChainID)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), param.ContractAddress...),
		toChainId...), cverBytes...))
	return verify(param.Signs, addrs, hash)
}

func verifyBtcTxParam(param *BtcTxParam, addrs []btcutil.Address) (map[string][]byte, error) {
	r := make([]byte, len(param.Redeem))
	copy(r, param.Redeem)
	fromChainId := utils.GetUint64Bytes(param.RedeemChainId)
	frBytes := utils.GetUint64Bytes(param.Detial.FeeRate)
	mcBytes := utils.GetUint64Bytes(param.Detial.MinChange)
	verBytes := utils.GetUint64Bytes(param.Detial.PVersion)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), frBytes...), mcBytes...), verBytes...))
	return verify(param.Sigs, addrs, hash)
}

func verify(sigs [][]byte, addrs []btcutil.Address, hash []byte) (map[string][]byte, error) {
	res := make(map[string][]byte)
	for i, sig := range sigs {
		if len(sig) < 1 {
			return nil, fmt.Errorf("length of no.%d sig is less than 1", i)
		}
		pSig, err := btcec.ParseDERSignature(sig, btcec.S256())
		if err != nil {
			return nil, fmt.Errorf("failed to parse no.%d sig: %v", i, err)
		}
		for _, addr := range addrs {
			if pSig.Verify(hash, addr.(*btcutil.AddressPubKey).PubKey()) {
				res[addr.EncodeAddress()] = sig
			}
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no sigs is verified")
	}
	return res, nil
}

func putBtcRedeemScript(native *native.NativeContract, redeemScriptKey string, redeemScriptBytes []byte, redeemChainId uint64) error {
	chainIDBytes := utils.GetUint64Bytes(redeemChainId)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(REDEEM_SCRIPT), chainIDBytes, []byte(redeemScriptKey))

	cls := txscript.GetScriptClass(redeemScriptBytes)
	if cls.String() != "multisig" {
		return fmt.Errorf("putBtcRedeemScript, wrong type of redeem: %s", cls)
	}
	native.GetCacheDB().Put(key, redeemScriptBytes)
	return nil
}

func GetBtcRedeemScriptBytes(native *native.NativeContract, redeemScriptKey string, redeemChainId uint64) ([]byte, error) {
	chainIDBytes := utils.GetUint64Bytes(redeemChainId)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(REDEEM_SCRIPT), chainIDBytes, []byte(redeemScriptKey))
	redeemStore, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("getBtcRedeemScript, get btcProofStore error: %v", err)
	}
	if redeemStore == nil {
		return nil, fmt.Errorf("getBtcRedeemScript, can not find any records")
	}
	return redeemStore, nil
}
