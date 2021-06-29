package backend

import (
	"bytes"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	hsc "github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	lru "github.com/hashicorp/golang-lru"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
	inmemoryPeers      = 1000
	inmemoryMessages   = 1024
)

// HotStuff protocol constants.
var (
	extraVanity = crypto.DigestLength    // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for signer seal

	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	uncleHash         = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	now               = time.Now

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.
)

func (s *backend) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, s.signatures)
}

func (s *backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return s.verifyHeader(chain, header, nil, seal)
}

func (s *backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			seal := false
			if seals != nil && len(seals) > i {
				seal = seals[i]
			}
			err := s.verifyHeader(chain, header, headers[:i], seal)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (s *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

func (s *backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = s.signer
	header.Nonce = emptyNonce
	header.MixDigest = types.HotstuffDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	valset := s.getValidatorsAddress()

	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, valset)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	header.Time = parent.Time + s.config.BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	return nil
}

func (s *backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
}

func (s *backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil)), nil
}

func (s *backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()
	// Bail out if we're unauthorized to sign a block

	snap := s.getValidators()
	if _, v := snap.GetByAddress(s.Address()); v == nil {
		return errUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	block, err = s.updateBlock(block)
	if err != nil {
		return err
	}

	delay := time.Unix(int64(block.Header().Time), 0).Sub(now())

	go func() {
		// wait for the timestamp of header, use this to adjust the block period
		select {
		case <-time.After(delay):
		case <-stop:
			results <- nil
			return
		}

		// get the proposed block hash and clear it if the seal() is completed.
		s.sealMu.Lock()
		s.proposedBlockHash = block.Hash()

		defer func() {
			s.proposedBlockHash = common.Hash{}
			s.sealMu.Unlock()
		}()
		// post block into Istanbul engine
		go s.EventMux().Post(hotstuff.RequestEvent{
			Proposal: block,
		})
		for {
			select {
			case result := <-s.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				// todo: if result != nil && block.Hash() == result.Hash() {
				if result != nil {
					results <- result
					return
				}
			case <-stop:
				results <- nil
				return
			}
		}
	}()
	return nil
}

func (s *backend) SealHash(header *types.Header) common.Hash {
	return basic.SigHash(header)
}

// useless
func (s *backend) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

func (s *backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "istanbul",
		Version:   "1.0",
		Service:   &API{chain: chain, hotstuff: s},
		Public:    true,
	}}
}

func (s *backend) Close() error {
	return nil
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := sigcache.Get(hash); ok {
		return addr.(common.Address), nil
	}

	// Retrieve the signature from the header extra-data
	hotstuffExtra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := basic.GetSignatureAddress(basic.SigHash(header).Bytes(), hotstuffExtra.Seal)
	if err != nil {
		return addr, err
	}
	sigcache.Add(hash, addr)
	return addr, nil
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.HotstuffExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.HotstuffExtraVanity])

	ist := &types.HotstuffExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.HotstuffExtraSeal != 0 {
		return errInvalidSignature
	}

	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}

	extra.Seal = seal
	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.HotstuffExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return errInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.HotstuffExtraSeal {
			return errInvalidCommittedSeals
		}
	}

	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}

	extra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(extra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.HotstuffExtraVanity], payload...)
	return nil
}

// verifySigner checks whether the signer is in parent's validator set
func (s *backend) verifySigner(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	snap := s.getValidators()

	// resolve the authorization key and check against signers
	signer, err := ecrecover(header, s.signatures)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (s *backend) verifyCommittedSeals(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return errEmptyCommittedSeals
	}

	snap := s.getValidators()
	validators := snap.Copy()
	// Check whether the committed seals are generated by parent's validators
	committers, err := s.Signers(header)
	if err != nil {
		return err
	}
	return s.checkValidatorQuorum(committers, validators)
}

func (s *backend) checkValidatorQuorum(committers []common.Address, validators hotstuff.ValidatorSet) error {
	validSeal := 0
	for _, addr := range committers {
		if validators.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return errInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= validators.Q() {
		return errInvalidCommittedSeals
	}
	return nil
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (s *backend) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, seal bool) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// todo: verify header time
	// Don't waste time checking blocks from the future (adjusting for allowed threshold)
	//adjustedTimeNow := now().Add(time.Duration(s.config.AllowedFutureBlockTime) * time.Second).Unix()
	//if header.Time > uint64(adjustedTimeNow) {
	//	return consensus.ErrFutureBlock
	//}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractHotstuffExtra(header); err != nil {
		return errInvalidExtraDataFormat
	}

	// todo: check nonce
	// Ensure that the coinbase is valid
	//if header.Nonce != (emptyNonce) && !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
	//	return errInvalidNonce
	//}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.HotstuffDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in Istanbul
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	return s.verifyCascadingFields(chain, header, parents, seal)
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func (s *backend) Signers(header *types.Header) ([]common.Address, error) {
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return []common.Address{}, err
	}

	return s.signersFromCommittedSeals(header.Hash(), extra.CommittedSeal)
}

func (s *backend) signersFromCommittedSeals(hash common.Hash, seals [][]byte) ([]common.Address, error) {
	var addrs []common.Address
	proposalSeal := hsc.PrepareCommittedSeal(hash)

	// 1. Get committed seals from current header
	for _, seal := range seals {
		// 2. Get the original address by seal and parent block hash
		addr, err := hsc.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			s.logger.Error("not a valid address", "err", err)
			return nil, errInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (s *backend) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, seal bool) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time+s.config.BlockPeriod > header.Time {
		return errInvalidTimestamp
	}
	// Verify validators in extraData. Validators in snapshot and extraData should be the same.
	snap := s.getValidators()
	validators := make([]byte, snap.Size()*common.AddressLength)
	for i, validator := range snap.List() {
		copy(validators[i*common.AddressLength:], validator.Address().Bytes()[:])
	}
	if err := s.verifySigner(chain, header, parents); err != nil {
		return err
	}

	// verify unsealed proposal
	if seal {
		return s.verifyCommittedSeals(chain, header, parents)
	}
	return nil
}

// update timestamp and signature of the block based on its number of transactions
func (s *backend) updateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	// sign the hash
	seal, err := s.sealHeader(header)
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

func (s *backend) sealHeader(header *types.Header) ([]byte, error) {
	return s.Sign(basic.SigHash(header).Bytes())
}

// Start and Stop invoked in worker.go - start and stop
// Start implements consensus.HotStuff.Start
// Start implements consensus.Istanbul.Start
func (s *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	s.coreMu.Lock()
	defer s.coreMu.Unlock()
	if s.coreStarted {
		return hotstuff.ErrStartedEngine
	}

	// clear previous data
	s.proposedBlockHash = common.Hash{}
	if s.commitCh != nil {
		close(s.commitCh)
	}
	s.commitCh = make(chan *types.Block, 1)

	s.chain = chain
	s.currentBlock = currentBlock
	s.hasBadBlock = hasBadBlock

	if err := s.core.Start(); err != nil {
		return err
	}

	s.coreStarted = true
	return nil
}

// Stop implements consensus.Istanbul.Stop
func (s *backend) Stop() error {
	s.coreMu.Lock()
	defer s.coreMu.Unlock()
	if !s.coreStarted {
		return hotstuff.ErrStoppedEngine
	}
	if err := s.core.Stop(); err != nil {
		return err
	}
	s.coreStarted = false
	return nil
}
