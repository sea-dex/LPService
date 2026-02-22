// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package poolquery

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

// UniswapV3PoolQueryMetaData contains all meta data concerning the UniswapV3PoolQuery contract.
var UniswapV3PoolQueryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"poolAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"words\",\"type\":\"uint256\"}],\"name\":\"getAllTicksInWord\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"internalType\":\"int24\",\"name\":\"tickCur\",\"type\":\"int24\"},{\"internalType\":\"int24[]\",\"name\":\"ticks\",\"type\":\"int24[]\"},{\"internalType\":\"int128[]\",\"name\":\"liquidityNets\",\"type\":\"int128[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"sdatas\",\"type\":\"bytes[]\"}],\"name\":\"getPopulatedTicksInWords\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"datas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"poolAddrs\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"tickDatas\",\"type\":\"bytes[]\"}],\"name\":\"getTickLiqs\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"datas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// UniswapV3PoolQueryABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV3PoolQueryMetaData.ABI instead.
var UniswapV3PoolQueryABI = UniswapV3PoolQueryMetaData.ABI

// UniswapV3PoolQuery is an auto generated Go binding around an Ethereum contract.
type UniswapV3PoolQuery struct {
	UniswapV3PoolQueryCaller     // Read-only binding to the contract
	UniswapV3PoolQueryTransactor // Write-only binding to the contract
	UniswapV3PoolQueryFilterer   // Log filterer for contract events
}

// UniswapV3PoolQueryCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV3PoolQueryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3PoolQueryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV3PoolQueryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3PoolQueryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV3PoolQueryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3PoolQuerySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV3PoolQuerySession struct {
	Contract     *UniswapV3PoolQuery // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// UniswapV3PoolQueryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV3PoolQueryCallerSession struct {
	Contract *UniswapV3PoolQueryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// UniswapV3PoolQueryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV3PoolQueryTransactorSession struct {
	Contract     *UniswapV3PoolQueryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// UniswapV3PoolQueryRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV3PoolQueryRaw struct {
	Contract *UniswapV3PoolQuery // Generic contract binding to access the raw methods on
}

// UniswapV3PoolQueryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV3PoolQueryCallerRaw struct {
	Contract *UniswapV3PoolQueryCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV3PoolQueryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV3PoolQueryTransactorRaw struct {
	Contract *UniswapV3PoolQueryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV3PoolQuery creates a new instance of UniswapV3PoolQuery, bound to a specific deployed contract.
func NewUniswapV3PoolQuery(address common.Address, backend bind.ContractBackend) (*UniswapV3PoolQuery, error) {
	contract, err := bindUniswapV3PoolQuery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV3PoolQuery{UniswapV3PoolQueryCaller: UniswapV3PoolQueryCaller{contract: contract}, UniswapV3PoolQueryTransactor: UniswapV3PoolQueryTransactor{contract: contract}, UniswapV3PoolQueryFilterer: UniswapV3PoolQueryFilterer{contract: contract}}, nil
}

// NewUniswapV3PoolQueryCaller creates a new read-only instance of UniswapV3PoolQuery, bound to a specific deployed contract.
func NewUniswapV3PoolQueryCaller(address common.Address, caller bind.ContractCaller) (*UniswapV3PoolQueryCaller, error) {
	contract, err := bindUniswapV3PoolQuery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3PoolQueryCaller{contract: contract}, nil
}

// NewUniswapV3PoolQueryTransactor creates a new write-only instance of UniswapV3PoolQuery, bound to a specific deployed contract.
func NewUniswapV3PoolQueryTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV3PoolQueryTransactor, error) {
	contract, err := bindUniswapV3PoolQuery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3PoolQueryTransactor{contract: contract}, nil
}

// NewUniswapV3PoolQueryFilterer creates a new log filterer instance of UniswapV3PoolQuery, bound to a specific deployed contract.
func NewUniswapV3PoolQueryFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV3PoolQueryFilterer, error) {
	contract, err := bindUniswapV3PoolQuery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV3PoolQueryFilterer{contract: contract}, nil
}

// bindUniswapV3PoolQuery binds a generic wrapper to an already deployed contract.
func bindUniswapV3PoolQuery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV3PoolQueryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3PoolQuery.Contract.UniswapV3PoolQueryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3PoolQuery.Contract.UniswapV3PoolQueryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3PoolQuery.Contract.UniswapV3PoolQueryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3PoolQuery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3PoolQuery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3PoolQuery *UniswapV3PoolQueryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3PoolQuery.Contract.contract.Transact(opts, method, params...)
}

// GetAllTicksInWord is a free data retrieval call binding the contract method 0x14ffd68e.
//
// Solidity: function getAllTicksInWord(address poolAddr, uint256 words) view returns(uint160 sqrtPriceX96, uint128 liquidity, int24 tickCur, int24[] ticks, int128[] liquidityNets)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCaller) GetAllTicksInWord(opts *bind.CallOpts, poolAddr common.Address, words *big.Int) (struct {
	SqrtPriceX96  *big.Int
	Liquidity     *big.Int
	TickCur       *big.Int
	Ticks         []*big.Int
	LiquidityNets []*big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3PoolQuery.contract.Call(opts, &out, "getAllTicksInWord", poolAddr, words)

	outstruct := new(struct {
		SqrtPriceX96  *big.Int
		Liquidity     *big.Int
		TickCur       *big.Int
		Ticks         []*big.Int
		LiquidityNets []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SqrtPriceX96 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Liquidity = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TickCur = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Ticks = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)
	outstruct.LiquidityNets = *abi.ConvertType(out[4], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetAllTicksInWord is a free data retrieval call binding the contract method 0x14ffd68e.
//
// Solidity: function getAllTicksInWord(address poolAddr, uint256 words) view returns(uint160 sqrtPriceX96, uint128 liquidity, int24 tickCur, int24[] ticks, int128[] liquidityNets)
func (_UniswapV3PoolQuery *UniswapV3PoolQuerySession) GetAllTicksInWord(poolAddr common.Address, words *big.Int) (struct {
	SqrtPriceX96  *big.Int
	Liquidity     *big.Int
	TickCur       *big.Int
	Ticks         []*big.Int
	LiquidityNets []*big.Int
}, error) {
	return _UniswapV3PoolQuery.Contract.GetAllTicksInWord(&_UniswapV3PoolQuery.CallOpts, poolAddr, words)
}

// GetAllTicksInWord is a free data retrieval call binding the contract method 0x14ffd68e.
//
// Solidity: function getAllTicksInWord(address poolAddr, uint256 words) view returns(uint160 sqrtPriceX96, uint128 liquidity, int24 tickCur, int24[] ticks, int128[] liquidityNets)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCallerSession) GetAllTicksInWord(poolAddr common.Address, words *big.Int) (struct {
	SqrtPriceX96  *big.Int
	Liquidity     *big.Int
	TickCur       *big.Int
	Ticks         []*big.Int
	LiquidityNets []*big.Int
}, error) {
	return _UniswapV3PoolQuery.Contract.GetAllTicksInWord(&_UniswapV3PoolQuery.CallOpts, poolAddr, words)
}

// GetPopulatedTicksInWords is a free data retrieval call binding the contract method 0x6ed31357.
//
// Solidity: function getPopulatedTicksInWords(bytes[] sdatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCaller) GetPopulatedTicksInWords(opts *bind.CallOpts, sdatas [][]byte) ([][]byte, error) {
	var out []interface{}
	err := _UniswapV3PoolQuery.contract.Call(opts, &out, "getPopulatedTicksInWords", sdatas)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// GetPopulatedTicksInWords is a free data retrieval call binding the contract method 0x6ed31357.
//
// Solidity: function getPopulatedTicksInWords(bytes[] sdatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQuerySession) GetPopulatedTicksInWords(sdatas [][]byte) ([][]byte, error) {
	return _UniswapV3PoolQuery.Contract.GetPopulatedTicksInWords(&_UniswapV3PoolQuery.CallOpts, sdatas)
}

// GetPopulatedTicksInWords is a free data retrieval call binding the contract method 0x6ed31357.
//
// Solidity: function getPopulatedTicksInWords(bytes[] sdatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCallerSession) GetPopulatedTicksInWords(sdatas [][]byte) ([][]byte, error) {
	return _UniswapV3PoolQuery.Contract.GetPopulatedTicksInWords(&_UniswapV3PoolQuery.CallOpts, sdatas)
}

// GetTickLiqs is a free data retrieval call binding the contract method 0x52746a2d.
//
// Solidity: function getTickLiqs(address[] poolAddrs, bytes[] tickDatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCaller) GetTickLiqs(opts *bind.CallOpts, poolAddrs []common.Address, tickDatas [][]byte) ([][]byte, error) {
	var out []interface{}
	err := _UniswapV3PoolQuery.contract.Call(opts, &out, "getTickLiqs", poolAddrs, tickDatas)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// GetTickLiqs is a free data retrieval call binding the contract method 0x52746a2d.
//
// Solidity: function getTickLiqs(address[] poolAddrs, bytes[] tickDatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQuerySession) GetTickLiqs(poolAddrs []common.Address, tickDatas [][]byte) ([][]byte, error) {
	return _UniswapV3PoolQuery.Contract.GetTickLiqs(&_UniswapV3PoolQuery.CallOpts, poolAddrs, tickDatas)
}

// GetTickLiqs is a free data retrieval call binding the contract method 0x52746a2d.
//
// Solidity: function getTickLiqs(address[] poolAddrs, bytes[] tickDatas) view returns(bytes[] datas)
func (_UniswapV3PoolQuery *UniswapV3PoolQueryCallerSession) GetTickLiqs(poolAddrs []common.Address, tickDatas [][]byte) ([][]byte, error) {
	return _UniswapV3PoolQuery.Contract.GetTickLiqs(&_UniswapV3PoolQuery.CallOpts, poolAddrs, tickDatas)
}
