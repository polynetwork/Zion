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

	MethodProposeCommunity = "proposeCommunity"

	MethodProposeConfig = "proposeConfig"

	MethodVoteProposal = "voteProposal"

	MethodGetCommunityProposalList = "getCommunityProposalList"

	MethodGetConfigProposalList = "getConfigProposalList"

	MethodGetProposal = "getProposal"

	MethodGetProposalList = "getProposalList"

	EventPropose = "Propose"

	EventProposeCommunity = "ProposeCommunity"

	EventProposeConfig = "ProposeConfig"

	EventVoteProposal = "VoteProposal"
)

// IProposalManagerABI is the input ABI used to generate the binding from.
const IProposalManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"stake\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"content\",\"type\":\"string\"}],\"name\":\"Propose\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"stake\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"content\",\"type\":\"string\"}],\"name\":\"ProposeCommunity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"stake\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"content\",\"type\":\"string\"}],\"name\":\"ProposeConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ID\",\"type\":\"string\"}],\"name\":\"VoteProposal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getCommunityProposalList\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigProposalList\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ID\",\"type\":\"int256\"}],\"name\":\"getProposal\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposalList\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content\",\"type\":\"bytes\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content\",\"type\":\"bytes\"}],\"name\":\"proposeCommunity\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content\",\"type\":\"bytes\"}],\"name\":\"proposeConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ID\",\"type\":\"int256\"}],\"name\":\"voteProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IProposalManagerFuncSigs maps the 4-byte function signature to its string representation.
var IProposalManagerFuncSigs = map[string]string{
	"0085b673": "getCommunityProposalList()",
	"de63d452": "getConfigProposalList()",
	"2a69c349": "getProposal(int256)",
	"346750f3": "getProposalList()",
	"37558af5": "propose(bytes)",
	"8682c1d0": "proposeCommunity(bytes)",
	"529aaa13": "proposeConfig(bytes)",
	"e3b917ca": "voteProposal(int256)",
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

// GetCommunityProposalList is a free data retrieval call binding the contract method 0x0085b673.
//
// Solidity: function getCommunityProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCaller) GetCommunityProposalList(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IProposalManager.contract.Call(opts, &out, "getCommunityProposalList")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetCommunityProposalList is a free data retrieval call binding the contract method 0x0085b673.
//
// Solidity: function getCommunityProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerSession) GetCommunityProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetCommunityProposalList(&_IProposalManager.CallOpts)
}

// GetCommunityProposalList is a free data retrieval call binding the contract method 0x0085b673.
//
// Solidity: function getCommunityProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCallerSession) GetCommunityProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetCommunityProposalList(&_IProposalManager.CallOpts)
}

// GetConfigProposalList is a free data retrieval call binding the contract method 0xde63d452.
//
// Solidity: function getConfigProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCaller) GetConfigProposalList(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IProposalManager.contract.Call(opts, &out, "getConfigProposalList")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigProposalList is a free data retrieval call binding the contract method 0xde63d452.
//
// Solidity: function getConfigProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerSession) GetConfigProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetConfigProposalList(&_IProposalManager.CallOpts)
}

// GetConfigProposalList is a free data retrieval call binding the contract method 0xde63d452.
//
// Solidity: function getConfigProposalList() view returns(bytes)
func (_IProposalManager *IProposalManagerCallerSession) GetConfigProposalList() ([]byte, error) {
	return _IProposalManager.Contract.GetConfigProposalList(&_IProposalManager.CallOpts)
}

// GetProposal is a free data retrieval call binding the contract method 0x2a69c349.
//
// Solidity: function getProposal(int256 ID) view returns(bytes)
func (_IProposalManager *IProposalManagerCaller) GetProposal(opts *bind.CallOpts, ID *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IProposalManager.contract.Call(opts, &out, "getProposal", ID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetProposal is a free data retrieval call binding the contract method 0x2a69c349.
//
// Solidity: function getProposal(int256 ID) view returns(bytes)
func (_IProposalManager *IProposalManagerSession) GetProposal(ID *big.Int) ([]byte, error) {
	return _IProposalManager.Contract.GetProposal(&_IProposalManager.CallOpts, ID)
}

// GetProposal is a free data retrieval call binding the contract method 0x2a69c349.
//
// Solidity: function getProposal(int256 ID) view returns(bytes)
func (_IProposalManager *IProposalManagerCallerSession) GetProposal(ID *big.Int) ([]byte, error) {
	return _IProposalManager.Contract.GetProposal(&_IProposalManager.CallOpts, ID)
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

// Propose is a paid mutator transaction binding the contract method 0x37558af5.
//
// Solidity: function propose(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) Propose(opts *bind.TransactOpts, content []byte) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "propose", content)
}

// Propose is a paid mutator transaction binding the contract method 0x37558af5.
//
// Solidity: function propose(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerSession) Propose(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.Propose(&_IProposalManager.TransactOpts, content)
}

// Propose is a paid mutator transaction binding the contract method 0x37558af5.
//
// Solidity: function propose(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) Propose(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.Propose(&_IProposalManager.TransactOpts, content)
}

// ProposeCommunity is a paid mutator transaction binding the contract method 0x8682c1d0.
//
// Solidity: function proposeCommunity(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) ProposeCommunity(opts *bind.TransactOpts, content []byte) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "proposeCommunity", content)
}

// ProposeCommunity is a paid mutator transaction binding the contract method 0x8682c1d0.
//
// Solidity: function proposeCommunity(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerSession) ProposeCommunity(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.ProposeCommunity(&_IProposalManager.TransactOpts, content)
}

// ProposeCommunity is a paid mutator transaction binding the contract method 0x8682c1d0.
//
// Solidity: function proposeCommunity(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) ProposeCommunity(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.ProposeCommunity(&_IProposalManager.TransactOpts, content)
}

// ProposeConfig is a paid mutator transaction binding the contract method 0x529aaa13.
//
// Solidity: function proposeConfig(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) ProposeConfig(opts *bind.TransactOpts, content []byte) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "proposeConfig", content)
}

// ProposeConfig is a paid mutator transaction binding the contract method 0x529aaa13.
//
// Solidity: function proposeConfig(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerSession) ProposeConfig(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.ProposeConfig(&_IProposalManager.TransactOpts, content)
}

// ProposeConfig is a paid mutator transaction binding the contract method 0x529aaa13.
//
// Solidity: function proposeConfig(bytes content) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) ProposeConfig(content []byte) (*types.Transaction, error) {
	return _IProposalManager.Contract.ProposeConfig(&_IProposalManager.TransactOpts, content)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xe3b917ca.
//
// Solidity: function voteProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) VoteProposal(opts *bind.TransactOpts, ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "voteProposal", ID)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xe3b917ca.
//
// Solidity: function voteProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerSession) VoteProposal(ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.VoteProposal(&_IProposalManager.TransactOpts, ID)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xe3b917ca.
//
// Solidity: function voteProposal(int256 ID) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) VoteProposal(ID *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.VoteProposal(&_IProposalManager.TransactOpts, ID)
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
	ID      string
	Caller  string
	Stake   string
	Content string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPropose is a free log retrieval operation binding the contract event 0x85dc3bd90ead16a6343614ed337d3c3e10994d4da6691ba4d8d14c78370260d6.
//
// Solidity: event Propose(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) FilterPropose(opts *bind.FilterOpts) (*IProposalManagerProposeIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "Propose")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerProposeIterator{contract: _IProposalManager.contract, event: "Propose", logs: logs, sub: sub}, nil
}

// WatchPropose is a free log subscription operation binding the contract event 0x85dc3bd90ead16a6343614ed337d3c3e10994d4da6691ba4d8d14c78370260d6.
//
// Solidity: event Propose(string ID, string caller, string stake, string content)
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

// ParsePropose is a log parse operation binding the contract event 0x85dc3bd90ead16a6343614ed337d3c3e10994d4da6691ba4d8d14c78370260d6.
//
// Solidity: event Propose(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) ParsePropose(log types.Log) (*IProposalManagerPropose, error) {
	event := new(IProposalManagerPropose)
	if err := _IProposalManager.contract.UnpackLog(event, "Propose", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IProposalManagerProposeCommunityIterator is returned from FilterProposeCommunity and is used to iterate over the raw logs and unpacked data for ProposeCommunity events raised by the IProposalManager contract.
type IProposalManagerProposeCommunityIterator struct {
	Event *IProposalManagerProposeCommunity // Event containing the contract specifics and raw log

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
func (it *IProposalManagerProposeCommunityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerProposeCommunity)
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
		it.Event = new(IProposalManagerProposeCommunity)
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
func (it *IProposalManagerProposeCommunityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerProposeCommunityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerProposeCommunity represents a ProposeCommunity event raised by the IProposalManager contract.
type IProposalManagerProposeCommunity struct {
	ID      string
	Caller  string
	Stake   string
	Content string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterProposeCommunity is a free log retrieval operation binding the contract event 0xac6d6135cfcf6030bf2699fe5e4e0fdea6ef460903bb389dede1fe67474049df.
//
// Solidity: event ProposeCommunity(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) FilterProposeCommunity(opts *bind.FilterOpts) (*IProposalManagerProposeCommunityIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "ProposeCommunity")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerProposeCommunityIterator{contract: _IProposalManager.contract, event: "ProposeCommunity", logs: logs, sub: sub}, nil
}

// WatchProposeCommunity is a free log subscription operation binding the contract event 0xac6d6135cfcf6030bf2699fe5e4e0fdea6ef460903bb389dede1fe67474049df.
//
// Solidity: event ProposeCommunity(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) WatchProposeCommunity(opts *bind.WatchOpts, sink chan<- *IProposalManagerProposeCommunity) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "ProposeCommunity")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerProposeCommunity)
				if err := _IProposalManager.contract.UnpackLog(event, "ProposeCommunity", log); err != nil {
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

// ParseProposeCommunity is a log parse operation binding the contract event 0xac6d6135cfcf6030bf2699fe5e4e0fdea6ef460903bb389dede1fe67474049df.
//
// Solidity: event ProposeCommunity(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) ParseProposeCommunity(log types.Log) (*IProposalManagerProposeCommunity, error) {
	event := new(IProposalManagerProposeCommunity)
	if err := _IProposalManager.contract.UnpackLog(event, "ProposeCommunity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IProposalManagerProposeConfigIterator is returned from FilterProposeConfig and is used to iterate over the raw logs and unpacked data for ProposeConfig events raised by the IProposalManager contract.
type IProposalManagerProposeConfigIterator struct {
	Event *IProposalManagerProposeConfig // Event containing the contract specifics and raw log

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
func (it *IProposalManagerProposeConfigIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerProposeConfig)
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
		it.Event = new(IProposalManagerProposeConfig)
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
func (it *IProposalManagerProposeConfigIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerProposeConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerProposeConfig represents a ProposeConfig event raised by the IProposalManager contract.
type IProposalManagerProposeConfig struct {
	ID      string
	Caller  string
	Stake   string
	Content string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterProposeConfig is a free log retrieval operation binding the contract event 0x74bbe87e16daa665c9fe316a858f4968102afdba7df9b3a9ffe75c705a8c7a86.
//
// Solidity: event ProposeConfig(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) FilterProposeConfig(opts *bind.FilterOpts) (*IProposalManagerProposeConfigIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "ProposeConfig")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerProposeConfigIterator{contract: _IProposalManager.contract, event: "ProposeConfig", logs: logs, sub: sub}, nil
}

// WatchProposeConfig is a free log subscription operation binding the contract event 0x74bbe87e16daa665c9fe316a858f4968102afdba7df9b3a9ffe75c705a8c7a86.
//
// Solidity: event ProposeConfig(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) WatchProposeConfig(opts *bind.WatchOpts, sink chan<- *IProposalManagerProposeConfig) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "ProposeConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerProposeConfig)
				if err := _IProposalManager.contract.UnpackLog(event, "ProposeConfig", log); err != nil {
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

// ParseProposeConfig is a log parse operation binding the contract event 0x74bbe87e16daa665c9fe316a858f4968102afdba7df9b3a9ffe75c705a8c7a86.
//
// Solidity: event ProposeConfig(string ID, string caller, string stake, string content)
func (_IProposalManager *IProposalManagerFilterer) ParseProposeConfig(log types.Log) (*IProposalManagerProposeConfig, error) {
	event := new(IProposalManagerProposeConfig)
	if err := _IProposalManager.contract.UnpackLog(event, "ProposeConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IProposalManagerVoteProposalIterator is returned from FilterVoteProposal and is used to iterate over the raw logs and unpacked data for VoteProposal events raised by the IProposalManager contract.
type IProposalManagerVoteProposalIterator struct {
	Event *IProposalManagerVoteProposal // Event containing the contract specifics and raw log

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
func (it *IProposalManagerVoteProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerVoteProposal)
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
		it.Event = new(IProposalManagerVoteProposal)
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
func (it *IProposalManagerVoteProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerVoteProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerVoteProposal represents a VoteProposal event raised by the IProposalManager contract.
type IProposalManagerVoteProposal struct {
	ID  string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterVoteProposal is a free log retrieval operation binding the contract event 0xc2dd9b8e110ee6b2faa8e933b5638ec5a62091253ebf357146d0637a9e30bafe.
//
// Solidity: event VoteProposal(string ID)
func (_IProposalManager *IProposalManagerFilterer) FilterVoteProposal(opts *bind.FilterOpts) (*IProposalManagerVoteProposalIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "VoteProposal")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerVoteProposalIterator{contract: _IProposalManager.contract, event: "VoteProposal", logs: logs, sub: sub}, nil
}

// WatchVoteProposal is a free log subscription operation binding the contract event 0xc2dd9b8e110ee6b2faa8e933b5638ec5a62091253ebf357146d0637a9e30bafe.
//
// Solidity: event VoteProposal(string ID)
func (_IProposalManager *IProposalManagerFilterer) WatchVoteProposal(opts *bind.WatchOpts, sink chan<- *IProposalManagerVoteProposal) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "VoteProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerVoteProposal)
				if err := _IProposalManager.contract.UnpackLog(event, "VoteProposal", log); err != nil {
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

// ParseVoteProposal is a log parse operation binding the contract event 0xc2dd9b8e110ee6b2faa8e933b5638ec5a62091253ebf357146d0637a9e30bafe.
//
// Solidity: event VoteProposal(string ID)
func (_IProposalManager *IProposalManagerFilterer) ParseVoteProposal(log types.Log) (*IProposalManagerVoteProposal, error) {
	event := new(IProposalManagerVoteProposal)
	if err := _IProposalManager.contract.UnpackLog(event, "VoteProposal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
