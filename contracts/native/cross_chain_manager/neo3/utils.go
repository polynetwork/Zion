package neo3

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/neo3_state_manager"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/mpt"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
)

func verifyCrossChainMsg(native *native.NativeContract, ccm *CrossChainMsg, magic uint32) error {
	// get neo3 state validator from native contract
	svListBytes, err := neo3_state_manager.GetCurrentStateValidator(native)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, neo3_state_manager.GetCurrentStateValidator error: %v", err)
	}
	svStrings, err := neo3_state_manager.DeserializeStringArray(svListBytes)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, neo3_state_manager.DeserializeStringArray error: %v", err)
	}
	pubKeys := make([]crypto.ECPoint, len(svStrings), len(svStrings))
	for i, v := range svStrings {
		pubKey, err := crypto.NewECPointFromString(v)
		if err != nil {
			return fmt.Errorf("verifyCrossChainMsg, crypto.NewECPointFromString error: %v", err)
		}
		pubKeys[i] = *pubKey
	}
	n := len(pubKeys)
	m := n - (n-1)/3
	msc, err := sc.CreateMultiSigContract(m, pubKeys) // sort public keys inside
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, sc.CreateMultiSigContract error: %v", err)
	}
	expected := msc.GetScriptHash()
	got, err := ccm.GetScriptHash()
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, getScripthash error: %v", err)
	}
	// compare state validator
	if !expected.Equals(got) {
		return fmt.Errorf("verifyCrossChainMsg, invalid script hash in NeoCrossChainMsg error, expected: %s, got: %s", expected.String(), got.String())
	}
	msg, err := ccm.GetMessage(magic)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, unable to get unsigned message of neo crossChainMsg")
	}
	// verify witness
	if len(ccm.Witnesses) == 0 {
		return fmt.Errorf("verifyCrossChainMsg, incorrect witness length")
	}
	invScript, err := crypto.Base64Decode(ccm.Witnesses[0].Invocation)
	if err != nil {
		return fmt.Errorf("crypto.Base64Decode, decode invocation script error: %v", err)
	}
	verScript, err := crypto.Base64Decode(ccm.Witnesses[0].Verification)
	if err != nil {
		return fmt.Errorf("crypto.Base64Decode, decode verification script error: %v", err)
	}
	witness := &tx.Witness{
		InvocationScript:   invScript,
		VerificationScript: verScript,
	}
	v1 := tx.VerifyMultiSignatureWitness(msg, witness)
	if !v1 {
		return fmt.Errorf("verifyCrossChainMsg, verify witness failed, height: %d", ccm.Index)
	}
	return nil
}

func verifyFromNeoTx(proof []byte, ccm *CrossChainMsg, contractId int) (*scom.MakeTxParam, error) {
	root, err := helper.UInt256FromString(ccm.RootHash)
	if err != nil {
		return nil, fmt.Errorf("UInt256FromString error: %v", err)
	}
	value, err := verifyNeoCrossChainProof(proof, root.ToByteArray(), contractId)
	if err != nil {
		return nil, fmt.Errorf("verifyNeoCrossChainProof error: %v", err)
	}

	txParam, err := DeserializeNeoMakeTxParam(value)
	if err != nil {
		return nil, fmt.Errorf("DecodeTxParam error: %v", err)
	}
	return txParam.MakeTxParam, nil
}

func verifyNeoCrossChainProof(proof []byte, stateRoot []byte, contractId int) ([]byte, error) {
	id, key, proofs, err := mpt.ResolveProof(proof)
	if err != nil {
		return nil, fmt.Errorf("VerifyNeoCrossChainProof, neo3-gogogo mpt.ResolveProof error: %v", err)
	}
	if id != contractId {
		return nil, fmt.Errorf("VerifyNeoCrossChainProof, error: id is not CCMC contract id, expected: %d, but got: %d", contractId, id)
	}
	root := helper.UInt256FromBytes(stateRoot)
	value, err := mpt.VerifyProof(root, contractId, key, proofs)
	if err != nil {
		return nil, fmt.Errorf("VerifyNeoCrossChainProof, neo3-gogogo mpt.VerifyProof error: %v", err)
	}
	return value, nil
}
