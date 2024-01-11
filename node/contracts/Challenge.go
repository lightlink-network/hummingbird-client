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

type ChallengeDaInfo struct{
	BlockIndex *big.Int `pretty:"Block Index"`
	Challenger string `pretty:"Challenger"`
	Expiry     *big.Int `pretty:"Expiry"`
	Status     uint8 	`pretty:"Status"`
}

// ChallengeDataAvailabilityChallengeDAProof is an auto generated low-level Go binding around an user-defined struct.
type ChallengeDataAvailabilityChallengeDAProof struct {
	RootNonce *big.Int
	Proof     BinaryMerkleProof
}

// ChallengeContractMetaData contains all meta data concerning the ChallengeContract contract.
var ChallengeContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_chain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_daOracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_mipsChallenge\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_blockIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"enumChallengeDataAvailability.ChallengeDAStatus\",\"name\":\"_status\",\"type\":\"uint8\"}],\"name\":\"ChallengeDAUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_blockIndex\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"enumChallengeHeader.InvalidHeaderReason\",\"name\":\"_reason\",\"type\":\"uint8\"}],\"name\":\"InvalidHeader\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"chain\",\"outputs\":[{\"internalType\":\"contractICanonicalStateChain\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blockIndex\",\"type\":\"uint256\"}],\"name\":\"challengeDataRootInclusion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengeReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengeWindow\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"daChallenges\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockIndex\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"challenger\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"enumChallengeDataAvailability.ChallengeDAStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daOracle\",\"outputs\":[{\"internalType\":\"contractIDAOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rootNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"sideNodes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numLeaves\",\"type\":\"uint256\"}],\"internalType\":\"structBinaryMerkleProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"internalType\":\"structChallengeDataAvailability.ChallengeDAProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"defendDataRootInclusion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"execChallenges\",\"outputs\":[{\"internalType\":\"enumChallengeExecution.ChallengeStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"headerIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"mipSteps\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"assertionRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"finalSystemState\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"mipsChallengeId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_headerIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_blockIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"_mipSteps\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"_assertionRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_finalSystemState\",\"type\":\"bytes32\"}],\"name\":\"initiateChallenge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blockIndex\",\"type\":\"uint256\"}],\"name\":\"invalidateHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mipsChallenge\",\"outputs\":[{\"internalType\":\"contractIMipsChallenge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_challengeFee\",\"type\":\"uint256\"}],\"name\":\"setChallengeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_challengePeriod\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_challengeReward\",\"type\":\"uint256\"}],\"name\":\"setChallengeReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_challengeWindow\",\"type\":\"uint256\"}],\"name\":\"setChallengeWindow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_defender\",\"type\":\"address\"}],\"name\":\"setDefender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_blockhash\",\"type\":\"bytes32\"}],\"name\":\"settleDataRootInclusion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"contractITreasury\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ChallengeContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ChallengeContractMetaData.ABI instead.
var ChallengeContractABI = ChallengeContractMetaData.ABI

// ChallengeContract is an auto generated Go binding around an Ethereum contract.
type ChallengeContract struct {
	ChallengeContractCaller     // Read-only binding to the contract
	ChallengeContractTransactor // Write-only binding to the contract
	ChallengeContractFilterer   // Log filterer for contract events
}

// ChallengeContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChallengeContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChallengeContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChallengeContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChallengeContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChallengeContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChallengeContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChallengeContractSession struct {
	Contract     *ChallengeContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ChallengeContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChallengeContractCallerSession struct {
	Contract *ChallengeContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ChallengeContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChallengeContractTransactorSession struct {
	Contract     *ChallengeContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ChallengeContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChallengeContractRaw struct {
	Contract *ChallengeContract // Generic contract binding to access the raw methods on
}

// ChallengeContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChallengeContractCallerRaw struct {
	Contract *ChallengeContractCaller // Generic read-only contract binding to access the raw methods on
}

// ChallengeContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChallengeContractTransactorRaw struct {
	Contract *ChallengeContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChallengeContract creates a new instance of ChallengeContract, bound to a specific deployed contract.
func NewChallengeContract(address common.Address, backend bind.ContractBackend) (*ChallengeContract, error) {
	contract, err := bindChallengeContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChallengeContract{ChallengeContractCaller: ChallengeContractCaller{contract: contract}, ChallengeContractTransactor: ChallengeContractTransactor{contract: contract}, ChallengeContractFilterer: ChallengeContractFilterer{contract: contract}}, nil
}

// NewChallengeContractCaller creates a new read-only instance of ChallengeContract, bound to a specific deployed contract.
func NewChallengeContractCaller(address common.Address, caller bind.ContractCaller) (*ChallengeContractCaller, error) {
	contract, err := bindChallengeContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractCaller{contract: contract}, nil
}

// NewChallengeContractTransactor creates a new write-only instance of ChallengeContract, bound to a specific deployed contract.
func NewChallengeContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ChallengeContractTransactor, error) {
	contract, err := bindChallengeContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractTransactor{contract: contract}, nil
}

// NewChallengeContractFilterer creates a new log filterer instance of ChallengeContract, bound to a specific deployed contract.
func NewChallengeContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ChallengeContractFilterer, error) {
	contract, err := bindChallengeContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractFilterer{contract: contract}, nil
}

// bindChallengeContract binds a generic wrapper to an already deployed contract.
func bindChallengeContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChallengeContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChallengeContract *ChallengeContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChallengeContract.Contract.ChallengeContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChallengeContract *ChallengeContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChallengeContract.Contract.ChallengeContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChallengeContract *ChallengeContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChallengeContract.Contract.ChallengeContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChallengeContract *ChallengeContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChallengeContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChallengeContract *ChallengeContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChallengeContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChallengeContract *ChallengeContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChallengeContract.Contract.contract.Transact(opts, method, params...)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) Chain(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "chain")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() view returns(address)
func (_ChallengeContract *ChallengeContractSession) Chain() (common.Address, error) {
	return _ChallengeContract.Contract.Chain(&_ChallengeContract.CallOpts)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) Chain() (common.Address, error) {
	return _ChallengeContract.Contract.Chain(&_ChallengeContract.CallOpts)
}

// ChallengeFee is a free data retrieval call binding the contract method 0x1bd8f9ca.
//
// Solidity: function challengeFee() view returns(uint256)
func (_ChallengeContract *ChallengeContractCaller) ChallengeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "challengeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChallengeFee is a free data retrieval call binding the contract method 0x1bd8f9ca.
//
// Solidity: function challengeFee() view returns(uint256)
func (_ChallengeContract *ChallengeContractSession) ChallengeFee() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeFee(&_ChallengeContract.CallOpts)
}

// ChallengeFee is a free data retrieval call binding the contract method 0x1bd8f9ca.
//
// Solidity: function challengeFee() view returns(uint256)
func (_ChallengeContract *ChallengeContractCallerSession) ChallengeFee() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeFee(&_ChallengeContract.CallOpts)
}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_ChallengeContract *ChallengeContractCaller) ChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "challengePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_ChallengeContract *ChallengeContractSession) ChallengePeriod() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengePeriod(&_ChallengeContract.CallOpts)
}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_ChallengeContract *ChallengeContractCallerSession) ChallengePeriod() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengePeriod(&_ChallengeContract.CallOpts)
}

// ChallengeReward is a free data retrieval call binding the contract method 0x3ea0c15e.
//
// Solidity: function challengeReward() view returns(uint256)
func (_ChallengeContract *ChallengeContractCaller) ChallengeReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "challengeReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChallengeReward is a free data retrieval call binding the contract method 0x3ea0c15e.
//
// Solidity: function challengeReward() view returns(uint256)
func (_ChallengeContract *ChallengeContractSession) ChallengeReward() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeReward(&_ChallengeContract.CallOpts)
}

// ChallengeReward is a free data retrieval call binding the contract method 0x3ea0c15e.
//
// Solidity: function challengeReward() view returns(uint256)
func (_ChallengeContract *ChallengeContractCallerSession) ChallengeReward() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeReward(&_ChallengeContract.CallOpts)
}

// ChallengeWindow is a free data retrieval call binding the contract method 0x861a1412.
//
// Solidity: function challengeWindow() view returns(uint256)
func (_ChallengeContract *ChallengeContractCaller) ChallengeWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "challengeWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChallengeWindow is a free data retrieval call binding the contract method 0x861a1412.
//
// Solidity: function challengeWindow() view returns(uint256)
func (_ChallengeContract *ChallengeContractSession) ChallengeWindow() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeWindow(&_ChallengeContract.CallOpts)
}

// ChallengeWindow is a free data retrieval call binding the contract method 0x861a1412.
//
// Solidity: function challengeWindow() view returns(uint256)
func (_ChallengeContract *ChallengeContractCallerSession) ChallengeWindow() (*big.Int, error) {
	return _ChallengeContract.Contract.ChallengeWindow(&_ChallengeContract.CallOpts)
}

// DaChallenges is a free data retrieval call binding the contract method 0x113e70fb.
//
// Solidity: function daChallenges(bytes32 ) view returns(uint256 blockIndex, address challenger, uint256 expiry, uint8 status)
func (_ChallengeContract *ChallengeContractCaller) DaChallenges(opts *bind.CallOpts, arg0 [32]byte) (struct {
	BlockIndex *big.Int
	Challenger common.Address
	Expiry     *big.Int
	Status     uint8
}, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "daChallenges", arg0)

	outstruct := new(struct {
		BlockIndex *big.Int
		Challenger common.Address
		Expiry     *big.Int
		Status     uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockIndex = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Challenger = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Expiry = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[3], new(uint8)).(*uint8)

	return *outstruct, err

}

// DaChallenges is a free data retrieval call binding the contract method 0x113e70fb.
//
// Solidity: function daChallenges(bytes32 ) view returns(uint256 blockIndex, address challenger, uint256 expiry, uint8 status)
func (_ChallengeContract *ChallengeContractSession) DaChallenges(arg0 [32]byte) (struct {
	BlockIndex *big.Int
	Challenger common.Address
	Expiry     *big.Int
	Status     uint8
}, error) {
	return _ChallengeContract.Contract.DaChallenges(&_ChallengeContract.CallOpts, arg0)
}

// DaChallenges is a free data retrieval call binding the contract method 0x113e70fb.
//
// Solidity: function daChallenges(bytes32 ) view returns(uint256 blockIndex, address challenger, uint256 expiry, uint8 status)
func (_ChallengeContract *ChallengeContractCallerSession) DaChallenges(arg0 [32]byte) (struct {
	BlockIndex *big.Int
	Challenger common.Address
	Expiry     *big.Int
	Status     uint8
}, error) {
	return _ChallengeContract.Contract.DaChallenges(&_ChallengeContract.CallOpts, arg0)
}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) DaOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "daOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChallengeContract *ChallengeContractSession) DaOracle() (common.Address, error) {
	return _ChallengeContract.Contract.DaOracle(&_ChallengeContract.CallOpts)
}

// DaOracle is a free data retrieval call binding the contract method 0xee223c02.
//
// Solidity: function daOracle() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) DaOracle() (common.Address, error) {
	return _ChallengeContract.Contract.DaOracle(&_ChallengeContract.CallOpts)
}

// Defender is a free data retrieval call binding the contract method 0x7f4c91c5.
//
// Solidity: function defender() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) Defender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "defender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Defender is a free data retrieval call binding the contract method 0x7f4c91c5.
//
// Solidity: function defender() view returns(address)
func (_ChallengeContract *ChallengeContractSession) Defender() (common.Address, error) {
	return _ChallengeContract.Contract.Defender(&_ChallengeContract.CallOpts)
}

// Defender is a free data retrieval call binding the contract method 0x7f4c91c5.
//
// Solidity: function defender() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) Defender() (common.Address, error) {
	return _ChallengeContract.Contract.Defender(&_ChallengeContract.CallOpts)
}

// ExecChallenges is a free data retrieval call binding the contract method 0x89f86d12.
//
// Solidity: function execChallenges(uint256 ) view returns(uint8 status, uint256 headerIndex, uint256 blockIndex, uint64 mipSteps, bytes32 assertionRoot, bytes32 finalSystemState, uint256 mipsChallengeId)
func (_ChallengeContract *ChallengeContractCaller) ExecChallenges(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Status           uint8
	HeaderIndex      *big.Int
	BlockIndex       *big.Int
	MipSteps         uint64
	AssertionRoot    [32]byte
	FinalSystemState [32]byte
	MipsChallengeId  *big.Int
}, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "execChallenges", arg0)

	outstruct := new(struct {
		Status           uint8
		HeaderIndex      *big.Int
		BlockIndex       *big.Int
		MipSteps         uint64
		AssertionRoot    [32]byte
		FinalSystemState [32]byte
		MipsChallengeId  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Status = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.HeaderIndex = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockIndex = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MipSteps = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.AssertionRoot = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.FinalSystemState = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.MipsChallengeId = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ExecChallenges is a free data retrieval call binding the contract method 0x89f86d12.
//
// Solidity: function execChallenges(uint256 ) view returns(uint8 status, uint256 headerIndex, uint256 blockIndex, uint64 mipSteps, bytes32 assertionRoot, bytes32 finalSystemState, uint256 mipsChallengeId)
func (_ChallengeContract *ChallengeContractSession) ExecChallenges(arg0 *big.Int) (struct {
	Status           uint8
	HeaderIndex      *big.Int
	BlockIndex       *big.Int
	MipSteps         uint64
	AssertionRoot    [32]byte
	FinalSystemState [32]byte
	MipsChallengeId  *big.Int
}, error) {
	return _ChallengeContract.Contract.ExecChallenges(&_ChallengeContract.CallOpts, arg0)
}

// ExecChallenges is a free data retrieval call binding the contract method 0x89f86d12.
//
// Solidity: function execChallenges(uint256 ) view returns(uint8 status, uint256 headerIndex, uint256 blockIndex, uint64 mipSteps, bytes32 assertionRoot, bytes32 finalSystemState, uint256 mipsChallengeId)
func (_ChallengeContract *ChallengeContractCallerSession) ExecChallenges(arg0 *big.Int) (struct {
	Status           uint8
	HeaderIndex      *big.Int
	BlockIndex       *big.Int
	MipSteps         uint64
	AssertionRoot    [32]byte
	FinalSystemState [32]byte
	MipsChallengeId  *big.Int
}, error) {
	return _ChallengeContract.Contract.ExecChallenges(&_ChallengeContract.CallOpts, arg0)
}

// MipsChallenge is a free data retrieval call binding the contract method 0xdca384be.
//
// Solidity: function mipsChallenge() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) MipsChallenge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "mipsChallenge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MipsChallenge is a free data retrieval call binding the contract method 0xdca384be.
//
// Solidity: function mipsChallenge() view returns(address)
func (_ChallengeContract *ChallengeContractSession) MipsChallenge() (common.Address, error) {
	return _ChallengeContract.Contract.MipsChallenge(&_ChallengeContract.CallOpts)
}

// MipsChallenge is a free data retrieval call binding the contract method 0xdca384be.
//
// Solidity: function mipsChallenge() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) MipsChallenge() (common.Address, error) {
	return _ChallengeContract.Contract.MipsChallenge(&_ChallengeContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChallengeContract *ChallengeContractSession) Owner() (common.Address, error) {
	return _ChallengeContract.Contract.Owner(&_ChallengeContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) Owner() (common.Address, error) {
	return _ChallengeContract.Contract.Owner(&_ChallengeContract.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_ChallengeContract *ChallengeContractCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChallengeContract.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_ChallengeContract *ChallengeContractSession) Treasury() (common.Address, error) {
	return _ChallengeContract.Contract.Treasury(&_ChallengeContract.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_ChallengeContract *ChallengeContractCallerSession) Treasury() (common.Address, error) {
	return _ChallengeContract.Contract.Treasury(&_ChallengeContract.CallOpts)
}

// ChallengeDataRootInclusion is a paid mutator transaction binding the contract method 0x7739f135.
//
// Solidity: function challengeDataRootInclusion(uint256 _blockIndex) payable returns(uint256)
func (_ChallengeContract *ChallengeContractTransactor) ChallengeDataRootInclusion(opts *bind.TransactOpts, _blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "challengeDataRootInclusion", _blockIndex)
}

// ChallengeDataRootInclusion is a paid mutator transaction binding the contract method 0x7739f135.
//
// Solidity: function challengeDataRootInclusion(uint256 _blockIndex) payable returns(uint256)
func (_ChallengeContract *ChallengeContractSession) ChallengeDataRootInclusion(_blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.ChallengeDataRootInclusion(&_ChallengeContract.TransactOpts, _blockIndex)
}

// ChallengeDataRootInclusion is a paid mutator transaction binding the contract method 0x7739f135.
//
// Solidity: function challengeDataRootInclusion(uint256 _blockIndex) payable returns(uint256)
func (_ChallengeContract *ChallengeContractTransactorSession) ChallengeDataRootInclusion(_blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.ChallengeDataRootInclusion(&_ChallengeContract.TransactOpts, _blockIndex)
}

// DefendDataRootInclusion is a paid mutator transaction binding the contract method 0xd898a5f4.
//
// Solidity: function defendDataRootInclusion(bytes32 _blockHash, (uint256,(bytes32[],uint256,uint256)) _proof) returns()
func (_ChallengeContract *ChallengeContractTransactor) DefendDataRootInclusion(opts *bind.TransactOpts, _blockHash [32]byte, _proof ChallengeDataAvailabilityChallengeDAProof) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "defendDataRootInclusion", _blockHash, _proof)
}

// DefendDataRootInclusion is a paid mutator transaction binding the contract method 0xd898a5f4.
//
// Solidity: function defendDataRootInclusion(bytes32 _blockHash, (uint256,(bytes32[],uint256,uint256)) _proof) returns()
func (_ChallengeContract *ChallengeContractSession) DefendDataRootInclusion(_blockHash [32]byte, _proof ChallengeDataAvailabilityChallengeDAProof) (*types.Transaction, error) {
	return _ChallengeContract.Contract.DefendDataRootInclusion(&_ChallengeContract.TransactOpts, _blockHash, _proof)
}

// DefendDataRootInclusion is a paid mutator transaction binding the contract method 0xd898a5f4.
//
// Solidity: function defendDataRootInclusion(bytes32 _blockHash, (uint256,(bytes32[],uint256,uint256)) _proof) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) DefendDataRootInclusion(_blockHash [32]byte, _proof ChallengeDataAvailabilityChallengeDAProof) (*types.Transaction, error) {
	return _ChallengeContract.Contract.DefendDataRootInclusion(&_ChallengeContract.TransactOpts, _blockHash, _proof)
}

// InitiateChallenge is a paid mutator transaction binding the contract method 0x78f142bb.
//
// Solidity: function initiateChallenge(uint256 _headerIndex, uint256 _blockIndex, uint64 _mipSteps, bytes32 _assertionRoot, bytes32 _finalSystemState) payable returns(uint256)
func (_ChallengeContract *ChallengeContractTransactor) InitiateChallenge(opts *bind.TransactOpts, _headerIndex *big.Int, _blockIndex *big.Int, _mipSteps uint64, _assertionRoot [32]byte, _finalSystemState [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "initiateChallenge", _headerIndex, _blockIndex, _mipSteps, _assertionRoot, _finalSystemState)
}

// InitiateChallenge is a paid mutator transaction binding the contract method 0x78f142bb.
//
// Solidity: function initiateChallenge(uint256 _headerIndex, uint256 _blockIndex, uint64 _mipSteps, bytes32 _assertionRoot, bytes32 _finalSystemState) payable returns(uint256)
func (_ChallengeContract *ChallengeContractSession) InitiateChallenge(_headerIndex *big.Int, _blockIndex *big.Int, _mipSteps uint64, _assertionRoot [32]byte, _finalSystemState [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.Contract.InitiateChallenge(&_ChallengeContract.TransactOpts, _headerIndex, _blockIndex, _mipSteps, _assertionRoot, _finalSystemState)
}

// InitiateChallenge is a paid mutator transaction binding the contract method 0x78f142bb.
//
// Solidity: function initiateChallenge(uint256 _headerIndex, uint256 _blockIndex, uint64 _mipSteps, bytes32 _assertionRoot, bytes32 _finalSystemState) payable returns(uint256)
func (_ChallengeContract *ChallengeContractTransactorSession) InitiateChallenge(_headerIndex *big.Int, _blockIndex *big.Int, _mipSteps uint64, _assertionRoot [32]byte, _finalSystemState [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.Contract.InitiateChallenge(&_ChallengeContract.TransactOpts, _headerIndex, _blockIndex, _mipSteps, _assertionRoot, _finalSystemState)
}

// InvalidateHeader is a paid mutator transaction binding the contract method 0x5dade412.
//
// Solidity: function invalidateHeader(uint256 _blockIndex) returns()
func (_ChallengeContract *ChallengeContractTransactor) InvalidateHeader(opts *bind.TransactOpts, _blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "invalidateHeader", _blockIndex)
}

// InvalidateHeader is a paid mutator transaction binding the contract method 0x5dade412.
//
// Solidity: function invalidateHeader(uint256 _blockIndex) returns()
func (_ChallengeContract *ChallengeContractSession) InvalidateHeader(_blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.InvalidateHeader(&_ChallengeContract.TransactOpts, _blockIndex)
}

// InvalidateHeader is a paid mutator transaction binding the contract method 0x5dade412.
//
// Solidity: function invalidateHeader(uint256 _blockIndex) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) InvalidateHeader(_blockIndex *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.InvalidateHeader(&_ChallengeContract.TransactOpts, _blockIndex)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ChallengeContract *ChallengeContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ChallengeContract *ChallengeContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _ChallengeContract.Contract.RenounceOwnership(&_ChallengeContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ChallengeContract *ChallengeContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ChallengeContract.Contract.RenounceOwnership(&_ChallengeContract.TransactOpts)
}

// SetChallengeFee is a paid mutator transaction binding the contract method 0x35bf82f6.
//
// Solidity: function setChallengeFee(uint256 _challengeFee) returns()
func (_ChallengeContract *ChallengeContractTransactor) SetChallengeFee(opts *bind.TransactOpts, _challengeFee *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "setChallengeFee", _challengeFee)
}

// SetChallengeFee is a paid mutator transaction binding the contract method 0x35bf82f6.
//
// Solidity: function setChallengeFee(uint256 _challengeFee) returns()
func (_ChallengeContract *ChallengeContractSession) SetChallengeFee(_challengeFee *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeFee(&_ChallengeContract.TransactOpts, _challengeFee)
}

// SetChallengeFee is a paid mutator transaction binding the contract method 0x35bf82f6.
//
// Solidity: function setChallengeFee(uint256 _challengeFee) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SetChallengeFee(_challengeFee *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeFee(&_ChallengeContract.TransactOpts, _challengeFee)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _challengePeriod) returns()
func (_ChallengeContract *ChallengeContractTransactor) SetChallengePeriod(opts *bind.TransactOpts, _challengePeriod *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "setChallengePeriod", _challengePeriod)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _challengePeriod) returns()
func (_ChallengeContract *ChallengeContractSession) SetChallengePeriod(_challengePeriod *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengePeriod(&_ChallengeContract.TransactOpts, _challengePeriod)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _challengePeriod) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SetChallengePeriod(_challengePeriod *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengePeriod(&_ChallengeContract.TransactOpts, _challengePeriod)
}

// SetChallengeReward is a paid mutator transaction binding the contract method 0x7d3020ad.
//
// Solidity: function setChallengeReward(uint256 _challengeReward) returns()
func (_ChallengeContract *ChallengeContractTransactor) SetChallengeReward(opts *bind.TransactOpts, _challengeReward *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "setChallengeReward", _challengeReward)
}

// SetChallengeReward is a paid mutator transaction binding the contract method 0x7d3020ad.
//
// Solidity: function setChallengeReward(uint256 _challengeReward) returns()
func (_ChallengeContract *ChallengeContractSession) SetChallengeReward(_challengeReward *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeReward(&_ChallengeContract.TransactOpts, _challengeReward)
}

// SetChallengeReward is a paid mutator transaction binding the contract method 0x7d3020ad.
//
// Solidity: function setChallengeReward(uint256 _challengeReward) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SetChallengeReward(_challengeReward *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeReward(&_ChallengeContract.TransactOpts, _challengeReward)
}

// SetChallengeWindow is a paid mutator transaction binding the contract method 0x01c1aa0d.
//
// Solidity: function setChallengeWindow(uint256 _challengeWindow) returns()
func (_ChallengeContract *ChallengeContractTransactor) SetChallengeWindow(opts *bind.TransactOpts, _challengeWindow *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "setChallengeWindow", _challengeWindow)
}

// SetChallengeWindow is a paid mutator transaction binding the contract method 0x01c1aa0d.
//
// Solidity: function setChallengeWindow(uint256 _challengeWindow) returns()
func (_ChallengeContract *ChallengeContractSession) SetChallengeWindow(_challengeWindow *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeWindow(&_ChallengeContract.TransactOpts, _challengeWindow)
}

// SetChallengeWindow is a paid mutator transaction binding the contract method 0x01c1aa0d.
//
// Solidity: function setChallengeWindow(uint256 _challengeWindow) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SetChallengeWindow(_challengeWindow *big.Int) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetChallengeWindow(&_ChallengeContract.TransactOpts, _challengeWindow)
}

// SetDefender is a paid mutator transaction binding the contract method 0x163a7177.
//
// Solidity: function setDefender(address _defender) returns()
func (_ChallengeContract *ChallengeContractTransactor) SetDefender(opts *bind.TransactOpts, _defender common.Address) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "setDefender", _defender)
}

// SetDefender is a paid mutator transaction binding the contract method 0x163a7177.
//
// Solidity: function setDefender(address _defender) returns()
func (_ChallengeContract *ChallengeContractSession) SetDefender(_defender common.Address) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetDefender(&_ChallengeContract.TransactOpts, _defender)
}

// SetDefender is a paid mutator transaction binding the contract method 0x163a7177.
//
// Solidity: function setDefender(address _defender) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SetDefender(_defender common.Address) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SetDefender(&_ChallengeContract.TransactOpts, _defender)
}

// SettleDataRootInclusion is a paid mutator transaction binding the contract method 0x5bba0ea9.
//
// Solidity: function settleDataRootInclusion(bytes32 _blockhash) returns()
func (_ChallengeContract *ChallengeContractTransactor) SettleDataRootInclusion(opts *bind.TransactOpts, _blockhash [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "settleDataRootInclusion", _blockhash)
}

// SettleDataRootInclusion is a paid mutator transaction binding the contract method 0x5bba0ea9.
//
// Solidity: function settleDataRootInclusion(bytes32 _blockhash) returns()
func (_ChallengeContract *ChallengeContractSession) SettleDataRootInclusion(_blockhash [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SettleDataRootInclusion(&_ChallengeContract.TransactOpts, _blockhash)
}

// SettleDataRootInclusion is a paid mutator transaction binding the contract method 0x5bba0ea9.
//
// Solidity: function settleDataRootInclusion(bytes32 _blockhash) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) SettleDataRootInclusion(_blockhash [32]byte) (*types.Transaction, error) {
	return _ChallengeContract.Contract.SettleDataRootInclusion(&_ChallengeContract.TransactOpts, _blockhash)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChallengeContract *ChallengeContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ChallengeContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChallengeContract *ChallengeContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ChallengeContract.Contract.TransferOwnership(&_ChallengeContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChallengeContract *ChallengeContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ChallengeContract.Contract.TransferOwnership(&_ChallengeContract.TransactOpts, newOwner)
}

// ChallengeContractChallengeDAUpdateIterator is returned from FilterChallengeDAUpdate and is used to iterate over the raw logs and unpacked data for ChallengeDAUpdate events raised by the ChallengeContract contract.
type ChallengeContractChallengeDAUpdateIterator struct {
	Event *ChallengeContractChallengeDAUpdate // Event containing the contract specifics and raw log

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
func (it *ChallengeContractChallengeDAUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChallengeContractChallengeDAUpdate)
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
		it.Event = new(ChallengeContractChallengeDAUpdate)
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
func (it *ChallengeContractChallengeDAUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChallengeContractChallengeDAUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChallengeContractChallengeDAUpdate represents a ChallengeDAUpdate event raised by the ChallengeContract contract.
type ChallengeContractChallengeDAUpdate struct {
	BlockHash  [32]byte
	BlockIndex *big.Int
	Expiry     *big.Int
	Status     uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChallengeDAUpdate is a free log retrieval operation binding the contract event 0xcdbf857f75b46d9b94fe7c990b2b0457f168aef7b089a0cb5d4e88074f2d88e8.
//
// Solidity: event ChallengeDAUpdate(bytes32 indexed _blockHash, uint256 indexed _blockIndex, uint256 _expiry, uint8 indexed _status)
func (_ChallengeContract *ChallengeContractFilterer) FilterChallengeDAUpdate(opts *bind.FilterOpts, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (*ChallengeContractChallengeDAUpdateIterator, error) {

	var _blockHashRule []interface{}
	for _, _blockHashItem := range _blockHash {
		_blockHashRule = append(_blockHashRule, _blockHashItem)
	}
	var _blockIndexRule []interface{}
	for _, _blockIndexItem := range _blockIndex {
		_blockIndexRule = append(_blockIndexRule, _blockIndexItem)
	}

	var _statusRule []interface{}
	for _, _statusItem := range _status {
		_statusRule = append(_statusRule, _statusItem)
	}

	logs, sub, err := _ChallengeContract.contract.FilterLogs(opts, "ChallengeDAUpdate", _blockHashRule, _blockIndexRule, _statusRule)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractChallengeDAUpdateIterator{contract: _ChallengeContract.contract, event: "ChallengeDAUpdate", logs: logs, sub: sub}, nil
}

// WatchChallengeDAUpdate is a free log subscription operation binding the contract event 0xcdbf857f75b46d9b94fe7c990b2b0457f168aef7b089a0cb5d4e88074f2d88e8.
//
// Solidity: event ChallengeDAUpdate(bytes32 indexed _blockHash, uint256 indexed _blockIndex, uint256 _expiry, uint8 indexed _status)
func (_ChallengeContract *ChallengeContractFilterer) WatchChallengeDAUpdate(opts *bind.WatchOpts, sink chan<- *ChallengeContractChallengeDAUpdate, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (event.Subscription, error) {

	var _blockHashRule []interface{}
	for _, _blockHashItem := range _blockHash {
		_blockHashRule = append(_blockHashRule, _blockHashItem)
	}
	var _blockIndexRule []interface{}
	for _, _blockIndexItem := range _blockIndex {
		_blockIndexRule = append(_blockIndexRule, _blockIndexItem)
	}

	var _statusRule []interface{}
	for _, _statusItem := range _status {
		_statusRule = append(_statusRule, _statusItem)
	}

	logs, sub, err := _ChallengeContract.contract.WatchLogs(opts, "ChallengeDAUpdate", _blockHashRule, _blockIndexRule, _statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChallengeContractChallengeDAUpdate)
				if err := _ChallengeContract.contract.UnpackLog(event, "ChallengeDAUpdate", log); err != nil {
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

// ParseChallengeDAUpdate is a log parse operation binding the contract event 0xcdbf857f75b46d9b94fe7c990b2b0457f168aef7b089a0cb5d4e88074f2d88e8.
//
// Solidity: event ChallengeDAUpdate(bytes32 indexed _blockHash, uint256 indexed _blockIndex, uint256 _expiry, uint8 indexed _status)
func (_ChallengeContract *ChallengeContractFilterer) ParseChallengeDAUpdate(log types.Log) (*ChallengeContractChallengeDAUpdate, error) {
	event := new(ChallengeContractChallengeDAUpdate)
	if err := _ChallengeContract.contract.UnpackLog(event, "ChallengeDAUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChallengeContractInvalidHeaderIterator is returned from FilterInvalidHeader and is used to iterate over the raw logs and unpacked data for InvalidHeader events raised by the ChallengeContract contract.
type ChallengeContractInvalidHeaderIterator struct {
	Event *ChallengeContractInvalidHeader // Event containing the contract specifics and raw log

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
func (it *ChallengeContractInvalidHeaderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChallengeContractInvalidHeader)
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
		it.Event = new(ChallengeContractInvalidHeader)
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
func (it *ChallengeContractInvalidHeaderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChallengeContractInvalidHeaderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChallengeContractInvalidHeader represents a InvalidHeader event raised by the ChallengeContract contract.
type ChallengeContractInvalidHeader struct {
	BlockIndex *big.Int
	Hash       [32]byte
	Reason     uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterInvalidHeader is a free log retrieval operation binding the contract event 0xea46f8ad2711844c28d6aa0fe8ed10b1ac38bdcdc6df7ba3b8f3bfc35232f31b.
//
// Solidity: event InvalidHeader(uint256 indexed _blockIndex, bytes32 indexed _hash, uint8 indexed _reason)
func (_ChallengeContract *ChallengeContractFilterer) FilterInvalidHeader(opts *bind.FilterOpts, _blockIndex []*big.Int, _hash [][32]byte, _reason []uint8) (*ChallengeContractInvalidHeaderIterator, error) {

	var _blockIndexRule []interface{}
	for _, _blockIndexItem := range _blockIndex {
		_blockIndexRule = append(_blockIndexRule, _blockIndexItem)
	}
	var _hashRule []interface{}
	for _, _hashItem := range _hash {
		_hashRule = append(_hashRule, _hashItem)
	}
	var _reasonRule []interface{}
	for _, _reasonItem := range _reason {
		_reasonRule = append(_reasonRule, _reasonItem)
	}

	logs, sub, err := _ChallengeContract.contract.FilterLogs(opts, "InvalidHeader", _blockIndexRule, _hashRule, _reasonRule)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractInvalidHeaderIterator{contract: _ChallengeContract.contract, event: "InvalidHeader", logs: logs, sub: sub}, nil
}

// WatchInvalidHeader is a free log subscription operation binding the contract event 0xea46f8ad2711844c28d6aa0fe8ed10b1ac38bdcdc6df7ba3b8f3bfc35232f31b.
//
// Solidity: event InvalidHeader(uint256 indexed _blockIndex, bytes32 indexed _hash, uint8 indexed _reason)
func (_ChallengeContract *ChallengeContractFilterer) WatchInvalidHeader(opts *bind.WatchOpts, sink chan<- *ChallengeContractInvalidHeader, _blockIndex []*big.Int, _hash [][32]byte, _reason []uint8) (event.Subscription, error) {

	var _blockIndexRule []interface{}
	for _, _blockIndexItem := range _blockIndex {
		_blockIndexRule = append(_blockIndexRule, _blockIndexItem)
	}
	var _hashRule []interface{}
	for _, _hashItem := range _hash {
		_hashRule = append(_hashRule, _hashItem)
	}
	var _reasonRule []interface{}
	for _, _reasonItem := range _reason {
		_reasonRule = append(_reasonRule, _reasonItem)
	}

	logs, sub, err := _ChallengeContract.contract.WatchLogs(opts, "InvalidHeader", _blockIndexRule, _hashRule, _reasonRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChallengeContractInvalidHeader)
				if err := _ChallengeContract.contract.UnpackLog(event, "InvalidHeader", log); err != nil {
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

// ParseInvalidHeader is a log parse operation binding the contract event 0xea46f8ad2711844c28d6aa0fe8ed10b1ac38bdcdc6df7ba3b8f3bfc35232f31b.
//
// Solidity: event InvalidHeader(uint256 indexed _blockIndex, bytes32 indexed _hash, uint8 indexed _reason)
func (_ChallengeContract *ChallengeContractFilterer) ParseInvalidHeader(log types.Log) (*ChallengeContractInvalidHeader, error) {
	event := new(ChallengeContractInvalidHeader)
	if err := _ChallengeContract.contract.UnpackLog(event, "InvalidHeader", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChallengeContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ChallengeContract contract.
type ChallengeContractOwnershipTransferredIterator struct {
	Event *ChallengeContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ChallengeContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChallengeContractOwnershipTransferred)
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
		it.Event = new(ChallengeContractOwnershipTransferred)
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
func (it *ChallengeContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChallengeContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChallengeContractOwnershipTransferred represents a OwnershipTransferred event raised by the ChallengeContract contract.
type ChallengeContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ChallengeContract *ChallengeContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ChallengeContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ChallengeContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ChallengeContractOwnershipTransferredIterator{contract: _ChallengeContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ChallengeContract *ChallengeContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChallengeContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ChallengeContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChallengeContractOwnershipTransferred)
				if err := _ChallengeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ChallengeContract *ChallengeContractFilterer) ParseOwnershipTransferred(log types.Log) (*ChallengeContractOwnershipTransferred, error) {
	event := new(ChallengeContractOwnershipTransferred)
	if err := _ChallengeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
