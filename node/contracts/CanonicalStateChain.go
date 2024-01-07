// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// CanonicalStateChainHeader is an auto generated low-level Go binding around an user-defined struct.
type CanonicalStateChainHeader struct {
	Epoch            uint64
	L2Height         uint64
	PrevHash         [32]byte
	TxRoot           [32]byte
	BlockRoot        [32]byte
	StateRoot        [32]byte
	CelestiaHeight   uint64
	CelestiaDataRoot [32]byte
}

// CanonicalStateChainContractMetaData contains all meta data concerning the CanonicalStateChainContract contract.
var CanonicalStateChainContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_publisher\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"l2Height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"celestiaHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"celestiaDataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structCanonicalStateChain.Header\",\"name\":\"_header\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"BlockAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"challenge\",\"type\":\"address\"}],\"name\":\"ChallengeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"}],\"name\":\"PublisherChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"RolledBack\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"chain\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainHead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challenge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"l2Height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"celestiaHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"celestiaDataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structCanonicalStateChain.Header\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHead\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"l2Height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"celestiaHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"celestiaDataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structCanonicalStateChain.Header\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"headerMetadata\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"headers\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"l2Height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"celestiaHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"celestiaDataRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publisher\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"l2Height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"celestiaHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"celestiaDataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structCanonicalStateChain.Header\",\"name\":\"_header\",\"type\":\"tuple\"}],\"name\":\"pushBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"name\":\"rollback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_challenge\",\"type\":\"address\"}],\"name\":\"setChallengeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_publisher\",\"type\":\"address\"}],\"name\":\"setPublisher\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CanonicalStateChainContractABI is the input ABI used to generate the binding from.
// Deprecated: Use CanonicalStateChainContractMetaData.ABI instead.
var CanonicalStateChainContractABI = CanonicalStateChainContractMetaData.ABI

// CanonicalStateChainContract is an auto generated Go binding around an Ethereum contract.
type CanonicalStateChainContract struct {
	CanonicalStateChainContractCaller     // Read-only binding to the contract
	CanonicalStateChainContractTransactor // Write-only binding to the contract
	CanonicalStateChainContractFilterer   // Log filterer for contract events
}

// CanonicalStateChainContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type CanonicalStateChainContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CanonicalStateChainContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CanonicalStateChainContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CanonicalStateChainContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CanonicalStateChainContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CanonicalStateChainContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CanonicalStateChainContractSession struct {
	Contract     *CanonicalStateChainContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// CanonicalStateChainContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CanonicalStateChainContractCallerSession struct {
	Contract *CanonicalStateChainContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// CanonicalStateChainContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CanonicalStateChainContractTransactorSession struct {
	Contract     *CanonicalStateChainContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// CanonicalStateChainContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type CanonicalStateChainContractRaw struct {
	Contract *CanonicalStateChainContract // Generic contract binding to access the raw methods on
}

// CanonicalStateChainContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CanonicalStateChainContractCallerRaw struct {
	Contract *CanonicalStateChainContractCaller // Generic read-only contract binding to access the raw methods on
}

// CanonicalStateChainContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CanonicalStateChainContractTransactorRaw struct {
	Contract *CanonicalStateChainContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCanonicalStateChainContract creates a new instance of CanonicalStateChainContract, bound to a specific deployed contract.
func NewCanonicalStateChainContract(address common.Address, backend bind.ContractBackend) (*CanonicalStateChainContract, error) {
	contract, err := bindCanonicalStateChainContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContract{CanonicalStateChainContractCaller: CanonicalStateChainContractCaller{contract: contract}, CanonicalStateChainContractTransactor: CanonicalStateChainContractTransactor{contract: contract}, CanonicalStateChainContractFilterer: CanonicalStateChainContractFilterer{contract: contract}}, nil
}

// NewCanonicalStateChainContractCaller creates a new read-only instance of CanonicalStateChainContract, bound to a specific deployed contract.
func NewCanonicalStateChainContractCaller(address common.Address, caller bind.ContractCaller) (*CanonicalStateChainContractCaller, error) {
	contract, err := bindCanonicalStateChainContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractCaller{contract: contract}, nil
}

// NewCanonicalStateChainContractTransactor creates a new write-only instance of CanonicalStateChainContract, bound to a specific deployed contract.
func NewCanonicalStateChainContractTransactor(address common.Address, transactor bind.ContractTransactor) (*CanonicalStateChainContractTransactor, error) {
	contract, err := bindCanonicalStateChainContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractTransactor{contract: contract}, nil
}

// NewCanonicalStateChainContractFilterer creates a new log filterer instance of CanonicalStateChainContract, bound to a specific deployed contract.
func NewCanonicalStateChainContractFilterer(address common.Address, filterer bind.ContractFilterer) (*CanonicalStateChainContractFilterer, error) {
	contract, err := bindCanonicalStateChainContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractFilterer{contract: contract}, nil
}

// bindCanonicalStateChainContract binds a generic wrapper to an already deployed contract.
func bindCanonicalStateChainContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CanonicalStateChainContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CanonicalStateChainContract *CanonicalStateChainContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CanonicalStateChainContract.Contract.CanonicalStateChainContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CanonicalStateChainContract *CanonicalStateChainContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.CanonicalStateChainContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CanonicalStateChainContract *CanonicalStateChainContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.CanonicalStateChainContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CanonicalStateChainContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.contract.Transact(opts, method, params...)
}

// Chain is a free data retrieval call binding the contract method 0x5852cc0c.
//
// Solidity: function chain(uint256 ) view returns(bytes32)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Chain(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "chain", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Chain is a free data retrieval call binding the contract method 0x5852cc0c.
//
// Solidity: function chain(uint256 ) view returns(bytes32)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Chain(arg0 *big.Int) ([32]byte, error) {
	return _CanonicalStateChainContract.Contract.Chain(&_CanonicalStateChainContract.CallOpts, arg0)
}

// Chain is a free data retrieval call binding the contract method 0x5852cc0c.
//
// Solidity: function chain(uint256 ) view returns(bytes32)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Chain(arg0 *big.Int) ([32]byte, error) {
	return _CanonicalStateChainContract.Contract.Chain(&_CanonicalStateChainContract.CallOpts, arg0)
}

// ChainHead is a free data retrieval call binding the contract method 0x008f51c6.
//
// Solidity: function chainHead() view returns(uint256)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) ChainHead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "chainHead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChainHead is a free data retrieval call binding the contract method 0x008f51c6.
//
// Solidity: function chainHead() view returns(uint256)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) ChainHead() (*big.Int, error) {
	return _CanonicalStateChainContract.Contract.ChainHead(&_CanonicalStateChainContract.CallOpts)
}

// ChainHead is a free data retrieval call binding the contract method 0x008f51c6.
//
// Solidity: function chainHead() view returns(uint256)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) ChainHead() (*big.Int, error) {
	return _CanonicalStateChainContract.Contract.ChainHead(&_CanonicalStateChainContract.CallOpts)
}

// Challenge is a free data retrieval call binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Challenge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "challenge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Challenge is a free data retrieval call binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Challenge() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Challenge(&_CanonicalStateChainContract.CallOpts)
}

// Challenge is a free data retrieval call binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Challenge() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Challenge(&_CanonicalStateChainContract.CallOpts)
}

// GetBlock is a free data retrieval call binding the contract method 0x04c07569.
//
// Solidity: function getBlock(uint256 _index) view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) GetBlock(opts *bind.CallOpts, _index *big.Int) (CanonicalStateChainHeader, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "getBlock", _index)

	if err != nil {
		return *new(CanonicalStateChainHeader), err
	}

	out0 := *abi.ConvertType(out[0], new(CanonicalStateChainHeader)).(*CanonicalStateChainHeader)

	return out0, err

}

// GetBlock is a free data retrieval call binding the contract method 0x04c07569.
//
// Solidity: function getBlock(uint256 _index) view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) GetBlock(_index *big.Int) (CanonicalStateChainHeader, error) {
	return _CanonicalStateChainContract.Contract.GetBlock(&_CanonicalStateChainContract.CallOpts, _index)
}

// GetBlock is a free data retrieval call binding the contract method 0x04c07569.
//
// Solidity: function getBlock(uint256 _index) view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) GetBlock(_index *big.Int) (CanonicalStateChainHeader, error) {
	return _CanonicalStateChainContract.Contract.GetBlock(&_CanonicalStateChainContract.CallOpts, _index)
}

// GetHead is a free data retrieval call binding the contract method 0xdc281aff.
//
// Solidity: function getHead() view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) GetHead(opts *bind.CallOpts) (CanonicalStateChainHeader, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "getHead")

	if err != nil {
		return *new(CanonicalStateChainHeader), err
	}

	out0 := *abi.ConvertType(out[0], new(CanonicalStateChainHeader)).(*CanonicalStateChainHeader)

	return out0, err

}

// GetHead is a free data retrieval call binding the contract method 0xdc281aff.
//
// Solidity: function getHead() view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) GetHead() (CanonicalStateChainHeader, error) {
	return _CanonicalStateChainContract.Contract.GetHead(&_CanonicalStateChainContract.CallOpts)
}

// GetHead is a free data retrieval call binding the contract method 0xdc281aff.
//
// Solidity: function getHead() view returns((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32))
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) GetHead() (CanonicalStateChainHeader, error) {
	return _CanonicalStateChainContract.Contract.GetHead(&_CanonicalStateChainContract.CallOpts)
}

// HeaderMetadata is a free data retrieval call binding the contract method 0x28a8d0e4.
//
// Solidity: function headerMetadata(bytes32 ) view returns(uint64 timestamp, address publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) HeaderMetadata(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Timestamp uint64
	Publisher common.Address
}, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "headerMetadata", arg0)

	outstruct := new(struct {
		Timestamp uint64
		Publisher common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Publisher = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// HeaderMetadata is a free data retrieval call binding the contract method 0x28a8d0e4.
//
// Solidity: function headerMetadata(bytes32 ) view returns(uint64 timestamp, address publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) HeaderMetadata(arg0 [32]byte) (struct {
	Timestamp uint64
	Publisher common.Address
}, error) {
	return _CanonicalStateChainContract.Contract.HeaderMetadata(&_CanonicalStateChainContract.CallOpts, arg0)
}

// HeaderMetadata is a free data retrieval call binding the contract method 0x28a8d0e4.
//
// Solidity: function headerMetadata(bytes32 ) view returns(uint64 timestamp, address publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) HeaderMetadata(arg0 [32]byte) (struct {
	Timestamp uint64
	Publisher common.Address
}, error) {
	return _CanonicalStateChainContract.Contract.HeaderMetadata(&_CanonicalStateChainContract.CallOpts, arg0)
}

// Headers is a free data retrieval call binding the contract method 0x9e7f2700.
//
// Solidity: function headers(bytes32 ) view returns(uint64 epoch, uint64 l2Height, bytes32 prevHash, bytes32 txRoot, bytes32 blockRoot, bytes32 stateRoot, uint64 celestiaHeight, bytes32 celestiaDataRoot)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Headers(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Epoch            uint64
	L2Height         uint64
	PrevHash         [32]byte
	TxRoot           [32]byte
	BlockRoot        [32]byte
	StateRoot        [32]byte
	CelestiaHeight   uint64
	CelestiaDataRoot [32]byte
}, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "headers", arg0)

	outstruct := new(struct {
		Epoch            uint64
		L2Height         uint64
		PrevHash         [32]byte
		TxRoot           [32]byte
		BlockRoot        [32]byte
		StateRoot        [32]byte
		CelestiaHeight   uint64
		CelestiaDataRoot [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Epoch = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.L2Height = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.PrevHash = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.TxRoot = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.BlockRoot = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.StateRoot = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.CelestiaHeight = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.CelestiaDataRoot = *abi.ConvertType(out[7], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Headers is a free data retrieval call binding the contract method 0x9e7f2700.
//
// Solidity: function headers(bytes32 ) view returns(uint64 epoch, uint64 l2Height, bytes32 prevHash, bytes32 txRoot, bytes32 blockRoot, bytes32 stateRoot, uint64 celestiaHeight, bytes32 celestiaDataRoot)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Headers(arg0 [32]byte) (struct {
	Epoch            uint64
	L2Height         uint64
	PrevHash         [32]byte
	TxRoot           [32]byte
	BlockRoot        [32]byte
	StateRoot        [32]byte
	CelestiaHeight   uint64
	CelestiaDataRoot [32]byte
}, error) {
	return _CanonicalStateChainContract.Contract.Headers(&_CanonicalStateChainContract.CallOpts, arg0)
}

// Headers is a free data retrieval call binding the contract method 0x9e7f2700.
//
// Solidity: function headers(bytes32 ) view returns(uint64 epoch, uint64 l2Height, bytes32 prevHash, bytes32 txRoot, bytes32 blockRoot, bytes32 stateRoot, uint64 celestiaHeight, bytes32 celestiaDataRoot)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Headers(arg0 [32]byte) (struct {
	Epoch            uint64
	L2Height         uint64
	PrevHash         [32]byte
	TxRoot           [32]byte
	BlockRoot        [32]byte
	StateRoot        [32]byte
	CelestiaHeight   uint64
	CelestiaDataRoot [32]byte
}, error) {
	return _CanonicalStateChainContract.Contract.Headers(&_CanonicalStateChainContract.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Owner() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Owner(&_CanonicalStateChainContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Owner() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Owner(&_CanonicalStateChainContract.CallOpts)
}

// Publisher is a free data retrieval call binding the contract method 0x8c72c54e.
//
// Solidity: function publisher() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Publisher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "publisher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Publisher is a free data retrieval call binding the contract method 0x8c72c54e.
//
// Solidity: function publisher() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Publisher() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Publisher(&_CanonicalStateChainContract.CallOpts)
}

// Publisher is a free data retrieval call binding the contract method 0x8c72c54e.
//
// Solidity: function publisher() view returns(address)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Publisher() (common.Address, error) {
	return _CanonicalStateChainContract.Contract.Publisher(&_CanonicalStateChainContract.CallOpts)
}

// Timestamps is a free data retrieval call binding the contract method 0xb5872958.
//
// Solidity: function timestamps(bytes32 ) view returns(uint64)
func (_CanonicalStateChainContract *CanonicalStateChainContractCaller) Timestamps(opts *bind.CallOpts, arg0 [32]byte) (uint64, error) {
	var out []interface{}
	err := _CanonicalStateChainContract.contract.Call(opts, &out, "timestamps", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Timestamps is a free data retrieval call binding the contract method 0xb5872958.
//
// Solidity: function timestamps(bytes32 ) view returns(uint64)
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Timestamps(arg0 [32]byte) (uint64, error) {
	return _CanonicalStateChainContract.Contract.Timestamps(&_CanonicalStateChainContract.CallOpts, arg0)
}

// Timestamps is a free data retrieval call binding the contract method 0xb5872958.
//
// Solidity: function timestamps(bytes32 ) view returns(uint64)
func (_CanonicalStateChainContract *CanonicalStateChainContractCallerSession) Timestamps(arg0 [32]byte) (uint64, error) {
	return _CanonicalStateChainContract.Contract.Timestamps(&_CanonicalStateChainContract.CallOpts, arg0)
}

// PushBlock is a paid mutator transaction binding the contract method 0xce71d01b.
//
// Solidity: function pushBlock((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32) _header) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) PushBlock(opts *bind.TransactOpts, _header CanonicalStateChainHeader) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "pushBlock", _header)
}

// PushBlock is a paid mutator transaction binding the contract method 0xce71d01b.
//
// Solidity: function pushBlock((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32) _header) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) PushBlock(_header CanonicalStateChainHeader) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.PushBlock(&_CanonicalStateChainContract.TransactOpts, _header)
}

// PushBlock is a paid mutator transaction binding the contract method 0xce71d01b.
//
// Solidity: function pushBlock((uint64,uint64,bytes32,bytes32,bytes32,bytes32,uint64,bytes32) _header) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) PushBlock(_header CanonicalStateChainHeader) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.PushBlock(&_CanonicalStateChainContract.TransactOpts, _header)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.RenounceOwnership(&_CanonicalStateChainContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.RenounceOwnership(&_CanonicalStateChainContract.TransactOpts)
}

// Rollback is a paid mutator transaction binding the contract method 0x0da9da20.
//
// Solidity: function rollback(uint256 _blockNumber) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) Rollback(opts *bind.TransactOpts, _blockNumber *big.Int) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "rollback", _blockNumber)
}

// Rollback is a paid mutator transaction binding the contract method 0x0da9da20.
//
// Solidity: function rollback(uint256 _blockNumber) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) Rollback(_blockNumber *big.Int) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.Rollback(&_CanonicalStateChainContract.TransactOpts, _blockNumber)
}

// Rollback is a paid mutator transaction binding the contract method 0x0da9da20.
//
// Solidity: function rollback(uint256 _blockNumber) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) Rollback(_blockNumber *big.Int) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.Rollback(&_CanonicalStateChainContract.TransactOpts, _blockNumber)
}

// SetChallengeContract is a paid mutator transaction binding the contract method 0xb37256b9.
//
// Solidity: function setChallengeContract(address _challenge) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) SetChallengeContract(opts *bind.TransactOpts, _challenge common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "setChallengeContract", _challenge)
}

// SetChallengeContract is a paid mutator transaction binding the contract method 0xb37256b9.
//
// Solidity: function setChallengeContract(address _challenge) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) SetChallengeContract(_challenge common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.SetChallengeContract(&_CanonicalStateChainContract.TransactOpts, _challenge)
}

// SetChallengeContract is a paid mutator transaction binding the contract method 0xb37256b9.
//
// Solidity: function setChallengeContract(address _challenge) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) SetChallengeContract(_challenge common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.SetChallengeContract(&_CanonicalStateChainContract.TransactOpts, _challenge)
}

// SetPublisher is a paid mutator transaction binding the contract method 0xcab63661.
//
// Solidity: function setPublisher(address _publisher) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) SetPublisher(opts *bind.TransactOpts, _publisher common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "setPublisher", _publisher)
}

// SetPublisher is a paid mutator transaction binding the contract method 0xcab63661.
//
// Solidity: function setPublisher(address _publisher) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) SetPublisher(_publisher common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.SetPublisher(&_CanonicalStateChainContract.TransactOpts, _publisher)
}

// SetPublisher is a paid mutator transaction binding the contract method 0xcab63661.
//
// Solidity: function setPublisher(address _publisher) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) SetPublisher(_publisher common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.SetPublisher(&_CanonicalStateChainContract.TransactOpts, _publisher)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.TransferOwnership(&_CanonicalStateChainContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CanonicalStateChainContract *CanonicalStateChainContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CanonicalStateChainContract.Contract.TransferOwnership(&_CanonicalStateChainContract.TransactOpts, newOwner)
}

// CanonicalStateChainContractBlockAddedIterator is returned from FilterBlockAdded and is used to iterate over the raw logs and unpacked data for BlockAdded events raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractBlockAddedIterator struct {
	Event *CanonicalStateChainContractBlockAdded // Event containing the contract specifics and raw log

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
func (it *CanonicalStateChainContractBlockAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CanonicalStateChainContractBlockAdded)
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
		it.Event = new(CanonicalStateChainContractBlockAdded)
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
func (it *CanonicalStateChainContractBlockAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CanonicalStateChainContractBlockAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CanonicalStateChainContractBlockAdded represents a BlockAdded event raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractBlockAdded struct {
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBlockAdded is a free log retrieval operation binding the contract event 0xa37f9fb2f8e66e6e5746e84c33d55fc62d920182d22358f2adc6855d3ac4d437.
//
// Solidity: event BlockAdded(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) FilterBlockAdded(opts *bind.FilterOpts, blockNumber []*big.Int) (*CanonicalStateChainContractBlockAddedIterator, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.FilterLogs(opts, "BlockAdded", blockNumberRule)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractBlockAddedIterator{contract: _CanonicalStateChainContract.contract, event: "BlockAdded", logs: logs, sub: sub}, nil
}

// WatchBlockAdded is a free log subscription operation binding the contract event 0xa37f9fb2f8e66e6e5746e84c33d55fc62d920182d22358f2adc6855d3ac4d437.
//
// Solidity: event BlockAdded(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) WatchBlockAdded(opts *bind.WatchOpts, sink chan<- *CanonicalStateChainContractBlockAdded, blockNumber []*big.Int) (event.Subscription, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.WatchLogs(opts, "BlockAdded", blockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CanonicalStateChainContractBlockAdded)
				if err := _CanonicalStateChainContract.contract.UnpackLog(event, "BlockAdded", log); err != nil {
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

// ParseBlockAdded is a log parse operation binding the contract event 0xa37f9fb2f8e66e6e5746e84c33d55fc62d920182d22358f2adc6855d3ac4d437.
//
// Solidity: event BlockAdded(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) ParseBlockAdded(log types.Log) (*CanonicalStateChainContractBlockAdded, error) {
	event := new(CanonicalStateChainContractBlockAdded)
	if err := _CanonicalStateChainContract.contract.UnpackLog(event, "BlockAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CanonicalStateChainContractChallengeChangedIterator is returned from FilterChallengeChanged and is used to iterate over the raw logs and unpacked data for ChallengeChanged events raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractChallengeChangedIterator struct {
	Event *CanonicalStateChainContractChallengeChanged // Event containing the contract specifics and raw log

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
func (it *CanonicalStateChainContractChallengeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CanonicalStateChainContractChallengeChanged)
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
		it.Event = new(CanonicalStateChainContractChallengeChanged)
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
func (it *CanonicalStateChainContractChallengeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CanonicalStateChainContractChallengeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CanonicalStateChainContractChallengeChanged represents a ChallengeChanged event raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractChallengeChanged struct {
	Challenge common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChallengeChanged is a free log retrieval operation binding the contract event 0xe06eac444661557e3ac16a5251a66b82c3f985c3e3b15eac7ea4b4fac6eeac2c.
//
// Solidity: event ChallengeChanged(address indexed challenge)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) FilterChallengeChanged(opts *bind.FilterOpts, challenge []common.Address) (*CanonicalStateChainContractChallengeChangedIterator, error) {

	var challengeRule []interface{}
	for _, challengeItem := range challenge {
		challengeRule = append(challengeRule, challengeItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.FilterLogs(opts, "ChallengeChanged", challengeRule)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractChallengeChangedIterator{contract: _CanonicalStateChainContract.contract, event: "ChallengeChanged", logs: logs, sub: sub}, nil
}

// WatchChallengeChanged is a free log subscription operation binding the contract event 0xe06eac444661557e3ac16a5251a66b82c3f985c3e3b15eac7ea4b4fac6eeac2c.
//
// Solidity: event ChallengeChanged(address indexed challenge)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) WatchChallengeChanged(opts *bind.WatchOpts, sink chan<- *CanonicalStateChainContractChallengeChanged, challenge []common.Address) (event.Subscription, error) {

	var challengeRule []interface{}
	for _, challengeItem := range challenge {
		challengeRule = append(challengeRule, challengeItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.WatchLogs(opts, "ChallengeChanged", challengeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CanonicalStateChainContractChallengeChanged)
				if err := _CanonicalStateChainContract.contract.UnpackLog(event, "ChallengeChanged", log); err != nil {
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

// ParseChallengeChanged is a log parse operation binding the contract event 0xe06eac444661557e3ac16a5251a66b82c3f985c3e3b15eac7ea4b4fac6eeac2c.
//
// Solidity: event ChallengeChanged(address indexed challenge)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) ParseChallengeChanged(log types.Log) (*CanonicalStateChainContractChallengeChanged, error) {
	event := new(CanonicalStateChainContractChallengeChanged)
	if err := _CanonicalStateChainContract.contract.UnpackLog(event, "ChallengeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CanonicalStateChainContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractOwnershipTransferredIterator struct {
	Event *CanonicalStateChainContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CanonicalStateChainContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CanonicalStateChainContractOwnershipTransferred)
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
		it.Event = new(CanonicalStateChainContractOwnershipTransferred)
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
func (it *CanonicalStateChainContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CanonicalStateChainContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CanonicalStateChainContractOwnershipTransferred represents a OwnershipTransferred event raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CanonicalStateChainContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractOwnershipTransferredIterator{contract: _CanonicalStateChainContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CanonicalStateChainContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CanonicalStateChainContractOwnershipTransferred)
				if err := _CanonicalStateChainContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) ParseOwnershipTransferred(log types.Log) (*CanonicalStateChainContractOwnershipTransferred, error) {
	event := new(CanonicalStateChainContractOwnershipTransferred)
	if err := _CanonicalStateChainContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CanonicalStateChainContractPublisherChangedIterator is returned from FilterPublisherChanged and is used to iterate over the raw logs and unpacked data for PublisherChanged events raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractPublisherChangedIterator struct {
	Event *CanonicalStateChainContractPublisherChanged // Event containing the contract specifics and raw log

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
func (it *CanonicalStateChainContractPublisherChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CanonicalStateChainContractPublisherChanged)
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
		it.Event = new(CanonicalStateChainContractPublisherChanged)
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
func (it *CanonicalStateChainContractPublisherChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CanonicalStateChainContractPublisherChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CanonicalStateChainContractPublisherChanged represents a PublisherChanged event raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractPublisherChanged struct {
	Publisher common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPublisherChanged is a free log retrieval operation binding the contract event 0x55eb99d77b0e1ed261c0a8d11f026f811b8af01455a2b45189bcc87b93dfbbb7.
//
// Solidity: event PublisherChanged(address indexed publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) FilterPublisherChanged(opts *bind.FilterOpts, publisher []common.Address) (*CanonicalStateChainContractPublisherChangedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.FilterLogs(opts, "PublisherChanged", publisherRule)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractPublisherChangedIterator{contract: _CanonicalStateChainContract.contract, event: "PublisherChanged", logs: logs, sub: sub}, nil
}

// WatchPublisherChanged is a free log subscription operation binding the contract event 0x55eb99d77b0e1ed261c0a8d11f026f811b8af01455a2b45189bcc87b93dfbbb7.
//
// Solidity: event PublisherChanged(address indexed publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) WatchPublisherChanged(opts *bind.WatchOpts, sink chan<- *CanonicalStateChainContractPublisherChanged, publisher []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.WatchLogs(opts, "PublisherChanged", publisherRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CanonicalStateChainContractPublisherChanged)
				if err := _CanonicalStateChainContract.contract.UnpackLog(event, "PublisherChanged", log); err != nil {
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

// ParsePublisherChanged is a log parse operation binding the contract event 0x55eb99d77b0e1ed261c0a8d11f026f811b8af01455a2b45189bcc87b93dfbbb7.
//
// Solidity: event PublisherChanged(address indexed publisher)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) ParsePublisherChanged(log types.Log) (*CanonicalStateChainContractPublisherChanged, error) {
	event := new(CanonicalStateChainContractPublisherChanged)
	if err := _CanonicalStateChainContract.contract.UnpackLog(event, "PublisherChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CanonicalStateChainContractRolledBackIterator is returned from FilterRolledBack and is used to iterate over the raw logs and unpacked data for RolledBack events raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractRolledBackIterator struct {
	Event *CanonicalStateChainContractRolledBack // Event containing the contract specifics and raw log

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
func (it *CanonicalStateChainContractRolledBackIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CanonicalStateChainContractRolledBack)
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
		it.Event = new(CanonicalStateChainContractRolledBack)
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
func (it *CanonicalStateChainContractRolledBackIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CanonicalStateChainContractRolledBackIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CanonicalStateChainContractRolledBack represents a RolledBack event raised by the CanonicalStateChainContract contract.
type CanonicalStateChainContractRolledBack struct {
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRolledBack is a free log retrieval operation binding the contract event 0xbd56a1ce5e71ef906a2c86c43372d012f8ab2422ff19bfdba9b686ac0936f86f.
//
// Solidity: event RolledBack(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) FilterRolledBack(opts *bind.FilterOpts, blockNumber []*big.Int) (*CanonicalStateChainContractRolledBackIterator, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.FilterLogs(opts, "RolledBack", blockNumberRule)
	if err != nil {
		return nil, err
	}
	return &CanonicalStateChainContractRolledBackIterator{contract: _CanonicalStateChainContract.contract, event: "RolledBack", logs: logs, sub: sub}, nil
}

// WatchRolledBack is a free log subscription operation binding the contract event 0xbd56a1ce5e71ef906a2c86c43372d012f8ab2422ff19bfdba9b686ac0936f86f.
//
// Solidity: event RolledBack(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) WatchRolledBack(opts *bind.WatchOpts, sink chan<- *CanonicalStateChainContractRolledBack, blockNumber []*big.Int) (event.Subscription, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _CanonicalStateChainContract.contract.WatchLogs(opts, "RolledBack", blockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CanonicalStateChainContractRolledBack)
				if err := _CanonicalStateChainContract.contract.UnpackLog(event, "RolledBack", log); err != nil {
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

// ParseRolledBack is a log parse operation binding the contract event 0xbd56a1ce5e71ef906a2c86c43372d012f8ab2422ff19bfdba9b686ac0936f86f.
//
// Solidity: event RolledBack(uint256 indexed blockNumber)
func (_CanonicalStateChainContract *CanonicalStateChainContractFilterer) ParseRolledBack(log types.Log) (*CanonicalStateChainContractRolledBack, error) {
	event := new(CanonicalStateChainContractRolledBack)
	if err := _CanonicalStateChainContract.contract.UnpackLog(event, "RolledBack", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
