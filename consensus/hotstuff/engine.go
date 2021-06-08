package hotstuff

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
)

var (
	blockPeriod uint64 = 10 // block period is 10 second

	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now

	inmemoryAddresses  = 20 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)
)

func NewHotstuffConsensusEngine(chainReader consensus.ChainReader, privateKey *ecdsa.PrivateKey) consensus.Engine {
	paceMaker := newPaceMaker(chainReader)
	roundState := newRoundState(privateKey, chainReader)
	paceMaker.core = roundState
	roundState.pace = paceMaker

	ctx := context.Background()
	go roundState.handleSelfMsg(ctx)
	go paceMaker.Start(ctx)

	return roundState
}

func (s *roundState) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header)
}

func (s *roundState) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return s.verifyHeader(chain, header, nil)
}

func (s *roundState) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := s.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (s *roundState) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

func (s *roundState) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = common.Address{}
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

	// Assemble the voting snapshot
	snap, err := s.snapshot(chain)
	if err != nil {
		return err
	}

	// todo: add candidate node as validator, set coinbase as new validator

	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, snap.validators())
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	// todo: use config.BlockPeriod
	header.Time = parent.Time + blockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	return nil
}

func (s *roundState) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
}

func (s *roundState) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil)), nil
}

func (s *roundState) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()
	// Bail out if we're unauthorized to sign a block
	snap, err := s.snapshot(chain)
	if err != nil {
		return err
	}
	if _, v := snap.ValSet.GetByAddress(s.address); v == nil {
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

		// todo(fuk): commit ch should be unbuffer channel, use mutex to lock the state
		// get the proposed block hash and clear it if the seal() is completed.
		//s.sealMu.Lock()
		//s.proposedBlockHash = block.Hash()
		//
		//defer func() {
		//	sb.proposedBlockHash = common.Hash{}
		//	sb.sealMu.Unlock()
		//}()

		// post block into Istanbul engine
		//go s.pace.istanbulEventMux.Post(RequestEvent{block: block})
		s.pace.requestCh <- block
		for {
			select {
			case result := <-s.pace.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				if result != nil && block.Hash() == result.Hash() {
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

func (s *roundState) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

// useless
func (s *roundState) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

func (s *roundState) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "istanbul",
		Version:   "1.0",
		Service:   &API{chain: chain, core: s},
		Public:    true,
	}}
}

func (s *roundState) Close() error {
	return nil
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.HotstuffFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}

	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := GetSignatureAddress(sigHash(header).Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
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

	HotstuffExtra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}

	HotstuffExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&HotstuffExtra)
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

// snapshot
// todo: retrieves the authorization snapshot at a given point in time.
func (s *roundState) snapshot(chain consensus.ChainHeaderReader) (*Snapshot, error) {
	if s.snap == nil {
		genesisHeader := chain.GetHeaderByNumber(0)
		extra, err := types.ExtractHotstuffExtra(genesisHeader)
		if err != nil {
			return nil, err
		}
		valset := newDefaultSet(extra.Validators)
		s.snap = newSnapshot(valset)
	}
	return s.snap.copy(), nil
}

// verifySigner checks whether the signer is in parent's validator set
func (s *roundState) verifySigner(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := s.snapshot(chain)
	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := ecrecover(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.ValSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (s *roundState) verifyCommittedSeals(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := s.snapshot(chain)
	if err != nil {
		return err
	}

	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return errEmptyCommittedSeals
	}

	// Check whether the committed seals are generated by parent's validators
	validSeal := 0
	validators := snap.ValSet.Copy()
	committers, err := Signers(header)
	if err != nil {
		return err
	}
	for _, addr := range committers {
		if validators.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return errInvalidCommittedSeals
	}

	// The length of validSeal should be large than the number of 2F(faulty node)
	if validSeal < snap.major() {
		return errInvalidCommittedSeals
	}

	return nil
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *roundState) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// todo: verify header time
	// Don't waste time checking blocks from the future (adjusting for allowed threshold)
	//adjustedTimeNow := now().Add(time.Duration(sb.config.AllowedFutureBlockTime) * time.Second).Unix()
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

	return sb.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (s *roundState) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
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
		// todo
		parent = s.store.GetHeader(header.ParentHash) //chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	// todo: set config.BlockPeriod
	//if parent.Time+sb.config.BlockPeriod > header.Time {
	//	return errInvalidTimestamp
	//}

	// Verify validators in extraData. Validators in snapshot and extraData should be the same.
	snap, err := s.snapshot(chain)
	if err != nil {
		return err
	}
	validators := make([]byte, len(snap.validators())*common.AddressLength)
	for i, validator := range snap.validators() {
		copy(validators[i*common.AddressLength:], validator[:])
	}
	if err := s.verifySigner(chain, header, parents); err != nil {
		return err
	}

	return s.verifyCommittedSeals(chain, header, parents)
}

// update timestamp and signature of the block based on its number of transactions
func (s *roundState) updateBlock(block *types.Block) (*types.Block, error) {
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

func (s *roundState) sealHeader(header *types.Header) ([]byte, error) {
	return s.Sign(sigHash(header).Bytes())
}

// Sign implements istanbul.Backend.Sign
func (s *roundState) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, s.privateKey)
}
