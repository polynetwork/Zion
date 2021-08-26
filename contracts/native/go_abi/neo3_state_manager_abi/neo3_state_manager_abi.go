// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package neo3_state_manager_abi

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

// Neo3StateManagerABI is the input ABI used to generate the binding from.
const Neo3StateManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"EventApproveRegisterStateValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"EventApproveRemoveStateValidator\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveRegisterStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveRemoveStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodGetCurrentStateValidator\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"Validator\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"StateValidators\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodRegisterStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"StateValidators\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodRemoveStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Neo3StateManagerFuncSigs maps the 4-byte function signature to its string representation.
var Neo3StateManagerFuncSigs = map[string]string{
	"db3a5951": "MethodApproveRegisterStateValidator(uint64,address)",
	"a27b72f9": "MethodApproveRemoveStateValidator(uint64,address)",
	"e50f8f44": "MethodContractName()",
	"495d18a4": "MethodGetCurrentStateValidator()",
	"42753232": "MethodRegisterStateValidator(string[],address)",
	"60d23093": "MethodRemoveStateValidator(string[],address)",
}

// Neo3StateManagerBin is the compiled bytecode used for deploying new contracts.
var Neo3StateManagerBin = "0x608060405234801561001057600080fd5b5061031d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80634275323214610067578063495d18a41461009257806360d2309314610067578063a27b72f9146100a1578063db3a5951146100a1578063e50f8f4414610092575b600080fd5b61007d6100753660046100cb565b600092915050565b60405190151581526020015b60405180910390f35b60606040516100899190610286565b61007d6100753660046101f7565b80356001600160a01b03811681146100c657600080fd5b919050565b60008060408084860312156100df57600080fd5b833567ffffffffffffffff808211156100f757600080fd5b818601915086601f83011261010b57600080fd5b813560208282111561011f5761011f6102d1565b8160051b61012e8282016102a0565b8381528281019086840183880185018d101561014957600080fd5b600093505b858410156101d85780358781111561016557600080fd5b8801603f81018e1361017657600080fd5b858101358881111561018a5761018a6102d1565b61019c601f8201601f191688016102a0565b8181528f8c8385010111156101b057600080fd5b818c84018983013760009181018801919091528452506001939093019291840191840161014e565b5098506101e99150508882016100af565b955050505050509250929050565b6000806040838503121561020a57600080fd5b823567ffffffffffffffff8116811461022257600080fd5b9150610230602084016100af565b90509250929050565b6000815180845260005b8181101561025f57602081850181015186830182015201610243565b81811115610271576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006102996020830184610239565b9392505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156102c9576102c96102d1565b604052919050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212208decd5be1b43984fd9d415af4f2132079b00c5e3190935fa4e7d700e3b3706d864736f6c63430008060033"

// DeployNeo3StateManager deploys a new Ethereum contract, binding an instance of Neo3StateManager to it.
func DeployNeo3StateManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Neo3StateManager, error) {
	parsed, err := abi.JSON(strings.NewReader(Neo3StateManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Neo3StateManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Neo3StateManager{Neo3StateManagerCaller: Neo3StateManagerCaller{contract: contract}, Neo3StateManagerTransactor: Neo3StateManagerTransactor{contract: contract}, Neo3StateManagerFilterer: Neo3StateManagerFilterer{contract: contract}}, nil
}

// Neo3StateManager is an auto generated Go binding around an Ethereum contract.
type Neo3StateManager struct {
	Neo3StateManagerCaller     // Read-only binding to the contract
	Neo3StateManagerTransactor // Write-only binding to the contract
	Neo3StateManagerFilterer   // Log filterer for contract events
}

// Neo3StateManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type Neo3StateManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Neo3StateManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Neo3StateManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Neo3StateManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Neo3StateManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Neo3StateManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Neo3StateManagerSession struct {
	Contract     *Neo3StateManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Neo3StateManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Neo3StateManagerCallerSession struct {
	Contract *Neo3StateManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// Neo3StateManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Neo3StateManagerTransactorSession struct {
	Contract     *Neo3StateManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// Neo3StateManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type Neo3StateManagerRaw struct {
	Contract *Neo3StateManager // Generic contract binding to access the raw methods on
}

// Neo3StateManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Neo3StateManagerCallerRaw struct {
	Contract *Neo3StateManagerCaller // Generic read-only contract binding to access the raw methods on
}

// Neo3StateManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Neo3StateManagerTransactorRaw struct {
	Contract *Neo3StateManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNeo3StateManager creates a new instance of Neo3StateManager, bound to a specific deployed contract.
func NewNeo3StateManager(address common.Address, backend bind.ContractBackend) (*Neo3StateManager, error) {
	contract, err := bindNeo3StateManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Neo3StateManager{Neo3StateManagerCaller: Neo3StateManagerCaller{contract: contract}, Neo3StateManagerTransactor: Neo3StateManagerTransactor{contract: contract}, Neo3StateManagerFilterer: Neo3StateManagerFilterer{contract: contract}}, nil
}

// NewNeo3StateManagerCaller creates a new read-only instance of Neo3StateManager, bound to a specific deployed contract.
func NewNeo3StateManagerCaller(address common.Address, caller bind.ContractCaller) (*Neo3StateManagerCaller, error) {
	contract, err := bindNeo3StateManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerCaller{contract: contract}, nil
}

// NewNeo3StateManagerTransactor creates a new write-only instance of Neo3StateManager, bound to a specific deployed contract.
func NewNeo3StateManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*Neo3StateManagerTransactor, error) {
	contract, err := bindNeo3StateManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerTransactor{contract: contract}, nil
}

// NewNeo3StateManagerFilterer creates a new log filterer instance of Neo3StateManager, bound to a specific deployed contract.
func NewNeo3StateManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*Neo3StateManagerFilterer, error) {
	contract, err := bindNeo3StateManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerFilterer{contract: contract}, nil
}

// bindNeo3StateManager binds a generic wrapper to an already deployed contract.
func bindNeo3StateManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Neo3StateManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Neo3StateManager *Neo3StateManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Neo3StateManager.Contract.Neo3StateManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Neo3StateManager *Neo3StateManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.Neo3StateManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Neo3StateManager *Neo3StateManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.Neo3StateManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Neo3StateManager *Neo3StateManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Neo3StateManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Neo3StateManager *Neo3StateManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Neo3StateManager *Neo3StateManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.contract.Transact(opts, method, params...)
}

// MethodApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xdb3a5951.
//
// Solidity: function MethodApproveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodApproveRegisterStateValidator(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodApproveRegisterStateValidator", ID, Address)
}

// MethodApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xdb3a5951.
//
// Solidity: function MethodApproveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) MethodApproveRegisterStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodApproveRegisterStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// MethodApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xdb3a5951.
//
// Solidity: function MethodApproveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodApproveRegisterStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodApproveRegisterStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// MethodApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0xa27b72f9.
//
// Solidity: function MethodApproveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodApproveRemoveStateValidator(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodApproveRemoveStateValidator", ID, Address)
}

// MethodApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0xa27b72f9.
//
// Solidity: function MethodApproveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) MethodApproveRemoveStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodApproveRemoveStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// MethodApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0xa27b72f9.
//
// Solidity: function MethodApproveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodApproveRemoveStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodApproveRemoveStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodContractName(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodContractName")
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerSession) MethodContractName() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodContractName(&_Neo3StateManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodContractName() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodContractName(&_Neo3StateManager.TransactOpts)
}

// MethodGetCurrentStateValidator is a paid mutator transaction binding the contract method 0x495d18a4.
//
// Solidity: function MethodGetCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodGetCurrentStateValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodGetCurrentStateValidator")
}

// MethodGetCurrentStateValidator is a paid mutator transaction binding the contract method 0x495d18a4.
//
// Solidity: function MethodGetCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerSession) MethodGetCurrentStateValidator() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodGetCurrentStateValidator(&_Neo3StateManager.TransactOpts)
}

// MethodGetCurrentStateValidator is a paid mutator transaction binding the contract method 0x495d18a4.
//
// Solidity: function MethodGetCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodGetCurrentStateValidator() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodGetCurrentStateValidator(&_Neo3StateManager.TransactOpts)
}

// MethodRegisterStateValidator is a paid mutator transaction binding the contract method 0x42753232.
//
// Solidity: function MethodRegisterStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodRegisterStateValidator(opts *bind.TransactOpts, StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodRegisterStateValidator", StateValidators, Address)
}

// MethodRegisterStateValidator is a paid mutator transaction binding the contract method 0x42753232.
//
// Solidity: function MethodRegisterStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) MethodRegisterStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodRegisterStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// MethodRegisterStateValidator is a paid mutator transaction binding the contract method 0x42753232.
//
// Solidity: function MethodRegisterStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodRegisterStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodRegisterStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// MethodRemoveStateValidator is a paid mutator transaction binding the contract method 0x60d23093.
//
// Solidity: function MethodRemoveStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) MethodRemoveStateValidator(opts *bind.TransactOpts, StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "MethodRemoveStateValidator", StateValidators, Address)
}

// MethodRemoveStateValidator is a paid mutator transaction binding the contract method 0x60d23093.
//
// Solidity: function MethodRemoveStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) MethodRemoveStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodRemoveStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// MethodRemoveStateValidator is a paid mutator transaction binding the contract method 0x60d23093.
//
// Solidity: function MethodRemoveStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) MethodRemoveStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.MethodRemoveStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// Neo3StateManagerEventApproveRegisterStateValidatorIterator is returned from FilterEventApproveRegisterStateValidator and is used to iterate over the raw logs and unpacked data for EventApproveRegisterStateValidator events raised by the Neo3StateManager contract.
type Neo3StateManagerEventApproveRegisterStateValidatorIterator struct {
	Event *Neo3StateManagerEventApproveRegisterStateValidator // Event containing the contract specifics and raw log

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
func (it *Neo3StateManagerEventApproveRegisterStateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Neo3StateManagerEventApproveRegisterStateValidator)
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
		it.Event = new(Neo3StateManagerEventApproveRegisterStateValidator)
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
func (it *Neo3StateManagerEventApproveRegisterStateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Neo3StateManagerEventApproveRegisterStateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Neo3StateManagerEventApproveRegisterStateValidator represents a EventApproveRegisterStateValidator event raised by the Neo3StateManager contract.
type Neo3StateManagerEventApproveRegisterStateValidator struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventApproveRegisterStateValidator is a free log retrieval operation binding the contract event 0x5f1a0b61a603a9ca51716619dd9f1da0e41ed981eaf2964dc19a807542aab580.
//
// Solidity: event EventApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) FilterEventApproveRegisterStateValidator(opts *bind.FilterOpts) (*Neo3StateManagerEventApproveRegisterStateValidatorIterator, error) {

	logs, sub, err := _Neo3StateManager.contract.FilterLogs(opts, "EventApproveRegisterStateValidator")
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerEventApproveRegisterStateValidatorIterator{contract: _Neo3StateManager.contract, event: "EventApproveRegisterStateValidator", logs: logs, sub: sub}, nil
}

// WatchEventApproveRegisterStateValidator is a free log subscription operation binding the contract event 0x5f1a0b61a603a9ca51716619dd9f1da0e41ed981eaf2964dc19a807542aab580.
//
// Solidity: event EventApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) WatchEventApproveRegisterStateValidator(opts *bind.WatchOpts, sink chan<- *Neo3StateManagerEventApproveRegisterStateValidator) (event.Subscription, error) {

	logs, sub, err := _Neo3StateManager.contract.WatchLogs(opts, "EventApproveRegisterStateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Neo3StateManagerEventApproveRegisterStateValidator)
				if err := _Neo3StateManager.contract.UnpackLog(event, "EventApproveRegisterStateValidator", log); err != nil {
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

// ParseEventApproveRegisterStateValidator is a log parse operation binding the contract event 0x5f1a0b61a603a9ca51716619dd9f1da0e41ed981eaf2964dc19a807542aab580.
//
// Solidity: event EventApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) ParseEventApproveRegisterStateValidator(log types.Log) (*Neo3StateManagerEventApproveRegisterStateValidator, error) {
	event := new(Neo3StateManagerEventApproveRegisterStateValidator)
	if err := _Neo3StateManager.contract.UnpackLog(event, "EventApproveRegisterStateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Neo3StateManagerEventApproveRemoveStateValidatorIterator is returned from FilterEventApproveRemoveStateValidator and is used to iterate over the raw logs and unpacked data for EventApproveRemoveStateValidator events raised by the Neo3StateManager contract.
type Neo3StateManagerEventApproveRemoveStateValidatorIterator struct {
	Event *Neo3StateManagerEventApproveRemoveStateValidator // Event containing the contract specifics and raw log

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
func (it *Neo3StateManagerEventApproveRemoveStateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Neo3StateManagerEventApproveRemoveStateValidator)
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
		it.Event = new(Neo3StateManagerEventApproveRemoveStateValidator)
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
func (it *Neo3StateManagerEventApproveRemoveStateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Neo3StateManagerEventApproveRemoveStateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Neo3StateManagerEventApproveRemoveStateValidator represents a EventApproveRemoveStateValidator event raised by the Neo3StateManager contract.
type Neo3StateManagerEventApproveRemoveStateValidator struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventApproveRemoveStateValidator is a free log retrieval operation binding the contract event 0xb3af0407c3807615e9840ef2aaeb42ba5174e45fce7f1793350965c6ecb13fa7.
//
// Solidity: event EventApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) FilterEventApproveRemoveStateValidator(opts *bind.FilterOpts) (*Neo3StateManagerEventApproveRemoveStateValidatorIterator, error) {

	logs, sub, err := _Neo3StateManager.contract.FilterLogs(opts, "EventApproveRemoveStateValidator")
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerEventApproveRemoveStateValidatorIterator{contract: _Neo3StateManager.contract, event: "EventApproveRemoveStateValidator", logs: logs, sub: sub}, nil
}

// WatchEventApproveRemoveStateValidator is a free log subscription operation binding the contract event 0xb3af0407c3807615e9840ef2aaeb42ba5174e45fce7f1793350965c6ecb13fa7.
//
// Solidity: event EventApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) WatchEventApproveRemoveStateValidator(opts *bind.WatchOpts, sink chan<- *Neo3StateManagerEventApproveRemoveStateValidator) (event.Subscription, error) {

	logs, sub, err := _Neo3StateManager.contract.WatchLogs(opts, "EventApproveRemoveStateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Neo3StateManagerEventApproveRemoveStateValidator)
				if err := _Neo3StateManager.contract.UnpackLog(event, "EventApproveRemoveStateValidator", log); err != nil {
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

// ParseEventApproveRemoveStateValidator is a log parse operation binding the contract event 0xb3af0407c3807615e9840ef2aaeb42ba5174e45fce7f1793350965c6ecb13fa7.
//
// Solidity: event EventApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) ParseEventApproveRemoveStateValidator(log types.Log) (*Neo3StateManagerEventApproveRemoveStateValidator, error) {
	event := new(Neo3StateManagerEventApproveRemoveStateValidator)
	if err := _Neo3StateManager.contract.UnpackLog(event, "EventApproveRemoveStateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
