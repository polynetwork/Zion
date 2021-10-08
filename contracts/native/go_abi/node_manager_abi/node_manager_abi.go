// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node_manager_abi

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

// node_managerConfiguration is an auto generated low-level Go binding around an user-defined struct.
type node_managerConfiguration struct {
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	MaxBlockChangeView   uint32
}

// node_managerUpdateConfigParam is an auto generated low-level Go binding around an user-defined struct.
type node_managerUpdateConfigParam struct {
	Config node_managerConfiguration
}

// node_managerVBFTPeerInfo is an auto generated low-level Go binding around an user-defined struct.
type node_managerVBFTPeerInfo struct {
	Index      uint32
	PeerPubkey string
	Address    string
}

var (
	MethodApproveCandidate = "approveCandidate"

	MethodBlackNode = "blackNode"

	MethodCommitDpos = "commitDpos"

	MethodInitConfig = "initConfig"

	MethodName = "name"

	MethodQuitNode = "quitNode"

	MethodRegisterCandidate = "registerCandidate"

	MethodUnRegisterCandidate = "unRegisterCandidate"

	MethodUpdateConfig = "updateConfig"

	MethodWhiteNode = "whiteNode"
)

// NodeManagerABI is the input ABI used to generate the binding from.
const NodeManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"signs\",\"type\":\"uint64\"}],\"name\":\"CheckConsensusSignsEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"evtApproveCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"PubkeyList\",\"type\":\"string[]\"}],\"name\":\"evtBlackNode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"evtCommitDpos\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"evtQuitNode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"evtRegisterCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"evtUnRegisterCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structnode_manager.Configuration\",\"name\":\"Config\",\"type\":\"tuple\"}],\"name\":\"evtUpdateConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"evtWhiteNode\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"PeerPubkeyList\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"blackNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commitDpos\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"VrfValue\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"VrfProof\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"Index\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Address\",\"type\":\"string\"}],\"internalType\":\"structnode_manager.VBFTPeerInfo\",\"name\":\"Peers\",\"type\":\"tuple\"}],\"name\":\"initConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"quitNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"registerCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"unRegisterCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"}],\"internalType\":\"structnode_manager.Configuration\",\"name\":\"Config\",\"type\":\"tuple\"}],\"internalType\":\"structnode_manager.UpdateConfigParam\",\"name\":\"ConfigParam\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"whiteNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NodeManagerFuncSigs maps the 4-byte function signature to its string representation.
var NodeManagerFuncSigs = map[string]string{
	"2437d889": "approveCandidate(string,address)",
	"1c21bf1b": "blackNode(string[],address)",
	"39f097d5": "commitDpos()",
	"5bc58f4f": "initConfig(uint32,uint32,uint32,uint32,string,string,(uint32,string,string))",
	"06fdde03": "name()",
	"5506c5bf": "quitNode(string,address)",
	"7a9b6fb9": "registerCandidate(string,address)",
	"e7472e0c": "unRegisterCandidate(string,address)",
	"b8e3e81b": "updateConfig(((uint32,uint32,uint32,uint32)))",
	"8141608b": "whiteNode(string,address)",
}

// NodeManagerBin is the compiled bytecode used for deploying new contracts.
var NodeManagerBin = "0x608060405234801561001057600080fd5b506105b5806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c80635bc58f4f116100665780635bc58f4f146100f65780637a9b6fb9146100e15780638141608b146100e1578063b8e3e81b14610111578063e7472e0c146100e157600080fd5b806306fdde03146100a35780631c21bf1b146100bb5780632437d889146100e157806339f097d5146100ef5780635506c5bf146100e1575b600080fd5b60606040516100b291906104ba565b60405180910390f35b6100d16100c93660046101c5565b600092915050565b60405190151581526020016100b2565b6100d16100c936600461029f565b60006100d1565b6100d161010436600461038a565b6000979650505050505050565b6100d161011f3660046102ed565b50600090565b80356001600160a01b038116811461013c57600080fd5b919050565b600082601f83011261015257600080fd5b813567ffffffffffffffff81111561016c5761016c610569565b61017f601f8201601f1916602001610538565b81815284602083860101111561019457600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff8116811461013c57600080fd5b600080604083850312156101d857600080fd5b823567ffffffffffffffff808211156101f057600080fd5b818501915085601f83011261020457600080fd5b813560208282111561021857610218610569565b8160051b610227828201610538565b8381528281019086840183880185018c101561024257600080fd5b60009350835b8681101561027f5781358881111561025e578586fd5b61026c8e88838d0101610141565b8552509285019290850190600101610248565b5050809850505050610292818801610125565b9450505050509250929050565b600080604083850312156102b257600080fd5b823567ffffffffffffffff8111156102c957600080fd5b6102d585828601610141565b9250506102e460208401610125565b90509250929050565b6000608082840312156102ff57600080fd5b6040516020810167ffffffffffffffff828210818311171561032357610323610569565b8160405260a08301828110828211171561033f5761033f610569565b6040525061034c846101b1565b815261035a602085016101b1565b604083015261036b604085016101b1565b606083015261037c606085016101b1565b608083015281529392505050565b600080600080600080600060e0888a0312156103a557600080fd5b6103ae886101b1565b96506103bc602089016101b1565b95506103ca604089016101b1565b94506103d8606089016101b1565b9350608088013567ffffffffffffffff808211156103f557600080fd5b6104018b838c01610141565b945060a08a013591508082111561041757600080fd5b6104238b838c01610141565b935060c08a013591508082111561043957600080fd5b908901906060828c03121561044d57600080fd5b61045561050f565b61045e836101b1565b815260208301358281111561047257600080fd5b61047e8d828601610141565b60208301525060408301358281111561049657600080fd5b6104a28d828601610141565b60408301525080935050505092959891949750929550565b600060208083528351808285015260005b818110156104e7578581018301518582016040015282016104cb565b818111156104f9576000604083870101525b50601f01601f1916929092016040019392505050565b6040516060810167ffffffffffffffff8111828210171561053257610532610569565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561056157610561610569565b604052919050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122088c676cde9d59b570f4911b54e947a82df34b5d31c1586d7a4785f913ed4fe9164736f6c63430008060033"

// DeployNodeManager deploys a new Ethereum contract, binding an instance of NodeManager to it.
func DeployNodeManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NodeManager, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NodeManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeManager{NodeManagerCaller: NodeManagerCaller{contract: contract}, NodeManagerTransactor: NodeManagerTransactor{contract: contract}, NodeManagerFilterer: NodeManagerFilterer{contract: contract}}, nil
}

// NodeManager is an auto generated Go binding around an Ethereum contract.
type NodeManager struct {
	NodeManagerCaller     // Read-only binding to the contract
	NodeManagerTransactor // Write-only binding to the contract
	NodeManagerFilterer   // Log filterer for contract events
}

// NodeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeManagerSession struct {
	Contract     *NodeManager      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeManagerCallerSession struct {
	Contract *NodeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// NodeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeManagerTransactorSession struct {
	Contract     *NodeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// NodeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeManagerRaw struct {
	Contract *NodeManager // Generic contract binding to access the raw methods on
}

// NodeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeManagerCallerRaw struct {
	Contract *NodeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// NodeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeManagerTransactorRaw struct {
	Contract *NodeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeManager creates a new instance of NodeManager, bound to a specific deployed contract.
func NewNodeManager(address common.Address, backend bind.ContractBackend) (*NodeManager, error) {
	contract, err := bindNodeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeManager{NodeManagerCaller: NodeManagerCaller{contract: contract}, NodeManagerTransactor: NodeManagerTransactor{contract: contract}, NodeManagerFilterer: NodeManagerFilterer{contract: contract}}, nil
}

// NewNodeManagerCaller creates a new read-only instance of NodeManager, bound to a specific deployed contract.
func NewNodeManagerCaller(address common.Address, caller bind.ContractCaller) (*NodeManagerCaller, error) {
	contract, err := bindNodeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeManagerCaller{contract: contract}, nil
}

// NewNodeManagerTransactor creates a new write-only instance of NodeManager, bound to a specific deployed contract.
func NewNodeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeManagerTransactor, error) {
	contract, err := bindNodeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeManagerTransactor{contract: contract}, nil
}

// NewNodeManagerFilterer creates a new log filterer instance of NodeManager, bound to a specific deployed contract.
func NewNodeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeManagerFilterer, error) {
	contract, err := bindNodeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeManagerFilterer{contract: contract}, nil
}

// bindNodeManager binds a generic wrapper to an already deployed contract.
func bindNodeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeManager *NodeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeManager.Contract.NodeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeManager *NodeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.Contract.NodeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeManager *NodeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeManager.Contract.NodeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeManager *NodeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeManager *NodeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeManager *NodeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeManager.Contract.contract.Transact(opts, method, params...)
}

// ApproveCandidate is a paid mutator transaction binding the contract method 0x2437d889.
//
// Solidity: function approveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) ApproveCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "approveCandidate", PeerPubkey, Address)
}

// ApproveCandidate is a paid mutator transaction binding the contract method 0x2437d889.
//
// Solidity: function approveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) ApproveCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.ApproveCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// ApproveCandidate is a paid mutator transaction binding the contract method 0x2437d889.
//
// Solidity: function approveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) ApproveCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.ApproveCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// BlackNode is a paid mutator transaction binding the contract method 0x1c21bf1b.
//
// Solidity: function blackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) BlackNode(opts *bind.TransactOpts, PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "blackNode", PeerPubkeyList, Address)
}

// BlackNode is a paid mutator transaction binding the contract method 0x1c21bf1b.
//
// Solidity: function blackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) BlackNode(PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.BlackNode(&_NodeManager.TransactOpts, PeerPubkeyList, Address)
}

// BlackNode is a paid mutator transaction binding the contract method 0x1c21bf1b.
//
// Solidity: function blackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) BlackNode(PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.BlackNode(&_NodeManager.TransactOpts, PeerPubkeyList, Address)
}

// CommitDpos is a paid mutator transaction binding the contract method 0x39f097d5.
//
// Solidity: function commitDpos() returns(bool success)
func (_NodeManager *NodeManagerTransactor) CommitDpos(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "commitDpos")
}

// CommitDpos is a paid mutator transaction binding the contract method 0x39f097d5.
//
// Solidity: function commitDpos() returns(bool success)
func (_NodeManager *NodeManagerSession) CommitDpos() (*types.Transaction, error) {
	return _NodeManager.Contract.CommitDpos(&_NodeManager.TransactOpts)
}

// CommitDpos is a paid mutator transaction binding the contract method 0x39f097d5.
//
// Solidity: function commitDpos() returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) CommitDpos() (*types.Transaction, error) {
	return _NodeManager.Contract.CommitDpos(&_NodeManager.TransactOpts)
}

// InitConfig is a paid mutator transaction binding the contract method 0x5bc58f4f.
//
// Solidity: function initConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerTransactor) InitConfig(opts *bind.TransactOpts, BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "initConfig", BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// InitConfig is a paid mutator transaction binding the contract method 0x5bc58f4f.
//
// Solidity: function initConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerSession) InitConfig(BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.Contract.InitConfig(&_NodeManager.TransactOpts, BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// InitConfig is a paid mutator transaction binding the contract method 0x5bc58f4f.
//
// Solidity: function initConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) InitConfig(BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.Contract.InitConfig(&_NodeManager.TransactOpts, BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_NodeManager *NodeManagerTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_NodeManager *NodeManagerSession) Name() (*types.Transaction, error) {
	return _NodeManager.Contract.Name(&_NodeManager.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_NodeManager *NodeManagerTransactorSession) Name() (*types.Transaction, error) {
	return _NodeManager.Contract.Name(&_NodeManager.TransactOpts)
}

// QuitNode is a paid mutator transaction binding the contract method 0x5506c5bf.
//
// Solidity: function quitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) QuitNode(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "quitNode", PeerPubkey, Address)
}

// QuitNode is a paid mutator transaction binding the contract method 0x5506c5bf.
//
// Solidity: function quitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) QuitNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.QuitNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// QuitNode is a paid mutator transaction binding the contract method 0x5506c5bf.
//
// Solidity: function quitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) QuitNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.QuitNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0x7a9b6fb9.
//
// Solidity: function registerCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) RegisterCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "registerCandidate", PeerPubkey, Address)
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0x7a9b6fb9.
//
// Solidity: function registerCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) RegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.RegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0x7a9b6fb9.
//
// Solidity: function registerCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) RegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.RegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// UnRegisterCandidate is a paid mutator transaction binding the contract method 0xe7472e0c.
//
// Solidity: function unRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) UnRegisterCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "unRegisterCandidate", PeerPubkey, Address)
}

// UnRegisterCandidate is a paid mutator transaction binding the contract method 0xe7472e0c.
//
// Solidity: function unRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) UnRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.UnRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// UnRegisterCandidate is a paid mutator transaction binding the contract method 0xe7472e0c.
//
// Solidity: function unRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) UnRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.UnRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0xb8e3e81b.
//
// Solidity: function updateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerTransactor) UpdateConfig(opts *bind.TransactOpts, ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "updateConfig", ConfigParam)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0xb8e3e81b.
//
// Solidity: function updateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerSession) UpdateConfig(ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.Contract.UpdateConfig(&_NodeManager.TransactOpts, ConfigParam)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0xb8e3e81b.
//
// Solidity: function updateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) UpdateConfig(ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.Contract.UpdateConfig(&_NodeManager.TransactOpts, ConfigParam)
}

// WhiteNode is a paid mutator transaction binding the contract method 0x8141608b.
//
// Solidity: function whiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) WhiteNode(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "whiteNode", PeerPubkey, Address)
}

// WhiteNode is a paid mutator transaction binding the contract method 0x8141608b.
//
// Solidity: function whiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) WhiteNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.WhiteNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// WhiteNode is a paid mutator transaction binding the contract method 0x8141608b.
//
// Solidity: function whiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) WhiteNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.WhiteNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// NodeManagerCheckConsensusSignsEventIterator is returned from FilterCheckConsensusSignsEvent and is used to iterate over the raw logs and unpacked data for CheckConsensusSignsEvent events raised by the NodeManager contract.
type NodeManagerCheckConsensusSignsEventIterator struct {
	Event *NodeManagerCheckConsensusSignsEvent // Event containing the contract specifics and raw log

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
func (it *NodeManagerCheckConsensusSignsEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerCheckConsensusSignsEvent)
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
		it.Event = new(NodeManagerCheckConsensusSignsEvent)
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
func (it *NodeManagerCheckConsensusSignsEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerCheckConsensusSignsEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerCheckConsensusSignsEvent represents a CheckConsensusSignsEvent event raised by the NodeManager contract.
type NodeManagerCheckConsensusSignsEvent struct {
	Signs uint64
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCheckConsensusSignsEvent is a free log retrieval operation binding the contract event 0x11a255420a0713251f48807b7c1efabe8e912db6332ec448cd07bf4ca381c0d5.
//
// Solidity: event CheckConsensusSignsEvent(uint64 signs)
func (_NodeManager *NodeManagerFilterer) FilterCheckConsensusSignsEvent(opts *bind.FilterOpts) (*NodeManagerCheckConsensusSignsEventIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "CheckConsensusSignsEvent")
	if err != nil {
		return nil, err
	}
	return &NodeManagerCheckConsensusSignsEventIterator{contract: _NodeManager.contract, event: "CheckConsensusSignsEvent", logs: logs, sub: sub}, nil
}

// WatchCheckConsensusSignsEvent is a free log subscription operation binding the contract event 0x11a255420a0713251f48807b7c1efabe8e912db6332ec448cd07bf4ca381c0d5.
//
// Solidity: event CheckConsensusSignsEvent(uint64 signs)
func (_NodeManager *NodeManagerFilterer) WatchCheckConsensusSignsEvent(opts *bind.WatchOpts, sink chan<- *NodeManagerCheckConsensusSignsEvent) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "CheckConsensusSignsEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerCheckConsensusSignsEvent)
				if err := _NodeManager.contract.UnpackLog(event, "CheckConsensusSignsEvent", log); err != nil {
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

// ParseCheckConsensusSignsEvent is a log parse operation binding the contract event 0x11a255420a0713251f48807b7c1efabe8e912db6332ec448cd07bf4ca381c0d5.
//
// Solidity: event CheckConsensusSignsEvent(uint64 signs)
func (_NodeManager *NodeManagerFilterer) ParseCheckConsensusSignsEvent(log types.Log) (*NodeManagerCheckConsensusSignsEvent, error) {
	event := new(NodeManagerCheckConsensusSignsEvent)
	if err := _NodeManager.contract.UnpackLog(event, "CheckConsensusSignsEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerApproveCandidateIterator is returned from FilterApproveCandidate and is used to iterate over the raw logs and unpacked data for ApproveCandidate events raised by the NodeManager contract.
type NodeManagerApproveCandidateIterator struct {
	Event *NodeManagerApproveCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerApproveCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerApproveCandidate)
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
		it.Event = new(NodeManagerApproveCandidate)
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
func (it *NodeManagerApproveCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerApproveCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerApproveCandidate represents a ApproveCandidate event raised by the NodeManager contract.
type NodeManagerApproveCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterApproveCandidate is a free log retrieval operation binding the contract event 0x902caf6ae405aec1d7aa7994a0c27c467dc04d5efc4b94188aeb483cdbce2276.
//
// Solidity: event evtApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterApproveCandidate(opts *bind.FilterOpts) (*NodeManagerApproveCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtApproveCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerApproveCandidateIterator{contract: _NodeManager.contract, event: "evtApproveCandidate", logs: logs, sub: sub}, nil
}

// WatchApproveCandidate is a free log subscription operation binding the contract event 0x902caf6ae405aec1d7aa7994a0c27c467dc04d5efc4b94188aeb483cdbce2276.
//
// Solidity: event evtApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchApproveCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerApproveCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtApproveCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerApproveCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "evtApproveCandidate", log); err != nil {
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

// ParseApproveCandidate is a log parse operation binding the contract event 0x902caf6ae405aec1d7aa7994a0c27c467dc04d5efc4b94188aeb483cdbce2276.
//
// Solidity: event evtApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseApproveCandidate(log types.Log) (*NodeManagerApproveCandidate, error) {
	event := new(NodeManagerApproveCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "evtApproveCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerBlackNodeIterator is returned from FilterBlackNode and is used to iterate over the raw logs and unpacked data for BlackNode events raised by the NodeManager contract.
type NodeManagerBlackNodeIterator struct {
	Event *NodeManagerBlackNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerBlackNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerBlackNode)
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
		it.Event = new(NodeManagerBlackNode)
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
func (it *NodeManagerBlackNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerBlackNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerBlackNode represents a BlackNode event raised by the NodeManager contract.
type NodeManagerBlackNode struct {
	PubkeyList []string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBlackNode is a free log retrieval operation binding the contract event 0x0812f652fc317d5144576e96dd27a9260a90bdd876cc69afbd2405d989ec9eb5.
//
// Solidity: event evtBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) FilterBlackNode(opts *bind.FilterOpts) (*NodeManagerBlackNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtBlackNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerBlackNodeIterator{contract: _NodeManager.contract, event: "evtBlackNode", logs: logs, sub: sub}, nil
}

// WatchBlackNode is a free log subscription operation binding the contract event 0x0812f652fc317d5144576e96dd27a9260a90bdd876cc69afbd2405d989ec9eb5.
//
// Solidity: event evtBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) WatchBlackNode(opts *bind.WatchOpts, sink chan<- *NodeManagerBlackNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtBlackNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerBlackNode)
				if err := _NodeManager.contract.UnpackLog(event, "evtBlackNode", log); err != nil {
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

// ParseBlackNode is a log parse operation binding the contract event 0x0812f652fc317d5144576e96dd27a9260a90bdd876cc69afbd2405d989ec9eb5.
//
// Solidity: event evtBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) ParseBlackNode(log types.Log) (*NodeManagerBlackNode, error) {
	event := new(NodeManagerBlackNode)
	if err := _NodeManager.contract.UnpackLog(event, "evtBlackNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerCommitDposIterator is returned from FilterCommitDpos and is used to iterate over the raw logs and unpacked data for CommitDpos events raised by the NodeManager contract.
type NodeManagerCommitDposIterator struct {
	Event *NodeManagerCommitDpos // Event containing the contract specifics and raw log

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
func (it *NodeManagerCommitDposIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerCommitDpos)
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
		it.Event = new(NodeManagerCommitDpos)
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
func (it *NodeManagerCommitDposIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerCommitDposIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerCommitDpos represents a CommitDpos event raised by the NodeManager contract.
type NodeManagerCommitDpos struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCommitDpos is a free log retrieval operation binding the contract event 0x4a1167039259b24583f23f1b361b44243af5c995970a48c79605418fd6319af7.
//
// Solidity: event evtCommitDpos()
func (_NodeManager *NodeManagerFilterer) FilterCommitDpos(opts *bind.FilterOpts) (*NodeManagerCommitDposIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtCommitDpos")
	if err != nil {
		return nil, err
	}
	return &NodeManagerCommitDposIterator{contract: _NodeManager.contract, event: "evtCommitDpos", logs: logs, sub: sub}, nil
}

// WatchCommitDpos is a free log subscription operation binding the contract event 0x4a1167039259b24583f23f1b361b44243af5c995970a48c79605418fd6319af7.
//
// Solidity: event evtCommitDpos()
func (_NodeManager *NodeManagerFilterer) WatchCommitDpos(opts *bind.WatchOpts, sink chan<- *NodeManagerCommitDpos) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtCommitDpos")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerCommitDpos)
				if err := _NodeManager.contract.UnpackLog(event, "evtCommitDpos", log); err != nil {
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

// ParseCommitDpos is a log parse operation binding the contract event 0x4a1167039259b24583f23f1b361b44243af5c995970a48c79605418fd6319af7.
//
// Solidity: event evtCommitDpos()
func (_NodeManager *NodeManagerFilterer) ParseCommitDpos(log types.Log) (*NodeManagerCommitDpos, error) {
	event := new(NodeManagerCommitDpos)
	if err := _NodeManager.contract.UnpackLog(event, "evtCommitDpos", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerQuitNodeIterator is returned from FilterQuitNode and is used to iterate over the raw logs and unpacked data for QuitNode events raised by the NodeManager contract.
type NodeManagerQuitNodeIterator struct {
	Event *NodeManagerQuitNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerQuitNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerQuitNode)
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
		it.Event = new(NodeManagerQuitNode)
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
func (it *NodeManagerQuitNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerQuitNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerQuitNode represents a QuitNode event raised by the NodeManager contract.
type NodeManagerQuitNode struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterQuitNode is a free log retrieval operation binding the contract event 0xd9960f1ffee6726984338fbb6b8243175d18b2a9b984b9b6c84f6b7d312b8aa8.
//
// Solidity: event evtQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterQuitNode(opts *bind.FilterOpts) (*NodeManagerQuitNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtQuitNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerQuitNodeIterator{contract: _NodeManager.contract, event: "evtQuitNode", logs: logs, sub: sub}, nil
}

// WatchQuitNode is a free log subscription operation binding the contract event 0xd9960f1ffee6726984338fbb6b8243175d18b2a9b984b9b6c84f6b7d312b8aa8.
//
// Solidity: event evtQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchQuitNode(opts *bind.WatchOpts, sink chan<- *NodeManagerQuitNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtQuitNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerQuitNode)
				if err := _NodeManager.contract.UnpackLog(event, "evtQuitNode", log); err != nil {
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

// ParseQuitNode is a log parse operation binding the contract event 0xd9960f1ffee6726984338fbb6b8243175d18b2a9b984b9b6c84f6b7d312b8aa8.
//
// Solidity: event evtQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseQuitNode(log types.Log) (*NodeManagerQuitNode, error) {
	event := new(NodeManagerQuitNode)
	if err := _NodeManager.contract.UnpackLog(event, "evtQuitNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerRegisterCandidateIterator is returned from FilterRegisterCandidate and is used to iterate over the raw logs and unpacked data for RegisterCandidate events raised by the NodeManager contract.
type NodeManagerRegisterCandidateIterator struct {
	Event *NodeManagerRegisterCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerRegisterCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerRegisterCandidate)
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
		it.Event = new(NodeManagerRegisterCandidate)
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
func (it *NodeManagerRegisterCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerRegisterCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerRegisterCandidate represents a RegisterCandidate event raised by the NodeManager contract.
type NodeManagerRegisterCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegisterCandidate is a free log retrieval operation binding the contract event 0x9e86b25c573201cdd348c36cc051e55aeeb2fdf6f054c4d52dd100a475350626.
//
// Solidity: event evtRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterRegisterCandidate(opts *bind.FilterOpts) (*NodeManagerRegisterCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerRegisterCandidateIterator{contract: _NodeManager.contract, event: "evtRegisterCandidate", logs: logs, sub: sub}, nil
}

// WatchRegisterCandidate is a free log subscription operation binding the contract event 0x9e86b25c573201cdd348c36cc051e55aeeb2fdf6f054c4d52dd100a475350626.
//
// Solidity: event evtRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchRegisterCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerRegisterCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerRegisterCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "evtRegisterCandidate", log); err != nil {
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

// ParseRegisterCandidate is a log parse operation binding the contract event 0x9e86b25c573201cdd348c36cc051e55aeeb2fdf6f054c4d52dd100a475350626.
//
// Solidity: event evtRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseRegisterCandidate(log types.Log) (*NodeManagerRegisterCandidate, error) {
	event := new(NodeManagerRegisterCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "evtRegisterCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerUnRegisterCandidateIterator is returned from FilterUnRegisterCandidate and is used to iterate over the raw logs and unpacked data for UnRegisterCandidate events raised by the NodeManager contract.
type NodeManagerUnRegisterCandidateIterator struct {
	Event *NodeManagerUnRegisterCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerUnRegisterCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerUnRegisterCandidate)
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
		it.Event = new(NodeManagerUnRegisterCandidate)
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
func (it *NodeManagerUnRegisterCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerUnRegisterCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerUnRegisterCandidate represents a UnRegisterCandidate event raised by the NodeManager contract.
type NodeManagerUnRegisterCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUnRegisterCandidate is a free log retrieval operation binding the contract event 0xb472901c46184dd851d1cc49bfb0c8a2e148ae4d8ae9961de16ec1451df3c557.
//
// Solidity: event evtUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterUnRegisterCandidate(opts *bind.FilterOpts) (*NodeManagerUnRegisterCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtUnRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerUnRegisterCandidateIterator{contract: _NodeManager.contract, event: "evtUnRegisterCandidate", logs: logs, sub: sub}, nil
}

// WatchUnRegisterCandidate is a free log subscription operation binding the contract event 0xb472901c46184dd851d1cc49bfb0c8a2e148ae4d8ae9961de16ec1451df3c557.
//
// Solidity: event evtUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchUnRegisterCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerUnRegisterCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtUnRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerUnRegisterCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "evtUnRegisterCandidate", log); err != nil {
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

// ParseUnRegisterCandidate is a log parse operation binding the contract event 0xb472901c46184dd851d1cc49bfb0c8a2e148ae4d8ae9961de16ec1451df3c557.
//
// Solidity: event evtUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseUnRegisterCandidate(log types.Log) (*NodeManagerUnRegisterCandidate, error) {
	event := new(NodeManagerUnRegisterCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "evtUnRegisterCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerUpdateConfigIterator is returned from FilterUpdateConfig and is used to iterate over the raw logs and unpacked data for UpdateConfig events raised by the NodeManager contract.
type NodeManagerUpdateConfigIterator struct {
	Event *NodeManagerUpdateConfig // Event containing the contract specifics and raw log

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
func (it *NodeManagerUpdateConfigIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerUpdateConfig)
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
		it.Event = new(NodeManagerUpdateConfig)
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
func (it *NodeManagerUpdateConfigIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerUpdateConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerUpdateConfig represents a UpdateConfig event raised by the NodeManager contract.
type NodeManagerUpdateConfig struct {
	Config node_managerConfiguration
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUpdateConfig is a free log retrieval operation binding the contract event 0x0ea859a9f788215a0a769be70546414d2c2bd5c6eac02eb711b7c47d0472bd76.
//
// Solidity: event evtUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) FilterUpdateConfig(opts *bind.FilterOpts) (*NodeManagerUpdateConfigIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtUpdateConfig")
	if err != nil {
		return nil, err
	}
	return &NodeManagerUpdateConfigIterator{contract: _NodeManager.contract, event: "evtUpdateConfig", logs: logs, sub: sub}, nil
}

// WatchUpdateConfig is a free log subscription operation binding the contract event 0x0ea859a9f788215a0a769be70546414d2c2bd5c6eac02eb711b7c47d0472bd76.
//
// Solidity: event evtUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) WatchUpdateConfig(opts *bind.WatchOpts, sink chan<- *NodeManagerUpdateConfig) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtUpdateConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerUpdateConfig)
				if err := _NodeManager.contract.UnpackLog(event, "evtUpdateConfig", log); err != nil {
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

// ParseUpdateConfig is a log parse operation binding the contract event 0x0ea859a9f788215a0a769be70546414d2c2bd5c6eac02eb711b7c47d0472bd76.
//
// Solidity: event evtUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) ParseUpdateConfig(log types.Log) (*NodeManagerUpdateConfig, error) {
	event := new(NodeManagerUpdateConfig)
	if err := _NodeManager.contract.UnpackLog(event, "evtUpdateConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerWhiteNodeIterator is returned from FilterWhiteNode and is used to iterate over the raw logs and unpacked data for WhiteNode events raised by the NodeManager contract.
type NodeManagerWhiteNodeIterator struct {
	Event *NodeManagerWhiteNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerWhiteNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerWhiteNode)
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
		it.Event = new(NodeManagerWhiteNode)
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
func (it *NodeManagerWhiteNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerWhiteNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerWhiteNode represents a WhiteNode event raised by the NodeManager contract.
type NodeManagerWhiteNode struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWhiteNode is a free log retrieval operation binding the contract event 0xb7d8dbc21c076ae2b1b72557e8ef8ad0a55e17e2eba9a97d5f132cae3e533cd1.
//
// Solidity: event evtWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterWhiteNode(opts *bind.FilterOpts) (*NodeManagerWhiteNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "evtWhiteNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerWhiteNodeIterator{contract: _NodeManager.contract, event: "evtWhiteNode", logs: logs, sub: sub}, nil
}

// WatchWhiteNode is a free log subscription operation binding the contract event 0xb7d8dbc21c076ae2b1b72557e8ef8ad0a55e17e2eba9a97d5f132cae3e533cd1.
//
// Solidity: event evtWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchWhiteNode(opts *bind.WatchOpts, sink chan<- *NodeManagerWhiteNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "evtWhiteNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerWhiteNode)
				if err := _NodeManager.contract.UnpackLog(event, "evtWhiteNode", log); err != nil {
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

// ParseWhiteNode is a log parse operation binding the contract event 0xb7d8dbc21c076ae2b1b72557e8ef8ad0a55e17e2eba9a97d5f132cae3e533cd1.
//
// Solidity: event evtWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseWhiteNode(log types.Log) (*NodeManagerWhiteNode, error) {
	event := new(NodeManagerWhiteNode)
	if err := _NodeManager.contract.UnpackLog(event, "evtWhiteNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
