// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chainloader

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AttestationProof is an auto generated low-level Go binding around an user-defined struct.
type AttestationProof struct {
	TupleRootNonce *big.Int
	Tuple          DataRootTuple
	Proof          BinaryMerkleProof
}

// BinaryMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type BinaryMerkleProof struct {
	SideNodes [][32]byte
	Key       *big.Int
	NumLeaves *big.Int
}

// DataRootTuple is an auto generated low-level Go binding around an user-defined struct.
type DataRootTuple struct {
	Height   *big.Int
	DataRoot [32]byte
}

// Namespace is an auto generated low-level Go binding around an user-defined struct.
type Namespace struct {
	Version [1]byte
	Id      [28]byte
}

// NamespaceMerkleMultiproof is an auto generated low-level Go binding around an user-defined struct.
type NamespaceMerkleMultiproof struct {
	BeginKey  *big.Int
	EndKey    *big.Int
	SideNodes []NamespaceNode
}

// NamespaceNode is an auto generated low-level Go binding around an user-defined struct.
type NamespaceNode struct {
	Min    Namespace
	Max    Namespace
	Digest [32]byte
}

// SharesProof is an auto generated low-level Go binding around an user-defined struct.
type SharesProof struct {
	Data             [][]byte
	ShareProofs      []NamespaceMerkleMultiproof
	Namespace        Namespace
	RowRoots         []NamespaceNode
	RowProofs        []BinaryMerkleProof
	AttestationProof AttestationProof
}

// ChainLoaderMetaData contains all meta data concerning the ChainLoader contract.
var ChainLoaderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_canonicalStateChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_daOracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_rblock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"_shareData\",\"type\":\"bytes[]\"}],\"name\":\"ShareKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canonicalStateChain\",\"outputs\":[{\"internalType\":\"contractICanonicalStateChain\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daOracle\",\"outputs\":[{\"internalType\":\"contractIDAOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_rblock\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"beginKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endKey\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes1\",\"name\":\"version\",\"type\":\"bytes1\"},{\"internalType\":\"bytes28\",\"name\":\"id\",\"type\":\"bytes28\"}],\"internalType\":\"structNamespace\",\"name\":\"min\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes1\",\"name\":\"version\",\"type\":\"bytes1\"},{\"internalType\":\"bytes28\",\"name\":\"id\",\"type\":\"bytes28\"}],\"internalType\":\"structNamespace\",\"name\":\"max\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structNamespaceNode[]\",\"name\":\"sideNodes\",\"type\":\"tuple[]\"}],\"internalType\":\"structNamespaceMerkleMultiproof[]\",\"name\":\"shareProofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes1\",\"name\":\"version\",\"type\":\"bytes1\"},{\"internalType\":\"bytes28\",\"name\":\"id\",\"type\":\"bytes28\"}],\"internalType\":\"structNamespace\",\"name\":\"namespace\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes1\",\"name\":\"version\",\"type\":\"bytes1\"},{\"internalType\":\"bytes28\",\"name\":\"id\",\"type\":\"bytes28\"}],\"internalType\":\"structNamespace\",\"name\":\"min\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes1\",\"name\":\"version\",\"type\":\"bytes1\"},{\"internalType\":\"bytes28\",\"name\":\"id\",\"type\":\"bytes28\"}],\"internalType\":\"structNamespace\",\"name\":\"max\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structNamespaceNode[]\",\"name\":\"rowRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"sideNodes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numLeaves\",\"type\":\"uint256\"}],\"internalType\":\"structBinaryMerkleProof[]\",\"name\":\"rowProofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tupleRootNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structDataRootTuple\",\"name\":\"tuple\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"sideNodes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numLeaves\",\"type\":\"uint256\"}],\"internalType\":\"structBinaryMerkleProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"internalType\":\"structAttestationProof\",\"name\":\"attestationProof\",\"type\":\"tuple\"}],\"internalType\":\"structSharesProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"loadShares\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"shares\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ChainLoaderABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainLoaderMetaData.ABI instead.
var ChainLoaderABI = ChainLoaderMetaData.ABI

// ChainLoader is an auto generated Go binding around an Ethereum contract.
type ChainLoader struct {
	ChainLoaderCaller     // Read-only binding to the contract
	ChainLoaderTransactor // Write-only binding to the contract
	ChainLoaderFilterer   // Log filterer for contract events
}

// ChainLoaderCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainLoaderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainLoaderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainLoaderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainLoaderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainLoaderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainLoaderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainLoaderSession struct {
	Contract     *ChainLoader      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainLoaderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainLoaderCallerSession struct {
	Contract *ChainLoaderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ChainLoaderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainLoaderTransactorSession struct {
	Contract     *ChainLoaderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ChainLoaderRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainLoaderRaw struct {
	Contract *ChainLoader // Generic contract binding to access the raw methods on
}

// ChainLoaderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainLoaderCallerRaw struct {
	Contract *ChainLoaderCaller // Generic read-only contract binding to access the raw methods on
}

// ChainLoaderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainLoaderTransactorRaw struct {
	Contract *ChainLoaderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainLoader creates a new instance of ChainLoader, bound to a specific deployed contract.
func NewChainLoader(address common.Address, backend bind.ContractBackend) (*ChainLoader, error) {
	contract, err := bindChainLoader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainLoader{ChainLoaderCaller: ChainLoaderCaller{contract: contract}, ChainLoaderTransactor: ChainLoaderTransactor{contract: contract}, ChainLoaderFilterer: ChainLoaderFilterer{contract: contract}}, nil
}

// NewChainLoaderCaller creates a new read-only instance of ChainLoader, bound to a specific deployed contract.
func NewChainLoaderCaller(address common.Address, caller bind.ContractCaller) (*ChainLoaderCaller, error) {
	contract, err := bindChainLoader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainLoaderCaller{contract: contract}, nil
}

// NewChainLoaderTransactor creates a new write-only instance of ChainLoader, bound to a specific deployed contract.
func NewChainLoaderTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainLoaderTransactor, error) {
	contract, err := bindChainLoader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainLoaderTransactor{contract: contract}, nil
}

// NewChainLoaderFilterer creates a new log filterer instance of ChainLoader, bound to a specific deployed contract.
func NewChainLoaderFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainLoaderFilterer, error) {
	contract, err := bindChainLoader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainLoaderFilterer{contract: contract}, nil
}

// bindChainLoader binds a generic wrapper to an already deployed contract.
func bindChainLoader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainLoaderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainLoader *ChainLoaderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainLoader.Contract.ChainLoaderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainLoader *ChainLoaderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainLoader.Contract.ChainLoaderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainLoader *ChainLoaderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainLoader.Contract.ChainLoaderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainLoader *ChainLoaderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainLoader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainLoader *ChainLoaderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainLoader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainLoader *ChainLoaderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainLoader.Contract.contract.Transact(opts, method, params...)
}

// ShareKey is a free data retrieval call binding the contract method 0x40a2ab6e.
//
// Solidity: function ShareKey(bytes32 _rblock, bytes[] _shareData) pure returns(bytes32)
func (_ChainLoader *ChainLoaderCaller) ShareKey(opts *bind.CallOpts, _rblock [32]byte, _shareData [][]byte) ([32]byte, error) {
	var out []interface{}
	err := _ChainLoader.contract.Call(opts, &out, "ShareKey", _rblock, _shareData)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ShareKey is a free data retrieval call binding the contract method 0x40a2ab6e.
//
// Solidity: function ShareKey(bytes32 _rblock, bytes[] _shareData) pure returns(bytes32)
func (_ChainLoader *ChainLoaderSession) ShareKey(_rblock [32]byte, _shareData [][]byte) ([32]byte, error) {
	return _ChainLoader.Contract.ShareKey(&_ChainLoader.CallOpts, _rblock, _shareData)
}

// ShareKey is a free data retrieval call binding the contract method 0x40a2ab6e.
//
// Solidity: function ShareKey(bytes32 _rblock, bytes[] _shareData) pure returns(bytes32)
func (_ChainLoader *ChainLoaderCallerSession) ShareKey(_rblock [32]byte, _shareData [][]byte) ([32]byte, error) {
	return _ChainLoader.Contract.ShareKey(&_ChainLoader.CallOpts, _rblock, _shareData)
}

// CanonicalStateChain is a free data retrieval call binding the contract method 0x8c69fa5d.
//
// Solidity: function canonicalStateChain() view returns(address)
func (_ChainLoader *ChainLoaderCaller) CanonicalStateChain(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChainLoader.contract.Call(opts, &out, "canonicalStateChain")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CanonicalStateChain is a free data retrieval call binding the contract method 0x8c69fa5d.
//
// Solidity: function canonicalStateChain() view returns(address)
func (_ChainLoader *ChainLoaderSession) CanonicalStateChain() (common.Address, error) {
	return _ChainLoader.Contract.CanonicalStateChain(&_ChainLoader.CallOpts)
}

// CanonicalStateChain is a free data retrieval call binding the contract method 0x8c69fa5d.
//
// Solidity: function canonicalStateChain() view returns(address)
func (_ChainLoader *ChainLoaderCallerSession) CanonicalStateChain() (common.Address, error) {
	return _ChainLoader.Contract.CanonicalStateChain(&_ChainLoader.CallOpts)
}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChainLoader *ChainLoaderCaller) DaOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChainLoader.contract.Call(opts, &out, "daOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChainLoader *ChainLoaderSession) DaOracle() (common.Address, error) {
	return _ChainLoader.Contract.DaOracle(&_ChainLoader.CallOpts)
}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChainLoader *ChainLoaderCallerSession) DaOracle() (common.Address, error) {
	return _ChainLoader.Contract.DaOracle(&_ChainLoader.CallOpts)
}

// Shares is a free data retrieval call binding the contract method 0x263d5f11.
//
// Solidity: function shares(bytes32 , uint256 ) view returns(bytes)
func (_ChainLoader *ChainLoaderCaller) Shares(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _ChainLoader.contract.Call(opts, &out, "shares", arg0, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shares is a free data retrieval call binding the contract method 0x263d5f11.
//
// Solidity: function shares(bytes32 , uint256 ) view returns(bytes)
func (_ChainLoader *ChainLoaderSession) Shares(arg0 [32]byte, arg1 *big.Int) ([]byte, error) {
	return _ChainLoader.Contract.Shares(&_ChainLoader.CallOpts, arg0, arg1)
}

// Shares is a free data retrieval call binding the contract method 0x263d5f11.
//
// Solidity: function shares(bytes32 , uint256 ) view returns(bytes)
func (_ChainLoader *ChainLoaderCallerSession) Shares(arg0 [32]byte, arg1 *big.Int) ([]byte, error) {
	return _ChainLoader.Contract.Shares(&_ChainLoader.CallOpts, arg0, arg1)
}

// LoadShares is a paid mutator transaction binding the contract method 0xde37a2b0.
//
// Solidity: function loadShares(bytes32 _rblock, (bytes[],(uint256,uint256,((bytes1,bytes28),(bytes1,bytes28),bytes32)[])[],(bytes1,bytes28),((bytes1,bytes28),(bytes1,bytes28),bytes32)[],(bytes32[],uint256,uint256)[],(uint256,(uint256,bytes32),(bytes32[],uint256,uint256))) _proof) returns(bytes32)
func (_ChainLoader *ChainLoaderTransactor) LoadShares(opts *bind.TransactOpts, _rblock [32]byte, _proof SharesProof) (*types.Transaction, error) {
	return _ChainLoader.contract.Transact(opts, "loadShares", _rblock, _proof)
}

// LoadShares is a paid mutator transaction binding the contract method 0xde37a2b0.
//
// Solidity: function loadShares(bytes32 _rblock, (bytes[],(uint256,uint256,((bytes1,bytes28),(bytes1,bytes28),bytes32)[])[],(bytes1,bytes28),((bytes1,bytes28),(bytes1,bytes28),bytes32)[],(bytes32[],uint256,uint256)[],(uint256,(uint256,bytes32),(bytes32[],uint256,uint256))) _proof) returns(bytes32)
func (_ChainLoader *ChainLoaderSession) LoadShares(_rblock [32]byte, _proof SharesProof) (*types.Transaction, error) {
	return _ChainLoader.Contract.LoadShares(&_ChainLoader.TransactOpts, _rblock, _proof)
}

// LoadShares is a paid mutator transaction binding the contract method 0xde37a2b0.
//
// Solidity: function loadShares(bytes32 _rblock, (bytes[],(uint256,uint256,((bytes1,bytes28),(bytes1,bytes28),bytes32)[])[],(bytes1,bytes28),((bytes1,bytes28),(bytes1,bytes28),bytes32)[],(bytes32[],uint256,uint256)[],(uint256,(uint256,bytes32),(bytes32[],uint256,uint256))) _proof) returns(bytes32)
func (_ChainLoader *ChainLoaderTransactorSession) LoadShares(_rblock [32]byte, _proof SharesProof) (*types.Transaction, error) {
	return _ChainLoader.Contract.LoadShares(&_ChainLoader.TransactOpts, _rblock, _proof)
}
