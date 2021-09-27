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

var (
	MethodApproveRegisterStateValidator = "approveRegisterStateValidator"

	MethodApproveRemoveStateValidator = "approveRemoveStateValidator"

	MethodGetCurrentStateValidator = "getCurrentStateValidator"

	MethodName = "name"

	MethodRegisterStateValidator = "registerStateValidator"

	MethodRemoveStateValidator = "removeStateValidator"
)

// Neo3StateManagerABI is the input ABI used to generate the binding from.
const Neo3StateManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"evtApproveRegisterStateValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"evtApproveRemoveStateValidator\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveRegisterStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveRemoveStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentStateValidator\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"Validator\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"StateValidators\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"registerStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"StateValidators\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"removeStateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Neo3StateManagerFuncSigs maps the 4-byte function signature to its string representation.
var Neo3StateManagerFuncSigs = map[string]string{
	"ca1c4d1b": "approveRegisterStateValidator(uint64,address)",
	"3473fd55": "approveRemoveStateValidator(uint64,address)",
	"770fa9ad": "getCurrentStateValidator()",
	"06fdde03": "name()",
	"f7531edd": "registerStateValidator(string[],address)",
	"d62c2f61": "removeStateValidator(string[],address)",
}

// Neo3StateManagerBin is the compiled bytecode used for deploying new contracts.
var Neo3StateManagerBin = "0x608060405234801561001057600080fd5b50610321806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806306fdde03146100675780633473fd551461007f578063770fa9ad14610067578063ca1c4d1b1461007f578063d62c2f61146100a5578063f7531edd146100a5575b600080fd5b6060604051610076919061028a565b60405180910390f35b61009561008d3660046101fb565b600092915050565b6040519015158152602001610076565b61009561008d3660046100cf565b80356001600160a01b03811681146100ca57600080fd5b919050565b60008060408084860312156100e357600080fd5b833567ffffffffffffffff808211156100fb57600080fd5b818601915086601f83011261010f57600080fd5b8135602082821115610123576101236102d5565b8160051b6101328282016102a4565b8381528281019086840183880185018d101561014d57600080fd5b600093505b858410156101dc5780358781111561016957600080fd5b8801603f81018e1361017a57600080fd5b858101358881111561018e5761018e6102d5565b6101a0601f8201601f191688016102a4565b8181528f8c8385010111156101b457600080fd5b818c840189830137600091810188019190915284525060019390930192918401918401610152565b5098506101ed9150508882016100b3565b955050505050509250929050565b6000806040838503121561020e57600080fd5b823567ffffffffffffffff8116811461022657600080fd5b9150610234602084016100b3565b90509250929050565b6000815180845260005b8181101561026357602081850181015186830182015201610247565b81811115610275576000602083870101525b50601f01601f19169290920160200192915050565b60208152600061029d602083018461023d565b9392505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156102cd576102cd6102d5565b604052919050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212202c6e8dc15980e12ec1b93cc0fd811e701b6b90ebd3fd4b8985934601617ca04e64736f6c63430008060033"

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

// ApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xca1c4d1b.
//
// Solidity: function approveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) ApproveRegisterStateValidator(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "approveRegisterStateValidator", ID, Address)
}

// ApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xca1c4d1b.
//
// Solidity: function approveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) ApproveRegisterStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.ApproveRegisterStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// ApproveRegisterStateValidator is a paid mutator transaction binding the contract method 0xca1c4d1b.
//
// Solidity: function approveRegisterStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) ApproveRegisterStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.ApproveRegisterStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// ApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0x3473fd55.
//
// Solidity: function approveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) ApproveRemoveStateValidator(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "approveRemoveStateValidator", ID, Address)
}

// ApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0x3473fd55.
//
// Solidity: function approveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) ApproveRemoveStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.ApproveRemoveStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// ApproveRemoveStateValidator is a paid mutator transaction binding the contract method 0x3473fd55.
//
// Solidity: function approveRemoveStateValidator(uint64 ID, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) ApproveRemoveStateValidator(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.ApproveRemoveStateValidator(&_Neo3StateManager.TransactOpts, ID, Address)
}

// GetCurrentStateValidator is a paid mutator transaction binding the contract method 0x770fa9ad.
//
// Solidity: function getCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerTransactor) GetCurrentStateValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "getCurrentStateValidator")
}

// GetCurrentStateValidator is a paid mutator transaction binding the contract method 0x770fa9ad.
//
// Solidity: function getCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerSession) GetCurrentStateValidator() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.GetCurrentStateValidator(&_Neo3StateManager.TransactOpts)
}

// GetCurrentStateValidator is a paid mutator transaction binding the contract method 0x770fa9ad.
//
// Solidity: function getCurrentStateValidator() returns(bytes Validator)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) GetCurrentStateValidator() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.GetCurrentStateValidator(&_Neo3StateManager.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerSession) Name() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.Name(&_Neo3StateManager.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) Name() (*types.Transaction, error) {
	return _Neo3StateManager.Contract.Name(&_Neo3StateManager.TransactOpts)
}

// RegisterStateValidator is a paid mutator transaction binding the contract method 0xf7531edd.
//
// Solidity: function registerStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) RegisterStateValidator(opts *bind.TransactOpts, StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "registerStateValidator", StateValidators, Address)
}

// RegisterStateValidator is a paid mutator transaction binding the contract method 0xf7531edd.
//
// Solidity: function registerStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) RegisterStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.RegisterStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// RegisterStateValidator is a paid mutator transaction binding the contract method 0xf7531edd.
//
// Solidity: function registerStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) RegisterStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.RegisterStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// RemoveStateValidator is a paid mutator transaction binding the contract method 0xd62c2f61.
//
// Solidity: function removeStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactor) RemoveStateValidator(opts *bind.TransactOpts, StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.contract.Transact(opts, "removeStateValidator", StateValidators, Address)
}

// RemoveStateValidator is a paid mutator transaction binding the contract method 0xd62c2f61.
//
// Solidity: function removeStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerSession) RemoveStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.RemoveStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// RemoveStateValidator is a paid mutator transaction binding the contract method 0xd62c2f61.
//
// Solidity: function removeStateValidator(string[] StateValidators, address Address) returns(bool success)
func (_Neo3StateManager *Neo3StateManagerTransactorSession) RemoveStateValidator(StateValidators []string, Address common.Address) (*types.Transaction, error) {
	return _Neo3StateManager.Contract.RemoveStateValidator(&_Neo3StateManager.TransactOpts, StateValidators, Address)
}

// Neo3StateManagerApproveRegisterStateValidatorIterator is returned from FilterApproveRegisterStateValidator and is used to iterate over the raw logs and unpacked data for ApproveRegisterStateValidator events raised by the Neo3StateManager contract.
type Neo3StateManagerApproveRegisterStateValidatorIterator struct {
	Event *Neo3StateManagerApproveRegisterStateValidator // Event containing the contract specifics and raw log

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
func (it *Neo3StateManagerApproveRegisterStateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Neo3StateManagerApproveRegisterStateValidator)
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
		it.Event = new(Neo3StateManagerApproveRegisterStateValidator)
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
func (it *Neo3StateManagerApproveRegisterStateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Neo3StateManagerApproveRegisterStateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Neo3StateManagerApproveRegisterStateValidator represents a ApproveRegisterStateValidator event raised by the Neo3StateManager contract.
type Neo3StateManagerApproveRegisterStateValidator struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproveRegisterStateValidator is a free log retrieval operation binding the contract event 0xa2d1b53f6d5ea7964c14d81acde072ef679b771780c18a93206cd3319c6621d6.
//
// Solidity: event evtApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) FilterApproveRegisterStateValidator(opts *bind.FilterOpts) (*Neo3StateManagerApproveRegisterStateValidatorIterator, error) {

	logs, sub, err := _Neo3StateManager.contract.FilterLogs(opts, "evtApproveRegisterStateValidator")
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerApproveRegisterStateValidatorIterator{contract: _Neo3StateManager.contract, event: "evtApproveRegisterStateValidator", logs: logs, sub: sub}, nil
}

// WatchApproveRegisterStateValidator is a free log subscription operation binding the contract event 0xa2d1b53f6d5ea7964c14d81acde072ef679b771780c18a93206cd3319c6621d6.
//
// Solidity: event evtApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) WatchApproveRegisterStateValidator(opts *bind.WatchOpts, sink chan<- *Neo3StateManagerApproveRegisterStateValidator) (event.Subscription, error) {

	logs, sub, err := _Neo3StateManager.contract.WatchLogs(opts, "evtApproveRegisterStateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Neo3StateManagerApproveRegisterStateValidator)
				if err := _Neo3StateManager.contract.UnpackLog(event, "evtApproveRegisterStateValidator", log); err != nil {
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

// ParseApproveRegisterStateValidator is a log parse operation binding the contract event 0xa2d1b53f6d5ea7964c14d81acde072ef679b771780c18a93206cd3319c6621d6.
//
// Solidity: event evtApproveRegisterStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) ParseApproveRegisterStateValidator(log types.Log) (*Neo3StateManagerApproveRegisterStateValidator, error) {
	event := new(Neo3StateManagerApproveRegisterStateValidator)
	if err := _Neo3StateManager.contract.UnpackLog(event, "evtApproveRegisterStateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Neo3StateManagerApproveRemoveStateValidatorIterator is returned from FilterApproveRemoveStateValidator and is used to iterate over the raw logs and unpacked data for ApproveRemoveStateValidator events raised by the Neo3StateManager contract.
type Neo3StateManagerApproveRemoveStateValidatorIterator struct {
	Event *Neo3StateManagerApproveRemoveStateValidator // Event containing the contract specifics and raw log

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
func (it *Neo3StateManagerApproveRemoveStateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Neo3StateManagerApproveRemoveStateValidator)
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
		it.Event = new(Neo3StateManagerApproveRemoveStateValidator)
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
func (it *Neo3StateManagerApproveRemoveStateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Neo3StateManagerApproveRemoveStateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Neo3StateManagerApproveRemoveStateValidator represents a ApproveRemoveStateValidator event raised by the Neo3StateManager contract.
type Neo3StateManagerApproveRemoveStateValidator struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproveRemoveStateValidator is a free log retrieval operation binding the contract event 0x32050489af2e3d6571e4c66ba0ca36ccf2b8feed7a7059d23aab90998d263604.
//
// Solidity: event evtApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) FilterApproveRemoveStateValidator(opts *bind.FilterOpts) (*Neo3StateManagerApproveRemoveStateValidatorIterator, error) {

	logs, sub, err := _Neo3StateManager.contract.FilterLogs(opts, "evtApproveRemoveStateValidator")
	if err != nil {
		return nil, err
	}
	return &Neo3StateManagerApproveRemoveStateValidatorIterator{contract: _Neo3StateManager.contract, event: "evtApproveRemoveStateValidator", logs: logs, sub: sub}, nil
}

// WatchApproveRemoveStateValidator is a free log subscription operation binding the contract event 0x32050489af2e3d6571e4c66ba0ca36ccf2b8feed7a7059d23aab90998d263604.
//
// Solidity: event evtApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) WatchApproveRemoveStateValidator(opts *bind.WatchOpts, sink chan<- *Neo3StateManagerApproveRemoveStateValidator) (event.Subscription, error) {

	logs, sub, err := _Neo3StateManager.contract.WatchLogs(opts, "evtApproveRemoveStateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Neo3StateManagerApproveRemoveStateValidator)
				if err := _Neo3StateManager.contract.UnpackLog(event, "evtApproveRemoveStateValidator", log); err != nil {
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

// ParseApproveRemoveStateValidator is a log parse operation binding the contract event 0x32050489af2e3d6571e4c66ba0ca36ccf2b8feed7a7059d23aab90998d263604.
//
// Solidity: event evtApproveRemoveStateValidator(uint64 ID)
func (_Neo3StateManager *Neo3StateManagerFilterer) ParseApproveRemoveStateValidator(log types.Log) (*Neo3StateManagerApproveRemoveStateValidator, error) {
	event := new(Neo3StateManagerApproveRemoveStateValidator)
	if err := _Neo3StateManager.contract.UnpackLog(event, "evtApproveRemoveStateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
