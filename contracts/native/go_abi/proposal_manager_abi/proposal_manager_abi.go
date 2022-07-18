// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proposal_manager_abi

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
	MethodPropose = "propose"

	MethodSetActiveProposal = "setActiveProposal"

	MethodVoteActiveProposal = "voteActiveProposal"

	MethodGetActiveProposal = "getActiveProposal"

	MethodGetProposalList = "getProposalList"

	EventPropose = "Propose"

	EventVoteActiveProposal = "VoteActiveProposal"
)

// IProposalManagerABI is the input ABI used to generate the binding from.
const IProposalManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"pType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"stake\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"content\",\"type\":\"string\"}],\"name\":\"Propose\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"pType\",\"type\":\"uint8\"}],\"name\":\"VoteActiveProposal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getActiveProposal\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposalList\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"pType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"content\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"stake\",\"type\":\"int256\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setActiveProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ID\",\"type\":\"int256\"}],\"name\":\"voteActiveProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IProposalManagerFuncSigs maps the 4-byte function signature to its string representation.
var IProposalManagerFuncSigs = map[string]string{
	"d7335d03": "getActiveProposal()",
	"346750f3": "getProposalList()",
	"a02a2b06": "propose(uint8,bytes,int256)",
	"cfd74457": "setActiveProposal()",
	"a8bcb673": "voteActiveProposal(int256)",
}

// IProposalManager is an auto generated Go binding around an Ethereum contract.
type IProposalManager struct {
	IProposalManagerCaller     // Read-only binding to the contract
	IProposalManagerTransactor // Write-only binding to the contract
	IProposalManagerFilterer   // Log filterer for contract events
}

// IProposalManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IProposalManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IProposalManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IProposalManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IProposalManagerSession struct {
	Contract     *IProposalManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IProposalManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IProposalManagerCallerSession struct {
	Contract *IProposalManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IProposalManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IProposalManagerTransactorSession struct {
	Contract     *IProposalManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IProposalManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IProposalManagerRaw struct {
	Contract *IProposalManager // Generic contract binding to access the raw methods on
}

// IProposalManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IProposalManagerCallerRaw struct {
	Contract *IProposalManagerCaller // Generic read-only contract binding to access the raw methods on
}

// IProposalManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IProposalManagerTransactorRaw struct {
	Contract *IProposalManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIProposalManager creates a new instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManager(address common.Address, backend bind.ContractBackend) (*IProposalManager, error) {
	contract, err := bindIProposalManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IProposalManager{IProposalManagerCaller: IProposalManagerCaller{contract: contract}, IProposalManagerTransactor: IProposalManagerTransactor{contract: contract}, IProposalManagerFilterer: IProposalManagerFilterer{contract: contract}}, nil
}

// NewIProposalManagerCaller creates a new read-only instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerCaller(address common.Address, caller bind.ContractCaller) (*IProposalManagerCaller, error) {
	contract, err := bindIProposalManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerCaller{contract: contract}, nil
}

// NewIProposalManagerTransactor creates a new write-only instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*IProposalManagerTransactor, error) {
	contract, err := bindIProposalManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerTransactor{contract: contract}, nil
}

// NewIProposalManagerFilterer creates a new log filterer instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*IProposalManagerFilterer, error) {
	contract, err := bindIProposalManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerFilterer{contract: contract}, nil
}

// bindIProposalManager binds a generic wrapper to an already deployed contract.
func bindIProposalManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IProposalManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IProposalManager *IProposalManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IProposalManager.Contract.IProposalManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IProposalManager *IProposalManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IProposalManager.Contract.IProposalManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IProposalManager *IProposalManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IProposalManager.Contract.IProposalManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IProposalManager *IProposalManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IProposalManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IProposalManager *IProposalManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IProposalManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IProposalManager *IProposalManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IProposalManager.Contract.contract.Transact(opts, method, params...)
}

// GetActiveProposal is a free data retrieval call binding the contract method 0xd7335d03.
//
// Solidity: function getActiveProposal() view returns(bytes)
func (_IProposalManager *IProposalManagerCaller) GetActiveProposal(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IProposalManager.contract.Call(opts, &out, "getActiveProposal")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetActiveProposal is a free data retrieval call binding the contract method 0xd7335d03.
//
// Solidity: function getActiveProposal() view returns(bytes)
func (_IProposalManager *IProposalManagerSession) GetActiveProposal() ([]byte, error) {
	return _IProposalManager.Contract.GetActiveProposal(&_IProposalManager.CallOpts)
}

// GetActiveProposal is a free data retrieval call binding the contract method 0xd7335d03.
//
// Solidity: function getActiveProposal() view returns(bytes)
func (_IProposalManager *IProposalManagerCallerSession) GetActiveProposal() ([]byte, error) {
	return _IProposalManager.Contract.GetActiveProposal(&_IProposalManager.CallOpts)
}

// GetProposalList is a free data retrieval call binding the contract method 0x346750f3.
//
// Solidity: function getProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCaller) GetProposalList(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IProposalManager.contract.Call(opts, &out, "getProposalList")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetProposalList is a free data retrieval call binding the contract method 0x346750f3.
//
// Solidity: function getProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerSession) GetProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetProposalList(&_IProposalManager.CallOpts)
}

// GetProposalList is a free data retrieval call binding the contract method 0x346750f3.
//
// Solidity: function getProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCallerSession) GetProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetProposalList(&_IProposalManager.CallOpts)
}

// Propose is a paid mutator transaction binding the contract method 0xa02a2b06.
//
// Solidity: function propose(uint8 pType, bytes content, int256 stake) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) Propose(opts *bind.TransactOpts, pType uint8, content []byte, stake *big.Int) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "propose", pType, content, stake)
}

// Propose is a paid mutator transaction binding the contract method 0xa02a2b06.
//
// Solidity: function propose(uint8 pType, bytes content, int256 stake) returns(bool success)
func (_IProposalManager *IProposalManagerSession) Propose(pType uint8, content []byte, stake *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.Propose(&_IProposalManager.TransactOpts, pType, content, stake)
}

// Propose is a paid mutator transaction binding the contract method 0xa02a2b06.
//
// Solidity: function propose(uint8 pType, bytes content, int256 stake) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) Propose(pType uint8, content []byte, stake *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.Propose(&_IProposalManager.TransactOpts, pType, content, stake)
}

// SetActiveProposal is a paid mutator transaction binding the contract method 0xcfd74457.
//
// Solidity: function setActiveProposal() returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) SetActiveProposal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "setActiveProposal")
}

// SetActiveProposal is a paid mutator transaction binding the contract method 0xcfd74457.
//
// Solidity: function setActiveProposal() returns(bool success)
func (_IProposalManager *IProposalManagerSession) SetActiveProposal() (*types.Transaction, error) {
	return _IProposalManager.Contract.SetActiveProposal(&_IProposalManager.TransactOpts)
}

// SetActiveProposal is a paid mutator transaction binding the contract method 0xcfd74457.
//
// Solidity: function setActiveProposal() returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) SetActiveProposal() (*types.Transaction, error) {
	return _IProposalManager.Contract.SetActiveProposal(&_IProposalManager.TransactOpts)
}

// VoteActiveProposal is a paid mutator transaction binding the contract method 0xa8bcb673.
//
// Solidity: function voteActiveProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) VoteActiveProposal(opts *bind.TransactOpts, ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "voteActiveProposal", ID)
}

// VoteActiveProposal is a paid mutator transaction binding the contract method 0xa8bcb673.
//
// Solidity: function voteActiveProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerSession) VoteActiveProposal(ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.VoteActiveProposal(&_IProposalManager.TransactOpts, ID)
}

// VoteActiveProposal is a paid mutator transaction binding the contract method 0xa8bcb673.
//
// Solidity: function voteActiveProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) VoteActiveProposal(ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.VoteActiveProposal(&_IProposalManager.TransactOpts, ID)
}

// IProposalManagerProposeIterator is returned from FilterPropose and is used to iterate over the raw logs and unpacked data for Propose events raised by the IProposalManager contract.
type IProposalManagerProposeIterator struct {
	Event *IProposalManagerPropose // Event containing the contract specifics and raw log

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
func (it *IProposalManagerProposeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerPropose)
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
		it.Event = new(IProposalManagerPropose)
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
func (it *IProposalManagerProposeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerProposeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerPropose represents a Propose event raised by the IProposalManager contract.
type IProposalManagerPropose struct {
	Caller  string
	PType   uint8
	Stake   string
	Content string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPropose is a free log retrieval operation binding the contract event 0x1eedae0efed6a9892a40aa176f002fab58d8b558e2f1330fbb125562b0e20af1.
//
// Solidity: event Propose(string caller, uint8 pType, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) FilterPropose(opts *bind.FilterOpts) (*IProposalManagerProposeIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "Propose")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerProposeIterator{contract: _IProposalManager.contract, event: "Propose", logs: logs, sub: sub}, nil
}

// WatchPropose is a free log subscription operation binding the contract event 0x1eedae0efed6a9892a40aa176f002fab58d8b558e2f1330fbb125562b0e20af1.
//
// Solidity: event Propose(string caller, uint8 pType, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) WatchPropose(opts *bind.WatchOpts, sink chan<- *IProposalManagerPropose) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "Propose")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerPropose)
				if err := _IProposalManager.contract.UnpackLog(event, "Propose", log); err != nil {
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

// ParsePropose is a log parse operation binding the contract event 0x1eedae0efed6a9892a40aa176f002fab58d8b558e2f1330fbb125562b0e20af1.
//
// Solidity: event Propose(string caller, uint8 pType, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) ParsePropose(log types.Log) (*IProposalManagerPropose, error) {
	event := new(IProposalManagerPropose)
	if err := _IProposalManager.contract.UnpackLog(event, "Propose", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IProposalManagerVoteActiveProposalIterator is returned from FilterVoteActiveProposal and is used to iterate over the raw logs and unpacked data for VoteActiveProposal events raised by the IProposalManager contract.
type IProposalManagerVoteActiveProposalIterator struct {
	Event *IProposalManagerVoteActiveProposal // Event containing the contract specifics and raw log

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
func (it *IProposalManagerVoteActiveProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerVoteActiveProposal)
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
		it.Event = new(IProposalManagerVoteActiveProposal)
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
func (it *IProposalManagerVoteActiveProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerVoteActiveProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerVoteActiveProposal represents a VoteActiveProposal event raised by the IProposalManager contract.
type IProposalManagerVoteActiveProposal struct {
	ID    string
	PType uint8
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVoteActiveProposal is a free log retrieval operation binding the contract event 0xeccf1292b1b5cc7f99be34efed92ae9b3412e330d6bb21cfadafbe8851727306.
//
// Solidity: event VoteActiveProposal(string ID, uint8 pType)
func (_IProposalManager *IProposalManagerFilterer) FilterVoteActiveProposal(opts *bind.FilterOpts) (*IProposalManagerVoteActiveProposalIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "VoteActiveProposal")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerVoteActiveProposalIterator{contract: _IProposalManager.contract, event: "VoteActiveProposal", logs: logs, sub: sub}, nil
}

// WatchVoteActiveProposal is a free log subscription operation binding the contract event 0xeccf1292b1b5cc7f99be34efed92ae9b3412e330d6bb21cfadafbe8851727306.
//
// Solidity: event VoteActiveProposal(string ID, uint8 pType)
func (_IProposalManager *IProposalManagerFilterer) WatchVoteActiveProposal(opts *bind.WatchOpts, sink chan<- *IProposalManagerVoteActiveProposal) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "VoteActiveProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerVoteActiveProposal)
				if err := _IProposalManager.contract.UnpackLog(event, "VoteActiveProposal", log); err != nil {
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

// ParseVoteActiveProposal is a log parse operation binding the contract event 0xeccf1292b1b5cc7f99be34efed92ae9b3412e330d6bb21cfadafbe8851727306.
//
// Solidity: event VoteActiveProposal(string ID, uint8 pType)
func (_IProposalManager *IProposalManagerFilterer) ParseVoteActiveProposal(log types.Log) (*IProposalManagerVoteActiveProposal, error) {
	event := new(IProposalManagerVoteActiveProposal)
	if err := _IProposalManager.contract.UnpackLog(event, "VoteActiveProposal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

