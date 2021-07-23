package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type treeNode struct {
	block *types.Block
	children []common.Hash
}

type blockTree struct {

}

type BlockPool struct {

}