// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegate_abi

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
	MethodApprove = "approve"

	MethodAllowance = "allowance"

	EventApproval = "Approval"
)

// IDelegateABI is the input ABI used to generate the binding from.
const IDelegateABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IDelegateFuncSigs maps the 4-byte function signature to its string representation.
var IDelegateFuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
}

// IDelegate is an auto generated Go binding around an Ethereum contract.
type IDelegate struct {
	IDelegateCaller     // Read-only binding to the contract
	IDelegateTransactor // Write-only binding to the contract
	IDelegateFilterer   // Log filterer for contract events
}

// IDelegateCaller is an auto generated read-only Go binding around an Ethereum contract.
type IDelegateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDelegateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IDelegateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDelegateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IDelegateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDelegateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IDelegateSession struct {
	Contract     *IDelegate        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IDelegateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IDelegateCallerSession struct {
	Contract *IDelegateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IDelegateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IDelegateTransactorSession struct {
	Contract     *IDelegateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IDelegateRaw is an auto generated low-level Go binding around an Ethereum contract.
type IDelegateRaw struct {
	Contract *IDelegate // Generic contract binding to access the raw methods on
}

// IDelegateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IDelegateCallerRaw struct {
	Contract *IDelegateCaller // Generic read-only contract binding to access the raw methods on
}

// IDelegateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IDelegateTransactorRaw struct {
	Contract *IDelegateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIDelegate creates a new instance of IDelegate, bound to a specific deployed contract.
func NewIDelegate(address common.Address, backend bind.ContractBackend) (*IDelegate, error) {
	contract, err := bindIDelegate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IDelegate{IDelegateCaller: IDelegateCaller{contract: contract}, IDelegateTransactor: IDelegateTransactor{contract: contract}, IDelegateFilterer: IDelegateFilterer{contract: contract}}, nil
}

// NewIDelegateCaller creates a new read-only instance of IDelegate, bound to a specific deployed contract.
func NewIDelegateCaller(address common.Address, caller bind.ContractCaller) (*IDelegateCaller, error) {
	contract, err := bindIDelegate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IDelegateCaller{contract: contract}, nil
}

// NewIDelegateTransactor creates a new write-only instance of IDelegate, bound to a specific deployed contract.
func NewIDelegateTransactor(address common.Address, transactor bind.ContractTransactor) (*IDelegateTransactor, error) {
	contract, err := bindIDelegate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IDelegateTransactor{contract: contract}, nil
}

// NewIDelegateFilterer creates a new log filterer instance of IDelegate, bound to a specific deployed contract.
func NewIDelegateFilterer(address common.Address, filterer bind.ContractFilterer) (*IDelegateFilterer, error) {
	contract, err := bindIDelegate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IDelegateFilterer{contract: contract}, nil
}

// bindIDelegate binds a generic wrapper to an already deployed contract.
func bindIDelegate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IDelegateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDelegate *IDelegateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDelegate.Contract.IDelegateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDelegate *IDelegateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDelegate.Contract.IDelegateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDelegate *IDelegateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDelegate.Contract.IDelegateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDelegate *IDelegateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDelegate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDelegate *IDelegateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDelegate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDelegate *IDelegateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDelegate.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IDelegate *IDelegateCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IDelegate.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IDelegate *IDelegateSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IDelegate.Contract.Allowance(&_IDelegate.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IDelegate *IDelegateCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IDelegate.Contract.Allowance(&_IDelegate.CallOpts, owner, spender)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IDelegate *IDelegateTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDelegate.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IDelegate *IDelegateSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDelegate.Contract.Approve(&_IDelegate.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IDelegate *IDelegateTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDelegate.Contract.Approve(&_IDelegate.TransactOpts, spender, amount)
}

// IDelegateApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IDelegate contract.
type IDelegateApprovalIterator struct {
	Event *IDelegateApproval // Event containing the contract specifics and raw log

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
func (it *IDelegateApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IDelegateApproval)
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
		it.Event = new(IDelegateApproval)
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
func (it *IDelegateApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IDelegateApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IDelegateApproval represents a Approval event raised by the IDelegate contract.
type IDelegateApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IDelegate *IDelegateFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IDelegateApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IDelegate.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IDelegateApprovalIterator{contract: _IDelegate.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IDelegate *IDelegateFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IDelegateApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IDelegate.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IDelegateApproval)
				if err := _IDelegate.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IDelegate *IDelegateFilterer) ParseApproval(log types.Log) (*IDelegateApproval, error) {
	event := new(IDelegateApproval)
	if err := _IDelegate.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

