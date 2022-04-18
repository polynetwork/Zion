package maas_config

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	resetTestContext()
	type TestCase struct {
		Key []byte
		Value []byte
		expect error
	}

	cases := []TestCase {
		{
			Key: getOwnerKey(),
			Value: []byte("0x2D3913c12ACa0E4A2278f829Fb78A682123c0129"),
			expect: nil,
		},
		{
			Key: []byte("testKey"),
			Value: []byte("testValue"),
			expect: errors.New("CacheDB should only be used for native contract storage"),
		},
		{
			Key: []byte("0x2D3913c12ACa0E4A2278f829Fb78A682123c0129"),
			Value: []byte("testValue"),
			expect: nil,
		},
	}
	for _, testCase := range cases {
		key := testCase.Key
		value := testCase.Value
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, err, testCase.expect.Error())
				t.Log("error:", err)
			}
		}()
		set(testEmptyCtx, key, value)
		res, err := get(testEmptyCtx, key)
		assert.NoError(t, err)
		t.Log(string(res))
		assert.Equal(t, testCase.Value, res)
	}
}

func TestSetDelAndGet(t *testing.T) {
	resetTestContext()
	type TestCase struct {
		Key []byte
		Value []byte
		BeforeHandler func(testCase *TestCase)
		AfterHandler func(testCase *TestCase)
		expect error
	}

	cases := []TestCase{
		{
			Key:    getOwnerKey(),
			Value:  []byte("0x2D3913c12ACa0E4A2278f829Fb78A682123c0129"),
			BeforeHandler: func(testCase *TestCase) {
				key := testCase.Key
				value := testCase.Value
				set(testEmptyCtx, key, value)
			},
			AfterHandler: func(testCase *TestCase) {
				key := testCase.Key
				res, err := get(testEmptyCtx, key)
				assert.NoError(t, err)
				t.Log(string(res))
				assert.Equal(t, res, testCase.Value)
				del(testEmptyCtx, key)
				res, err = get(testEmptyCtx, key)
				assert.NoError(t, err)
				t.Log(string(res))
				assert.Equal(t, res, []byte(nil))
			},
			expect: nil,
		},
		{
			Key: []byte("testKey"),
			Value: []byte("testValue"),
			BeforeHandler: func(testCase *TestCase) {
				set(testEmptyCtx, testCase.Key, testCase.Value)
			},
			AfterHandler: nil,
			expect: errors.New("CacheDB should only be used for native contract storage"),
		},
	}

	for _, testCase := range cases {
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, err, testCase.expect.Error())
				t.Log("error:", err)
			}
		}()
		testCase.BeforeHandler(&testCase)
		testCase.AfterHandler(&testCase)
	}


}