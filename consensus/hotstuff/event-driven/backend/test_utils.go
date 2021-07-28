package backend

import (
	"bytes"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	snr "github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	elog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	testLogger = elog.New()
)

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey, hotstuff.ValidatorSet) {
	// Setup validators
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}

	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig
	// force enable Istanbul engine
	genesis.Config.HotStuff = &params.HotStuffConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = defaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.HotstuffDigest

	appendValidators(genesis, addrs)
	valset := makeValSet(addrs)
	return genesis, nodeKeys, valset
}

func appendValidators(genesis *core.Genesis, addrs []common.Address) {

	if len(genesis.ExtraData) < types.HotstuffExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.HotstuffExtraVanity]

	ist := &types.HotstuffExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func makeHeader(parent *types.Block, config *hotstuff.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       parent.Time() + config.BlockPeriod,

		Difficulty: defaultDifficulty,
	}
	return header
}

func makeBlock(t *testing.T, chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	block := makeBlockWithoutSeal(chain, engine, parent)
	header := block.Header()

	assert.NoError(t, engine.signer.SealBeforeCommit(header))
	expectBlock := block.WithSeal(header)

	resultCh := make(chan *types.Block, 10)
	go func() {
		if err := engine.Seal(chain, expectBlock, resultCh, make(chan struct{})); err != nil {
			t.Errorf("seal block failed, err: %s", err)
		}
	}()

	return <-resultCh
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	header := makeHeader(parent, engine.config)
	engine.Prepare(chain, header)
	state, _ := chain.StateAt(parent.Root())
	block, _ := engine.FinalizeAndAssemble(chain, header, state, nil, nil, nil)
	return block
}

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func singleNodeChain() (*core.BlockChain, *backend) {
	testLogger.SetHandler(elog.StdoutHandler)

	genesis, nodeKeys, valset := getGenesisAndKeys(1)
	memDB := rawdb.NewMemoryDatabase()
	config := hotstuff.DefaultBasicConfig
	// Use the first key as private key
	b, _ := New(config, nodeKeys[0], memDB, valset).(*backend)
	genesis.MustCommit(memDB)

	txLookUpLimit := uint64(100)
	cacheConfig := &core.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
	}
	blockchain, err := core.NewBlockChain(memDB, cacheConfig, genesis.Config, b, vm.Config{}, nil, &txLookUpLimit)
	if err != nil {
		panic(err)
	}

	b.Start(blockchain, blockchain.CurrentBlock, nil)

	return blockchain, b
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestValidatorSet(n int) (hotstuff.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, hotstuff.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}

func postAndWait(backend *backend, block *types.Block, t *testing.T) {
	eventSub := backend.EventMux().Subscribe(hotstuff.RequestEvent{})
	defer eventSub.Unsubscribe()
	stop := make(chan struct{}, 1)
	eventLoop := func() {
		<-eventSub.Chan()
		stop <- struct{}{}
	}
	go eventLoop()
	if err := backend.EventMux().Post(hotstuff.RequestEvent{
		Proposal: block,
	}); err != nil {
		t.Fatalf("%s", err)
	}
	<-stop
}

func buildArbitraryP2PNewBlockMessage(t *testing.T, invalidMsg bool) (*types.Block, p2p.Msg) {
	arbitraryBlock := types.NewBlock(&types.Header{
		Number:    big.NewInt(1),
		GasLimit:  0,
		MixDigest: types.HotstuffDigest,
	}, nil, nil, nil, nil)
	request := []interface{}{&arbitraryBlock, big.NewInt(1)}
	if invalidMsg {
		request = []interface{}{"invalid msg"}
	}
	size, r, err := rlp.EncodeToReader(request)
	if err != nil {
		t.Fatalf("can't encode due to %s", err)
	}
	payload, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("can't read payload due to %s", err)
	}
	arbitraryP2PMessage := p2p.Msg{Code: NewBlockMsg, Size: uint32(size), Payload: bytes.NewReader(payload)}
	return arbitraryBlock, arbitraryP2PMessage
}

var emptySigner = &snr.SignerImpl{}

func (s *backend) UpdateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	if err := s.signer.SealBeforeCommit(header); err != nil {
		return nil, err
	}
	newBlock := block.WithSeal(header)
	return newBlock, nil
}

func makeValSet(validators []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(validators, hotstuff.RoundRobin)
}

func newTestSigner() hotstuff.Signer {
	key, _ := generatePrivateKey()
	return snr.NewSigner(key, 3)
}
