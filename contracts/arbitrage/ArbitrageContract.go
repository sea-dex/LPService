// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrage

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

// ArbitrageContractMetaData contains all meta data concerning the ArbitrageContract contract.
var ArbitrageContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hook\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pancakeCall\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pancakeV3FlashCallback\",\"inputs\":[{\"name\":\"fee0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pancakeV3SwapCallback\",\"inputs\":[{\"name\":\"amount0Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"amount1Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWhitelist\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allow\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapFromV2\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapFromV3\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwner\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniswapV2Call\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniswapV3FlashCallback\",\"inputs\":[{\"name\":\"fee0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniswapV3SwapCallback\",\"inputs\":[{\"name\":\"amount0Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"amount1Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// ArbitrageContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ArbitrageContractMetaData.ABI instead.
var ArbitrageContractABI = ArbitrageContractMetaData.ABI

// ArbitrageContract is an auto generated Go binding around an Ethereum contract.
type ArbitrageContract struct {
	ArbitrageContractCaller     // Read-only binding to the contract
	ArbitrageContractTransactor // Write-only binding to the contract
	ArbitrageContractFilterer   // Log filterer for contract events
}

// ArbitrageContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArbitrageContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArbitrageContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArbitrageContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArbitrageContractSession struct {
	Contract     *ArbitrageContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ArbitrageContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArbitrageContractCallerSession struct {
	Contract *ArbitrageContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ArbitrageContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArbitrageContractTransactorSession struct {
	Contract     *ArbitrageContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ArbitrageContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArbitrageContractRaw struct {
	Contract *ArbitrageContract // Generic contract binding to access the raw methods on
}

// ArbitrageContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArbitrageContractCallerRaw struct {
	Contract *ArbitrageContractCaller // Generic read-only contract binding to access the raw methods on
}

// ArbitrageContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArbitrageContractTransactorRaw struct {
	Contract *ArbitrageContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArbitrageContract creates a new instance of ArbitrageContract, bound to a specific deployed contract.
func NewArbitrageContract(address common.Address, backend bind.ContractBackend) (*ArbitrageContract, error) {
	contract, err := bindArbitrageContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrageContract{ArbitrageContractCaller: ArbitrageContractCaller{contract: contract}, ArbitrageContractTransactor: ArbitrageContractTransactor{contract: contract}, ArbitrageContractFilterer: ArbitrageContractFilterer{contract: contract}}, nil
}

// NewArbitrageContractCaller creates a new read-only instance of ArbitrageContract, bound to a specific deployed contract.
func NewArbitrageContractCaller(address common.Address, caller bind.ContractCaller) (*ArbitrageContractCaller, error) {
	contract, err := bindArbitrageContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrageContractCaller{contract: contract}, nil
}

// NewArbitrageContractTransactor creates a new write-only instance of ArbitrageContract, bound to a specific deployed contract.
func NewArbitrageContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrageContractTransactor, error) {
	contract, err := bindArbitrageContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrageContractTransactor{contract: contract}, nil
}

// NewArbitrageContractFilterer creates a new log filterer instance of ArbitrageContract, bound to a specific deployed contract.
func NewArbitrageContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrageContractFilterer, error) {
	contract, err := bindArbitrageContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrageContractFilterer{contract: contract}, nil
}

// bindArbitrageContract binds a generic wrapper to an already deployed contract.
func bindArbitrageContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrageContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArbitrageContract *ArbitrageContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrageContract.Contract.ArbitrageContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArbitrageContract *ArbitrageContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.ArbitrageContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArbitrageContract *ArbitrageContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.ArbitrageContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArbitrageContract *ArbitrageContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrageContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArbitrageContract *ArbitrageContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArbitrageContract *ArbitrageContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbitrageContract *ArbitrageContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrageContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbitrageContract *ArbitrageContractSession) Owner() (common.Address, error) {
	return _ArbitrageContract.Contract.Owner(&_ArbitrageContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbitrageContract *ArbitrageContractCallerSession) Owner() (common.Address, error) {
	return _ArbitrageContract.Contract.Owner(&_ArbitrageContract.CallOpts)
}

// Hook is a paid mutator transaction binding the contract method 0x9a7bff79.
//
// Solidity: function hook(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) Hook(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "hook", sender, amount0, amount1, data)
}

// Hook is a paid mutator transaction binding the contract method 0x9a7bff79.
//
// Solidity: function hook(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) Hook(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.Hook(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// Hook is a paid mutator transaction binding the contract method 0x9a7bff79.
//
// Solidity: function hook(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) Hook(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.Hook(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) PancakeCall(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "pancakeCall", sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) PancakeCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeCall(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) PancakeCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeCall(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// PancakeV3FlashCallback is a paid mutator transaction binding the contract method 0xa1d48336.
//
// Solidity: function pancakeV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) PancakeV3FlashCallback(opts *bind.TransactOpts, fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "pancakeV3FlashCallback", fee0, fee1, data)
}

// PancakeV3FlashCallback is a paid mutator transaction binding the contract method 0xa1d48336.
//
// Solidity: function pancakeV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) PancakeV3FlashCallback(fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeV3FlashCallback(&_ArbitrageContract.TransactOpts, fee0, fee1, data)
}

// PancakeV3FlashCallback is a paid mutator transaction binding the contract method 0xa1d48336.
//
// Solidity: function pancakeV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) PancakeV3FlashCallback(fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeV3FlashCallback(&_ArbitrageContract.TransactOpts, fee0, fee1, data)
}

// PancakeV3SwapCallback is a paid mutator transaction binding the contract method 0x23a69e75.
//
// Solidity: function pancakeV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) PancakeV3SwapCallback(opts *bind.TransactOpts, amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "pancakeV3SwapCallback", amount0Delta, amount1Delta, data)
}

// PancakeV3SwapCallback is a paid mutator transaction binding the contract method 0x23a69e75.
//
// Solidity: function pancakeV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) PancakeV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeV3SwapCallback(&_ArbitrageContract.TransactOpts, amount0Delta, amount1Delta, data)
}

// PancakeV3SwapCallback is a paid mutator transaction binding the contract method 0x23a69e75.
//
// Solidity: function pancakeV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) PancakeV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.PancakeV3SwapCallback(&_ArbitrageContract.TransactOpts, amount0Delta, amount1Delta, data)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address addr, bool allow) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) SetWhitelist(opts *bind.TransactOpts, addr common.Address, allow bool) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "setWhitelist", addr, allow)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address addr, bool allow) returns()
func (_ArbitrageContract *ArbitrageContractSession) SetWhitelist(addr common.Address, allow bool) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SetWhitelist(&_ArbitrageContract.TransactOpts, addr, allow)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address addr, bool allow) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) SetWhitelist(addr common.Address, allow bool) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SetWhitelist(&_ArbitrageContract.TransactOpts, addr, allow)
}

// SwapFromV2 is a paid mutator transaction binding the contract method 0x477da876.
//
// Solidity: function swapFromV2(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) SwapFromV2(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "swapFromV2", data)
}

// SwapFromV2 is a paid mutator transaction binding the contract method 0x477da876.
//
// Solidity: function swapFromV2(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) SwapFromV2(data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SwapFromV2(&_ArbitrageContract.TransactOpts, data)
}

// SwapFromV2 is a paid mutator transaction binding the contract method 0x477da876.
//
// Solidity: function swapFromV2(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) SwapFromV2(data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SwapFromV2(&_ArbitrageContract.TransactOpts, data)
}

// SwapFromV3 is a paid mutator transaction binding the contract method 0x575b89aa.
//
// Solidity: function swapFromV3(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) SwapFromV3(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "swapFromV3", data)
}

// SwapFromV3 is a paid mutator transaction binding the contract method 0x575b89aa.
//
// Solidity: function swapFromV3(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) SwapFromV3(data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SwapFromV3(&_ArbitrageContract.TransactOpts, data)
}

// SwapFromV3 is a paid mutator transaction binding the contract method 0x575b89aa.
//
// Solidity: function swapFromV3(bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) SwapFromV3(data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.SwapFromV3(&_ArbitrageContract.TransactOpts, data)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x4fb2e45d.
//
// Solidity: function transferOwner(address _owner) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) TransferOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "transferOwner", _owner)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x4fb2e45d.
//
// Solidity: function transferOwner(address _owner) returns()
func (_ArbitrageContract *ArbitrageContractSession) TransferOwner(_owner common.Address) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.TransferOwner(&_ArbitrageContract.TransactOpts, _owner)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x4fb2e45d.
//
// Solidity: function transferOwner(address _owner) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) TransferOwner(_owner common.Address) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.TransferOwner(&_ArbitrageContract.TransactOpts, _owner)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) UniswapV2Call(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "uniswapV2Call", sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV2Call(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV2Call(&_ArbitrageContract.TransactOpts, sender, amount0, amount1, data)
}

// UniswapV3FlashCallback is a paid mutator transaction binding the contract method 0xe9cbafb0.
//
// Solidity: function uniswapV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) UniswapV3FlashCallback(opts *bind.TransactOpts, fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "uniswapV3FlashCallback", fee0, fee1, data)
}

// UniswapV3FlashCallback is a paid mutator transaction binding the contract method 0xe9cbafb0.
//
// Solidity: function uniswapV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) UniswapV3FlashCallback(fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV3FlashCallback(&_ArbitrageContract.TransactOpts, fee0, fee1, data)
}

// UniswapV3FlashCallback is a paid mutator transaction binding the contract method 0xe9cbafb0.
//
// Solidity: function uniswapV3FlashCallback(uint256 fee0, uint256 fee1, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) UniswapV3FlashCallback(fee0 *big.Int, fee1 *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV3FlashCallback(&_ArbitrageContract.TransactOpts, fee0, fee1, data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) UniswapV3SwapCallback(opts *bind.TransactOpts, amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "uniswapV3SwapCallback", amount0Delta, amount1Delta, data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV3SwapCallback(&_ArbitrageContract.TransactOpts, amount0Delta, amount1Delta, data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.UniswapV3SwapCallback(&_ArbitrageContract.TransactOpts, amount0Delta, amount1Delta, data)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_ArbitrageContract *ArbitrageContractTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrageContract.contract.Transact(opts, "withdrawToken", token, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_ArbitrageContract *ArbitrageContractSession) WithdrawToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.WithdrawToken(&_ArbitrageContract.TransactOpts, token, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_ArbitrageContract *ArbitrageContractTransactorSession) WithdrawToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrageContract.Contract.WithdrawToken(&_ArbitrageContract.TransactOpts, token, amount)
}
