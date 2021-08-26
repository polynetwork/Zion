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
package utils

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// MethodID only used for register method handler and prepare native contract context ref.
func MethodID(ab *abi.ABI, name string) string {
	m, ok := ab.Methods[name]
	if !ok {
		panic(fmt.Sprintf("method name %s not exist", name))
	}
	return hexutil.Encode(m.ID)
}

func PackMethodWithStruct(ab *abi.ABI, name string, data interface{}) ([]byte, error) {

	value := reflect.ValueOf(data).Elem()

	var args []interface{}
	n := value.NumField()
	for i := 0; i < n; i++ {
		fv := value.Field(i)
		args = append(args, fv.Interface())
	}

	return PackMethod(ab, name, args...)
}

func PackMethod(ab *abi.ABI, name string, args ...interface{}) ([]byte, error) {
	method, exist := ab.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(method.ID, arguments...), nil
}

func UnpackMethod(ab *abi.ABI, name string, data interface{}, payload []byte) error {
	mth, ok := ab.Methods[name]
	if !ok {
		return fmt.Errorf("abi method %s not exist", name)
	}

	if len(payload) < 4 || len(payload[4:])%32 != 0 {
		return fmt.Errorf("invalid payload")
	}

	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("method input should be pointer")
	}

	args := mth.Inputs
	unpacked, err := args.Unpack(payload[4:])
	if err != nil {
		return err
	}
	return args.Copy(data, unpacked)
}

func PackOutputs(ab *abi.ABI, method string, args ...interface{}) ([]byte, error) {
	mth, exist := ab.Methods[method]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", method)
	}
	return mth.Outputs.Pack(args...)
}

func UnpackOutputs(ab *abi.ABI, name string, data interface{}, payload []byte) error {
	mth, ok := ab.Methods[name]
	if !ok {
		return fmt.Errorf("abi method %s not exist", name)
	}

	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("method output should be pointer")
	}

	args := mth.Outputs
	unpacked, err := args.Unpack(payload)
	if err != nil {
		return err
	}
	return args.Copy(data, unpacked)
}

func PackEvents(ab *abi.ABI, event string, args ...interface{}) ([]byte, error) {
	evt, exist := ab.Events[event]
	if !exist {
		return nil, fmt.Errorf("event '%s' not found", event)
	}
	return evt.Inputs.Pack(args...)
}
