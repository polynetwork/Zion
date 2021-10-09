// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package header_sync_abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

var (
	MethodName = "name"

	MethodSyncBlockHeader = "syncBlockHeader"

	MethodSyncCrossChainMsg = "syncCrossChainMsg"

	MethodSyncGenesisHeader = "syncGenesisHeader"
)

// HeaderSyncABI is the input ABI used to generate the binding from.
const HeaderSyncABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"BlockHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"NextValidatorsHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"InfoChainID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlockHeight\",\"type\":\"uint64\"}],\"name\":\"OKEpochSwitchInfoEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"blockHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"BlockHeight\",\"type\":\"uint256\"}],\"name\":\"syncHeader\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"Headers\",\"type\":\"bytes[]\"}],\"name\":\"syncBlockHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"CrossChainMsgs\",\"type\":\"bytes[]\"}],\"name\":\"syncCrossChainMsg\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"GenesisHeader\",\"type\":\"bytes\"}],\"name\":\"syncGenesisHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// HeaderSyncFuncSigs maps the 4-byte function signature to its string representation.
var HeaderSyncFuncSigs = map[string]string{
	"06fdde03": "name()",
	"72ce6700": "syncBlockHeader(uint64,address,bytes[])",
	"21b5cff5": "syncCrossChainMsg(uint64,address,bytes[])",
	"b5ace618": "syncGenesisHeader(uint64,bytes)",
}

// HeaderSyncBin is the compiled bytecode used for deploying new contracts.
var HeaderSyncBin = "0x608060405234801561001057600080fd5b5061034b806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806306fdde031461005157806321b5cff51461006957806372ce670014610069578063b5ace61814610090575b600080fd5b60606040516100609190610279565b60405180910390f35b610080610077366004610133565b60009392505050565b6040519015158152602001610060565b61008061009e36600461022b565b600092915050565b600082601f8301126100b757600080fd5b813567ffffffffffffffff8111156100d1576100d16102ff565b6100e4601f8201601f19166020016102ce565b8181528460208386010111156100f957600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461012e57600080fd5b919050565b60008060006060848603121561014857600080fd5b61015184610116565b92506020848101356001600160a01b038116811461016e57600080fd5b9250604085013567ffffffffffffffff8082111561018b57600080fd5b818701915087601f83011261019f57600080fd5b8135818111156101b1576101b16102ff565b8060051b6101c08582016102ce565b8281528581019085870183870188018d10156101db57600080fd5b60009350835b85811015610218578135878111156101f7578586fd5b6102058f8b838c01016100a6565b85525092880192908801906001016101e1565b5050809750505050505050509250925092565b6000806040838503121561023e57600080fd5b61024783610116565b9150602083013567ffffffffffffffff81111561026357600080fd5b61026f858286016100a6565b9150509250929050565b600060208083528351808285015260005b818110156102a65785810183015185820160400152820161028a565b818111156102b8576000604083870101525b50601f01601f1916929092016040019392505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156102f7576102f76102ff565b604052919050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212204231d0ee59fc1838c92d148dc40fade4e810e808d3b6d562c1891ea89e6b013a64736f6c63430008060033"

// DeployHeaderSync deploys a new Ethereum contract, binding an instance of HeaderSync to it.
func DeployHeaderSync(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *HeaderSync, error) {
	parsed, err := abi.JSON(strings.NewReader(HeaderSyncABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HeaderSyncBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HeaderSync{HeaderSyncCaller: HeaderSyncCaller{contract: contract}, HeaderSyncTransactor: HeaderSyncTransactor{contract: contract}, HeaderSyncFilterer: HeaderSyncFilterer{contract: contract}}, nil
}

// HeaderSync is an auto generated Go binding around an Ethereum contract.
type HeaderSync struct {
	HeaderSyncCaller     // Read-only binding to the contract
	HeaderSyncTransactor // Write-only binding to the contract
	HeaderSyncFilterer   // Log filterer for contract events
}

// HeaderSyncCaller is an auto generated read-only Go binding around an Ethereum contract.
type HeaderSyncCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeaderSyncTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HeaderSyncTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeaderSyncFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HeaderSyncFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeaderSyncSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HeaderSyncSession struct {
	Contract     *HeaderSync       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HeaderSyncCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HeaderSyncCallerSession struct {
	Contract *HeaderSyncCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// HeaderSyncTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HeaderSyncTransactorSession struct {
	Contract     *HeaderSyncTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// HeaderSyncRaw is an auto generated low-level Go binding around an Ethereum contract.
type HeaderSyncRaw struct {
	Contract *HeaderSync // Generic contract binding to access the raw methods on
}

// HeaderSyncCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HeaderSyncCallerRaw struct {
	Contract *HeaderSyncCaller // Generic read-only contract binding to access the raw methods on
}

// HeaderSyncTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HeaderSyncTransactorRaw struct {
	Contract *HeaderSyncTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHeaderSync creates a new instance of HeaderSync, bound to a specific deployed contract.
func NewHeaderSync(address common.Address, backend bind.ContractBackend) (*HeaderSync, error) {
	contract, err := bindHeaderSync(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HeaderSync{HeaderSyncCaller: HeaderSyncCaller{contract: contract}, HeaderSyncTransactor: HeaderSyncTransactor{contract: contract}, HeaderSyncFilterer: HeaderSyncFilterer{contract: contract}}, nil
}

// NewHeaderSyncCaller creates a new read-only instance of HeaderSync, bound to a specific deployed contract.
func NewHeaderSyncCaller(address common.Address, caller bind.ContractCaller) (*HeaderSyncCaller, error) {
	contract, err := bindHeaderSync(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HeaderSyncCaller{contract: contract}, nil
}

// NewHeaderSyncTransactor creates a new write-only instance of HeaderSync, bound to a specific deployed contract.
func NewHeaderSyncTransactor(address common.Address, transactor bind.ContractTransactor) (*HeaderSyncTransactor, error) {
	contract, err := bindHeaderSync(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HeaderSyncTransactor{contract: contract}, nil
}

// NewHeaderSyncFilterer creates a new log filterer instance of HeaderSync, bound to a specific deployed contract.
func NewHeaderSyncFilterer(address common.Address, filterer bind.ContractFilterer) (*HeaderSyncFilterer, error) {
	contract, err := bindHeaderSync(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HeaderSyncFilterer{contract: contract}, nil
}

// bindHeaderSync binds a generic wrapper to an already deployed contract.
func bindHeaderSync(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HeaderSyncABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HeaderSync *HeaderSyncRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HeaderSync.Contract.HeaderSyncCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HeaderSync *HeaderSyncRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HeaderSync.Contract.HeaderSyncTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HeaderSync *HeaderSyncRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HeaderSync.Contract.HeaderSyncTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HeaderSync *HeaderSyncCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HeaderSync.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HeaderSync *HeaderSyncTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HeaderSync.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HeaderSync *HeaderSyncTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HeaderSync.Contract.contract.Transact(opts, method, params...)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_HeaderSync *HeaderSyncTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HeaderSync.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_HeaderSync *HeaderSyncSession) Name() (*types.Transaction, error) {
	return _HeaderSync.Contract.Name(&_HeaderSync.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_HeaderSync *HeaderSyncTransactorSession) Name() (*types.Transaction, error) {
	return _HeaderSync.Contract.Name(&_HeaderSync.TransactOpts)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x72ce6700.
//
// Solidity: function syncBlockHeader(uint64 ChainID, address Address, bytes[] Headers) returns(bool success)
func (_HeaderSync *HeaderSyncTransactor) SyncBlockHeader(opts *bind.TransactOpts, ChainID uint64, Address common.Address, Headers [][]byte) (*types.Transaction, error) {
	return _HeaderSync.contract.Transact(opts, "syncBlockHeader", ChainID, Address, Headers)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x72ce6700.
//
// Solidity: function syncBlockHeader(uint64 ChainID, address Address, bytes[] Headers) returns(bool success)
func (_HeaderSync *HeaderSyncSession) SyncBlockHeader(ChainID uint64, Address common.Address, Headers [][]byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncBlockHeader(&_HeaderSync.TransactOpts, ChainID, Address, Headers)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x72ce6700.
//
// Solidity: function syncBlockHeader(uint64 ChainID, address Address, bytes[] Headers) returns(bool success)
func (_HeaderSync *HeaderSyncTransactorSession) SyncBlockHeader(ChainID uint64, Address common.Address, Headers [][]byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncBlockHeader(&_HeaderSync.TransactOpts, ChainID, Address, Headers)
}

// SyncCrossChainMsg is a paid mutator transaction binding the contract method 0x21b5cff5.
//
// Solidity: function syncCrossChainMsg(uint64 ChainID, address Address, bytes[] CrossChainMsgs) returns(bool success)
func (_HeaderSync *HeaderSyncTransactor) SyncCrossChainMsg(opts *bind.TransactOpts, ChainID uint64, Address common.Address, CrossChainMsgs [][]byte) (*types.Transaction, error) {
	return _HeaderSync.contract.Transact(opts, "syncCrossChainMsg", ChainID, Address, CrossChainMsgs)
}

// SyncCrossChainMsg is a paid mutator transaction binding the contract method 0x21b5cff5.
//
// Solidity: function syncCrossChainMsg(uint64 ChainID, address Address, bytes[] CrossChainMsgs) returns(bool success)
func (_HeaderSync *HeaderSyncSession) SyncCrossChainMsg(ChainID uint64, Address common.Address, CrossChainMsgs [][]byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncCrossChainMsg(&_HeaderSync.TransactOpts, ChainID, Address, CrossChainMsgs)
}

// SyncCrossChainMsg is a paid mutator transaction binding the contract method 0x21b5cff5.
//
// Solidity: function syncCrossChainMsg(uint64 ChainID, address Address, bytes[] CrossChainMsgs) returns(bool success)
func (_HeaderSync *HeaderSyncTransactorSession) SyncCrossChainMsg(ChainID uint64, Address common.Address, CrossChainMsgs [][]byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncCrossChainMsg(&_HeaderSync.TransactOpts, ChainID, Address, CrossChainMsgs)
}

// SyncGenesisHeader is a paid mutator transaction binding the contract method 0xb5ace618.
//
// Solidity: function syncGenesisHeader(uint64 ChainID, bytes GenesisHeader) returns(bool success)
func (_HeaderSync *HeaderSyncTransactor) SyncGenesisHeader(opts *bind.TransactOpts, ChainID uint64, GenesisHeader []byte) (*types.Transaction, error) {
	return _HeaderSync.contract.Transact(opts, "syncGenesisHeader", ChainID, GenesisHeader)
}

// SyncGenesisHeader is a paid mutator transaction binding the contract method 0xb5ace618.
//
// Solidity: function syncGenesisHeader(uint64 ChainID, bytes GenesisHeader) returns(bool success)
func (_HeaderSync *HeaderSyncSession) SyncGenesisHeader(ChainID uint64, GenesisHeader []byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncGenesisHeader(&_HeaderSync.TransactOpts, ChainID, GenesisHeader)
}

// SyncGenesisHeader is a paid mutator transaction binding the contract method 0xb5ace618.
//
// Solidity: function syncGenesisHeader(uint64 ChainID, bytes GenesisHeader) returns(bool success)
func (_HeaderSync *HeaderSyncTransactorSession) SyncGenesisHeader(ChainID uint64, GenesisHeader []byte) (*types.Transaction, error) {
	return _HeaderSync.Contract.SyncGenesisHeader(&_HeaderSync.TransactOpts, ChainID, GenesisHeader)
}

// HeaderSyncOKEpochSwitchInfoEventIterator is returned from FilterOKEpochSwitchInfoEvent and is used to iterate over the raw logs and unpacked data for OKEpochSwitchInfoEvent events raised by the HeaderSync contract.
type HeaderSyncOKEpochSwitchInfoEventIterator struct {
	Event *HeaderSyncOKEpochSwitchInfoEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HeaderSyncOKEpochSwitchInfoEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HeaderSyncOKEpochSwitchInfoEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HeaderSyncOKEpochSwitchInfoEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HeaderSyncOKEpochSwitchInfoEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HeaderSyncOKEpochSwitchInfoEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HeaderSyncOKEpochSwitchInfoEvent represents a OKEpochSwitchInfoEvent event raised by the HeaderSync contract.
type HeaderSyncOKEpochSwitchInfoEvent struct {
	ChainID            uint64
	BlockHash          string
	Height             uint64
	NextValidatorsHash string
	InfoChainID        string
	BlockHeight        uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOKEpochSwitchInfoEvent is a free log retrieval operation binding the contract event 0xbfd2d7144ec37c6f85850914f6d172957dd090c508f40d540062f1cd06c0852a.
//
// Solidity: event OKEpochSwitchInfoEvent(uint64 chainID, string BlockHash, uint64 Height, string NextValidatorsHash, string InfoChainID, uint64 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) FilterOKEpochSwitchInfoEvent(opts *bind.FilterOpts) (*HeaderSyncOKEpochSwitchInfoEventIterator, error) {

	logs, sub, err := _HeaderSync.contract.FilterLogs(opts, "OKEpochSwitchInfoEvent")
	if err != nil {
		return nil, err
	}
	return &HeaderSyncOKEpochSwitchInfoEventIterator{contract: _HeaderSync.contract, event: "OKEpochSwitchInfoEvent", logs: logs, sub: sub}, nil
}

// WatchOKEpochSwitchInfoEvent is a free log subscription operation binding the contract event 0xbfd2d7144ec37c6f85850914f6d172957dd090c508f40d540062f1cd06c0852a.
//
// Solidity: event OKEpochSwitchInfoEvent(uint64 chainID, string BlockHash, uint64 Height, string NextValidatorsHash, string InfoChainID, uint64 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) WatchOKEpochSwitchInfoEvent(opts *bind.WatchOpts, sink chan<- *HeaderSyncOKEpochSwitchInfoEvent) (event.Subscription, error) {

	logs, sub, err := _HeaderSync.contract.WatchLogs(opts, "OKEpochSwitchInfoEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HeaderSyncOKEpochSwitchInfoEvent)
				if err := _HeaderSync.contract.UnpackLog(event, "OKEpochSwitchInfoEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOKEpochSwitchInfoEvent is a log parse operation binding the contract event 0xbfd2d7144ec37c6f85850914f6d172957dd090c508f40d540062f1cd06c0852a.
//
// Solidity: event OKEpochSwitchInfoEvent(uint64 chainID, string BlockHash, uint64 Height, string NextValidatorsHash, string InfoChainID, uint64 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) ParseOKEpochSwitchInfoEvent(log types.Log) (*HeaderSyncOKEpochSwitchInfoEvent, error) {
	event := new(HeaderSyncOKEpochSwitchInfoEvent)
	if err := _HeaderSync.contract.UnpackLog(event, "OKEpochSwitchInfoEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HeaderSyncSyncHeaderIterator is returned from FilterSyncHeader and is used to iterate over the raw logs and unpacked data for SyncHeader events raised by the HeaderSync contract.
type HeaderSyncSyncHeaderIterator struct {
	Event *HeaderSyncSyncHeader // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HeaderSyncSyncHeaderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HeaderSyncSyncHeader)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HeaderSyncSyncHeader)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HeaderSyncSyncHeaderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HeaderSyncSyncHeaderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HeaderSyncSyncHeader represents a SyncHeader event raised by the HeaderSync contract.
type HeaderSyncSyncHeader struct {
	ChainID     uint64
	Height      uint64
	BlockHash   string
	BlockHeight *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSyncHeader is a free log retrieval operation binding the contract event 0xe4d5dbebcfbd7358435d9d612bc7b584bc51faf456160631d2537a0de2e1a236.
//
// Solidity: event syncHeader(uint64 chainID, uint64 height, string blockHash, uint256 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) FilterSyncHeader(opts *bind.FilterOpts) (*HeaderSyncSyncHeaderIterator, error) {

	logs, sub, err := _HeaderSync.contract.FilterLogs(opts, "syncHeader")
	if err != nil {
		return nil, err
	}
	return &HeaderSyncSyncHeaderIterator{contract: _HeaderSync.contract, event: "syncHeader", logs: logs, sub: sub}, nil
}

// WatchSyncHeader is a free log subscription operation binding the contract event 0xe4d5dbebcfbd7358435d9d612bc7b584bc51faf456160631d2537a0de2e1a236.
//
// Solidity: event syncHeader(uint64 chainID, uint64 height, string blockHash, uint256 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) WatchSyncHeader(opts *bind.WatchOpts, sink chan<- *HeaderSyncSyncHeader) (event.Subscription, error) {

	logs, sub, err := _HeaderSync.contract.WatchLogs(opts, "syncHeader")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HeaderSyncSyncHeader)
				if err := _HeaderSync.contract.UnpackLog(event, "syncHeader", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSyncHeader is a log parse operation binding the contract event 0xe4d5dbebcfbd7358435d9d612bc7b584bc51faf456160631d2537a0de2e1a236.
//
// Solidity: event syncHeader(uint64 chainID, uint64 height, string blockHash, uint256 BlockHeight)
func (_HeaderSync *HeaderSyncFilterer) ParseSyncHeader(log types.Log) (*HeaderSyncSyncHeader, error) {
	event := new(HeaderSyncSyncHeader)
	if err := _HeaderSync.contract.UnpackLog(event, "syncHeader", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
