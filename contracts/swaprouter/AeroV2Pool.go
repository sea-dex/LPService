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

// IPoolObservation is an auto generated low-level Go binding around an user-defined struct.
type IPoolObservation struct {
	Timestamp          *big.Int
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
}

// AeroV2PoolMetaData contains all meta data concerning the AeroV2Pool contract.
var AeroV2PoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BelowMinimumK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DepositsNotEqual\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FactoryAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientInputAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidityBurned\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidityMinted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientOutputAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTo\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IsPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"K\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEmergencyCouncil\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Claim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Fees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reserve0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reserve1\",\"type\":\"uint256\"}],\"name\":\"Sync\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockTimestampLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"claimed0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimed1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"claimable0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"claimable1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentCumulativePrices\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reserve0Cumulative\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve1Cumulative\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"}],\"name\":\"getAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getK\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reserve0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserve1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_blockTimestampLast\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"index0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"index1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token1\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_stable\",\"type\":\"bool\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastObservation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve0Cumulative\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve1Cumulative\",\"type\":\"uint256\"}],\"internalType\":\"structIPool.Observation\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"metadata\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"dec0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dec1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"r0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"r1\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"st\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"t0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"t1\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"observationLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"observations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve0Cumulative\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve1Cumulative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"periodSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolFees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"points\",\"type\":\"uint256\"}],\"name\":\"prices\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"granularity\",\"type\":\"uint256\"}],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve0CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve1CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"points\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"window\",\"type\":\"uint256\"}],\"name\":\"sample\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"__name\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"__symbol\",\"type\":\"string\"}],\"name\":\"setSymbol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"skim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"supplyIndex0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"supplyIndex1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sync\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AeroV2PoolABI is the input ABI used to generate the binding from.
// Deprecated: Use AeroV2PoolMetaData.ABI instead.
var AeroV2PoolABI = AeroV2PoolMetaData.ABI

// AeroV2Pool is an auto generated Go binding around an Ethereum contract.
type AeroV2Pool struct {
	AeroV2PoolCaller     // Read-only binding to the contract
	AeroV2PoolTransactor // Write-only binding to the contract
	AeroV2PoolFilterer   // Log filterer for contract events
}

// AeroV2PoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type AeroV2PoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AeroV2PoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AeroV2PoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AeroV2PoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AeroV2PoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AeroV2PoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AeroV2PoolSession struct {
	Contract     *AeroV2Pool       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AeroV2PoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AeroV2PoolCallerSession struct {
	Contract *AeroV2PoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AeroV2PoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AeroV2PoolTransactorSession struct {
	Contract     *AeroV2PoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AeroV2PoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type AeroV2PoolRaw struct {
	Contract *AeroV2Pool // Generic contract binding to access the raw methods on
}

// AeroV2PoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AeroV2PoolCallerRaw struct {
	Contract *AeroV2PoolCaller // Generic read-only contract binding to access the raw methods on
}

// AeroV2PoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AeroV2PoolTransactorRaw struct {
	Contract *AeroV2PoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAeroV2Pool creates a new instance of AeroV2Pool, bound to a specific deployed contract.
func NewAeroV2Pool(address common.Address, backend bind.ContractBackend) (*AeroV2Pool, error) {
	contract, err := bindAeroV2Pool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AeroV2Pool{AeroV2PoolCaller: AeroV2PoolCaller{contract: contract}, AeroV2PoolTransactor: AeroV2PoolTransactor{contract: contract}, AeroV2PoolFilterer: AeroV2PoolFilterer{contract: contract}}, nil
}

// NewAeroV2PoolCaller creates a new read-only instance of AeroV2Pool, bound to a specific deployed contract.
func NewAeroV2PoolCaller(address common.Address, caller bind.ContractCaller) (*AeroV2PoolCaller, error) {
	contract, err := bindAeroV2Pool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolCaller{contract: contract}, nil
}

// NewAeroV2PoolTransactor creates a new write-only instance of AeroV2Pool, bound to a specific deployed contract.
func NewAeroV2PoolTransactor(address common.Address, transactor bind.ContractTransactor) (*AeroV2PoolTransactor, error) {
	contract, err := bindAeroV2Pool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolTransactor{contract: contract}, nil
}

// NewAeroV2PoolFilterer creates a new log filterer instance of AeroV2Pool, bound to a specific deployed contract.
func NewAeroV2PoolFilterer(address common.Address, filterer bind.ContractFilterer) (*AeroV2PoolFilterer, error) {
	contract, err := bindAeroV2Pool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolFilterer{contract: contract}, nil
}

// bindAeroV2Pool binds a generic wrapper to an already deployed contract.
func bindAeroV2Pool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AeroV2PoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AeroV2Pool *AeroV2PoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AeroV2Pool.Contract.AeroV2PoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AeroV2Pool *AeroV2PoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.AeroV2PoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AeroV2Pool *AeroV2PoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.AeroV2PoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AeroV2Pool *AeroV2PoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AeroV2Pool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AeroV2Pool *AeroV2PoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AeroV2Pool *AeroV2PoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AeroV2Pool *AeroV2PoolCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AeroV2Pool *AeroV2PoolSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _AeroV2Pool.Contract.DOMAINSEPARATOR(&_AeroV2Pool.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AeroV2Pool *AeroV2PoolCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _AeroV2Pool.Contract.DOMAINSEPARATOR(&_AeroV2Pool.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Allowance(&_AeroV2Pool.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Allowance(&_AeroV2Pool.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.BalanceOf(&_AeroV2Pool.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.BalanceOf(&_AeroV2Pool.CallOpts, account)
}

// BlockTimestampLast is a free data retrieval call binding the contract method 0xc5700a02.
//
// Solidity: function blockTimestampLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) BlockTimestampLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "blockTimestampLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlockTimestampLast is a free data retrieval call binding the contract method 0xc5700a02.
//
// Solidity: function blockTimestampLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) BlockTimestampLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.BlockTimestampLast(&_AeroV2Pool.CallOpts)
}

// BlockTimestampLast is a free data retrieval call binding the contract method 0xc5700a02.
//
// Solidity: function blockTimestampLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) BlockTimestampLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.BlockTimestampLast(&_AeroV2Pool.CallOpts)
}

// Claimable0 is a free data retrieval call binding the contract method 0x4d5a9f8a.
//
// Solidity: function claimable0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Claimable0(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "claimable0", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Claimable0 is a free data retrieval call binding the contract method 0x4d5a9f8a.
//
// Solidity: function claimable0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Claimable0(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Claimable0(&_AeroV2Pool.CallOpts, arg0)
}

// Claimable0 is a free data retrieval call binding the contract method 0x4d5a9f8a.
//
// Solidity: function claimable0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Claimable0(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Claimable0(&_AeroV2Pool.CallOpts, arg0)
}

// Claimable1 is a free data retrieval call binding the contract method 0xa1ac4d13.
//
// Solidity: function claimable1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Claimable1(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "claimable1", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Claimable1 is a free data retrieval call binding the contract method 0xa1ac4d13.
//
// Solidity: function claimable1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Claimable1(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Claimable1(&_AeroV2Pool.CallOpts, arg0)
}

// Claimable1 is a free data retrieval call binding the contract method 0xa1ac4d13.
//
// Solidity: function claimable1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Claimable1(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Claimable1(&_AeroV2Pool.CallOpts, arg0)
}

// CurrentCumulativePrices is a free data retrieval call binding the contract method 0x1df8c717.
//
// Solidity: function currentCumulativePrices() view returns(uint256 reserve0Cumulative, uint256 reserve1Cumulative, uint256 blockTimestamp)
func (_AeroV2Pool *AeroV2PoolCaller) CurrentCumulativePrices(opts *bind.CallOpts) (struct {
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
	BlockTimestamp     *big.Int
}, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "currentCumulativePrices")

	outstruct := new(struct {
		Reserve0Cumulative *big.Int
		Reserve1Cumulative *big.Int
		BlockTimestamp     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0Cumulative = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1Cumulative = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CurrentCumulativePrices is a free data retrieval call binding the contract method 0x1df8c717.
//
// Solidity: function currentCumulativePrices() view returns(uint256 reserve0Cumulative, uint256 reserve1Cumulative, uint256 blockTimestamp)
func (_AeroV2Pool *AeroV2PoolSession) CurrentCumulativePrices() (struct {
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
	BlockTimestamp     *big.Int
}, error) {
	return _AeroV2Pool.Contract.CurrentCumulativePrices(&_AeroV2Pool.CallOpts)
}

// CurrentCumulativePrices is a free data retrieval call binding the contract method 0x1df8c717.
//
// Solidity: function currentCumulativePrices() view returns(uint256 reserve0Cumulative, uint256 reserve1Cumulative, uint256 blockTimestamp)
func (_AeroV2Pool *AeroV2PoolCallerSession) CurrentCumulativePrices() (struct {
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
	BlockTimestamp     *big.Int
}, error) {
	return _AeroV2Pool.Contract.CurrentCumulativePrices(&_AeroV2Pool.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AeroV2Pool *AeroV2PoolCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AeroV2Pool *AeroV2PoolSession) Decimals() (uint8, error) {
	return _AeroV2Pool.Contract.Decimals(&_AeroV2Pool.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AeroV2Pool *AeroV2PoolCallerSession) Decimals() (uint8, error) {
	return _AeroV2Pool.Contract.Decimals(&_AeroV2Pool.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AeroV2Pool *AeroV2PoolCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AeroV2Pool *AeroV2PoolSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AeroV2Pool.Contract.Eip712Domain(&_AeroV2Pool.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AeroV2Pool *AeroV2PoolCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AeroV2Pool.Contract.Eip712Domain(&_AeroV2Pool.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_AeroV2Pool *AeroV2PoolCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_AeroV2Pool *AeroV2PoolSession) Factory() (common.Address, error) {
	return _AeroV2Pool.Contract.Factory(&_AeroV2Pool.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_AeroV2Pool *AeroV2PoolCallerSession) Factory() (common.Address, error) {
	return _AeroV2Pool.Contract.Factory(&_AeroV2Pool.CallOpts)
}

// GetAmountOut is a free data retrieval call binding the contract method 0xf140a35a.
//
// Solidity: function getAmountOut(uint256 amountIn, address tokenIn) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) GetAmountOut(opts *bind.CallOpts, amountIn *big.Int, tokenIn common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "getAmountOut", amountIn, tokenIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOut is a free data retrieval call binding the contract method 0xf140a35a.
//
// Solidity: function getAmountOut(uint256 amountIn, address tokenIn) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) GetAmountOut(amountIn *big.Int, tokenIn common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.GetAmountOut(&_AeroV2Pool.CallOpts, amountIn, tokenIn)
}

// GetAmountOut is a free data retrieval call binding the contract method 0xf140a35a.
//
// Solidity: function getAmountOut(uint256 amountIn, address tokenIn) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) GetAmountOut(amountIn *big.Int, tokenIn common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.GetAmountOut(&_AeroV2Pool.CallOpts, amountIn, tokenIn)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint256 _blockTimestampLast)
func (_AeroV2Pool *AeroV2PoolCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast *big.Int
}, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint256 _blockTimestampLast)
func (_AeroV2Pool *AeroV2PoolSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast *big.Int
}, error) {
	return _AeroV2Pool.Contract.GetReserves(&_AeroV2Pool.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint256 _blockTimestampLast)
func (_AeroV2Pool *AeroV2PoolCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast *big.Int
}, error) {
	return _AeroV2Pool.Contract.GetReserves(&_AeroV2Pool.CallOpts)
}

// Index0 is a free data retrieval call binding the contract method 0x32c0defd.
//
// Solidity: function index0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Index0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "index0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Index0 is a free data retrieval call binding the contract method 0x32c0defd.
//
// Solidity: function index0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Index0() (*big.Int, error) {
	return _AeroV2Pool.Contract.Index0(&_AeroV2Pool.CallOpts)
}

// Index0 is a free data retrieval call binding the contract method 0x32c0defd.
//
// Solidity: function index0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Index0() (*big.Int, error) {
	return _AeroV2Pool.Contract.Index0(&_AeroV2Pool.CallOpts)
}

// Index1 is a free data retrieval call binding the contract method 0xbda39cad.
//
// Solidity: function index1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Index1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "index1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Index1 is a free data retrieval call binding the contract method 0xbda39cad.
//
// Solidity: function index1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Index1() (*big.Int, error) {
	return _AeroV2Pool.Contract.Index1(&_AeroV2Pool.CallOpts)
}

// Index1 is a free data retrieval call binding the contract method 0xbda39cad.
//
// Solidity: function index1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Index1() (*big.Int, error) {
	return _AeroV2Pool.Contract.Index1(&_AeroV2Pool.CallOpts)
}

// LastObservation is a free data retrieval call binding the contract method 0x8a7b8cf2.
//
// Solidity: function lastObservation() view returns((uint256,uint256,uint256))
func (_AeroV2Pool *AeroV2PoolCaller) LastObservation(opts *bind.CallOpts) (IPoolObservation, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "lastObservation")

	if err != nil {
		return *new(IPoolObservation), err
	}

	out0 := *abi.ConvertType(out[0], new(IPoolObservation)).(*IPoolObservation)

	return out0, err

}

// LastObservation is a free data retrieval call binding the contract method 0x8a7b8cf2.
//
// Solidity: function lastObservation() view returns((uint256,uint256,uint256))
func (_AeroV2Pool *AeroV2PoolSession) LastObservation() (IPoolObservation, error) {
	return _AeroV2Pool.Contract.LastObservation(&_AeroV2Pool.CallOpts)
}

// LastObservation is a free data retrieval call binding the contract method 0x8a7b8cf2.
//
// Solidity: function lastObservation() view returns((uint256,uint256,uint256))
func (_AeroV2Pool *AeroV2PoolCallerSession) LastObservation() (IPoolObservation, error) {
	return _AeroV2Pool.Contract.LastObservation(&_AeroV2Pool.CallOpts)
}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint256 dec0, uint256 dec1, uint256 r0, uint256 r1, bool st, address t0, address t1)
func (_AeroV2Pool *AeroV2PoolCaller) Metadata(opts *bind.CallOpts) (struct {
	Dec0 *big.Int
	Dec1 *big.Int
	R0   *big.Int
	R1   *big.Int
	St   bool
	T0   common.Address
	T1   common.Address
}, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "metadata")

	outstruct := new(struct {
		Dec0 *big.Int
		Dec1 *big.Int
		R0   *big.Int
		R1   *big.Int
		St   bool
		T0   common.Address
		T1   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Dec0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Dec1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.R0 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.R1 = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.St = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.T0 = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.T1 = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint256 dec0, uint256 dec1, uint256 r0, uint256 r1, bool st, address t0, address t1)
func (_AeroV2Pool *AeroV2PoolSession) Metadata() (struct {
	Dec0 *big.Int
	Dec1 *big.Int
	R0   *big.Int
	R1   *big.Int
	St   bool
	T0   common.Address
	T1   common.Address
}, error) {
	return _AeroV2Pool.Contract.Metadata(&_AeroV2Pool.CallOpts)
}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint256 dec0, uint256 dec1, uint256 r0, uint256 r1, bool st, address t0, address t1)
func (_AeroV2Pool *AeroV2PoolCallerSession) Metadata() (struct {
	Dec0 *big.Int
	Dec1 *big.Int
	R0   *big.Int
	R1   *big.Int
	St   bool
	T0   common.Address
	T1   common.Address
}, error) {
	return _AeroV2Pool.Contract.Metadata(&_AeroV2Pool.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AeroV2Pool *AeroV2PoolCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AeroV2Pool *AeroV2PoolSession) Name() (string, error) {
	return _AeroV2Pool.Contract.Name(&_AeroV2Pool.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AeroV2Pool *AeroV2PoolCallerSession) Name() (string, error) {
	return _AeroV2Pool.Contract.Name(&_AeroV2Pool.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Nonces(owner common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Nonces(&_AeroV2Pool.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.Nonces(&_AeroV2Pool.CallOpts, owner)
}

// ObservationLength is a free data retrieval call binding the contract method 0xebeb31db.
//
// Solidity: function observationLength() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) ObservationLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "observationLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ObservationLength is a free data retrieval call binding the contract method 0xebeb31db.
//
// Solidity: function observationLength() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) ObservationLength() (*big.Int, error) {
	return _AeroV2Pool.Contract.ObservationLength(&_AeroV2Pool.CallOpts)
}

// ObservationLength is a free data retrieval call binding the contract method 0xebeb31db.
//
// Solidity: function observationLength() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) ObservationLength() (*big.Int, error) {
	return _AeroV2Pool.Contract.ObservationLength(&_AeroV2Pool.CallOpts)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint256 timestamp, uint256 reserve0Cumulative, uint256 reserve1Cumulative)
func (_AeroV2Pool *AeroV2PoolCaller) Observations(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Timestamp          *big.Int
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
}, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "observations", arg0)

	outstruct := new(struct {
		Timestamp          *big.Int
		Reserve0Cumulative *big.Int
		Reserve1Cumulative *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve0Cumulative = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Reserve1Cumulative = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint256 timestamp, uint256 reserve0Cumulative, uint256 reserve1Cumulative)
func (_AeroV2Pool *AeroV2PoolSession) Observations(arg0 *big.Int) (struct {
	Timestamp          *big.Int
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
}, error) {
	return _AeroV2Pool.Contract.Observations(&_AeroV2Pool.CallOpts, arg0)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint256 timestamp, uint256 reserve0Cumulative, uint256 reserve1Cumulative)
func (_AeroV2Pool *AeroV2PoolCallerSession) Observations(arg0 *big.Int) (struct {
	Timestamp          *big.Int
	Reserve0Cumulative *big.Int
	Reserve1Cumulative *big.Int
}, error) {
	return _AeroV2Pool.Contract.Observations(&_AeroV2Pool.CallOpts, arg0)
}

// PeriodSize is a free data retrieval call binding the contract method 0xe4463eb2.
//
// Solidity: function periodSize() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) PeriodSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "periodSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PeriodSize is a free data retrieval call binding the contract method 0xe4463eb2.
//
// Solidity: function periodSize() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) PeriodSize() (*big.Int, error) {
	return _AeroV2Pool.Contract.PeriodSize(&_AeroV2Pool.CallOpts)
}

// PeriodSize is a free data retrieval call binding the contract method 0xe4463eb2.
//
// Solidity: function periodSize() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) PeriodSize() (*big.Int, error) {
	return _AeroV2Pool.Contract.PeriodSize(&_AeroV2Pool.CallOpts)
}

// PoolFees is a free data retrieval call binding the contract method 0x33580959.
//
// Solidity: function poolFees() view returns(address)
func (_AeroV2Pool *AeroV2PoolCaller) PoolFees(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "poolFees")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolFees is a free data retrieval call binding the contract method 0x33580959.
//
// Solidity: function poolFees() view returns(address)
func (_AeroV2Pool *AeroV2PoolSession) PoolFees() (common.Address, error) {
	return _AeroV2Pool.Contract.PoolFees(&_AeroV2Pool.CallOpts)
}

// PoolFees is a free data retrieval call binding the contract method 0x33580959.
//
// Solidity: function poolFees() view returns(address)
func (_AeroV2Pool *AeroV2PoolCallerSession) PoolFees() (common.Address, error) {
	return _AeroV2Pool.Contract.PoolFees(&_AeroV2Pool.CallOpts)
}

// Prices is a free data retrieval call binding the contract method 0x5881c475.
//
// Solidity: function prices(address tokenIn, uint256 amountIn, uint256 points) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolCaller) Prices(opts *bind.CallOpts, tokenIn common.Address, amountIn *big.Int, points *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "prices", tokenIn, amountIn, points)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// Prices is a free data retrieval call binding the contract method 0x5881c475.
//
// Solidity: function prices(address tokenIn, uint256 amountIn, uint256 points) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolSession) Prices(tokenIn common.Address, amountIn *big.Int, points *big.Int) ([]*big.Int, error) {
	return _AeroV2Pool.Contract.Prices(&_AeroV2Pool.CallOpts, tokenIn, amountIn, points)
}

// Prices is a free data retrieval call binding the contract method 0x5881c475.
//
// Solidity: function prices(address tokenIn, uint256 amountIn, uint256 points) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolCallerSession) Prices(tokenIn common.Address, amountIn *big.Int, points *big.Int) ([]*big.Int, error) {
	return _AeroV2Pool.Contract.Prices(&_AeroV2Pool.CallOpts, tokenIn, amountIn, points)
}

// Quote is a free data retrieval call binding the contract method 0x9e8cc04b.
//
// Solidity: function quote(address tokenIn, uint256 amountIn, uint256 granularity) view returns(uint256 amountOut)
func (_AeroV2Pool *AeroV2PoolCaller) Quote(opts *bind.CallOpts, tokenIn common.Address, amountIn *big.Int, granularity *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "quote", tokenIn, amountIn, granularity)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Quote is a free data retrieval call binding the contract method 0x9e8cc04b.
//
// Solidity: function quote(address tokenIn, uint256 amountIn, uint256 granularity) view returns(uint256 amountOut)
func (_AeroV2Pool *AeroV2PoolSession) Quote(tokenIn common.Address, amountIn *big.Int, granularity *big.Int) (*big.Int, error) {
	return _AeroV2Pool.Contract.Quote(&_AeroV2Pool.CallOpts, tokenIn, amountIn, granularity)
}

// Quote is a free data retrieval call binding the contract method 0x9e8cc04b.
//
// Solidity: function quote(address tokenIn, uint256 amountIn, uint256 granularity) view returns(uint256 amountOut)
func (_AeroV2Pool *AeroV2PoolCallerSession) Quote(tokenIn common.Address, amountIn *big.Int, granularity *big.Int) (*big.Int, error) {
	return _AeroV2Pool.Contract.Quote(&_AeroV2Pool.CallOpts, tokenIn, amountIn, granularity)
}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Reserve0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "reserve0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Reserve0() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve0(&_AeroV2Pool.CallOpts)
}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Reserve0() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve0(&_AeroV2Pool.CallOpts)
}

// Reserve0CumulativeLast is a free data retrieval call binding the contract method 0xbf944dbc.
//
// Solidity: function reserve0CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Reserve0CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "reserve0CumulativeLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve0CumulativeLast is a free data retrieval call binding the contract method 0xbf944dbc.
//
// Solidity: function reserve0CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Reserve0CumulativeLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve0CumulativeLast(&_AeroV2Pool.CallOpts)
}

// Reserve0CumulativeLast is a free data retrieval call binding the contract method 0xbf944dbc.
//
// Solidity: function reserve0CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Reserve0CumulativeLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve0CumulativeLast(&_AeroV2Pool.CallOpts)
}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Reserve1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "reserve1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Reserve1() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve1(&_AeroV2Pool.CallOpts)
}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Reserve1() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve1(&_AeroV2Pool.CallOpts)
}

// Reserve1CumulativeLast is a free data retrieval call binding the contract method 0xc245febc.
//
// Solidity: function reserve1CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) Reserve1CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "reserve1CumulativeLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve1CumulativeLast is a free data retrieval call binding the contract method 0xc245febc.
//
// Solidity: function reserve1CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) Reserve1CumulativeLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve1CumulativeLast(&_AeroV2Pool.CallOpts)
}

// Reserve1CumulativeLast is a free data retrieval call binding the contract method 0xc245febc.
//
// Solidity: function reserve1CumulativeLast() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) Reserve1CumulativeLast() (*big.Int, error) {
	return _AeroV2Pool.Contract.Reserve1CumulativeLast(&_AeroV2Pool.CallOpts)
}

// Sample is a free data retrieval call binding the contract method 0x13345fe1.
//
// Solidity: function sample(address tokenIn, uint256 amountIn, uint256 points, uint256 window) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolCaller) Sample(opts *bind.CallOpts, tokenIn common.Address, amountIn *big.Int, points *big.Int, window *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "sample", tokenIn, amountIn, points, window)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// Sample is a free data retrieval call binding the contract method 0x13345fe1.
//
// Solidity: function sample(address tokenIn, uint256 amountIn, uint256 points, uint256 window) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolSession) Sample(tokenIn common.Address, amountIn *big.Int, points *big.Int, window *big.Int) ([]*big.Int, error) {
	return _AeroV2Pool.Contract.Sample(&_AeroV2Pool.CallOpts, tokenIn, amountIn, points, window)
}

// Sample is a free data retrieval call binding the contract method 0x13345fe1.
//
// Solidity: function sample(address tokenIn, uint256 amountIn, uint256 points, uint256 window) view returns(uint256[])
func (_AeroV2Pool *AeroV2PoolCallerSession) Sample(tokenIn common.Address, amountIn *big.Int, points *big.Int, window *big.Int) ([]*big.Int, error) {
	return _AeroV2Pool.Contract.Sample(&_AeroV2Pool.CallOpts, tokenIn, amountIn, points, window)
}

// Stable is a free data retrieval call binding the contract method 0x22be3de1.
//
// Solidity: function stable() view returns(bool)
func (_AeroV2Pool *AeroV2PoolCaller) Stable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "stable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Stable is a free data retrieval call binding the contract method 0x22be3de1.
//
// Solidity: function stable() view returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) Stable() (bool, error) {
	return _AeroV2Pool.Contract.Stable(&_AeroV2Pool.CallOpts)
}

// Stable is a free data retrieval call binding the contract method 0x22be3de1.
//
// Solidity: function stable() view returns(bool)
func (_AeroV2Pool *AeroV2PoolCallerSession) Stable() (bool, error) {
	return _AeroV2Pool.Contract.Stable(&_AeroV2Pool.CallOpts)
}

// SupplyIndex0 is a free data retrieval call binding the contract method 0x9f767c88.
//
// Solidity: function supplyIndex0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) SupplyIndex0(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "supplyIndex0", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyIndex0 is a free data retrieval call binding the contract method 0x9f767c88.
//
// Solidity: function supplyIndex0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) SupplyIndex0(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.SupplyIndex0(&_AeroV2Pool.CallOpts, arg0)
}

// SupplyIndex0 is a free data retrieval call binding the contract method 0x9f767c88.
//
// Solidity: function supplyIndex0(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) SupplyIndex0(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.SupplyIndex0(&_AeroV2Pool.CallOpts, arg0)
}

// SupplyIndex1 is a free data retrieval call binding the contract method 0x205aabf1.
//
// Solidity: function supplyIndex1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) SupplyIndex1(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "supplyIndex1", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyIndex1 is a free data retrieval call binding the contract method 0x205aabf1.
//
// Solidity: function supplyIndex1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) SupplyIndex1(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.SupplyIndex1(&_AeroV2Pool.CallOpts, arg0)
}

// SupplyIndex1 is a free data retrieval call binding the contract method 0x205aabf1.
//
// Solidity: function supplyIndex1(address ) view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) SupplyIndex1(arg0 common.Address) (*big.Int, error) {
	return _AeroV2Pool.Contract.SupplyIndex1(&_AeroV2Pool.CallOpts, arg0)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AeroV2Pool *AeroV2PoolCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AeroV2Pool *AeroV2PoolSession) Symbol() (string, error) {
	return _AeroV2Pool.Contract.Symbol(&_AeroV2Pool.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AeroV2Pool *AeroV2PoolCallerSession) Symbol() (string, error) {
	return _AeroV2Pool.Contract.Symbol(&_AeroV2Pool.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_AeroV2Pool *AeroV2PoolCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_AeroV2Pool *AeroV2PoolSession) Token0() (common.Address, error) {
	return _AeroV2Pool.Contract.Token0(&_AeroV2Pool.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_AeroV2Pool *AeroV2PoolCallerSession) Token0() (common.Address, error) {
	return _AeroV2Pool.Contract.Token0(&_AeroV2Pool.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_AeroV2Pool *AeroV2PoolCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_AeroV2Pool *AeroV2PoolSession) Token1() (common.Address, error) {
	return _AeroV2Pool.Contract.Token1(&_AeroV2Pool.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_AeroV2Pool *AeroV2PoolCallerSession) Token1() (common.Address, error) {
	return _AeroV2Pool.Contract.Token1(&_AeroV2Pool.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() view returns(address, address)
func (_AeroV2Pool *AeroV2PoolCaller) Tokens(opts *bind.CallOpts) (common.Address, common.Address, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "tokens")

	if err != nil {
		return *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return out0, out1, err

}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() view returns(address, address)
func (_AeroV2Pool *AeroV2PoolSession) Tokens() (common.Address, common.Address, error) {
	return _AeroV2Pool.Contract.Tokens(&_AeroV2Pool.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() view returns(address, address)
func (_AeroV2Pool *AeroV2PoolCallerSession) Tokens() (common.Address, common.Address, error) {
	return _AeroV2Pool.Contract.Tokens(&_AeroV2Pool.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AeroV2Pool.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) TotalSupply() (*big.Int, error) {
	return _AeroV2Pool.Contract.TotalSupply(&_AeroV2Pool.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AeroV2Pool *AeroV2PoolCallerSession) TotalSupply() (*big.Int, error) {
	return _AeroV2Pool.Contract.TotalSupply(&_AeroV2Pool.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Approve(&_AeroV2Pool.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Approve(&_AeroV2Pool.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolTransactor) Burn(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "burn", to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolSession) Burn(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Burn(&_AeroV2Pool.TransactOpts, to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolTransactorSession) Burn(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Burn(&_AeroV2Pool.TransactOpts, to)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns(uint256 claimed0, uint256 claimed1)
func (_AeroV2Pool *AeroV2PoolTransactor) ClaimFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "claimFees")
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns(uint256 claimed0, uint256 claimed1)
func (_AeroV2Pool *AeroV2PoolSession) ClaimFees() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.ClaimFees(&_AeroV2Pool.TransactOpts)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns(uint256 claimed0, uint256 claimed1)
func (_AeroV2Pool *AeroV2PoolTransactorSession) ClaimFees() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.ClaimFees(&_AeroV2Pool.TransactOpts)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.DecreaseAllowance(&_AeroV2Pool.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.DecreaseAllowance(&_AeroV2Pool.TransactOpts, spender, subtractedValue)
}

// GetK is a paid mutator transaction binding the contract method 0xee39e7a0.
//
// Solidity: function getK() returns(uint256)
func (_AeroV2Pool *AeroV2PoolTransactor) GetK(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "getK")
}

// GetK is a paid mutator transaction binding the contract method 0xee39e7a0.
//
// Solidity: function getK() returns(uint256)
func (_AeroV2Pool *AeroV2PoolSession) GetK() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.GetK(&_AeroV2Pool.TransactOpts)
}

// GetK is a paid mutator transaction binding the contract method 0xee39e7a0.
//
// Solidity: function getK() returns(uint256)
func (_AeroV2Pool *AeroV2PoolTransactorSession) GetK() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.GetK(&_AeroV2Pool.TransactOpts)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.IncreaseAllowance(&_AeroV2Pool.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.IncreaseAllowance(&_AeroV2Pool.TransactOpts, spender, addedValue)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4bbb5a8.
//
// Solidity: function initialize(address _token0, address _token1, bool _stable) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) Initialize(opts *bind.TransactOpts, _token0 common.Address, _token1 common.Address, _stable bool) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "initialize", _token0, _token1, _stable)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4bbb5a8.
//
// Solidity: function initialize(address _token0, address _token1, bool _stable) returns()
func (_AeroV2Pool *AeroV2PoolSession) Initialize(_token0 common.Address, _token1 common.Address, _stable bool) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Initialize(&_AeroV2Pool.TransactOpts, _token0, _token1, _stable)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4bbb5a8.
//
// Solidity: function initialize(address _token0, address _token1, bool _stable) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) Initialize(_token0 common.Address, _token1 common.Address, _stable bool) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Initialize(&_AeroV2Pool.TransactOpts, _token0, _token1, _stable)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_AeroV2Pool *AeroV2PoolTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_AeroV2Pool *AeroV2PoolSession) Mint(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Mint(&_AeroV2Pool.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_AeroV2Pool *AeroV2PoolTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Mint(&_AeroV2Pool.TransactOpts, to)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AeroV2Pool *AeroV2PoolSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Permit(&_AeroV2Pool.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Permit(&_AeroV2Pool.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string __name) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) SetName(opts *bind.TransactOpts, __name string) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "setName", __name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string __name) returns()
func (_AeroV2Pool *AeroV2PoolSession) SetName(__name string) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.SetName(&_AeroV2Pool.TransactOpts, __name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string __name) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) SetName(__name string) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.SetName(&_AeroV2Pool.TransactOpts, __name)
}

// SetSymbol is a paid mutator transaction binding the contract method 0xb84c8246.
//
// Solidity: function setSymbol(string __symbol) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) SetSymbol(opts *bind.TransactOpts, __symbol string) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "setSymbol", __symbol)
}

// SetSymbol is a paid mutator transaction binding the contract method 0xb84c8246.
//
// Solidity: function setSymbol(string __symbol) returns()
func (_AeroV2Pool *AeroV2PoolSession) SetSymbol(__symbol string) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.SetSymbol(&_AeroV2Pool.TransactOpts, __symbol)
}

// SetSymbol is a paid mutator transaction binding the contract method 0xb84c8246.
//
// Solidity: function setSymbol(string __symbol) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) SetSymbol(__symbol string) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.SetSymbol(&_AeroV2Pool.TransactOpts, __symbol)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) Skim(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "skim", to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_AeroV2Pool *AeroV2PoolSession) Skim(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Skim(&_AeroV2Pool.TransactOpts, to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) Skim(to common.Address) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Skim(&_AeroV2Pool.TransactOpts, to)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_AeroV2Pool *AeroV2PoolTransactor) Swap(opts *bind.TransactOpts, amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "swap", amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_AeroV2Pool *AeroV2PoolSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Swap(&_AeroV2Pool.TransactOpts, amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Swap(&_AeroV2Pool.TransactOpts, amount0Out, amount1Out, to, data)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_AeroV2Pool *AeroV2PoolTransactor) Sync(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "sync")
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_AeroV2Pool *AeroV2PoolSession) Sync() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Sync(&_AeroV2Pool.TransactOpts)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_AeroV2Pool *AeroV2PoolTransactorSession) Sync() (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Sync(&_AeroV2Pool.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Transfer(&_AeroV2Pool.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.Transfer(&_AeroV2Pool.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.TransferFrom(&_AeroV2Pool.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_AeroV2Pool *AeroV2PoolTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AeroV2Pool.Contract.TransferFrom(&_AeroV2Pool.TransactOpts, from, to, amount)
}

// AeroV2PoolApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AeroV2Pool contract.
type AeroV2PoolApprovalIterator struct {
	Event *AeroV2PoolApproval // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolApproval)
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
		it.Event = new(AeroV2PoolApproval)
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
func (it *AeroV2PoolApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolApproval represents a Approval event raised by the AeroV2Pool contract.
type AeroV2PoolApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AeroV2PoolApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolApprovalIterator{contract: _AeroV2Pool.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AeroV2PoolApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolApproval)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseApproval(log types.Log) (*AeroV2PoolApproval, error) {
	event := new(AeroV2PoolApproval)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the AeroV2Pool contract.
type AeroV2PoolBurnIterator struct {
	Event *AeroV2PoolBurn // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolBurn)
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
		it.Event = new(AeroV2PoolBurn)
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
func (it *AeroV2PoolBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolBurn represents a Burn event raised by the AeroV2Pool contract.
type AeroV2PoolBurn struct {
	Sender  common.Address
	To      common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed sender, address indexed to, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterBurn(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*AeroV2PoolBurnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolBurnIterator{contract: _AeroV2Pool.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed sender, address indexed to, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *AeroV2PoolBurn, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolBurn)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed sender, address indexed to, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseBurn(log types.Log) (*AeroV2PoolBurn, error) {
	event := new(AeroV2PoolBurn)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the AeroV2Pool contract.
type AeroV2PoolClaimIterator struct {
	Event *AeroV2PoolClaim // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolClaim)
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
		it.Event = new(AeroV2PoolClaim)
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
func (it *AeroV2PoolClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolClaim represents a Claim event raised by the AeroV2Pool contract.
type AeroV2PoolClaim struct {
	Sender    common.Address
	Recipient common.Address
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterClaim(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*AeroV2PoolClaimIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Claim", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolClaimIterator{contract: _AeroV2Pool.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *AeroV2PoolClaim, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Claim", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolClaim)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Claim", log); err != nil {
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

// ParseClaim is a log parse operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseClaim(log types.Log) (*AeroV2PoolClaim, error) {
	event := new(AeroV2PoolClaim)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the AeroV2Pool contract.
type AeroV2PoolEIP712DomainChangedIterator struct {
	Event *AeroV2PoolEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolEIP712DomainChanged)
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
		it.Event = new(AeroV2PoolEIP712DomainChanged)
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
func (it *AeroV2PoolEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolEIP712DomainChanged represents a EIP712DomainChanged event raised by the AeroV2Pool contract.
type AeroV2PoolEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AeroV2Pool *AeroV2PoolFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*AeroV2PoolEIP712DomainChangedIterator, error) {

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolEIP712DomainChangedIterator{contract: _AeroV2Pool.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AeroV2Pool *AeroV2PoolFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *AeroV2PoolEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolEIP712DomainChanged)
				if err := _AeroV2Pool.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AeroV2Pool *AeroV2PoolFilterer) ParseEIP712DomainChanged(log types.Log) (*AeroV2PoolEIP712DomainChanged, error) {
	event := new(AeroV2PoolEIP712DomainChanged)
	if err := _AeroV2Pool.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolFeesIterator is returned from FilterFees and is used to iterate over the raw logs and unpacked data for Fees events raised by the AeroV2Pool contract.
type AeroV2PoolFeesIterator struct {
	Event *AeroV2PoolFees // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolFeesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolFees)
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
		it.Event = new(AeroV2PoolFees)
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
func (it *AeroV2PoolFeesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolFeesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolFees represents a Fees event raised by the AeroV2Pool contract.
type AeroV2PoolFees struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFees is a free log retrieval operation binding the contract event 0x112c256902bf554b6ed882d2936687aaeb4225e8cd5b51303c90ca6cf43a8602.
//
// Solidity: event Fees(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterFees(opts *bind.FilterOpts, sender []common.Address) (*AeroV2PoolFeesIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Fees", senderRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolFeesIterator{contract: _AeroV2Pool.contract, event: "Fees", logs: logs, sub: sub}, nil
}

// WatchFees is a free log subscription operation binding the contract event 0x112c256902bf554b6ed882d2936687aaeb4225e8cd5b51303c90ca6cf43a8602.
//
// Solidity: event Fees(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchFees(opts *bind.WatchOpts, sink chan<- *AeroV2PoolFees, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Fees", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolFees)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Fees", log); err != nil {
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

// ParseFees is a log parse operation binding the contract event 0x112c256902bf554b6ed882d2936687aaeb4225e8cd5b51303c90ca6cf43a8602.
//
// Solidity: event Fees(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseFees(log types.Log) (*AeroV2PoolFees, error) {
	event := new(AeroV2PoolFees)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Fees", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the AeroV2Pool contract.
type AeroV2PoolMintIterator struct {
	Event *AeroV2PoolMint // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolMint)
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
		it.Event = new(AeroV2PoolMint)
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
func (it *AeroV2PoolMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolMint represents a Mint event raised by the AeroV2Pool contract.
type AeroV2PoolMint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterMint(opts *bind.FilterOpts, sender []common.Address) (*AeroV2PoolMintIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolMintIterator{contract: _AeroV2Pool.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *AeroV2PoolMint, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolMint)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseMint(log types.Log) (*AeroV2PoolMint, error) {
	event := new(AeroV2PoolMint)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the AeroV2Pool contract.
type AeroV2PoolSwapIterator struct {
	Event *AeroV2PoolSwap // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolSwap)
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
		it.Event = new(AeroV2PoolSwap)
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
func (it *AeroV2PoolSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolSwap represents a Swap event raised by the AeroV2Pool contract.
type AeroV2PoolSwap struct {
	Sender     common.Address
	To         common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xb3e2773606abfd36b5bd91394b3a54d1398336c65005baf7bf7a05efeffaf75b.
//
// Solidity: event Swap(address indexed sender, address indexed to, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*AeroV2PoolSwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolSwapIterator{contract: _AeroV2Pool.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xb3e2773606abfd36b5bd91394b3a54d1398336c65005baf7bf7a05efeffaf75b.
//
// Solidity: event Swap(address indexed sender, address indexed to, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *AeroV2PoolSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolSwap)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0xb3e2773606abfd36b5bd91394b3a54d1398336c65005baf7bf7a05efeffaf75b.
//
// Solidity: event Swap(address indexed sender, address indexed to, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseSwap(log types.Log) (*AeroV2PoolSwap, error) {
	event := new(AeroV2PoolSwap)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolSyncIterator is returned from FilterSync and is used to iterate over the raw logs and unpacked data for Sync events raised by the AeroV2Pool contract.
type AeroV2PoolSyncIterator struct {
	Event *AeroV2PoolSync // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolSyncIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolSync)
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
		it.Event = new(AeroV2PoolSync)
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
func (it *AeroV2PoolSyncIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolSyncIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolSync represents a Sync event raised by the AeroV2Pool contract.
type AeroV2PoolSync struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSync is a free log retrieval operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 reserve0, uint256 reserve1)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterSync(opts *bind.FilterOpts) (*AeroV2PoolSyncIterator, error) {

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolSyncIterator{contract: _AeroV2Pool.contract, event: "Sync", logs: logs, sub: sub}, nil
}

// WatchSync is a free log subscription operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 reserve0, uint256 reserve1)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchSync(opts *bind.WatchOpts, sink chan<- *AeroV2PoolSync) (event.Subscription, error) {

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolSync)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Sync", log); err != nil {
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

// ParseSync is a log parse operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 reserve0, uint256 reserve1)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseSync(log types.Log) (*AeroV2PoolSync, error) {
	event := new(AeroV2PoolSync)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Sync", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AeroV2PoolTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AeroV2Pool contract.
type AeroV2PoolTransferIterator struct {
	Event *AeroV2PoolTransfer // Event containing the contract specifics and raw log

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
func (it *AeroV2PoolTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AeroV2PoolTransfer)
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
		it.Event = new(AeroV2PoolTransfer)
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
func (it *AeroV2PoolTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AeroV2PoolTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AeroV2PoolTransfer represents a Transfer event raised by the AeroV2Pool contract.
type AeroV2PoolTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AeroV2PoolTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AeroV2PoolTransferIterator{contract: _AeroV2Pool.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AeroV2PoolTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AeroV2Pool.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AeroV2PoolTransfer)
				if err := _AeroV2Pool.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AeroV2Pool *AeroV2PoolFilterer) ParseTransfer(log types.Log) (*AeroV2PoolTransfer, error) {
	event := new(AeroV2PoolTransfer)
	if err := _AeroV2Pool.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
