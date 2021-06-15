package core

import "github.com/ethereum/go-ethereum/common"

type core struct {

}

func (c *core) Start() error {
	return nil
}

func (c *core) Stop() error {
	return nil
}

func (c *core) IsProposer() bool {
	return false
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return false
}

//func (c *core) CurrentRoundState() *roundState {
//	return nil
//}