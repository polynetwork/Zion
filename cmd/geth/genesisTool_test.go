package main

import (
	"errors"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"testing"
)

func TestGenesisTool(t *testing.T){
	type TestCase struct {
		AllArgs []string
		Expect error
		AfterHandler func(c *TestCase)
	}
	testCases := []*TestCase{
		{
			[]string{"genesisTool", "generate", "3"},
			errors.New("Fatal: got 3 nodes, but hotstuff BFT requires at least 4 nodes"),
			nil,
		},
		{
			[]string{"genesisTool", "generate", "5", "-basePath", "./temp/"},
			nil,
			func(c *TestCase) {
				basePath := c.AllArgs[4]
				utils.DeleteBasePath(basePath)
			},
		},
	}

	for _, testCase := range testCases {
		geth := runGeth(t, testCase.AllArgs...)
		if testCase.Expect != nil {
			geth.ExpectRegexp(testCase.Expect.Error())
		}
		geth.ExpectRegexp("")
		if testCase.AfterHandler != nil {
			testCase.AfterHandler(testCase)
		}
	}
}
