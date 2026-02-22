// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swaprouter

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

// AggregatedSwapRouterMetaData contains all meta data concerning the AggregatedSwapRouter contract.
var AggregatedSwapRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"add_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"weth_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumAggregatedSwapRouter.EventType\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"SwapEvent\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"CALLSWAP\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_weth\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"creater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"defiCallBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"defiSwap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"defiSwapForEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"defiSwapFromEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"defiSync\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"swapForEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"callSwapAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"datas\",\"type\":\"bytes\"}],\"name\":\"swapFromEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// AggregatedSwapRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatedSwapRouterMetaData.ABI instead.
var AggregatedSwapRouterABI = AggregatedSwapRouterMetaData.ABI

// AggregatedSwapRouter is an auto generated Go binding around an Ethereum contract.
type AggregatedSwapRouter struct {
	AggregatedSwapRouterCaller     // Read-only binding to the contract
	AggregatedSwapRouterTransactor // Write-only binding to the contract
	AggregatedSwapRouterFilterer   // Log filterer for contract events
}

// AggregatedSwapRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatedSwapRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatedSwapRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatedSwapRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatedSwapRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatedSwapRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatedSwapRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatedSwapRouterSession struct {
	Contract     *AggregatedSwapRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AggregatedSwapRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatedSwapRouterCallerSession struct {
	Contract *AggregatedSwapRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// AggregatedSwapRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatedSwapRouterTransactorSession struct {
	Contract     *AggregatedSwapRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// AggregatedSwapRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatedSwapRouterRaw struct {
	Contract *AggregatedSwapRouter // Generic contract binding to access the raw methods on
}

// AggregatedSwapRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatedSwapRouterCallerRaw struct {
	Contract *AggregatedSwapRouterCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatedSwapRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatedSwapRouterTransactorRaw struct {
	Contract *AggregatedSwapRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatedSwapRouter creates a new instance of AggregatedSwapRouter, bound to a specific deployed contract.
func NewAggregatedSwapRouter(address common.Address, backend bind.ContractBackend) (*AggregatedSwapRouter, error) {
	contract, err := bindAggregatedSwapRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatedSwapRouter{AggregatedSwapRouterCaller: AggregatedSwapRouterCaller{contract: contract}, AggregatedSwapRouterTransactor: AggregatedSwapRouterTransactor{contract: contract}, AggregatedSwapRouterFilterer: AggregatedSwapRouterFilterer{contract: contract}}, nil
}

// NewAggregatedSwapRouterCaller creates a new read-only instance of AggregatedSwapRouter, bound to a specific deployed contract.
func NewAggregatedSwapRouterCaller(address common.Address, caller bind.ContractCaller) (*AggregatedSwapRouterCaller, error) {
	contract, err := bindAggregatedSwapRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatedSwapRouterCaller{contract: contract}, nil
}

// NewAggregatedSwapRouterTransactor creates a new write-only instance of AggregatedSwapRouter, bound to a specific deployed contract.
func NewAggregatedSwapRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatedSwapRouterTransactor, error) {
	contract, err := bindAggregatedSwapRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatedSwapRouterTransactor{contract: contract}, nil
}

// NewAggregatedSwapRouterFilterer creates a new log filterer instance of AggregatedSwapRouter, bound to a specific deployed contract.
func NewAggregatedSwapRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatedSwapRouterFilterer, error) {
	contract, err := bindAggregatedSwapRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatedSwapRouterFilterer{contract: contract}, nil
}

// bindAggregatedSwapRouter binds a generic wrapper to an already deployed contract.
func bindAggregatedSwapRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatedSwapRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatedSwapRouter *AggregatedSwapRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatedSwapRouter.Contract.AggregatedSwapRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatedSwapRouter *AggregatedSwapRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.AggregatedSwapRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatedSwapRouter *AggregatedSwapRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.AggregatedSwapRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatedSwapRouter *AggregatedSwapRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatedSwapRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.contract.Transact(opts, method, params...)
}

// CALLSWAP is a free data retrieval call binding the contract method 0x57d00ca7.
//
// Solidity: function CALLSWAP() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCaller) CALLSWAP(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatedSwapRouter.contract.Call(opts, &out, "CALLSWAP")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CALLSWAP is a free data retrieval call binding the contract method 0x57d00ca7.
//
// Solidity: function CALLSWAP() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) CALLSWAP() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.CALLSWAP(&_AggregatedSwapRouter.CallOpts)
}

// CALLSWAP is a free data retrieval call binding the contract method 0x57d00ca7.
//
// Solidity: function CALLSWAP() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCallerSession) CALLSWAP() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.CALLSWAP(&_AggregatedSwapRouter.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0xa1764595.
//
// Solidity: function _weth() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCaller) Weth(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatedSwapRouter.contract.Call(opts, &out, "_weth")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Weth is a free data retrieval call binding the contract method 0xa1764595.
//
// Solidity: function _weth() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) Weth() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.Weth(&_AggregatedSwapRouter.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0xa1764595.
//
// Solidity: function _weth() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCallerSession) Weth() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.Weth(&_AggregatedSwapRouter.CallOpts)
}

// Creater is a free data retrieval call binding the contract method 0x45653a6d.
//
// Solidity: function creater() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCaller) Creater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatedSwapRouter.contract.Call(opts, &out, "creater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creater is a free data retrieval call binding the contract method 0x45653a6d.
//
// Solidity: function creater() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) Creater() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.Creater(&_AggregatedSwapRouter.CallOpts)
}

// Creater is a free data retrieval call binding the contract method 0x45653a6d.
//
// Solidity: function creater() view returns(address)
func (_AggregatedSwapRouter *AggregatedSwapRouterCallerSession) Creater() (common.Address, error) {
	return _AggregatedSwapRouter.Contract.Creater(&_AggregatedSwapRouter.CallOpts)
}

// DefiCallBack is a paid mutator transaction binding the contract method 0x676f45a9.
//
// Solidity: function defiCallBack(address to, uint256 value) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) DefiCallBack(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "defiCallBack", to, value)
}

// DefiCallBack is a paid mutator transaction binding the contract method 0x676f45a9.
//
// Solidity: function defiCallBack(address to, uint256 value) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) DefiCallBack(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiCallBack(&_AggregatedSwapRouter.TransactOpts, to, value)
}

// DefiCallBack is a paid mutator transaction binding the contract method 0x676f45a9.
//
// Solidity: function defiCallBack(address to, uint256 value) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) DefiCallBack(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiCallBack(&_AggregatedSwapRouter.TransactOpts, to, value)
}

// DefiSwap is a paid mutator transaction binding the contract method 0x10793b38.
//
// Solidity: function defiSwap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) DefiSwap(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "defiSwap", amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSwap is a paid mutator transaction binding the contract method 0x10793b38.
//
// Solidity: function defiSwap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) DefiSwap(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwap(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSwap is a paid mutator transaction binding the contract method 0x10793b38.
//
// Solidity: function defiSwap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) DefiSwap(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwap(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSwapForEth is a paid mutator transaction binding the contract method 0xef7a21b7.
//
// Solidity: function defiSwapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) DefiSwapForEth(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "defiSwapForEth", amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// DefiSwapForEth is a paid mutator transaction binding the contract method 0xef7a21b7.
//
// Solidity: function defiSwapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) DefiSwapForEth(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwapForEth(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// DefiSwapForEth is a paid mutator transaction binding the contract method 0xef7a21b7.
//
// Solidity: function defiSwapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) DefiSwapForEth(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwapForEth(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// DefiSwapFromEth is a paid mutator transaction binding the contract method 0x6ff05b17.
//
// Solidity: function defiSwapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) DefiSwapFromEth(opts *bind.TransactOpts, amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "defiSwapFromEth", amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSwapFromEth is a paid mutator transaction binding the contract method 0x6ff05b17.
//
// Solidity: function defiSwapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) DefiSwapFromEth(amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwapFromEth(&_AggregatedSwapRouter.TransactOpts, amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSwapFromEth is a paid mutator transaction binding the contract method 0x6ff05b17.
//
// Solidity: function defiSwapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) DefiSwapFromEth(amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSwapFromEth(&_AggregatedSwapRouter.TransactOpts, amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// DefiSync is a paid mutator transaction binding the contract method 0x202a4e57.
//
// Solidity: function defiSync(address token) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) DefiSync(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "defiSync", token)
}

// DefiSync is a paid mutator transaction binding the contract method 0x202a4e57.
//
// Solidity: function defiSync(address token) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) DefiSync(token common.Address) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSync(&_AggregatedSwapRouter.TransactOpts, token)
}

// DefiSync is a paid mutator transaction binding the contract method 0x202a4e57.
//
// Solidity: function defiSync(address token) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) DefiSync(token common.Address) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.DefiSync(&_AggregatedSwapRouter.TransactOpts, token)
}

// Swap is a paid mutator transaction binding the contract method 0xe5488f5e.
//
// Solidity: function swap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) Swap(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "swap", amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// Swap is a paid mutator transaction binding the contract method 0xe5488f5e.
//
// Solidity: function swap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) Swap(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Swap(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// Swap is a paid mutator transaction binding the contract method 0xe5488f5e.
//
// Solidity: function swap(uint256 amountIn, uint256 amountOutMin, address tokenIn, address tokenOut, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) Swap(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Swap(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, tokenOut, receiver, callSwapAddr, datas)
}

// SwapForEth is a paid mutator transaction binding the contract method 0xd8a34ee4.
//
// Solidity: function swapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) SwapForEth(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "swapForEth", amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// SwapForEth is a paid mutator transaction binding the contract method 0xd8a34ee4.
//
// Solidity: function swapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) SwapForEth(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.SwapForEth(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// SwapForEth is a paid mutator transaction binding the contract method 0xd8a34ee4.
//
// Solidity: function swapForEth(uint256 amountIn, uint256 amountOutMin, address tokenIn, address receiver, address callSwapAddr, bytes datas) returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) SwapForEth(amountIn *big.Int, amountOutMin *big.Int, tokenIn common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.SwapForEth(&_AggregatedSwapRouter.TransactOpts, amountIn, amountOutMin, tokenIn, receiver, callSwapAddr, datas)
}

// SwapFromEth is a paid mutator transaction binding the contract method 0x5cd07242.
//
// Solidity: function swapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) SwapFromEth(opts *bind.TransactOpts, amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.Transact(opts, "swapFromEth", amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// SwapFromEth is a paid mutator transaction binding the contract method 0x5cd07242.
//
// Solidity: function swapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) SwapFromEth(amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.SwapFromEth(&_AggregatedSwapRouter.TransactOpts, amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// SwapFromEth is a paid mutator transaction binding the contract method 0x5cd07242.
//
// Solidity: function swapFromEth(uint256 amountOutMin, address tokenOut, address receiver, address callSwapAddr, bytes datas) payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) SwapFromEth(amountOutMin *big.Int, tokenOut common.Address, receiver common.Address, callSwapAddr common.Address, datas []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.SwapFromEth(&_AggregatedSwapRouter.TransactOpts, amountOutMin, tokenOut, receiver, callSwapAddr, datas)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Fallback(&_AggregatedSwapRouter.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Fallback(&_AggregatedSwapRouter.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatedSwapRouter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterSession) Receive() (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Receive(&_AggregatedSwapRouter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AggregatedSwapRouter *AggregatedSwapRouterTransactorSession) Receive() (*types.Transaction, error) {
	return _AggregatedSwapRouter.Contract.Receive(&_AggregatedSwapRouter.TransactOpts)
}

// AggregatedSwapRouterSwapEventIterator is returned from FilterSwapEvent and is used to iterate over the raw logs and unpacked data for SwapEvent events raised by the AggregatedSwapRouter contract.
type AggregatedSwapRouterSwapEventIterator struct {
	Event *AggregatedSwapRouterSwapEvent // Event containing the contract specifics and raw log

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
func (it *AggregatedSwapRouterSwapEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatedSwapRouterSwapEvent)
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
		it.Event = new(AggregatedSwapRouterSwapEvent)
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
func (it *AggregatedSwapRouterSwapEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatedSwapRouterSwapEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatedSwapRouterSwapEvent represents a SwapEvent event raised by the AggregatedSwapRouter contract.
type AggregatedSwapRouterSwapEvent struct {
	Arg0 common.Address
	Arg1 *big.Int
	Arg2 common.Address
	Arg3 *big.Int
	Arg4 common.Address
	Arg5 uint8
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSwapEvent is a free log retrieval operation binding the contract event 0x14b237993bd63b85f70e0dfd4a4acb834dd2109270a8142fe3d7ae371f66307a.
//
// Solidity: event SwapEvent(address arg0, uint256 arg1, address arg2, uint256 arg3, address arg4, uint8 arg5)
func (_AggregatedSwapRouter *AggregatedSwapRouterFilterer) FilterSwapEvent(opts *bind.FilterOpts) (*AggregatedSwapRouterSwapEventIterator, error) {

	logs, sub, err := _AggregatedSwapRouter.contract.FilterLogs(opts, "SwapEvent")
	if err != nil {
		return nil, err
	}
	return &AggregatedSwapRouterSwapEventIterator{contract: _AggregatedSwapRouter.contract, event: "SwapEvent", logs: logs, sub: sub}, nil
}

// WatchSwapEvent is a free log subscription operation binding the contract event 0x14b237993bd63b85f70e0dfd4a4acb834dd2109270a8142fe3d7ae371f66307a.
//
// Solidity: event SwapEvent(address arg0, uint256 arg1, address arg2, uint256 arg3, address arg4, uint8 arg5)
func (_AggregatedSwapRouter *AggregatedSwapRouterFilterer) WatchSwapEvent(opts *bind.WatchOpts, sink chan<- *AggregatedSwapRouterSwapEvent) (event.Subscription, error) {

	logs, sub, err := _AggregatedSwapRouter.contract.WatchLogs(opts, "SwapEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatedSwapRouterSwapEvent)
				if err := _AggregatedSwapRouter.contract.UnpackLog(event, "SwapEvent", log); err != nil {
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

// ParseSwapEvent is a log parse operation binding the contract event 0x14b237993bd63b85f70e0dfd4a4acb834dd2109270a8142fe3d7ae371f66307a.
//
// Solidity: event SwapEvent(address arg0, uint256 arg1, address arg2, uint256 arg3, address arg4, uint8 arg5)
func (_AggregatedSwapRouter *AggregatedSwapRouterFilterer) ParseSwapEvent(log types.Log) (*AggregatedSwapRouterSwapEvent, error) {
	event := new(AggregatedSwapRouterSwapEvent)
	if err := _AggregatedSwapRouter.contract.UnpackLog(event, "SwapEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
