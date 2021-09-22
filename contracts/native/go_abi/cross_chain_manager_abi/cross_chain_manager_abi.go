// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cross_chain_manager_abi

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

// CrossChainManagerABI is the input ABI used to generate the binding from.
const CrossChainManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"merkleValueHex\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"NOTIFY_MAKE_PROOF_EVENT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"TxHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"MultiSign\",\"type\":\"bytes\"}],\"name\":\"btcTxMultiSignEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"FromChainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"buf\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"FromTxHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"RedeemKey\",\"type\":\"string\"}],\"name\":\"btcTxToRelayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rk\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"buf\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64[]\",\"name\":\"amts\",\"type\":\"uint64[]\"}],\"name\":\"makeBtcTxEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"}],\"name\":\"MethodBlackChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"SourceChainID\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"Height\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"Proof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"RelayerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"Extra\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"HeaderOrCrossChainMsg\",\"type\":\"bytes\"}],\"name\":\"MethodImportOuterTransfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"RedeemKey\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"TxHash\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"Address\",\"type\":\"string\"},{\"internalType\":\"bytes[]\",\"name\":\"Signs\",\"type\":\"bytes[]\"}],\"name\":\"MethodMultiSign\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ChainID\",\"type\":\"uint64\"}],\"name\":\"MethodWhiteChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// CrossChainManagerFuncSigs maps the 4-byte function signature to its string representation.
var CrossChainManagerFuncSigs = map[string]string{
	"0d2ff800": "MethodBlackChain(uint64)",
	"e50f8f44": "MethodContractName()",
	"9921164e": "MethodImportOuterTransfer(uint64,uint32,bytes,bytes,bytes,bytes)",
	"4dbb371c": "MethodMultiSign(uint64,string,bytes,string,bytes[])",
	"2ae0b4bd": "MethodWhiteChain(uint64)",
}

// CrossChainManagerBin is the compiled bytecode used for deploying new contracts.
var CrossChainManagerBin = "0x608060405234801561001057600080fd5b50610471806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630d2ff8001461005c5780632ae0b4bd1461005c5780634dbb371c146100855780639921164e1461009e578063e50f8f44146100b8575b600080fd5b61007061006a366004610168565b50600090565b60405190151581526020015b60405180910390f35b61007061009336600461018a565b600095945050505050565b6100706100ac3660046102d2565b60009695505050505050565b606060405161007c919061039f565b600082601f8301126100d857600080fd5b813567ffffffffffffffff8111156100f2576100f2610425565b610105601f8201601f19166020016103f4565b81815284602083860101111561011a57600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff8116811461014b57600080fd5b919050565b803567ffffffffffffffff8116811461014b57600080fd5b60006020828403121561017a57600080fd5b61018382610150565b9392505050565b600080600080600060a086880312156101a257600080fd5b6101ab86610150565b945060208087013567ffffffffffffffff808211156101c957600080fd5b6101d58a838b016100c7565b965060408901359150808211156101eb57600080fd5b6101f78a838b016100c7565b9550606089013591508082111561020d57600080fd5b6102198a838b016100c7565b9450608089013591508082111561022f57600080fd5b818901915089601f83011261024357600080fd5b81358181111561025557610255610425565b8060051b6102648582016103f4565b8281528581019085870183870188018f101561027f57600080fd5b600093505b848410156102bd57858135111561029a57600080fd5b6102a98f8983358a01016100c7565b835260019390930192918701918701610284565b50809750505050505050509295509295909350565b60008060008060008060c087890312156102eb57600080fd5b6102f487610150565b955061030260208801610137565b9450604087013567ffffffffffffffff8082111561031f57600080fd5b61032b8a838b016100c7565b9550606089013591508082111561034157600080fd5b61034d8a838b016100c7565b9450608089013591508082111561036357600080fd5b61036f8a838b016100c7565b935060a089013591508082111561038557600080fd5b5061039289828a016100c7565b9150509295509295509295565b600060208083528351808285015260005b818110156103cc578581018301518582016040015282016103b0565b818111156103de576000604083870101525b50601f01601f1916929092016040019392505050565b604051601f8201601f1916810167ffffffffffffffff8111828210171561041d5761041d610425565b604052919050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212201a059f44d01f274f24f461ee3c212b8be9d81765a8258a100dc0f0b31ae8579164736f6c63430008060033"

// DeployCrossChainManager deploys a new Ethereum contract, binding an instance of CrossChainManager to it.
func DeployCrossChainManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainManager, error) {
	parsed, err := abi.JSON(strings.NewReader(CrossChainManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CrossChainManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainManager{CrossChainManagerCaller: CrossChainManagerCaller{contract: contract}, CrossChainManagerTransactor: CrossChainManagerTransactor{contract: contract}, CrossChainManagerFilterer: CrossChainManagerFilterer{contract: contract}}, nil
}

// CrossChainManager is an auto generated Go binding around an Ethereum contract.
type CrossChainManager struct {
	CrossChainManagerCaller     // Read-only binding to the contract
	CrossChainManagerTransactor // Write-only binding to the contract
	CrossChainManagerFilterer   // Log filterer for contract events
}

// CrossChainManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainManagerSession struct {
	Contract     *CrossChainManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CrossChainManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainManagerCallerSession struct {
	Contract *CrossChainManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// CrossChainManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainManagerTransactorSession struct {
	Contract     *CrossChainManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// CrossChainManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainManagerRaw struct {
	Contract *CrossChainManager // Generic contract binding to access the raw methods on
}

// CrossChainManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainManagerCallerRaw struct {
	Contract *CrossChainManagerCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainManagerTransactorRaw struct {
	Contract *CrossChainManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainManager creates a new instance of CrossChainManager, bound to a specific deployed contract.
func NewCrossChainManager(address common.Address, backend bind.ContractBackend) (*CrossChainManager, error) {
	contract, err := bindCrossChainManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainManager{CrossChainManagerCaller: CrossChainManagerCaller{contract: contract}, CrossChainManagerTransactor: CrossChainManagerTransactor{contract: contract}, CrossChainManagerFilterer: CrossChainManagerFilterer{contract: contract}}, nil
}

// NewCrossChainManagerCaller creates a new read-only instance of CrossChainManager, bound to a specific deployed contract.
func NewCrossChainManagerCaller(address common.Address, caller bind.ContractCaller) (*CrossChainManagerCaller, error) {
	contract, err := bindCrossChainManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerCaller{contract: contract}, nil
}

// NewCrossChainManagerTransactor creates a new write-only instance of CrossChainManager, bound to a specific deployed contract.
func NewCrossChainManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainManagerTransactor, error) {
	contract, err := bindCrossChainManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerTransactor{contract: contract}, nil
}

// NewCrossChainManagerFilterer creates a new log filterer instance of CrossChainManager, bound to a specific deployed contract.
func NewCrossChainManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainManagerFilterer, error) {
	contract, err := bindCrossChainManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerFilterer{contract: contract}, nil
}

// bindCrossChainManager binds a generic wrapper to an already deployed contract.
func bindCrossChainManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CrossChainManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainManager *CrossChainManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainManager.Contract.CrossChainManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainManager *CrossChainManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainManager.Contract.CrossChainManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainManager *CrossChainManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainManager.Contract.CrossChainManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainManager *CrossChainManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainManager *CrossChainManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainManager *CrossChainManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainManager.Contract.contract.Transact(opts, method, params...)
}

// MethodBlackChain is a paid mutator transaction binding the contract method 0x0d2ff800.
//
// Solidity: function MethodBlackChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactor) MethodBlackChain(opts *bind.TransactOpts, ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.contract.Transact(opts, "MethodBlackChain", ChainID)
}

// MethodBlackChain is a paid mutator transaction binding the contract method 0x0d2ff800.
//
// Solidity: function MethodBlackChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerSession) MethodBlackChain(ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodBlackChain(&_CrossChainManager.TransactOpts, ChainID)
}

// MethodBlackChain is a paid mutator transaction binding the contract method 0x0d2ff800.
//
// Solidity: function MethodBlackChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactorSession) MethodBlackChain(ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodBlackChain(&_CrossChainManager.TransactOpts, ChainID)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_CrossChainManager *CrossChainManagerTransactor) MethodContractName(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainManager.contract.Transact(opts, "MethodContractName")
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_CrossChainManager *CrossChainManagerSession) MethodContractName() (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodContractName(&_CrossChainManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_CrossChainManager *CrossChainManagerTransactorSession) MethodContractName() (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodContractName(&_CrossChainManager.TransactOpts)
}

// MethodImportOuterTransfer is a paid mutator transaction binding the contract method 0x9921164e.
//
// Solidity: function MethodImportOuterTransfer(uint64 SourceChainID, uint32 Height, bytes Proof, bytes RelayerAddress, bytes Extra, bytes HeaderOrCrossChainMsg) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactor) MethodImportOuterTransfer(opts *bind.TransactOpts, SourceChainID uint64, Height uint32, Proof []byte, RelayerAddress []byte, Extra []byte, HeaderOrCrossChainMsg []byte) (*types.Transaction, error) {
	return _CrossChainManager.contract.Transact(opts, "MethodImportOuterTransfer", SourceChainID, Height, Proof, RelayerAddress, Extra, HeaderOrCrossChainMsg)
}

// MethodImportOuterTransfer is a paid mutator transaction binding the contract method 0x9921164e.
//
// Solidity: function MethodImportOuterTransfer(uint64 SourceChainID, uint32 Height, bytes Proof, bytes RelayerAddress, bytes Extra, bytes HeaderOrCrossChainMsg) returns(bool success)
func (_CrossChainManager *CrossChainManagerSession) MethodImportOuterTransfer(SourceChainID uint64, Height uint32, Proof []byte, RelayerAddress []byte, Extra []byte, HeaderOrCrossChainMsg []byte) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodImportOuterTransfer(&_CrossChainManager.TransactOpts, SourceChainID, Height, Proof, RelayerAddress, Extra, HeaderOrCrossChainMsg)
}

// MethodImportOuterTransfer is a paid mutator transaction binding the contract method 0x9921164e.
//
// Solidity: function MethodImportOuterTransfer(uint64 SourceChainID, uint32 Height, bytes Proof, bytes RelayerAddress, bytes Extra, bytes HeaderOrCrossChainMsg) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactorSession) MethodImportOuterTransfer(SourceChainID uint64, Height uint32, Proof []byte, RelayerAddress []byte, Extra []byte, HeaderOrCrossChainMsg []byte) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodImportOuterTransfer(&_CrossChainManager.TransactOpts, SourceChainID, Height, Proof, RelayerAddress, Extra, HeaderOrCrossChainMsg)
}

// MethodMultiSign is a paid mutator transaction binding the contract method 0x4dbb371c.
//
// Solidity: function MethodMultiSign(uint64 ChainID, string RedeemKey, bytes TxHash, string Address, bytes[] Signs) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactor) MethodMultiSign(opts *bind.TransactOpts, ChainID uint64, RedeemKey string, TxHash []byte, Address string, Signs [][]byte) (*types.Transaction, error) {
	return _CrossChainManager.contract.Transact(opts, "MethodMultiSign", ChainID, RedeemKey, TxHash, Address, Signs)
}

// MethodMultiSign is a paid mutator transaction binding the contract method 0x4dbb371c.
//
// Solidity: function MethodMultiSign(uint64 ChainID, string RedeemKey, bytes TxHash, string Address, bytes[] Signs) returns(bool success)
func (_CrossChainManager *CrossChainManagerSession) MethodMultiSign(ChainID uint64, RedeemKey string, TxHash []byte, Address string, Signs [][]byte) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodMultiSign(&_CrossChainManager.TransactOpts, ChainID, RedeemKey, TxHash, Address, Signs)
}

// MethodMultiSign is a paid mutator transaction binding the contract method 0x4dbb371c.
//
// Solidity: function MethodMultiSign(uint64 ChainID, string RedeemKey, bytes TxHash, string Address, bytes[] Signs) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactorSession) MethodMultiSign(ChainID uint64, RedeemKey string, TxHash []byte, Address string, Signs [][]byte) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodMultiSign(&_CrossChainManager.TransactOpts, ChainID, RedeemKey, TxHash, Address, Signs)
}

// MethodWhiteChain is a paid mutator transaction binding the contract method 0x2ae0b4bd.
//
// Solidity: function MethodWhiteChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactor) MethodWhiteChain(opts *bind.TransactOpts, ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.contract.Transact(opts, "MethodWhiteChain", ChainID)
}

// MethodWhiteChain is a paid mutator transaction binding the contract method 0x2ae0b4bd.
//
// Solidity: function MethodWhiteChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerSession) MethodWhiteChain(ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodWhiteChain(&_CrossChainManager.TransactOpts, ChainID)
}

// MethodWhiteChain is a paid mutator transaction binding the contract method 0x2ae0b4bd.
//
// Solidity: function MethodWhiteChain(uint64 ChainID) returns(bool success)
func (_CrossChainManager *CrossChainManagerTransactorSession) MethodWhiteChain(ChainID uint64) (*types.Transaction, error) {
	return _CrossChainManager.Contract.MethodWhiteChain(&_CrossChainManager.TransactOpts, ChainID)
}

// CrossChainManagerNOTIFYMAKEPROOFEVENTIterator is returned from FilterNOTIFYMAKEPROOFEVENT and is used to iterate over the raw logs and unpacked data for NOTIFYMAKEPROOFEVENT events raised by the CrossChainManager contract.
type CrossChainManagerNOTIFYMAKEPROOFEVENTIterator struct {
	Event *CrossChainManagerNOTIFYMAKEPROOFEVENT // Event containing the contract specifics and raw log

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
func (it *CrossChainManagerNOTIFYMAKEPROOFEVENTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainManagerNOTIFYMAKEPROOFEVENT)
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
		it.Event = new(CrossChainManagerNOTIFYMAKEPROOFEVENT)
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
func (it *CrossChainManagerNOTIFYMAKEPROOFEVENTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainManagerNOTIFYMAKEPROOFEVENTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainManagerNOTIFYMAKEPROOFEVENT represents a NOTIFYMAKEPROOFEVENT event raised by the CrossChainManager contract.
type CrossChainManagerNOTIFYMAKEPROOFEVENT struct {
	MerkleValueHex string
	BlockHeight    uint64
	Key            string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNOTIFYMAKEPROOFEVENT is a free log retrieval operation binding the contract event 0xa01d1b38777d1877e6b1d5aac0794b6b3e700651b963abcfd2849d4cc843f186.
//
// Solidity: event NOTIFY_MAKE_PROOF_EVENT(string merkleValueHex, uint64 BlockHeight, string key)
func (_CrossChainManager *CrossChainManagerFilterer) FilterNOTIFYMAKEPROOFEVENT(opts *bind.FilterOpts) (*CrossChainManagerNOTIFYMAKEPROOFEVENTIterator, error) {

	logs, sub, err := _CrossChainManager.contract.FilterLogs(opts, "NOTIFY_MAKE_PROOF_EVENT")
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerNOTIFYMAKEPROOFEVENTIterator{contract: _CrossChainManager.contract, event: "NOTIFY_MAKE_PROOF_EVENT", logs: logs, sub: sub}, nil
}

// WatchNOTIFYMAKEPROOFEVENT is a free log subscription operation binding the contract event 0xa01d1b38777d1877e6b1d5aac0794b6b3e700651b963abcfd2849d4cc843f186.
//
// Solidity: event NOTIFY_MAKE_PROOF_EVENT(string merkleValueHex, uint64 BlockHeight, string key)
func (_CrossChainManager *CrossChainManagerFilterer) WatchNOTIFYMAKEPROOFEVENT(opts *bind.WatchOpts, sink chan<- *CrossChainManagerNOTIFYMAKEPROOFEVENT) (event.Subscription, error) {

	logs, sub, err := _CrossChainManager.contract.WatchLogs(opts, "NOTIFY_MAKE_PROOF_EVENT")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainManagerNOTIFYMAKEPROOFEVENT)
				if err := _CrossChainManager.contract.UnpackLog(event, "NOTIFY_MAKE_PROOF_EVENT", log); err != nil {
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

// ParseNOTIFYMAKEPROOFEVENT is a log parse operation binding the contract event 0xa01d1b38777d1877e6b1d5aac0794b6b3e700651b963abcfd2849d4cc843f186.
//
// Solidity: event NOTIFY_MAKE_PROOF_EVENT(string merkleValueHex, uint64 BlockHeight, string key)
func (_CrossChainManager *CrossChainManagerFilterer) ParseNOTIFYMAKEPROOFEVENT(log types.Log) (*CrossChainManagerNOTIFYMAKEPROOFEVENT, error) {
	event := new(CrossChainManagerNOTIFYMAKEPROOFEVENT)
	if err := _CrossChainManager.contract.UnpackLog(event, "NOTIFY_MAKE_PROOF_EVENT", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainManagerBtcTxMultiSignEventIterator is returned from FilterBtcTxMultiSignEvent and is used to iterate over the raw logs and unpacked data for BtcTxMultiSignEvent events raised by the CrossChainManager contract.
type CrossChainManagerBtcTxMultiSignEventIterator struct {
	Event *CrossChainManagerBtcTxMultiSignEvent // Event containing the contract specifics and raw log

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
func (it *CrossChainManagerBtcTxMultiSignEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainManagerBtcTxMultiSignEvent)
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
		it.Event = new(CrossChainManagerBtcTxMultiSignEvent)
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
func (it *CrossChainManagerBtcTxMultiSignEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainManagerBtcTxMultiSignEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainManagerBtcTxMultiSignEvent represents a BtcTxMultiSignEvent event raised by the CrossChainManager contract.
type CrossChainManagerBtcTxMultiSignEvent struct {
	TxHash    []byte
	MultiSign []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBtcTxMultiSignEvent is a free log retrieval operation binding the contract event 0x62fb550ff7fa48f759b0e56ea24757e77b5612d11efbcbdbed9545982cfe1770.
//
// Solidity: event btcTxMultiSignEvent(bytes TxHash, bytes MultiSign)
func (_CrossChainManager *CrossChainManagerFilterer) FilterBtcTxMultiSignEvent(opts *bind.FilterOpts) (*CrossChainManagerBtcTxMultiSignEventIterator, error) {

	logs, sub, err := _CrossChainManager.contract.FilterLogs(opts, "btcTxMultiSignEvent")
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerBtcTxMultiSignEventIterator{contract: _CrossChainManager.contract, event: "btcTxMultiSignEvent", logs: logs, sub: sub}, nil
}

// WatchBtcTxMultiSignEvent is a free log subscription operation binding the contract event 0x62fb550ff7fa48f759b0e56ea24757e77b5612d11efbcbdbed9545982cfe1770.
//
// Solidity: event btcTxMultiSignEvent(bytes TxHash, bytes MultiSign)
func (_CrossChainManager *CrossChainManagerFilterer) WatchBtcTxMultiSignEvent(opts *bind.WatchOpts, sink chan<- *CrossChainManagerBtcTxMultiSignEvent) (event.Subscription, error) {

	logs, sub, err := _CrossChainManager.contract.WatchLogs(opts, "btcTxMultiSignEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainManagerBtcTxMultiSignEvent)
				if err := _CrossChainManager.contract.UnpackLog(event, "btcTxMultiSignEvent", log); err != nil {
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

// ParseBtcTxMultiSignEvent is a log parse operation binding the contract event 0x62fb550ff7fa48f759b0e56ea24757e77b5612d11efbcbdbed9545982cfe1770.
//
// Solidity: event btcTxMultiSignEvent(bytes TxHash, bytes MultiSign)
func (_CrossChainManager *CrossChainManagerFilterer) ParseBtcTxMultiSignEvent(log types.Log) (*CrossChainManagerBtcTxMultiSignEvent, error) {
	event := new(CrossChainManagerBtcTxMultiSignEvent)
	if err := _CrossChainManager.contract.UnpackLog(event, "btcTxMultiSignEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainManagerBtcTxToRelayEventIterator is returned from FilterBtcTxToRelayEvent and is used to iterate over the raw logs and unpacked data for BtcTxToRelayEvent events raised by the CrossChainManager contract.
type CrossChainManagerBtcTxToRelayEventIterator struct {
	Event *CrossChainManagerBtcTxToRelayEvent // Event containing the contract specifics and raw log

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
func (it *CrossChainManagerBtcTxToRelayEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainManagerBtcTxToRelayEvent)
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
		it.Event = new(CrossChainManagerBtcTxToRelayEvent)
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
func (it *CrossChainManagerBtcTxToRelayEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainManagerBtcTxToRelayEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainManagerBtcTxToRelayEvent represents a BtcTxToRelayEvent event raised by the CrossChainManager contract.
type CrossChainManagerBtcTxToRelayEvent struct {
	FromChainID uint64
	ChainID     uint64
	Buf         string
	FromTxHash  string
	RedeemKey   string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBtcTxToRelayEvent is a free log retrieval operation binding the contract event 0x59c070ab5215dda625f463061a0ad421505b2cca1066a8597411c72a5ecac51b.
//
// Solidity: event btcTxToRelayEvent(uint64 FromChainID, uint64 ChainID, string buf, string FromTxHash, string RedeemKey)
func (_CrossChainManager *CrossChainManagerFilterer) FilterBtcTxToRelayEvent(opts *bind.FilterOpts) (*CrossChainManagerBtcTxToRelayEventIterator, error) {

	logs, sub, err := _CrossChainManager.contract.FilterLogs(opts, "btcTxToRelayEvent")
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerBtcTxToRelayEventIterator{contract: _CrossChainManager.contract, event: "btcTxToRelayEvent", logs: logs, sub: sub}, nil
}

// WatchBtcTxToRelayEvent is a free log subscription operation binding the contract event 0x59c070ab5215dda625f463061a0ad421505b2cca1066a8597411c72a5ecac51b.
//
// Solidity: event btcTxToRelayEvent(uint64 FromChainID, uint64 ChainID, string buf, string FromTxHash, string RedeemKey)
func (_CrossChainManager *CrossChainManagerFilterer) WatchBtcTxToRelayEvent(opts *bind.WatchOpts, sink chan<- *CrossChainManagerBtcTxToRelayEvent) (event.Subscription, error) {

	logs, sub, err := _CrossChainManager.contract.WatchLogs(opts, "btcTxToRelayEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainManagerBtcTxToRelayEvent)
				if err := _CrossChainManager.contract.UnpackLog(event, "btcTxToRelayEvent", log); err != nil {
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

// ParseBtcTxToRelayEvent is a log parse operation binding the contract event 0x59c070ab5215dda625f463061a0ad421505b2cca1066a8597411c72a5ecac51b.
//
// Solidity: event btcTxToRelayEvent(uint64 FromChainID, uint64 ChainID, string buf, string FromTxHash, string RedeemKey)
func (_CrossChainManager *CrossChainManagerFilterer) ParseBtcTxToRelayEvent(log types.Log) (*CrossChainManagerBtcTxToRelayEvent, error) {
	event := new(CrossChainManagerBtcTxToRelayEvent)
	if err := _CrossChainManager.contract.UnpackLog(event, "btcTxToRelayEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainManagerMakeBtcTxEventIterator is returned from FilterMakeBtcTxEvent and is used to iterate over the raw logs and unpacked data for MakeBtcTxEvent events raised by the CrossChainManager contract.
type CrossChainManagerMakeBtcTxEventIterator struct {
	Event *CrossChainManagerMakeBtcTxEvent // Event containing the contract specifics and raw log

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
func (it *CrossChainManagerMakeBtcTxEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainManagerMakeBtcTxEvent)
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
		it.Event = new(CrossChainManagerMakeBtcTxEvent)
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
func (it *CrossChainManagerMakeBtcTxEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainManagerMakeBtcTxEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainManagerMakeBtcTxEvent represents a MakeBtcTxEvent event raised by the CrossChainManager contract.
type CrossChainManagerMakeBtcTxEvent struct {
	Rk   string
	Buf  string
	Amts []uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMakeBtcTxEvent is a free log retrieval operation binding the contract event 0xa00d721fd040a2b479d1adf886244de04bdbb3e3e310dc75f1036e2602726234.
//
// Solidity: event makeBtcTxEvent(string rk, string buf, uint64[] amts)
func (_CrossChainManager *CrossChainManagerFilterer) FilterMakeBtcTxEvent(opts *bind.FilterOpts) (*CrossChainManagerMakeBtcTxEventIterator, error) {

	logs, sub, err := _CrossChainManager.contract.FilterLogs(opts, "makeBtcTxEvent")
	if err != nil {
		return nil, err
	}
	return &CrossChainManagerMakeBtcTxEventIterator{contract: _CrossChainManager.contract, event: "makeBtcTxEvent", logs: logs, sub: sub}, nil
}

// WatchMakeBtcTxEvent is a free log subscription operation binding the contract event 0xa00d721fd040a2b479d1adf886244de04bdbb3e3e310dc75f1036e2602726234.
//
// Solidity: event makeBtcTxEvent(string rk, string buf, uint64[] amts)
func (_CrossChainManager *CrossChainManagerFilterer) WatchMakeBtcTxEvent(opts *bind.WatchOpts, sink chan<- *CrossChainManagerMakeBtcTxEvent) (event.Subscription, error) {

	logs, sub, err := _CrossChainManager.contract.WatchLogs(opts, "makeBtcTxEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainManagerMakeBtcTxEvent)
				if err := _CrossChainManager.contract.UnpackLog(event, "makeBtcTxEvent", log); err != nil {
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

// ParseMakeBtcTxEvent is a log parse operation binding the contract event 0xa00d721fd040a2b479d1adf886244de04bdbb3e3e310dc75f1036e2602726234.
//
// Solidity: event makeBtcTxEvent(string rk, string buf, uint64[] amts)
func (_CrossChainManager *CrossChainManagerFilterer) ParseMakeBtcTxEvent(log types.Log) (*CrossChainManagerMakeBtcTxEvent, error) {
	event := new(CrossChainManagerMakeBtcTxEvent)
	if err := _CrossChainManager.contract.UnpackLog(event, "makeBtcTxEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
