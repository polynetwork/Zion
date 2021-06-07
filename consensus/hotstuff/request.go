package hotstuff

import "github.com/ethereum/go-ethereum/core/types"

type RequestEvent struct {
	block *types.Block
}
