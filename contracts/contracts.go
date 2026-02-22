package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// https://docs.base.org/docs/contracts

func Pack(f func() (*abi.ABI, error), method string, params ...interface{}) []byte {
	a, err := f()
	if err != nil {
		panic(err)
	}

	// Otherwise pack up the parameters and invoke the contract
	input, err := a.Pack(method, params...)
	if err != nil {
		panic(err)
	}

	return input
}

func EncodeFun(abiStr string, method string, params ...interface{}) []byte {
	return Pack(func() (*abi.ABI, error) {
		metaData := bind.MetaData{
			ABI: abiStr,
		}

		return metaData.GetAbi()
	},
		method,
		params...,
	)
}
