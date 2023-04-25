/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package ripple

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/cross_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/polynetwork/ripple-sdk/types"
	"github.com/rubblelabs/ripple/data"
)

type RippleHandler struct {
}

func NewRippleHandler() *RippleHandler {
	return &RippleHandler{}
}

func (this *RippleHandler) MakeDepositProposal(service *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//verify signature
	digest, err := params.Digest()
	if err != nil {
		return nil, fmt.Errorf("ripple MakeDepositProposal, digest input param error: %v", err)
	}
	pub, err := crypto.SigToPub(digest, params.Signature)
	if err != nil {
		return nil, fmt.Errorf("ripple MakeDepositProposal, crypto.SigToPub error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pub)

	ok, err := node_manager.CheckConsensusSigns(service, scom.MethodImportOuterTransfer, digest, addr, node_manager.Voter)
	if err != nil {
		return nil, fmt.Errorf("ripple MakeDepositProposal, node_manager.CheckConsensusSigns error: %v", err)
	}
	if ok {
		txParam := new(scom.MakeTxParam)
		if err := rlp.DecodeBytes(params.Extra, txParam); err != nil {
			return nil, fmt.Errorf("ripple MakeDepositProposal, deserialize MakeTxParam error:%s", err)
		}
		if err := scom.CheckDoneTx(service, txParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("ripple MakeDepositProposal, check done transaction error:%s", err)
		}
		if err := scom.PutDoneTx(service, txParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("ripple MakeDepositProposal, PutDoneTx error:%s", err)
		}

		//fulfill to contract address
		assetBind, err := side_chain_manager.GetAssetBind(service, params.SourceChainID)
		if err != nil {
			return nil, fmt.Errorf("ripple MakeDepositProposal, side_chain_manager.GetAssetBind error:%s", err)
		}
		txParam.ToContractAddress, ok = assetBind.LockProxyMap[txParam.ToChainID]
		if !ok {
			return nil, fmt.Errorf("ripple MakeDepositProposal, assetBind.LockProxyMap of %d not exist", txParam.ToChainID)
		}

		args := &scom.RippleTxArgs{
			Amount: new(big.Int),
		}
		err = rlp.DecodeBytes(txParam.Args, args)
		if err != nil {
			return nil, fmt.Errorf("ripple MakeDepositProposal, rlp.DecodeBytes error: %s", err)
		}
		b, err := scom.EncodeRippleTxArgs(args)
		txParam.Args = b

		return txParam, nil
	}
	return nil, nil
}

func (this *RippleHandler) MultiSign(service *native.NativeContract) error {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.MultiSignParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodMultiSignRipple, params, ctx.Payload); err != nil {
		return fmt.Errorf("MultiSign, contract params deserialize error: %v", err)
	}

	// get rippleExtraInfo
	rippleExtraInfo, err := side_chain_manager.GetRippleExtraInfo(service, params.ToChainId)
	if err != nil {
		return fmt.Errorf("MultiSign, side_chain_manager.GetRippleExtraInfo error: %v", err)
	}

	// get raw txJsonInfo
	raw, err := GetTxJsonInfo(service, params.FromChainId, params.TxHash)
	if err != nil {
		return fmt.Errorf("MultiSign, get txJsonInfo error: %v", err)
	}

	// check if aleady done
	multisignInfo, err := GetMultisignInfo(service, raw)
	if err != nil {
		return fmt.Errorf("MultiSign, GetMultisignInfo error: %v", err)
	}
	if multisignInfo.Status {
		return nil
	}

	// check if signature is valid
	txJson := new(types.MultisignPayment)
	err = json.Unmarshal([]byte(params.TxJson), txJson)
	if err != nil {
		return fmt.Errorf("MultiSign, unmarshal signed txjson error: %s", err)
	}
	for _, s := range txJson.Signers {
		signerAccount, err := data.NewAccountFromAddress(s.Signer.Account)
		if err != nil {
			return fmt.Errorf("MultiSign, data.NewAccountFromAddress error: %s", err)
		}
		signerPk, err := hex.DecodeString(s.Signer.SigningPubKey)
		if err != nil {
			return fmt.Errorf("MultiSign, hex.DecodeString signer pk error: %s", err)
		}
		signature, err := hex.DecodeString(s.Signer.TxnSignature)
		if err != nil {
			return fmt.Errorf("MultiSign, hex.DecodeString signature error: %s", err)
		}

		// check if valid signer
		flag := false
		for _, v := range rippleExtraInfo.Pks {
			if fmt.Sprintf("%X", v) == s.Signer.SigningPubKey {
				flag = true
				break
			}
		}
		if !flag {
			return fmt.Errorf("MultiSign, signer is not multisign account")
		}

		//check if valid signature
		err = types.CheckMultiSign(raw, *signerAccount, signerPk, signature)
		if err != nil {
			return fmt.Errorf("MultiSign, types.CheckMultiSign error: %s", err)
		}
		signer := &Signer{
			Account:       signerAccount.Bytes(),
			TxnSignature:  signature,
			SigningPubKey: signerPk,
		}
		blob, err := rlp.EncodeToBytes(signer)
		if err != nil {
			return fmt.Errorf("MultiSign, rlp.EncodeToBytes signer error: %v", err)
		}
		multisignInfo.SigMap[hex.EncodeToString(blob)] = true
	}

	if uint64(len(multisignInfo.SigMap)) >= rippleExtraInfo.Quorum {
		payment, err := types.DeserializeRawMultiSignTx(raw)
		if err != nil {
			return fmt.Errorf("MultiSign, types.DeserializeRawMultiSignTx error")
		}
		for s := range multisignInfo.SigMap {
			signerBytes, err := hex.DecodeString(s)
			if err != nil {
				return fmt.Errorf("MultiSign, hex.DecodeString signer bytes error")
			}
			signer := new(Signer)
			err = rlp.DecodeBytes(signerBytes, signer)
			if err != nil {
				return fmt.Errorf("MultiSign, deserialization signer bytes error")
			}
			sig := data.Signer{}
			sig.Signer.SigningPubKey = new(data.PublicKey)
			sig.Signer.TxnSignature = new(data.VariableLength)
			*sig.Signer.TxnSignature = signer.TxnSignature
			copy(sig.Signer.SigningPubKey[:], signer.SigningPubKey)
			acc := data.Account{}
			copy(acc[:], signer.Account)
			sig.Signer.Account = acc
			payment.Signers = append(payment.Signers, sig)
		}

		finalPayment, err := json.Marshal(payment)
		if err != nil {
			return fmt.Errorf("MultiSign, json.Marshal final payment error: %s", err)
		}
		err = service.AddNotify(scom.ABI, []string{cross_chain_manager_abi.EventMultiSign}, params.FromChainId, params.ToChainId,
			hex.EncodeToString(params.TxHash), string(finalPayment), payment.Sequence)
		if err != nil {
			return fmt.Errorf("MultiSign, AddNotify error: %v", err)
		}
		multisignInfo.Status = true
	}
	if err := PutMultisignInfo(service, raw, multisignInfo); err != nil {
		return fmt.Errorf("MultiSign, PutMultisignInfo error: %s", err)
	}
	return nil
}

func (this *RippleHandler) MakeTransaction(service *native.NativeContract, param *scom.MakeTxParam,
	fromChainID uint64) error {
	args, err := scom.DecodeRippleTxArgs(param.Args)
	if err != nil {
		fmt.Errorf("ripple MakeTransaction, deserialize asset hash error")
	}
	toAddrBytes := args.ToAddress
	amount_temp := args.Amount.Uint64()
	amount, err := data.NewAmount(new(big.Int).SetUint64(amount_temp).String())
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, data.NewAmount error: %s", err)
	}

	//get asset map
	assetBind, err := side_chain_manager.GetAssetBind(service, param.ToChainID)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, get asset map error: %s", err)
	}
	lockProxyAddress, ok := assetBind.LockProxyMap[param.ToChainID]
	if !ok {
		return fmt.Errorf("ripple MakeTransaction, lock proxy map of chain %d is not registered", param.ToChainID)
	}
	assetAddress, ok := assetBind.AssetMap[param.ToChainID]
	if !ok {
		return fmt.Errorf("ripple MakeTransaction, asset map of chain %d is not registered", param.ToChainID)
	}
	if hex.EncodeToString(assetAddress) != hex.EncodeToString(param.ToContractAddress) ||
		hex.EncodeToString(assetAddress) != hex.EncodeToString(lockProxyAddress) {
		return fmt.Errorf("ripple MakeTransaction, asset address is not match, assetAddress %x, "+
			"toContractAddress: %x, lockProxyAddress: %x", assetAddress, param.ToContractAddress, lockProxyAddress)
	}

	// get rippleExtraInfo
	rippleExtraInfo, err := side_chain_manager.GetRippleExtraInfo(service, param.ToChainID)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, side_chain_manager.GetRippleExtraInfo error")
	}

	//get fee
	baseFee, err := side_chain_manager.GetFeeObj(service, param.ToChainID)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, side_chain_manager.GetFee error: %v", err)
	}
	if baseFee.View == 0 {
		return fmt.Errorf("ripple MakeTransaction, base fee is not initialized")
	}

	//fee = baseFee * signerNum
	fee_temp := new(big.Int).Mul(baseFee.Fee, new(big.Int).SetUint64(rippleExtraInfo.SignerNum))
	fee, err := data.NewValue(ToStringByPrecise(fee_temp, 6), true)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, data.NewValue fee error: %s", err)
	}
	feeAmount, err := data.NewAmount(fee_temp.String())
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, data.NewAmount fee error: %s", err)
	}
	amountD, err := amount.Subtract(feeAmount)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, amount.Subtract fee error: %s", err)
	}
	reserveAmount, err := data.NewValue(rippleExtraInfo.ReserveAmount.String(), false)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, side_chain_manager.GetFee error: %v", err)
	}
	if amountD.Compare(*reserveAmount) < 0 {
		return fmt.Errorf("ripple MakeTransaction, amount is less than reserveAmount")
	}

	from := new(data.Account)
	to := new(data.Account)
	copy(from[:], assetAddress)
	copy(to[:], toAddrBytes)

	payment := types.GeneratePayment(*from, *to, *amountD, *fee, uint32(rippleExtraInfo.Sequence))
	_, raw, err := data.Raw(payment)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, data.Raw error: %s", err)
	}
	err = service.AddNotify(scom.ABI, []string{cross_chain_manager_abi.EventRippleTx}, fromChainID, param.ToChainID,
		hex.EncodeToString(param.TxHash), hex.EncodeToString(raw), payment.Sequence)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, AddNotify error: %v", err)
	}

	//sequence + 1
	rippleExtraInfo.Sequence = rippleExtraInfo.Sequence + 1
	err = side_chain_manager.PutRippleExtraInfo(service, param.ToChainID, rippleExtraInfo)
	if err != nil {
		return fmt.Errorf("ripple MakeTransaction, side_chain_manager.PutRippleExtraInfo error: %s", err)
	}

	//store txJson info
	PutTxJsonInfo(service, fromChainID, param.TxHash, hex.EncodeToString(raw))
	return nil
}

func (this *RippleHandler) ReconstructTx(service *native.NativeContract) error {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.ReconstructTxParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodReconstructRippleTx, params, ctx.Payload); err != nil {
		return fmt.Errorf("ReconstructTx, contract params deserialize error: %v", err)
	}

	//get tx json info
	raw, err := GetTxJsonInfo(service, params.FromChainId, params.TxHash)
	if err != nil {
		return fmt.Errorf("ReconstructTx, GetTxJsonInfo error: %v", err)
	}

	//get fee
	baseFee, err := side_chain_manager.GetFeeObj(service, params.ToChainId)
	if err != nil {
		return fmt.Errorf("ReconstructTx, side_chain_manager.GetFee error: %v", err)
	}
	if baseFee.View == 0 {
		return fmt.Errorf("ReconstructTx, base fee is not initialized")
	}

	//get ripple extra info
	rippleExtraInfo, err := side_chain_manager.GetRippleExtraInfo(service, params.ToChainId)
	if err != nil {
		return fmt.Errorf("ReconstructTx, side_chain_manager.GetRippleExtraInfo error: %v", err)
	}

	payment, err := types.DeserializeRawMultiSignTx(raw)
	if err != nil {
		return fmt.Errorf("ReconstructTx, types.DeserializeRawMultiSignTx error")
	}

	//fee = baseFee * signerNum
	fee_temp := new(big.Int).Mul(baseFee.Fee, new(big.Int).SetUint64(rippleExtraInfo.SignerNum))
	fee, err := data.NewValue(ToStringByPrecise(fee_temp, 6), true)
	if err != nil {
		return fmt.Errorf("ReconstructTx, data.NewValue fee error: %s", err)
	}

	payment.Fee = *fee
	_, newRaw, err := data.Raw(payment)
	if err != nil {
		return fmt.Errorf("ReconstructTx, data.Raw error: %s", err)
	}

	//store txJson info
	PutTxJsonInfo(service, params.FromChainId, params.TxHash, hex.EncodeToString(newRaw))

	err = service.AddNotify(scom.ABI, []string{cross_chain_manager_abi.EventRippleTx}, params.FromChainId, params.ToChainId,
		hex.EncodeToString(params.TxHash), hex.EncodeToString(newRaw), payment.Sequence)
	if err != nil {
		return fmt.Errorf("ReconstructTx, AddNotify error: %v", err)
	}
	return nil
}
