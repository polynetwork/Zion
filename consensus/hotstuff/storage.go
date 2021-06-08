package hotstuff

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"sync"
)

type Storage struct {
	mtx *sync.Mutex
	consensus.ChainReader
	list []common.Hash
	cache *lru.ARCCache
}

func newStorage(chain consensus.ChainReader) *Storage {
	cache, _ := lru.NewARC(100)
	s := new(Storage)
	s.ChainReader = chain
	s.cache = cache
	s.mtx = new(sync.Mutex)
	return s
}

func (s *Storage) Push(block *types.Block) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.list = append(s.list, block.Hash())
	s.set(block)
}

func (s *Storage) Pop() *types.Block {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.list) == 0 {
		return nil
	}
	hash := s.list[0]
	s.list = s.list[1:]
	block := s.get(hash)
	return block
}

func (s *Storage) Put(block *types.Block) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.set(block)
}

func (s *Storage) Get(hash common.Hash) *types.Block {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.get(hash)
}

func (s *Storage) GetHeader(hash common.Hash) *types.Header {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	block := s.get(hash)
	if block == nil {
		return nil
	}
	return block.Header()
}

func (s *Storage) set(block *types.Block) {
	key := block.Hash()
	s.cache.Add(key, block)
}

func (s *Storage) get(hash common.Hash) *types.Block {
	data, ok := s.cache.Get(hash)
	if ok {
		return data.(*types.Block)
	}
	header := s.GetHeaderByHash(hash)
	if header == nil {
		return nil
	}
	block := s.GetBlock(hash, header.Number.Uint64())
	return block
}
