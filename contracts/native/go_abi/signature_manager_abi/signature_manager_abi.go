// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package signature_manager_abi

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
	MethodAddSignature = "addSignature"

	EventAddSignatureQuorumEvent = "AddSignatureQuorumEvent"
)

// SignatureManagerABI is the input ABI used to generate the binding from.
const SignatureManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"subject\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sideChainID\",\"type\":\"uint256\"}],\"name\":\"AddSignatureQuorumEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sideChainID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"subject\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SignatureManagerFuncSigs maps the 4-byte function signature to its string representation.
var SignatureManagerFuncSigs = map[string]string{
	"29d75da9": "addSignature(address,uint256,bytes,bytes)",
}

// SignatureManagerBin is the compiled bytecode used for deploying new contracts.
var SignatureManagerBin = "0x608060405234801561001057600080fd5b506101ab806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806329d75da914610030575b600080fd5b61004461003e3660046100e9565b50505050565b005b634e487b7160e01b600052604160045260246000fd5b600082601f83011261006d57600080fd5b813567ffffffffffffffff8082111561008857610088610046565b604051601f8301601f19908116603f011681019082821181831017156100b0576100b0610046565b816040528381528660208588010111156100c957600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156100ff57600080fd5b84356001600160a01b038116811461011657600080fd5b935060208501359250604085013567ffffffffffffffff8082111561013a57600080fd5b6101468883890161005c565b9350606087013591508082111561015c57600080fd5b506101698782880161005c565b9150509295919450925056fea26469706673582212200d6198af87133340319c304c5fc322d59bdf351f86dc89e471340a954557c57364736f6c634300080e0033"

// DeploySignatureManager deploys a new Ethereum contract, binding an instance of SignatureManager to it.
func DeploySignatureManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SignatureManager, error) {
	parsed, err := abi.JSON(strings.NewReader(SignatureManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SignatureManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignatureManager{SignatureManagerCaller: SignatureManagerCaller{contract: contract}, SignatureManagerTransactor: SignatureManagerTransactor{contract: contract}, SignatureManagerFilterer: SignatureManagerFilterer{contract: contract}}, nil
}

// SignatureManager is an auto generated Go binding around an Ethereum contract.
type SignatureManager struct {
	SignatureManagerCaller     // Read-only binding to the contract
	SignatureManagerTransactor // Write-only binding to the contract
	SignatureManagerFilterer   // Log filterer for contract events
}

// SignatureManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignatureManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignatureManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignatureManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignatureManagerSession struct {
	Contract     *SignatureManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignatureManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignatureManagerCallerSession struct {
	Contract *SignatureManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SignatureManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignatureManagerTransactorSession struct {
	Contract     *SignatureManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SignatureManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignatureManagerRaw struct {
	Contract *SignatureManager // Generic contract binding to access the raw methods on
}

// SignatureManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignatureManagerCallerRaw struct {
	Contract *SignatureManagerCaller // Generic read-only contract binding to access the raw methods on
}

// SignatureManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignatureManagerTransactorRaw struct {
	Contract *SignatureManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignatureManager creates a new instance of SignatureManager, bound to a specific deployed contract.
func NewSignatureManager(address common.Address, backend bind.ContractBackend) (*SignatureManager, error) {
	contract, err := bindSignatureManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignatureManager{SignatureManagerCaller: SignatureManagerCaller{contract: contract}, SignatureManagerTransactor: SignatureManagerTransactor{contract: contract}, SignatureManagerFilterer: SignatureManagerFilterer{contract: contract}}, nil
}

// NewSignatureManagerCaller creates a new read-only instance of SignatureManager, bound to a specific deployed contract.
func NewSignatureManagerCaller(address common.Address, caller bind.ContractCaller) (*SignatureManagerCaller, error) {
	contract, err := bindSignatureManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureManagerCaller{contract: contract}, nil
}

// NewSignatureManagerTransactor creates a new write-only instance of SignatureManager, bound to a specific deployed contract.
func NewSignatureManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*SignatureManagerTransactor, error) {
	contract, err := bindSignatureManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureManagerTransactor{contract: contract}, nil
}

// NewSignatureManagerFilterer creates a new log filterer instance of SignatureManager, bound to a specific deployed contract.
func NewSignatureManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*SignatureManagerFilterer, error) {
	contract, err := bindSignatureManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignatureManagerFilterer{contract: contract}, nil
}

// bindSignatureManager binds a generic wrapper to an already deployed contract.
func bindSignatureManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignatureManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureManager *SignatureManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureManager.Contract.SignatureManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureManager *SignatureManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureManager.Contract.SignatureManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureManager *SignatureManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureManager.Contract.SignatureManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureManager *SignatureManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureManager *SignatureManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureManager *SignatureManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureManager.Contract.contract.Transact(opts, method, params...)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns()
func (_SignatureManager *SignatureManagerTransactor) AddSignature(opts *bind.TransactOpts, addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _SignatureManager.contract.Transact(opts, "addSignature", addr, sideChainID, subject, signature)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns()
func (_SignatureManager *SignatureManagerSession) AddSignature(addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _SignatureManager.Contract.AddSignature(&_SignatureManager.TransactOpts, addr, sideChainID, subject, signature)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns()
func (_SignatureManager *SignatureManagerTransactorSession) AddSignature(addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _SignatureManager.Contract.AddSignature(&_SignatureManager.TransactOpts, addr, sideChainID, subject, signature)
}

// SignatureManagerAddSignatureQuorumEventIterator is returned from FilterAddSignatureQuorumEvent and is used to iterate over the raw logs and unpacked data for AddSignatureQuorumEvent events raised by the SignatureManager contract.
type SignatureManagerAddSignatureQuorumEventIterator struct {
	Event *SignatureManagerAddSignatureQuorumEvent // Event containing the contract specifics and raw log

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
func (it *SignatureManagerAddSignatureQuorumEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignatureManagerAddSignatureQuorumEvent)
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
		it.Event = new(SignatureManagerAddSignatureQuorumEvent)
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
func (it *SignatureManagerAddSignatureQuorumEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignatureManagerAddSignatureQuorumEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignatureManagerAddSignatureQuorumEvent represents a AddSignatureQuorumEvent event raised by the SignatureManager contract.
type SignatureManagerAddSignatureQuorumEvent struct {
	Id         []byte
	Subject    []byte
	SideChainID *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddSignatureQuorumEvent is a free log retrieval operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_SignatureManager *SignatureManagerFilterer) FilterAddSignatureQuorumEvent(opts *bind.FilterOpts) (*SignatureManagerAddSignatureQuorumEventIterator, error) {

	logs, sub, err := _SignatureManager.contract.FilterLogs(opts, "AddSignatureQuorumEvent")
	if err != nil {
		return nil, err
	}
	return &SignatureManagerAddSignatureQuorumEventIterator{contract: _SignatureManager.contract, event: "AddSignatureQuorumEvent", logs: logs, sub: sub}, nil
}

// WatchAddSignatureQuorumEvent is a free log subscription operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_SignatureManager *SignatureManagerFilterer) WatchAddSignatureQuorumEvent(opts *bind.WatchOpts, sink chan<- *SignatureManagerAddSignatureQuorumEvent) (event.Subscription, error) {

	logs, sub, err := _SignatureManager.contract.WatchLogs(opts, "AddSignatureQuorumEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignatureManagerAddSignatureQuorumEvent)
				if err := _SignatureManager.contract.UnpackLog(event, "AddSignatureQuorumEvent", log); err != nil {
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

// ParseAddSignatureQuorumEvent is a log parse operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_SignatureManager *SignatureManagerFilterer) ParseAddSignatureQuorumEvent(log types.Log) (*SignatureManagerAddSignatureQuorumEvent, error) {
	event := new(SignatureManagerAddSignatureQuorumEvent)
	if err := _SignatureManager.contract.UnpackLog(event, "AddSignatureQuorumEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

