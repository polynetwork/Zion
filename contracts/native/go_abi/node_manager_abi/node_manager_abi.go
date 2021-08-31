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

// NodeManagerABI is the input ABI used to generate the binding from.
const NodeManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"signs\",\"type\":\"uint64\"}],\"name\":\"CheckConsensusSignsEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"EventApproveCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"PubkeyList\",\"type\":\"string[]\"}],\"name\":\"EventBlackNode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EventCommitDpos\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"EventQuitNode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"EventRegisterCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"EventUnRegisterCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structnode_manager.Configuration\",\"name\":\"Config\",\"type\":\"tuple\"}],\"name\":\"EventUpdateConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Pubkey\",\"type\":\"string\"}],\"name\":\"EventWhiteNode\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"PeerPubkeyList\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodBlackNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodCommitDpos\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"VrfValue\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"VrfProof\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"Index\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Address\",\"type\":\"string\"}],\"internalType\":\"structnode_manager.VBFTPeerInfo\",\"name\":\"Peers\",\"type\":\"tuple\"}],\"name\":\"MethodInitConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodQuitNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodRegisterCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodUnRegisterCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"BlockMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"HashMsgDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"PeerHandshakeTimeout\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"MaxBlockChangeView\",\"type\":\"uint32\"}],\"internalType\":\"structnode_manager.Configuration\",\"name\":\"Config\",\"type\":\"tuple\"}],\"internalType\":\"structnode_manager.UpdateConfigParam\",\"name\":\"ConfigParam\",\"type\":\"tuple\"}],\"name\":\"MethodUpdateConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"PeerPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodWhiteNode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NodeManagerFuncSigs maps the 4-byte function signature to its string representation.
var NodeManagerFuncSigs = map[string]string{
	"17fcf5b0": "MethodApproveCandidate(string,address)",
	"39d38786": "MethodBlackNode(string[],address)",
	"b6c90d19": "MethodCommitDpos()",
	"e50f8f44": "MethodContractName()",
	"f663417a": "MethodInitConfig(uint32,uint32,uint32,uint32,string,string,(uint32,string,string))",
	"a61bc535": "MethodQuitNode(string,address)",
	"943822fb": "MethodRegisterCandidate(string,address)",
	"54b017db": "MethodUnRegisterCandidate(string,address)",
	"c28ced43": "MethodUpdateConfig(((uint32,uint32,uint32,uint32)))",
	"fe3c2fb0": "MethodWhiteNode(string,address)",
}

// NodeManagerBin is the compiled bytecode used for deploying new contracts.
var NodeManagerBin = "0x608060405234801561001057600080fd5b506105b1806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063b6c90d1911610066578063b6c90d19146100dc578063c28ced43146100e3578063e50f8f44146100f7578063f663417a14610106578063fe3c2fb0146100a357600080fd5b806317fcf5b0146100a357806339d38786146100ce57806354b017db146100a3578063943822fb146100a3578063a61bc535146100a3575b600080fd5b6100b96100b136600461029b565b600092915050565b60405190151581526020015b60405180910390f35b6100b96100b13660046101c1565b60006100b9565b6100b96100f13660046102e9565b50600090565b60606040516100c591906104b6565b6100b9610114366004610386565b6000979650505050505050565b80356001600160a01b038116811461013857600080fd5b919050565b600082601f83011261014e57600080fd5b813567ffffffffffffffff81111561016857610168610565565b61017b601f8201601f1916602001610534565b81815284602083860101111561019057600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff8116811461013857600080fd5b600080604083850312156101d457600080fd5b823567ffffffffffffffff808211156101ec57600080fd5b818501915085601f83011261020057600080fd5b813560208282111561021457610214610565565b8160051b610223828201610534565b8381528281019086840183880185018c101561023e57600080fd5b60009350835b8681101561027b5781358881111561025a578586fd5b6102688e88838d010161013d565b8552509285019290850190600101610244565b505080985050505061028e818801610121565b9450505050509250929050565b600080604083850312156102ae57600080fd5b823567ffffffffffffffff8111156102c557600080fd5b6102d18582860161013d565b9250506102e060208401610121565b90509250929050565b6000608082840312156102fb57600080fd5b6040516020810167ffffffffffffffff828210818311171561031f5761031f610565565b8160405260a08301828110828211171561033b5761033b610565565b60405250610348846101ad565b8152610356602085016101ad565b6040830152610367604085016101ad565b6060830152610378606085016101ad565b608083015281529392505050565b600080600080600080600060e0888a0312156103a157600080fd5b6103aa886101ad565b96506103b8602089016101ad565b95506103c6604089016101ad565b94506103d4606089016101ad565b9350608088013567ffffffffffffffff808211156103f157600080fd5b6103fd8b838c0161013d565b945060a08a013591508082111561041357600080fd5b61041f8b838c0161013d565b935060c08a013591508082111561043557600080fd5b908901906060828c03121561044957600080fd5b61045161050b565b61045a836101ad565b815260208301358281111561046e57600080fd5b61047a8d82860161013d565b60208301525060408301358281111561049257600080fd5b61049e8d82860161013d565b60408301525080935050505092959891949750929550565b600060208083528351808285015260005b818110156104e3578581018301518582016040015282016104c7565b818111156104f5576000604083870101525b50601f01601f1916929092016040019392505050565b6040516060810167ffffffffffffffff8111828210171561052e5761052e610565565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561055d5761055d610565565b604052919050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122011425e946423e7a7cdd626df5d4d065bc4c2c50e7927f31f939d7ca57f15891b64736f6c63430008060033"

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

// MethodApproveCandidate is a paid mutator transaction binding the contract method 0x17fcf5b0.
//
// Solidity: function MethodApproveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodApproveCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodApproveCandidate", PeerPubkey, Address)
}

// MethodApproveCandidate is a paid mutator transaction binding the contract method 0x17fcf5b0.
//
// Solidity: function MethodApproveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodApproveCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodApproveCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodApproveCandidate is a paid mutator transaction binding the contract method 0x17fcf5b0.
//
// Solidity: function MethodApproveCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodApproveCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodApproveCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodBlackNode is a paid mutator transaction binding the contract method 0x39d38786.
//
// Solidity: function MethodBlackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodBlackNode(opts *bind.TransactOpts, PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodBlackNode", PeerPubkeyList, Address)
}

// MethodBlackNode is a paid mutator transaction binding the contract method 0x39d38786.
//
// Solidity: function MethodBlackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodBlackNode(PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodBlackNode(&_NodeManager.TransactOpts, PeerPubkeyList, Address)
}

// MethodBlackNode is a paid mutator transaction binding the contract method 0x39d38786.
//
// Solidity: function MethodBlackNode(string[] PeerPubkeyList, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodBlackNode(PeerPubkeyList []string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodBlackNode(&_NodeManager.TransactOpts, PeerPubkeyList, Address)
}

// MethodCommitDpos is a paid mutator transaction binding the contract method 0xb6c90d19.
//
// Solidity: function MethodCommitDpos() returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodCommitDpos(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodCommitDpos")
}

// MethodCommitDpos is a paid mutator transaction binding the contract method 0xb6c90d19.
//
// Solidity: function MethodCommitDpos() returns(bool success)
func (_NodeManager *NodeManagerSession) MethodCommitDpos() (*types.Transaction, error) {
	return _NodeManager.Contract.MethodCommitDpos(&_NodeManager.TransactOpts)
}

// MethodCommitDpos is a paid mutator transaction binding the contract method 0xb6c90d19.
//
// Solidity: function MethodCommitDpos() returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodCommitDpos() (*types.Transaction, error) {
	return _NodeManager.Contract.MethodCommitDpos(&_NodeManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_NodeManager *NodeManagerTransactor) MethodContractName(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodContractName")
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_NodeManager *NodeManagerSession) MethodContractName() (*types.Transaction, error) {
	return _NodeManager.Contract.MethodContractName(&_NodeManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_NodeManager *NodeManagerTransactorSession) MethodContractName() (*types.Transaction, error) {
	return _NodeManager.Contract.MethodContractName(&_NodeManager.TransactOpts)
}

// MethodInitConfig is a paid mutator transaction binding the contract method 0xf663417a.
//
// Solidity: function MethodInitConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodInitConfig(opts *bind.TransactOpts, BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodInitConfig", BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// MethodInitConfig is a paid mutator transaction binding the contract method 0xf663417a.
//
// Solidity: function MethodInitConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodInitConfig(BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodInitConfig(&_NodeManager.TransactOpts, BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// MethodInitConfig is a paid mutator transaction binding the contract method 0xf663417a.
//
// Solidity: function MethodInitConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string VrfValue, string VrfProof, (uint32,string,string) Peers) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodInitConfig(BlockMsgDelay uint32, HashMsgDelay uint32, PeerHandshakeTimeout uint32, MaxBlockChangeView uint32, VrfValue string, VrfProof string, Peers node_managerVBFTPeerInfo) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodInitConfig(&_NodeManager.TransactOpts, BlockMsgDelay, HashMsgDelay, PeerHandshakeTimeout, MaxBlockChangeView, VrfValue, VrfProof, Peers)
}

// MethodQuitNode is a paid mutator transaction binding the contract method 0xa61bc535.
//
// Solidity: function MethodQuitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodQuitNode(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodQuitNode", PeerPubkey, Address)
}

// MethodQuitNode is a paid mutator transaction binding the contract method 0xa61bc535.
//
// Solidity: function MethodQuitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodQuitNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodQuitNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodQuitNode is a paid mutator transaction binding the contract method 0xa61bc535.
//
// Solidity: function MethodQuitNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodQuitNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodQuitNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodRegisterCandidate is a paid mutator transaction binding the contract method 0x943822fb.
//
// Solidity: function MethodRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodRegisterCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodRegisterCandidate", PeerPubkey, Address)
}

// MethodRegisterCandidate is a paid mutator transaction binding the contract method 0x943822fb.
//
// Solidity: function MethodRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodRegisterCandidate is a paid mutator transaction binding the contract method 0x943822fb.
//
// Solidity: function MethodRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodUnRegisterCandidate is a paid mutator transaction binding the contract method 0x54b017db.
//
// Solidity: function MethodUnRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodUnRegisterCandidate(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodUnRegisterCandidate", PeerPubkey, Address)
}

// MethodUnRegisterCandidate is a paid mutator transaction binding the contract method 0x54b017db.
//
// Solidity: function MethodUnRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodUnRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodUnRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodUnRegisterCandidate is a paid mutator transaction binding the contract method 0x54b017db.
//
// Solidity: function MethodUnRegisterCandidate(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodUnRegisterCandidate(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodUnRegisterCandidate(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodUpdateConfig is a paid mutator transaction binding the contract method 0xc28ced43.
//
// Solidity: function MethodUpdateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodUpdateConfig(opts *bind.TransactOpts, ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodUpdateConfig", ConfigParam)
}

// MethodUpdateConfig is a paid mutator transaction binding the contract method 0xc28ced43.
//
// Solidity: function MethodUpdateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodUpdateConfig(ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodUpdateConfig(&_NodeManager.TransactOpts, ConfigParam)
}

// MethodUpdateConfig is a paid mutator transaction binding the contract method 0xc28ced43.
//
// Solidity: function MethodUpdateConfig(((uint32,uint32,uint32,uint32)) ConfigParam) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodUpdateConfig(ConfigParam node_managerUpdateConfigParam) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodUpdateConfig(&_NodeManager.TransactOpts, ConfigParam)
}

// MethodWhiteNode is a paid mutator transaction binding the contract method 0xfe3c2fb0.
//
// Solidity: function MethodWhiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactor) MethodWhiteNode(opts *bind.TransactOpts, PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.contract.Transact(opts, "MethodWhiteNode", PeerPubkey, Address)
}

// MethodWhiteNode is a paid mutator transaction binding the contract method 0xfe3c2fb0.
//
// Solidity: function MethodWhiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerSession) MethodWhiteNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodWhiteNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
}

// MethodWhiteNode is a paid mutator transaction binding the contract method 0xfe3c2fb0.
//
// Solidity: function MethodWhiteNode(string PeerPubkey, address Address) returns(bool success)
func (_NodeManager *NodeManagerTransactorSession) MethodWhiteNode(PeerPubkey string, Address common.Address) (*types.Transaction, error) {
	return _NodeManager.Contract.MethodWhiteNode(&_NodeManager.TransactOpts, PeerPubkey, Address)
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

// NodeManagerEventApproveCandidateIterator is returned from FilterEventApproveCandidate and is used to iterate over the raw logs and unpacked data for EventApproveCandidate events raised by the NodeManager contract.
type NodeManagerEventApproveCandidateIterator struct {
	Event *NodeManagerEventApproveCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventApproveCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventApproveCandidate)
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
		it.Event = new(NodeManagerEventApproveCandidate)
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
func (it *NodeManagerEventApproveCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventApproveCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventApproveCandidate represents a EventApproveCandidate event raised by the NodeManager contract.
type NodeManagerEventApproveCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventApproveCandidate is a free log retrieval operation binding the contract event 0x10ae7b365012d4a928487ad6fce7c22cac52806e831e3f1b6ec16784e2409a86.
//
// Solidity: event EventApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterEventApproveCandidate(opts *bind.FilterOpts) (*NodeManagerEventApproveCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventApproveCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventApproveCandidateIterator{contract: _NodeManager.contract, event: "EventApproveCandidate", logs: logs, sub: sub}, nil
}

// WatchEventApproveCandidate is a free log subscription operation binding the contract event 0x10ae7b365012d4a928487ad6fce7c22cac52806e831e3f1b6ec16784e2409a86.
//
// Solidity: event EventApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchEventApproveCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerEventApproveCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventApproveCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventApproveCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "EventApproveCandidate", log); err != nil {
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

// ParseEventApproveCandidate is a log parse operation binding the contract event 0x10ae7b365012d4a928487ad6fce7c22cac52806e831e3f1b6ec16784e2409a86.
//
// Solidity: event EventApproveCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseEventApproveCandidate(log types.Log) (*NodeManagerEventApproveCandidate, error) {
	event := new(NodeManagerEventApproveCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "EventApproveCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventBlackNodeIterator is returned from FilterEventBlackNode and is used to iterate over the raw logs and unpacked data for EventBlackNode events raised by the NodeManager contract.
type NodeManagerEventBlackNodeIterator struct {
	Event *NodeManagerEventBlackNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventBlackNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventBlackNode)
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
		it.Event = new(NodeManagerEventBlackNode)
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
func (it *NodeManagerEventBlackNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventBlackNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventBlackNode represents a EventBlackNode event raised by the NodeManager contract.
type NodeManagerEventBlackNode struct {
	PubkeyList []string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEventBlackNode is a free log retrieval operation binding the contract event 0x82751ad26bdf58cab168fa6640b43bfe4c17c54165b7bc3bf1fbbe4325b8ce52.
//
// Solidity: event EventBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) FilterEventBlackNode(opts *bind.FilterOpts) (*NodeManagerEventBlackNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventBlackNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventBlackNodeIterator{contract: _NodeManager.contract, event: "EventBlackNode", logs: logs, sub: sub}, nil
}

// WatchEventBlackNode is a free log subscription operation binding the contract event 0x82751ad26bdf58cab168fa6640b43bfe4c17c54165b7bc3bf1fbbe4325b8ce52.
//
// Solidity: event EventBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) WatchEventBlackNode(opts *bind.WatchOpts, sink chan<- *NodeManagerEventBlackNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventBlackNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventBlackNode)
				if err := _NodeManager.contract.UnpackLog(event, "EventBlackNode", log); err != nil {
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

// ParseEventBlackNode is a log parse operation binding the contract event 0x82751ad26bdf58cab168fa6640b43bfe4c17c54165b7bc3bf1fbbe4325b8ce52.
//
// Solidity: event EventBlackNode(string[] PubkeyList)
func (_NodeManager *NodeManagerFilterer) ParseEventBlackNode(log types.Log) (*NodeManagerEventBlackNode, error) {
	event := new(NodeManagerEventBlackNode)
	if err := _NodeManager.contract.UnpackLog(event, "EventBlackNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventCommitDposIterator is returned from FilterEventCommitDpos and is used to iterate over the raw logs and unpacked data for EventCommitDpos events raised by the NodeManager contract.
type NodeManagerEventCommitDposIterator struct {
	Event *NodeManagerEventCommitDpos // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventCommitDposIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventCommitDpos)
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
		it.Event = new(NodeManagerEventCommitDpos)
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
func (it *NodeManagerEventCommitDposIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventCommitDposIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventCommitDpos represents a EventCommitDpos event raised by the NodeManager contract.
type NodeManagerEventCommitDpos struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventCommitDpos is a free log retrieval operation binding the contract event 0x2aea5e6af017cf88ef818d37c5d229c9504a3da7fdc284f6bd12974b97b4a175.
//
// Solidity: event EventCommitDpos()
func (_NodeManager *NodeManagerFilterer) FilterEventCommitDpos(opts *bind.FilterOpts) (*NodeManagerEventCommitDposIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventCommitDpos")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventCommitDposIterator{contract: _NodeManager.contract, event: "EventCommitDpos", logs: logs, sub: sub}, nil
}

// WatchEventCommitDpos is a free log subscription operation binding the contract event 0x2aea5e6af017cf88ef818d37c5d229c9504a3da7fdc284f6bd12974b97b4a175.
//
// Solidity: event EventCommitDpos()
func (_NodeManager *NodeManagerFilterer) WatchEventCommitDpos(opts *bind.WatchOpts, sink chan<- *NodeManagerEventCommitDpos) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventCommitDpos")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventCommitDpos)
				if err := _NodeManager.contract.UnpackLog(event, "EventCommitDpos", log); err != nil {
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

// ParseEventCommitDpos is a log parse operation binding the contract event 0x2aea5e6af017cf88ef818d37c5d229c9504a3da7fdc284f6bd12974b97b4a175.
//
// Solidity: event EventCommitDpos()
func (_NodeManager *NodeManagerFilterer) ParseEventCommitDpos(log types.Log) (*NodeManagerEventCommitDpos, error) {
	event := new(NodeManagerEventCommitDpos)
	if err := _NodeManager.contract.UnpackLog(event, "EventCommitDpos", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventQuitNodeIterator is returned from FilterEventQuitNode and is used to iterate over the raw logs and unpacked data for EventQuitNode events raised by the NodeManager contract.
type NodeManagerEventQuitNodeIterator struct {
	Event *NodeManagerEventQuitNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventQuitNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventQuitNode)
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
		it.Event = new(NodeManagerEventQuitNode)
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
func (it *NodeManagerEventQuitNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventQuitNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventQuitNode represents a EventQuitNode event raised by the NodeManager contract.
type NodeManagerEventQuitNode struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventQuitNode is a free log retrieval operation binding the contract event 0x19b3b65e4221ad42d479fe6f763e22ec86cdcd3a1fbafea72a3c82f921fc0dfc.
//
// Solidity: event EventQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterEventQuitNode(opts *bind.FilterOpts) (*NodeManagerEventQuitNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventQuitNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventQuitNodeIterator{contract: _NodeManager.contract, event: "EventQuitNode", logs: logs, sub: sub}, nil
}

// WatchEventQuitNode is a free log subscription operation binding the contract event 0x19b3b65e4221ad42d479fe6f763e22ec86cdcd3a1fbafea72a3c82f921fc0dfc.
//
// Solidity: event EventQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchEventQuitNode(opts *bind.WatchOpts, sink chan<- *NodeManagerEventQuitNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventQuitNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventQuitNode)
				if err := _NodeManager.contract.UnpackLog(event, "EventQuitNode", log); err != nil {
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

// ParseEventQuitNode is a log parse operation binding the contract event 0x19b3b65e4221ad42d479fe6f763e22ec86cdcd3a1fbafea72a3c82f921fc0dfc.
//
// Solidity: event EventQuitNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseEventQuitNode(log types.Log) (*NodeManagerEventQuitNode, error) {
	event := new(NodeManagerEventQuitNode)
	if err := _NodeManager.contract.UnpackLog(event, "EventQuitNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventRegisterCandidateIterator is returned from FilterEventRegisterCandidate and is used to iterate over the raw logs and unpacked data for EventRegisterCandidate events raised by the NodeManager contract.
type NodeManagerEventRegisterCandidateIterator struct {
	Event *NodeManagerEventRegisterCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventRegisterCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventRegisterCandidate)
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
		it.Event = new(NodeManagerEventRegisterCandidate)
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
func (it *NodeManagerEventRegisterCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventRegisterCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventRegisterCandidate represents a EventRegisterCandidate event raised by the NodeManager contract.
type NodeManagerEventRegisterCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventRegisterCandidate is a free log retrieval operation binding the contract event 0x2cba64bf604c9ddd7bd252d3a2fa0a801f94cfb81eca6ec45ef727db4b753d4a.
//
// Solidity: event EventRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterEventRegisterCandidate(opts *bind.FilterOpts) (*NodeManagerEventRegisterCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventRegisterCandidateIterator{contract: _NodeManager.contract, event: "EventRegisterCandidate", logs: logs, sub: sub}, nil
}

// WatchEventRegisterCandidate is a free log subscription operation binding the contract event 0x2cba64bf604c9ddd7bd252d3a2fa0a801f94cfb81eca6ec45ef727db4b753d4a.
//
// Solidity: event EventRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchEventRegisterCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerEventRegisterCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventRegisterCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "EventRegisterCandidate", log); err != nil {
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

// ParseEventRegisterCandidate is a log parse operation binding the contract event 0x2cba64bf604c9ddd7bd252d3a2fa0a801f94cfb81eca6ec45ef727db4b753d4a.
//
// Solidity: event EventRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseEventRegisterCandidate(log types.Log) (*NodeManagerEventRegisterCandidate, error) {
	event := new(NodeManagerEventRegisterCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "EventRegisterCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventUnRegisterCandidateIterator is returned from FilterEventUnRegisterCandidate and is used to iterate over the raw logs and unpacked data for EventUnRegisterCandidate events raised by the NodeManager contract.
type NodeManagerEventUnRegisterCandidateIterator struct {
	Event *NodeManagerEventUnRegisterCandidate // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventUnRegisterCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventUnRegisterCandidate)
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
		it.Event = new(NodeManagerEventUnRegisterCandidate)
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
func (it *NodeManagerEventUnRegisterCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventUnRegisterCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventUnRegisterCandidate represents a EventUnRegisterCandidate event raised by the NodeManager contract.
type NodeManagerEventUnRegisterCandidate struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventUnRegisterCandidate is a free log retrieval operation binding the contract event 0x4a71b1f4ce29ab3184ff7bcdb316efd204ec12b3a308fcb2f095d7afa969bf04.
//
// Solidity: event EventUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterEventUnRegisterCandidate(opts *bind.FilterOpts) (*NodeManagerEventUnRegisterCandidateIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventUnRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventUnRegisterCandidateIterator{contract: _NodeManager.contract, event: "EventUnRegisterCandidate", logs: logs, sub: sub}, nil
}

// WatchEventUnRegisterCandidate is a free log subscription operation binding the contract event 0x4a71b1f4ce29ab3184ff7bcdb316efd204ec12b3a308fcb2f095d7afa969bf04.
//
// Solidity: event EventUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchEventUnRegisterCandidate(opts *bind.WatchOpts, sink chan<- *NodeManagerEventUnRegisterCandidate) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventUnRegisterCandidate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventUnRegisterCandidate)
				if err := _NodeManager.contract.UnpackLog(event, "EventUnRegisterCandidate", log); err != nil {
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

// ParseEventUnRegisterCandidate is a log parse operation binding the contract event 0x4a71b1f4ce29ab3184ff7bcdb316efd204ec12b3a308fcb2f095d7afa969bf04.
//
// Solidity: event EventUnRegisterCandidate(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseEventUnRegisterCandidate(log types.Log) (*NodeManagerEventUnRegisterCandidate, error) {
	event := new(NodeManagerEventUnRegisterCandidate)
	if err := _NodeManager.contract.UnpackLog(event, "EventUnRegisterCandidate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventUpdateConfigIterator is returned from FilterEventUpdateConfig and is used to iterate over the raw logs and unpacked data for EventUpdateConfig events raised by the NodeManager contract.
type NodeManagerEventUpdateConfigIterator struct {
	Event *NodeManagerEventUpdateConfig // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventUpdateConfigIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventUpdateConfig)
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
		it.Event = new(NodeManagerEventUpdateConfig)
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
func (it *NodeManagerEventUpdateConfigIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventUpdateConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventUpdateConfig represents a EventUpdateConfig event raised by the NodeManager contract.
type NodeManagerEventUpdateConfig struct {
	Config node_managerConfiguration
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventUpdateConfig is a free log retrieval operation binding the contract event 0x75f3adbd00bdc1752246184a67d7e5bad67bfa179acf61f021a1bf68798f90e8.
//
// Solidity: event EventUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) FilterEventUpdateConfig(opts *bind.FilterOpts) (*NodeManagerEventUpdateConfigIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventUpdateConfig")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventUpdateConfigIterator{contract: _NodeManager.contract, event: "EventUpdateConfig", logs: logs, sub: sub}, nil
}

// WatchEventUpdateConfig is a free log subscription operation binding the contract event 0x75f3adbd00bdc1752246184a67d7e5bad67bfa179acf61f021a1bf68798f90e8.
//
// Solidity: event EventUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) WatchEventUpdateConfig(opts *bind.WatchOpts, sink chan<- *NodeManagerEventUpdateConfig) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventUpdateConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventUpdateConfig)
				if err := _NodeManager.contract.UnpackLog(event, "EventUpdateConfig", log); err != nil {
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

// ParseEventUpdateConfig is a log parse operation binding the contract event 0x75f3adbd00bdc1752246184a67d7e5bad67bfa179acf61f021a1bf68798f90e8.
//
// Solidity: event EventUpdateConfig((uint32,uint32,uint32,uint32) Config)
func (_NodeManager *NodeManagerFilterer) ParseEventUpdateConfig(log types.Log) (*NodeManagerEventUpdateConfig, error) {
	event := new(NodeManagerEventUpdateConfig)
	if err := _NodeManager.contract.UnpackLog(event, "EventUpdateConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeManagerEventWhiteNodeIterator is returned from FilterEventWhiteNode and is used to iterate over the raw logs and unpacked data for EventWhiteNode events raised by the NodeManager contract.
type NodeManagerEventWhiteNodeIterator struct {
	Event *NodeManagerEventWhiteNode // Event containing the contract specifics and raw log

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
func (it *NodeManagerEventWhiteNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeManagerEventWhiteNode)
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
		it.Event = new(NodeManagerEventWhiteNode)
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
func (it *NodeManagerEventWhiteNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeManagerEventWhiteNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeManagerEventWhiteNode represents a EventWhiteNode event raised by the NodeManager contract.
type NodeManagerEventWhiteNode struct {
	Pubkey string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEventWhiteNode is a free log retrieval operation binding the contract event 0xa7d39a1a4b073df4b9c4a6233bdd982e098ad340a878fd67fff1476e2a049b41.
//
// Solidity: event EventWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) FilterEventWhiteNode(opts *bind.FilterOpts) (*NodeManagerEventWhiteNodeIterator, error) {

	logs, sub, err := _NodeManager.contract.FilterLogs(opts, "EventWhiteNode")
	if err != nil {
		return nil, err
	}
	return &NodeManagerEventWhiteNodeIterator{contract: _NodeManager.contract, event: "EventWhiteNode", logs: logs, sub: sub}, nil
}

// WatchEventWhiteNode is a free log subscription operation binding the contract event 0xa7d39a1a4b073df4b9c4a6233bdd982e098ad340a878fd67fff1476e2a049b41.
//
// Solidity: event EventWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) WatchEventWhiteNode(opts *bind.WatchOpts, sink chan<- *NodeManagerEventWhiteNode) (event.Subscription, error) {

	logs, sub, err := _NodeManager.contract.WatchLogs(opts, "EventWhiteNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeManagerEventWhiteNode)
				if err := _NodeManager.contract.UnpackLog(event, "EventWhiteNode", log); err != nil {
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

// ParseEventWhiteNode is a log parse operation binding the contract event 0xa7d39a1a4b073df4b9c4a6233bdd982e098ad340a878fd67fff1476e2a049b41.
//
// Solidity: event EventWhiteNode(string Pubkey)
func (_NodeManager *NodeManagerFilterer) ParseEventWhiteNode(log types.Log) (*NodeManagerEventWhiteNode, error) {
	event := new(NodeManagerEventWhiteNode)
	if err := _NodeManager.contract.UnpackLog(event, "EventWhiteNode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
