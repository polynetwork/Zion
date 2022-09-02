package neo3_state_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/joeqian10/neo3-gogogo/io"
)

func SerializeStringArray(data []string) ([]byte, error) {
	bw := io.NewBufBinaryWriter()
	bw.WriteVarUInt(uint64(len(data)))
	if bw.Err != nil {
		return nil, fmt.Errorf("WriteVarUInt error: %v", bw.Err)
	}
	for _, v := range data {
		bw.WriteVarString(v)
	}
	if bw.Err != nil {
		return nil, fmt.Errorf("WriteVarString error: %v", bw.Err)
	}
	return bw.Bytes(), nil
}

func DeserializeStringArray(source []byte) ([]string, error) {
	if len(source) == 0 {
		return []string{}, nil
	}
	br := io.NewBinaryReaderFromBuf(source)
	n := br.ReadVarUInt()
	if br.Err != nil {
		return nil, fmt.Errorf("ReadVarUInt error: %v", br.Err)
	}
	result := make([]string, 0, n)
	for i := 0; uint64(i) < n; i++ {
		s := br.ReadVarString(128)
		result = append(result, s)
	}
	return result, nil
}

func getStateValidators(native *native.NativeContract) ([]byte, error) {
	contract := utils.Neo3StateManagerContractAddress
	svBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(STATE_VALIDATOR)))
	if err != nil {
		return nil, fmt.Errorf("CacheDB.Get error: %v", err)
	}
	if svBytes == nil {
		return []byte{}, nil
	}

	return svBytes, nil
}

func putStateValidators(native *native.NativeContract, stateValidators []string) error {
	contract := utils.Neo3StateManagerContractAddress
	// get current stored value
	oldSvBytes, err := getStateValidators(native)
	if err != nil {
		return fmt.Errorf("getStateValidators error: %v", err)
	}
	oldSVs, err := DeserializeStringArray(oldSvBytes)
	if err != nil {
		return fmt.Errorf("DeserializeStringArray error: %v", err)
	}
	// max capacity = len(oldSVs)+len(stateValidators)
	newSVs := make([]string, 0, len(oldSVs)+len(stateValidators))
	newSVs = append(newSVs, oldSVs...)
	// filter duplicate svs
	for _, sv := range stateValidators {
		isInOld := false
		for _, oldSv := range oldSVs {
			if sv == oldSv {
				isInOld = true
				break
			}
		}
		if !isInOld {
			newSVs = append(newSVs, sv)
		}
	}
	// convert back to []byte
	data, err := SerializeStringArray(newSVs)
	if err != nil {
		return fmt.Errorf("SerializeStringArray error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR)), data)
	return nil
}

func removeStateValidators(native *native.NativeContract, stateValidators []string) error {
	contract := utils.Neo3StateManagerContractAddress
	// get current stored value
	oldSvBytes, err := getStateValidators(native)
	if err != nil {
		return fmt.Errorf("getStateValidators error: %v", err)
	}
	oldSVs, err := DeserializeStringArray(oldSvBytes)
	if err != nil {
		return fmt.Errorf("DeserializeStringArray error: %v", err)
	}
	// remove in the slice
	for _, sv := range stateValidators {
		for i, oldSv := range oldSVs {
			if sv == oldSv {
				oldSVs = append(oldSVs[:i], oldSVs[i+1:]...)
				break
			}
		}
	}
	// if no sv left, delete the storage, else put remaining back
	if len(oldSVs) == 0 {
		native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(STATE_VALIDATOR)))
		return nil
	}
	// convert back to []byte
	data, err := SerializeStringArray(oldSVs)
	if err != nil {
		return fmt.Errorf("SerializeStringArray error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR)), data)
	return nil
}

//-----
//below methods are for state validator apply
//-----

func getStateValidatorApply(native *native.NativeContract, applyID uint64) (*StateValidatorListParam, error) {
	contract := utils.Neo3StateManagerContractAddress
	svListParamBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_APPLY), utils.GetUint64Bytes(applyID)))
	if err != nil {
		return nil, fmt.Errorf("CacheDB.Get error: %v", err)
	}
	if len(svListParamBytes) == 0 {
		return nil, fmt.Errorf("cannot find any record")
	}
	svListParam := new(StateValidatorListParam)
	err = rlp.DecodeBytes(svListParamBytes, svListParam)
	if err != nil {
		return nil, fmt.Errorf("rlp.DecodeBytes error: %v", err)
	}
	return svListParam, nil
}

func putStateValidatorApply(native *native.NativeContract, stateValidatorListParam *StateValidatorListParam) error {
	contract := utils.Neo3StateManagerContractAddress
	applyID, err := getStateValidatorApplyID(native)
	if err != nil {
		return fmt.Errorf("getStateValidatorApplyID error: %v", err)
	}
	newApplyID := applyID + 1
	err = putStateValidatorApplyID(native, newApplyID)
	if err != nil {
		return fmt.Errorf("putStateValidatorApplyID error: %v", err)
	}

	blob, err := rlp.EncodeToBytes(stateValidatorListParam)
	if err != nil {
		return fmt.Errorf("rlp.EncodeToBytes error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_APPLY), utils.GetUint64Bytes(applyID)), blob)

	err = native.AddNotify(ABI, []string{EventRegisterStateValidator}, applyID)
	if err != nil {
		return fmt.Errorf("AddNotify error: %v", err)
	}
	return nil
}

func getStateValidatorApplyID(native *native.NativeContract) (uint64, error) {
	contract := utils.Neo3StateManagerContractAddress
	applyIDBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_APPLY_ID)))
	if err != nil {
		return 0, fmt.Errorf("CacheDB.Get error: %v", err)
	}
	var applyID uint64 = 0
	if len(applyIDBytes) != 0 {
		applyID = utils.GetBytesUint64(applyIDBytes)
	}
	return applyID, nil
}

func putStateValidatorApplyID(native *native.NativeContract, applyID uint64) error {
	contract := utils.Neo3StateManagerContractAddress
	applyIDBytes := utils.GetUint64Bytes(applyID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_APPLY_ID)), applyIDBytes)
	return nil
}

//-----
//below methods are for state validator removal
//-----

func getStateValidatorRemove(native *native.NativeContract, removeID uint64) (*StateValidatorListParam, error) {
	contract := utils.Neo3StateManagerContractAddress
	svListParamBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_REMOVE), utils.GetUint64Bytes(removeID)))
	if err != nil {
		return nil, fmt.Errorf("CacheDB.Get error: %v", err)
	}
	if len(svListParamBytes) == 0 {
		return nil, fmt.Errorf("cannot find any record")
	}

	svListParam := new(StateValidatorListParam)
	err = rlp.DecodeBytes(svListParamBytes, svListParam)
	if err != nil {
		return nil, fmt.Errorf("rlp.DecodeBytes error: %v", err)
	}
	return svListParam, nil
}

func putStateValidatorRemove(native *native.NativeContract, svListParam *StateValidatorListParam) error {
	contract := utils.Neo3StateManagerContractAddress
	removeID, err := getStateValidatorRemoveID(native)
	if err != nil {
		return fmt.Errorf("getStateValidatorRemoveID error: %v", err)
	}
	newRemoveID := removeID + 1
	err = putStateValidatorRemoveID(native, newRemoveID)
	if err != nil {
		return fmt.Errorf("putStateValidatorRemoveID error: %v", err)
	}

	blob, err := rlp.EncodeToBytes(svListParam)
	if err != nil {
		return fmt.Errorf("rlp.EncodeToBytes error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_REMOVE), utils.GetUint64Bytes(removeID)), blob)

	err = native.AddNotify(ABI, []string{EventRemoveStateValidator}, removeID)
	if err != nil {
		return fmt.Errorf("AddNofity error: %v", err)
	}
	return nil
}

func getStateValidatorRemoveID(native *native.NativeContract) (uint64, error) {
	contract := utils.Neo3StateManagerContractAddress
	removeIDBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_REMOVE_ID)))
	if err != nil {
		return 0, fmt.Errorf("CacheDB.Get error: %v", err)
	}
	var removeID uint64 = 0
	if len(removeIDBytes) != 0 {
		removeID = utils.GetBytesUint64(removeIDBytes)
	}
	return removeID, nil
}

func putStateValidatorRemoveID(native *native.NativeContract, removeID uint64) error {
	contract := utils.Neo3StateManagerContractAddress
	removeIDBytes := utils.GetUint64Bytes(removeID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(STATE_VALIDATOR_REMOVE_ID)), removeIDBytes)
	return nil
}
