package pool

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	defiabi "github.com/defiweb/go-eth/abi"
	defitypes "github.com/defiweb/go-eth/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	balanceOfMethodABI = defiabi.MustParseMethod("balanceOf(address)(uint256)")
	factoryMethodABI   = defiabi.MustParseMethod(`function factory() external view returns (address)`)
	reservesMethodABI  = defiabi.MustParseMethod(`function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast);`)
	token0MethodABI    = defiabi.MustParseMethod(`function token0() external view returns (address)`)
	token1MethodABI    = defiabi.MustParseMethod(`function token1() external view returns (address)`)

	aeroStableABI      = defiabi.MustParseMethod(`function stable() external view returns (bool)`)
	aeroGetFeeABI      = defiabi.MustParseMethod(`function getFee(address pool, bool _stable) public view returns (uint256)`)
	aeroGetReservesABI = defiabi.MustParseMethod(`function getReserves() external view returns (uint256 _reserve0, uint256 _reserve1, uint256 _blockTimestampLast);`)

	infusionGetReservesABI = defiabi.MustParseMethod(`function getReserves()
        external
        view
        returns (uint _reserve0, uint _reserve1, uint _blockTimestampLast);`)
	/*
			// AERO CLPool.
			feeMethodABI         = defiabi.MustParseMethod(`function fee() external view returns (uint24)`)
			tickSpacingMethodABI = defiabi.MustParseMethod(`function tickSpacing() external view returns (int24)`)
			liquidityMethodABI   = defiabi.MustParseMethod(`function liquidity() external view returns (uint128)`)
			// Uniswap V3 slot0.
			slot0MethodABI = defiabi.MustParseMethod(`function slot0()
		        external
		        view
		        returns (
		            uint160 sqrtPriceX96,
		            int24 tick,
		            uint16 observationIndex,
		            uint16 observationCardinality,
		            uint16 observationCardinalityNext,
		            uint8 feeProtocol,
		            bool unlocked
		        );`)

			// AERO slot0 https://github.com/aerodrome-finance/slipstream/blob/main/contracts/core/interfaces/pool/ICLPoolState.sol
			slot0AEROMethodABI = defiabi.MustParseMethod(`function slot0()
		        external
		        view
		        returns (
		            uint160 sqrtPriceX96,
		            int24 tick,
		            uint16 observationIndex,
		            uint16 observationCardinality,
		            uint16 observationCardinalityNext,
		            bool unlocked
		        );`)

			// PancakeSwap v3 slot0 https://github.com/pancakeswap/pancake-v3-contracts/tree/main/projects/v3-core/contracts/interfaces
			slot0PancakeMethodABI = defiabi.MustParseMethod(`function slot0()
			    external
			    view
			    returns (
			        uint160 sqrtPriceX96,
			        int24 tick,
			        uint16 observationIndex,
			        uint16 observationCardinality,
			        uint16 observationCardinalityNext,
			        uint32 feeProtocol,
			        bool unlocked
			    );`)
	*/

	// uniswapQueryABI      = defiabi.MustParseJSON([]byte(poolquery.UniswapV3PoolQueryABI)).
	getAllTicksInWordABI = defiabi.MustParseMethod(`function getAllTicksInWord2(
        address poolAddr,
        uint32 words
    )
        external
        view
        returns (
            uint160 sqrtPriceX96,
            uint128 liquidity,
            int24 tickCur,
            int24 tickSpacing,
            uint24 fee,
            int24[] memory ticks,
            int128[] memory liquidityNets,
            uint128[] memory liquidityGross
        )`)

	// syncEventABI = defiabi.MustParseEvent(`event Sync(uint112 reserve0, uint112 reserve1)`).

	multicallAddr = defitypes.MustAddressFromHex("0xcA11bde05977b3631167028862bE2a173976CA11")
	multicallABI  = defiabi.MustParseSignatures(
		"struct Call { address target; bytes callData; }",
		"struct Call3 { address target; bool allowFailure; bytes callData; }",
		"struct Call3Value { address target; bool allowFailure; uint256 value; bytes callData; }",
		"struct Result { bool success; bytes returnData; }",
		"function aggregate(Call[] calldata calls) public payable returns (uint256 blockNumber, bytes[] memory returnData)",
		"function aggregate3(Call3[] calldata calls) public payable returns (Result[] memory returnData)",
		"function aggregate3Value(Call3Value[] calldata calls) public payable returns (Result[] memory returnData)",
		"function blockAndAggregate(Call[] calldata calls) public payable returns (uint256 blockNumber, bytes32 blockHash, Result[] memory returnData)",
		"function getBasefee() view returns (uint256 basefee)",
		"function getBlockHash(uint256 blockNumber) view returns (bytes32 blockHash)",
		"function getBlockNumber() view returns (uint256 blockNumber)",
		"function getChainId() view returns (uint256 chainid)",
		"function getCurrentBlockCoinbase() view returns (address coinbase)",
		"function getCurrentBlockDifficulty() view returns (uint256 difficulty)",
		"function getCurrentBlockGasLimit() view returns (uint256 gaslimit)",
		"function getCurrentBlockTimestamp() view returns (uint256 timestamp)",
		"function getEthBalance(address addr) view returns (uint256 balance)",
		"function getLastBlockHash() view returns (bytes32 blockHash)",
		"function tryAggregate(bool requireSuccess, Call[] calldata calls) public payable returns (Result[] memory returnData)",
		"function tryBlockAndAggregate(bool requireSuccess, Call[] calldata calls) public payable returns (uint256 blockNumber, bytes32 blockHash, Result[] memory returnData)",
	)
)

// Call3 multicall3.
type Call3 struct {
	Target       defitypes.Address `abi:"target"`
	AllowFailure bool              `abi:"allowFailure"`
	CallData     []byte            `abi:"callData"`
}

// Result multicall3 result.
type Result struct {
	Success    bool   `abi:"success"`
	ReturnData []byte `abi:"returnData"`
}

// UniswapSlot0Result Uniswap v3 pool slot0() return.
type UniswapSlot0Result struct {
	SqrtPriceX96               *big.Int
	Tick                       int32
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

// AeroSlot0Result Aero pool slot0() return.
type AeroSlot0Result struct {
	SqrtPriceX96               *big.Int
	Tick                       int32
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	Unlocked                   bool
}

// PancakeSlot0Result Pancake swap v3 pool slot0() return.
type PancakeSlot0Result struct {
	SqrtPriceX96               *big.Int
	Tick                       int32
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint32
	Unlocked                   bool
}

// ReloadV3PoolReserves reload v3 pool reserves.
func (pp *ProviderPool) ReloadV3PoolReserves(pls []*Pool) {
	params := []Call3{}
	for _, pl := range pls {
		params = append(params,
			Call3{
				Target:       defitypes.MustAddressFromHex(pl.Token0),
				AllowFailure: true,
				CallData:     balanceOfMethodABI.MustEncodeArgs(pl.Address),
			},
			// 3 token1 balance
			Call3{
				Target:       defitypes.MustAddressFromHex(pl.Token1),
				AllowFailure: true,
				CallData:     balanceOfMethodABI.MustEncodeArgs(pl.Address),
			})
	}

	params = append(params, Call3{
		Target:   multicallAddr,
		CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
	})

	results, err := pp.Multicall(params, 3, "ReloadV3PoolReserves", false)
	if err != nil {
		logger.Warn().Err(err).Msgf("ReloadV3PoolReserves failed, pools: %d", len(pls))
		return
	}

	var blocknumber uint64

	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[len(pls)*2].ReturnData, &blocknumber)

	for i := 0; i < len(pls); i++ {
		var (
			err0     error
			err1     error
			reserve0 = big.NewInt(0)
			reserve1 = big.NewInt(0)
			pl       = pls[i]
		)

		if results[i*2].Success && results[i*2+1].Success {
			if err0 = balanceOfMethodABI.DecodeValues(results[i*2].ReturnData, &reserve0); err0 != nil {
				logger.Warn().Err(err0).Str("pool", pl.Address).Str("token0", pl.Token0).Msg("get v3 pool reserve0 failed")
			}

			if err1 = balanceOfMethodABI.DecodeValues(results[i*2+1].ReturnData, &reserve1); err1 != nil {
				logger.Warn().Err(err1).Str("pool", pl.Address).Str("token1", pl.Token1).Msg("get v3 pool reserve1 failed")
			}

			if err0 == nil && err1 == nil {
				// only blocknumber equal can compare
				if pl.LastBlockUpdated == blocknumber {
					if pl.Reserve0.Cmp(reserve0) != 0 {
						// logger.Warn().Str("pool", pl.Address).
						// 	Str("poolType", pl.Typ.String()).
						// 	Msgf("v3 pool reserve0 diff: %v onchain=%v block=%d", pl.Reserve0, reserve0, blocknumber)

						pl.Reserve0 = reserve0
					}
					//  else {
					// 	logger.Info().Str("pool", pl.Address).
					// 		Str("poolType", pl.Typ.String()).
					// 		Msgf("v3 pool reserve0 equals: %v block=%d", reserve0, blocknumber)
					// }
					//  else {
					// 	logger.Info().Str("pool", pl.Address).
					// 		Str("poolType", pl.Typ.String()).
					// 		Msgf("v3 pool reserve0 equals: %v block=%d", reserve0, blocknumber)
					// }

					//  else {
					// 	logger.Info().Str("pool", pl.Address).
					// 		Str("poolType", pl.Typ.String()).
					// 		Msgf("v3 pool reserve1 equals: %v block=%d", reserve1, blocknumber)
					// }
					if pl.Reserve1.Cmp(reserve1) != 0 {
						// logger.Warn().Str("pool", pl.Address).
						// 	Str("poolType", pl.Typ.String()).
						// 	Msgf("v3 pool reserve1 diff: %v onchain=%v block=%d", pl.Reserve1, reserve1, blocknumber)

						pl.Reserve1 = reserve1
					}
				}
			}
		} else {
			logger.Warn().Str("pool", pl.Address).Msg("pool reserve0 or reserve1 call NOT success")
		}
	}
}

// GetPoolLiquidity get v3 pool liquidity from contract.
func (pp *ProviderPool) GetV3PoolLiquidity(pool *Pool, queryAddr string, words uint) error {
	// if words == 0 {
	// 	// [MIN_TICK, MAX_TICK]
	// 	tick := pool.Tick
	// 	tickSpacing := pool.TickSpacing
	// 	lwords := (tick-MIN_TICK)/tickSpacing/256 + 1
	// 	rwords := (MAX_TICK-tick)/tickSpacing/256 + 1
	// 	words = uint((lwords)<<16 | rwords)
	// 	println("lwords:", lwords, "rwords:", rwords, "words:", words)
	// }
	query := defitypes.MustAddressFromHex(queryAddr)
	params := []Call3{
		// 0 block number
		{
			Target:   multicallAddr,
			CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
		},
		// 1 getAllTicksInWord
		{
			Target:   query,
			CallData: getAllTicksInWordABI.MustEncodeArgs(pool.Address, words),
		},
		// 2 token0 balance
		{
			Target:   defitypes.MustAddressFromHex(pool.Token0),
			CallData: balanceOfMethodABI.MustEncodeArgs(pool.Address),
		},
		// 3 token1 balance
		{
			Target:   defitypes.MustAddressFromHex(pool.Token1),
			CallData: balanceOfMethodABI.MustEncodeArgs(pool.Address),
		},
	}

	results, err := pp.Multicall(params, 3, "GetV3PoolLiquidity", true)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("pool", pool.Address).
			Str("vendor", pool.Vendor).
			Uint("fee", pool.Fee).
			Uint("words", words).
			Msg("GetV3PoolLiquidity failed")

		return err
	}

	var (
		blocknumber    uint64
		tick           int
		tickSpacing    int
		fee            uint
		liquidity      = big.NewInt(0)
		sqrtPrice      = big.NewInt(0)
		reserve0       = big.NewInt(0)
		reserve1       = big.NewInt(0)
		tickList       []int
		liquidityNet   []*big.Int
		liquidityGross []*big.Int
	)

	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[0].ReturnData, &blocknumber)
	// println(results[1].Success)
	getAllTicksInWordABI.MustDecodeValues(
		results[1].ReturnData,
		&sqrtPrice,
		&liquidity,
		&tick,
		&tickSpacing,
		&fee,
		&tickList,
		&liquidityNet,
		&liquidityGross,
	)

	if err = balanceOfMethodABI.DecodeValues(results[2].ReturnData, &reserve0); err != nil {
		logger.Error().Err(err).Str("pool", pool.Address).Str("token0", pool.Token0).Msg("get pool reserve0 failed")
	}

	if err = balanceOfMethodABI.DecodeValues(results[3].ReturnData, &reserve1); err != nil {
		logger.Error().Err(err).Str("pool", pool.Address).Str("token1", pool.Token1).Msg("get pool reserve1 failed")
	}

	if len(tickList) != len(liquidityNet) {
		return fmt.Errorf("ticks list len NOT equal liquidityNet len")
	}

	pool.Reserve0 = reserve0
	pool.Reserve1 = reserve1
	pool.Liquidity = liquidity
	pool.Tick = tick
	pool.Fee = fee
	pool.TickSpacing = tickSpacing
	// pool.TickList = tickList
	pool.SqrtPriceX96 = sqrtPrice
	pool.Synced = false
	pool.InitBlock = blocknumber
	// pool.LastBlockUpdated = blocknumber
	pool.tickBitmap = map[int16]*big.Int{}
	pool.Initialized = true

	tickInfo := map[int]*TickInfo{}

	for i := 0; i < len(tickList); i++ {
		tick := tickList[i]
		if _, ok := tickInfo[tick]; !ok {
			tickInfo[tick] = &TickInfo{
				Tick:           tick,
				LiquidityNet:   liquidityNet[i],
				LiquidityGross: liquidityGross[i],
			}

			flipTickBitmap(pool.tickBitmap, tick, pool.TickSpacing)
		}
	}

	pool.Ticks = tickInfo

	return nil
}

// GetV2PoolLiquid get v2 pool liquidity: reserves.
func (pp *ProviderPool) GetV2PoolLiquid(pl *Pool) (err error) {
	addr := defitypes.MustAddressFromHex(pl.Address)

	var callParams []Call3

	if pl.Typ == common.PoolTypeAMM {
		callParams = []Call3{
			// 0 reserves
			{
				Target:   addr,
				CallData: reservesMethodABI.MustEncodeArgs(),
			},
		}
	} else if pl.Typ == common.PoolTypeAeroAMM {
		callParams = []Call3{
			// 0 reserves
			{
				Target:   addr,
				CallData: aeroGetReservesABI.MustEncodeArgs(),
			},
		}
	} else if pl.Typ == common.PoolTypeInfusionAMM {
		callParams = []Call3{
			// 0 reserves
			{
				Target:   addr,
				CallData: infusionGetReservesABI.MustEncodeArgs(),
			},
		}
	} else {
		logger.Fatal().Msgf("invalid pool type: %s %s", pl.Address, pl.Typ)
	}

	callParams = append(callParams, Call3{
		Target:   multicallAddr,
		CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
	})

	results, err := pp.Multicall(callParams, 3, "GetV2PoolLiquid", false)
	if err != nil {
		logger.Warn().Err(err).
			Str("pool", pl.Address).
			Str("venddor", pl.Vendor).
			Msg("get pool v2 reserves failed")

		return err
	}

	// Decode the result.
	var (
		reserve0    = big.NewInt(0)
		reserve1    = big.NewInt(0)
		blockNumber = big.NewInt(0)
	)

	if pl.Typ == common.PoolTypeAMM {
		reservesMethodABI.MustDecodeValues(results[0].ReturnData, &reserve0, &reserve1, nil)
	} else if pl.Typ == common.PoolTypeAeroAMM {
		aeroGetReservesABI.MustDecodeValues(results[0].ReturnData, &reserve0, &reserve1, nil)
	} else if pl.Typ == common.PoolTypeInfusionAMM {
		infusionGetReservesABI.MustDecodeValues(results[0].ReturnData, &reserve0, &reserve1, nil)
	}

	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[1].ReturnData, &blockNumber)

	pl.Reserve0 = reserve0
	pl.Reserve1 = reserve1
	pl.InitBlock = blockNumber.Uint64()
	pl.Synced = true // v2 pool synced is always true

	return
}

func (pp *ProviderPool) getPoolsBasicInfo(poolAddrs []string, maxRetry uint, pools map[string]*Pool) ([]string, error) {
	failed := []string{}
	params := []Call3{}

	for _, paddr := range poolAddrs {
		addr := defitypes.MustAddressFromHex(paddr)
		params = append(params, Call3{
			Target:       addr,
			AllowFailure: true,
			CallData:     factoryMethodABI.MustEncodeArgs(),
		},
			// token0
			Call3{
				Target:       addr,
				AllowFailure: true,
				CallData:     token0MethodABI.MustEncodeArgs(),
			},
			// token1
			Call3{
				Target:       addr,
				AllowFailure: true,
				CallData:     token1MethodABI.MustEncodeArgs(),
			})
	}

	params = append(params, Call3{
		Target:   multicallAddr,
		CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
	})

	results, err := pp.Multicall(params, maxRetry, "getPoolsBasicInfo", false)
	if err != nil {
		logger.Warn().Err(err).Int("pools", len(poolAddrs)).Msg("GetPoolsBasicInfo: multicall failed")
		return nil, err
	}

	blockNumber := big.NewInt(0)
	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[len(results)-1].ReturnData, &blockNumber)

	for i := 0; i < len(poolAddrs); i++ {
		var (
			factoryAddr string
			token0Addr  string
			token1Addr  string
		)

		addr := poolAddrs[i]

		if results[i*3].Success && results[i*3+1].Success && results[i*3+2].Success {
			if len(results[i*3].ReturnData) > 0 {
				factoryMethodABI.MustDecodeValues(results[i*3].ReturnData, &factoryAddr)
			}

			if len(results[i*3+1].ReturnData) > 0 {
				token0MethodABI.MustDecodeValues(results[i*3+1].ReturnData, &token0Addr)
			}

			if len(results[i*3+2].ReturnData) > 0 {
				token1MethodABI.MustDecodeValues(results[i*3+2].ReturnData, &token1Addr)
			}

			pool := &Pool{
				PoolInfo: PoolInfo{
					Address:   addr,
					Factory:   strings.ToLower(factoryAddr),
					Token0:    strings.ToLower(token0Addr),
					Token1:    strings.ToLower(token1Addr),
					InitBlock: blockNumber.Uint64(),
				},
			}
			pool.UpdateKey()
			pools[addr] = pool
		} else {
			logger.Warn().
				Str("pool", addr).
				Str("factoryData", hex.EncodeToString(results[i*3].ReturnData)).
				Str("token0Data", hex.EncodeToString(results[i*3+1].ReturnData)).
				Str("token1Data", hex.EncodeToString(results[i*3+2].ReturnData)).
				Msg("GetPoolsBasicInfo: get pool basic info failed")

			failed = append(failed, addr)
		}
	}

	return failed, nil
}

// GetPoolsBasicInfo get pools basic info: factory, token0, token1.
func (pp *ProviderPool) GetPoolsBasicInfo(poolAddresses []string, maxRetry uint) (map[string]*Pool, []string, error) {
	ts1 := time.Now().UnixMilli()
	pools := map[string]*Pool{}
	failed := []string{}
	step := 5000

	for k := 0; k < len(poolAddresses); k += step {
		end := k + step
		if end > len(poolAddresses) {
			end = len(poolAddresses)
		}

		poolAddrs := poolAddresses[k:end]

		stepFailed, err := pp.getPoolsBasicInfo(poolAddrs, maxRetry, pools)
		if err == nil {
			failed = append(failed, stepFailed...)
		}
	}

	ts2 := time.Now().UnixMilli()
	logger.Info().
		Int("pools", len(poolAddresses)).
		Int("failed", len(failed)).
		Msgf("++++++++++ GetPoolsBasicInfo: get pools basic info used: %d", ts2-ts1)

	return pools, failed, nil
}

// GetPoolBasicInfo get pool basic info: factory, token0, token1.
func (pp *ProviderPool) GetPoolBasicInfo(poolAddr string, maxRetry uint) (*Pool, error) {
	addr := defitypes.MustAddressFromHex(poolAddr)
	params := []Call3{
		// factory
		{
			Target:       addr,
			AllowFailure: true,
			CallData:     factoryMethodABI.MustEncodeArgs(),
		},
		// token0
		{
			Target:   addr,
			CallData: token0MethodABI.MustEncodeArgs(),
		},
		// token1
		{
			Target:   addr,
			CallData: token1MethodABI.MustEncodeArgs(),
		},
	}

	results, err := pp.Multicall(params, maxRetry, "GetPoolBasicInfo", false)
	if err != nil {
		logger.Warn().Err(err).Str("pool", poolAddr).Msg("GetPoolBasicInfo: get pool basic info failed")
		return nil, err
	}

	// Decode the result.
	var (
		factoryAddr string
		token0Addr  string
		token1Addr  string
	)

	if results[0].Success && len(results[0].ReturnData) > 0 {
		factoryMethodABI.MustDecodeValues(results[0].ReturnData, &factoryAddr)
	} else {
		logger.Warn().Msgf("get pool factory failed: pool=%s", poolAddr)
	}

	token0MethodABI.MustDecodeValues(results[1].ReturnData, &token0Addr)
	token1MethodABI.MustDecodeValues(results[2].ReturnData, &token1Addr)

	pool := &Pool{
		PoolInfo: PoolInfo{
			Address: poolAddr,
			Factory: strings.ToLower(factoryAddr),
			Token0:  strings.ToLower(token0Addr),
			Token1:  strings.ToLower(token1Addr),
		},
	}

	return pool, nil
}

// GetAeroPoolV2Stable get pool is stable pool.
func (pp *ProviderPool) GetAeroPoolV2Stable(pls []*Pool, maxRetry uint) error {
	callParams := []Call3{}

	for _, pl := range pls {
		callParams = append(callParams, Call3{
			Target:   defitypes.MustAddressFromHex(pl.Address),
			CallData: aeroStableABI.MustEncodeArgs(),
		})
	}

	results, err := pp.Multicall(callParams, maxRetry, "GetAeroPoolV2Stable", false)
	if err != nil {
		logger.Warn().Err(err).Msg("get aero pool v2 stable failed")
		return err
	}

	for i, result := range results {
		// pls[i].Version = LatestVersion
		// fmt.Println(pls[i].Address, ": ", hex.EncodeToString(result.ReturnData))
		aeroStableABI.MustDecodeValues(result.ReturnData, &pls[i].Stable)
	}

	return nil
}
