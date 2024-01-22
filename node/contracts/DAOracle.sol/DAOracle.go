// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package daOracle

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

// DAOracleContractMetaData contains all meta data concerning the DAOracleContract contract.
var DAOracleContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tupleRootNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structDataRootTuple\",\"name\":\"_tuple\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"sideNodes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numLeaves\",\"type\":\"uint256\"}],\"internalType\":\"structBinaryMerkleProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"verifyAttestation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DAOracleContractABI is the input ABI used to generate the binding from.
// Deprecated: Use DAOracleContractMetaData.ABI instead.
var DAOracleContractABI = DAOracleContractMetaData.ABI

// DAOracleContract is an auto generated Go binding around an Ethereum contract.
type DAOracleContract struct {
	DAOracleContractCaller     // Read-only binding to the contract
	DAOracleContractTransactor // Write-only binding to the contract
	DAOracleContractFilterer   // Log filterer for contract events
}

// DAOracleContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type DAOracleContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAOracleContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DAOracleContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAOracleContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DAOracleContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAOracleContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DAOracleContractSession struct {
	Contract     *DAOracleContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DAOracleContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DAOracleContractCallerSession struct {
	Contract *DAOracleContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DAOracleContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DAOracleContractTransactorSession struct {
	Contract     *DAOracleContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DAOracleContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type DAOracleContractRaw struct {
	Contract *DAOracleContract // Generic contract binding to access the raw methods on
}

// DAOracleContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DAOracleContractCallerRaw struct {
	Contract *DAOracleContractCaller // Generic read-only contract binding to access the raw methods on
}

// DAOracleContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DAOracleContractTransactorRaw struct {
	Contract *DAOracleContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDAOracleContract creates a new instance of DAOracleContract, bound to a specific deployed contract.
func NewDAOracleContract(address common.Address, backend bind.ContractBackend) (*DAOracleContract, error) {
	contract, err := bindDAOracleContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DAOracleContract{DAOracleContractCaller: DAOracleContractCaller{contract: contract}, DAOracleContractTransactor: DAOracleContractTransactor{contract: contract}, DAOracleContractFilterer: DAOracleContractFilterer{contract: contract}}, nil
}

// NewDAOracleContractCaller creates a new read-only instance of DAOracleContract, bound to a specific deployed contract.
func NewDAOracleContractCaller(address common.Address, caller bind.ContractCaller) (*DAOracleContractCaller, error) {
	contract, err := bindDAOracleContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DAOracleContractCaller{contract: contract}, nil
}

// NewDAOracleContractTransactor creates a new write-only instance of DAOracleContract, bound to a specific deployed contract.
func NewDAOracleContractTransactor(address common.Address, transactor bind.ContractTransactor) (*DAOracleContractTransactor, error) {
	contract, err := bindDAOracleContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DAOracleContractTransactor{contract: contract}, nil
}

// NewDAOracleContractFilterer creates a new log filterer instance of DAOracleContract, bound to a specific deployed contract.
func NewDAOracleContractFilterer(address common.Address, filterer bind.ContractFilterer) (*DAOracleContractFilterer, error) {
	contract, err := bindDAOracleContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DAOracleContractFilterer{contract: contract}, nil
}

// bindDAOracleContract binds a generic wrapper to an already deployed contract.
func bindDAOracleContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DAOracleContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DAOracleContract *DAOracleContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DAOracleContract.Contract.DAOracleContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DAOracleContract *DAOracleContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAOracleContract.Contract.DAOracleContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DAOracleContract *DAOracleContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DAOracleContract.Contract.DAOracleContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DAOracleContract *DAOracleContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DAOracleContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DAOracleContract *DAOracleContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAOracleContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DAOracleContract *DAOracleContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DAOracleContract.Contract.contract.Transact(opts, method, params...)
}

// VerifyAttestation is a free data retrieval call binding the contract method 0x1f3302a9.
//
// Solidity: function verifyAttestation(uint256 _tupleRootNonce, (uint256,bytes32) _tuple, (bytes32[],uint256,uint256) _proof) view returns(bool)
func (_DAOracleContract *DAOracleContractCaller) VerifyAttestation(opts *bind.CallOpts, _tupleRootNonce *big.Int, _tuple DataRootTuple, _proof BinaryMerkleProof) (bool, error) {
	var out []interface{}
	err := _DAOracleContract.contract.Call(opts, &out, "verifyAttestation", _tupleRootNonce, _tuple, _proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyAttestation is a free data retrieval call binding the contract method 0x1f3302a9.
//
// Solidity: function verifyAttestation(uint256 _tupleRootNonce, (uint256,bytes32) _tuple, (bytes32[],uint256,uint256) _proof) view returns(bool)
func (_DAOracleContract *DAOracleContractSession) VerifyAttestation(_tupleRootNonce *big.Int, _tuple DataRootTuple, _proof BinaryMerkleProof) (bool, error) {
	return _DAOracleContract.Contract.VerifyAttestation(&_DAOracleContract.CallOpts, _tupleRootNonce, _tuple, _proof)
}

// VerifyAttestation is a free data retrieval call binding the contract method 0x1f3302a9.
//
// Solidity: function verifyAttestation(uint256 _tupleRootNonce, (uint256,bytes32) _tuple, (bytes32[],uint256,uint256) _proof) view returns(bool)
func (_DAOracleContract *DAOracleContractCallerSession) VerifyAttestation(_tupleRootNonce *big.Int, _tuple DataRootTuple, _proof BinaryMerkleProof) (bool, error) {
	return _DAOracleContract.Contract.VerifyAttestation(&_DAOracleContract.CallOpts, _tupleRootNonce, _tuple, _proof)
}
