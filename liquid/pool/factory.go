package pool

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	defiabi "github.com/defiweb/go-eth/abi"
	defitypes "github.com/defiweb/go-eth/types"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	allPairsMethodABI = defiabi.MustParseMethod("function allPairs(uint) external view returns (address pair);")

	allPairsLengthMethodABI = defiabi.MustParseMethod("function allPairsLength() external view returns (uint)")
	allPoolsLengthMethodABI = defiabi.MustParseMethod("function allPoolsLength() external view returns (uint256)")

	getV3PositionsMethodABI = defiabi.MustParseMethod(`function getV3Positions(
        address positionManager
    ) public view returns (uint256)`)
	getV3PoolsMethodABI = defiabi.MustParseMethod(`function getV3Pools(
        address factory,
        address positionManager,
        uint from,
        uint to
    ) public view returns (address[] memory pools)`)
	getAeroV3PoolsMethodABI = defiabi.MustParseMethod(`function getAeroV3pools(
        address factory,
        uint from,
        uint howmany
    ) public view returns (address[] memory pools)`)

	getPairsInfoMethodABI = defiabi.MustParseMethod(`function getPairsInfo(
        address factory,
        uint from,
        uint howmany
    )
        public
        view
        returns (
            address[] memory pairAddrs,
            address[] memory token0s,
            address[] memory token1s,
            uint[] memory reserve0s,
            uint[] memory reserve1s
        )`)

	getAeroPairsInfoMethodABI = defiabi.MustParseMethod(`function getAeroPairsInfo(
        address factory,
        uint from,
        uint howmany
    )
        public
        view
        returns (
            address[] memory pairAddrs,
            address[] memory token0s,
            address[] memory token1s,
            uint[] memory reserve0s,
            uint[] memory reserve1s,
            bool[] memory stables,
            uint256[] memory fees
        )`)

	getV2PairsLiquidityMethodABI = defiabi.MustParseMethod(`function getV2PairsLiquidity(
        address[] calldata pairs
    ) public view returns (uint[] memory reserve0s, uint[] memory reserve1s)`)
	getAeroV2PoolsLiquidityMethodABI = defiabi.MustParseMethod(`function getAeroV2PoolsLiquidity(
        address[] calldata pairs
    ) public view returns (uint[] memory reserve0s, uint[] memory reserve1s)`)

	getInfusionPoolMethodABI = defiabi.MustParseMethod(`function metadata()
        external
        view
        returns (
            uint dec0,
            uint dec1,
            uint r0,
            uint r1,
            bool st,
            address t0,
            address t1
        );`)

	getInfusionFeeMethodABI = defiabi.MustParseMethod(`function getFee(bool _stable) public view returns (uint256)`)

	// infusion pair created.
	infusionPairCreatedEventABI = defiabi.MustParseEvent(`PairCreated(
        address indexed token0,
        address indexed token1,
        bool stable,
        address pair,
        uint
    )`)
	// aerov2PoolCreatedEventABI = defiabi.MustParseEvent("PoolCreated(address indexed token0, address indexed token1, bool indexed stable, address pool, uint256)").

	setCustomeFeeV2EventABI = defiabi.MustParseEvent(`SetCustomFee(address indexed,uint256)`)
	setCustomeFeeV3EventABI = defiabi.MustParseEvent(`SetCustomFee(address indexed,uint24 indexed)`)
	ZERO_FEE_INDICATOR      = 420
)

// GetV2PoolInfos get uniswap v2 all pools, every pool info and reserves.
func (pp *ProviderPool) GetV2PoolInfos(query string,
	factory *common.SwapFactory,
	step int,
	parallel bool,
) (map[string]*Pool, error) {
	if !factory.Typ.IsAMMVariety() {
		return nil, fmt.Errorf("factory is NOT AMM variety")
	}

	factory.Address = strings.ToLower(factory.Address)
	ts1 := time.Now().UnixMilli()
	provider := pp.Get()

	pairsLen, err := getPairLength(provider, factory.Address, factory.Typ)
	if err != nil {
		logger.Error().Err(err).Str("factory", factory.Address).Msg("GetV2PoolInfos: get pair length failed")
		return nil, err
	}

	logger.Info().Str("vendor", factory.Name).Msgf("load v2 pools: %d", pairsLen)

	if step == 0 {
		step = 3000
	}

	if factory.Typ == common.PoolTypeInfusionAMM {
		return pp.getInfusionPools(factory, pairsLen, uint(step)) // nolint
	}

	tokensMap := map[string]bool{}
	pools := map[string]*Pool{}
	lk := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 100)
	routines := 0
	maxRetry := uint(5)
	queryAddr := defitypes.MustAddressFromHex(query)

	maxRoutines := provider.Tps
	if maxRoutines > 100 {
		maxRoutines = 100
	}

	callFn := func(from, howmany int) {
		stokens, spools, err := pp.getV2Pools(factory, queryAddr, from, howmany, maxRetry)
		if err != nil {
			errCh <- err
			return
		}

		if parallel {
			lk.Lock()
		}

		for addr := range stokens {
			tokensMap[addr] = true
		}

		for addr, p := range spools {
			pools[addr] = p
		}

		if parallel {
			lk.Unlock()
			wg.Done()
		}
	}

	for i := 0; i < int(pairsLen); i += step { // nolint
		howmany := step
		if i+step >= int(pairsLen) { // nolint
			howmany = int(pairsLen) - i - 1 // nolint
		}

		if parallel {
			wg.Add(1)

			go callFn(i, howmany)
		} else {
			callFn(i, howmany)
		}

		routines++
		if routines >= int(maxRoutines) { // nolint
			wg.Wait() // no parallel wait is also ok, just make code clean

			routines = 0
			// sleep 1 second whatever
			time.Sleep(time.Second)
		}

		if len(errCh) > 0 {
			err := <-errCh
			return nil, err
		}
	}

	wg.Wait() // no parallel wait is also ok, just make code clean

	ts2 := time.Now().UnixMilli()
	logger.Info().
		Str("factory", factory.Name).
		Msgf("load all v2 pairs info and reserves used: %d ms, pairs: %d tokens: %d",
			ts2-ts1, pairsLen, len(tokensMap))

	return pools, nil
}

// uniswap v2 or aerodrome v1 pools.
func (pp *ProviderPool) getV2Pools(factory *common.SwapFactory,
	queryAddr defitypes.Address,
	start, howmany int,
	maxRetry uint,
) (map[string]bool, map[string]*Pool, error) {
	tokensMap := map[string]bool{}
	pools := map[string]*Pool{}

	// Call method.
	params := []Call3{
		// 0 block number
		{
			Target:   multicallAddr,
			CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
		},
	}

	if factory.Typ == common.PoolTypeAeroAMM {
		params = append(params, Call3{
			Target:   queryAddr,
			CallData: getAeroPairsInfoMethodABI.MustEncodeArgs(factory.Address, start, howmany),
		})
	} else {
		params = append(params, Call3{
			Target:   queryAddr,
			CallData: getPairsInfoMethodABI.MustEncodeArgs(factory.Address, start, howmany),
		})
	}

	results, err := pp.Multicall(params, maxRetry, fmt.Sprintf("GetV2PoolInfos-%d-%d", start, howmany), false)
	if err != nil {
		logger.Error().
			Str("vendor", factory.Name).
			Err(err).
			Msgf("GetV2PoolInfos failed: from: %d howmany: %d", start, howmany)
		time.Sleep(time.Second)

		return nil, nil, err
	}

	var (
		blocknumber uint64
		pairs       []string
		token0s     []string
		token1s     []string
		reserve0s   []*big.Int
		reserve1s   []*big.Int
		stables     []bool
		fees        []*big.Int
	)

	// Decode the results.
	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[0].ReturnData, &blocknumber)

	if factory.Typ == common.PoolTypeAeroAMM {
		getAeroPairsInfoMethodABI.MustDecodeValues(results[1].ReturnData,
			&pairs, &token0s, &token1s, &reserve0s, &reserve1s, &stables, &fees)
	} else {
		getPairsInfoMethodABI.MustDecodeValues(results[1].ReturnData,
			&pairs, &token0s, &token1s, &reserve0s, &reserve1s)
	}

	for i := range pairs {
		pairAddr := strings.ToLower(pairs[i])
		token0 := strings.ToLower(token0s[i])
		token1 := strings.ToLower(token1s[i])
		pool := &Pool{
			PoolInfo: PoolInfo{
				Address:          pairAddr,
				Token0:           token0,
				Token1:           token1,
				Factory:          factory.Address,
				Vendor:           factory.Name,
				Typ:              factory.Typ,
				Fee:              uint(factory.Fee), // nolint
				TickSpacing:      0,
				LastBlockUpdated: 0,
				InitBlock:        blocknumber,
			},
			Reserve0: reserve0s[i],
			Reserve1: reserve1s[i],
		}
		tokensMap[token0] = true
		tokensMap[token1] = true

		if factory.Typ == common.PoolTypeAeroAMM {
			pool.Stable = stables[i]
			pool.Fee = uint(fees[i].Uint64() * 100) // convert to 1000000 denominator
		}

		pools[pairAddr] = pool
	}

	return tokensMap, pools, nil
}

// infusion pools.
func (pp *ProviderPool) getInfusionPools(factory *common.SwapFactory, pairs, step uint) (map[string]*Pool, error) {
	pools := map[string]*Pool{}

	pp.GetInfusionSwapFee(factory)

	factoryAddr := defitypes.MustAddressFromHex(factory.Address)

	for i := uint(0); i < pairs; {
		end := i + step
		if end >= pairs {
			end = pairs
		}

		params := []Call3{}
		for j := i; j < end; j++ {
			params = append(params, Call3{
				Target:   factoryAddr,
				CallData: allPairsMethodABI.MustEncodeArgs(j),
			})
		}

		results, err := pp.Multicall(params, 3, fmt.Sprintf("getInfusionPools-%d-%d", i, end), false)
		if err != nil {
			logger.Warn().
				Err(err).
				Msgf("getInfusionPools failed: from: %d end: %d", i, end)
			time.Sleep(time.Second)

			continue
		}

		pairs := []string{}

		for k := 0; k < len(results); k++ {
			var pairAddr string

			allPairsMethodABI.MustDecodeValues(results[k].ReturnData, &pairAddr)
			pairs = append(pairs, pairAddr)
		}

		poolsInfo, err := pp.getInfusionPoolInfo(factory, pairs, 3)
		if err == nil {
			for _, item := range poolsInfo {
				pools[item.Address] = item
			}
		}

		i += step
	}

	return pools, nil
}

func (pp *ProviderPool) getInfusionPoolInfo(factory *common.SwapFactory, addrs []string, maxRetry uint) ([]*Pool, error) {
	params := []Call3{}
	for _, addr := range addrs {
		params = append(params, Call3{
			Target:   defitypes.MustAddressFromHex(addr),
			CallData: getInfusionPoolMethodABI.MustEncodeArgs(),
		})
	}
	// block number
	params = append(params, Call3{
		Target:   multicallAddr,
		CallData: multicallABI.Methods["getBlockNumber"].MustEncodeArgs(),
	})

	results, err := pp.Multicall(params, maxRetry, "getInfusionPoolInfo", false)
	if err != nil {
		logger.Warn().
			Err(err).
			Msgf("getInfusionPoolInfo failed")

		return nil, err
	}

	var (
		blocknumber uint64
		pairs       = []*Pool{}
	)

	// Decode the results.
	multicallABI.Methods["getBlockNumber"].MustDecodeValues(results[len(results)-1].ReturnData, &blocknumber)

	for k := 0; k < len(results)-1; k++ {
		var (
			dec0, dec1, r0, r1 *big.Int
			st                 bool
			t0, t1             string
		)

		getInfusionPoolMethodABI.MustDecodeValues(results[k].ReturnData,
			&dec0, &dec1, &r0, &r1, &st, &t0, &t1)

		pl := &Pool{
			PoolInfo: PoolInfo{
				Typ:       factory.Typ,
				Factory:   factory.Address,
				Address:   addrs[k],
				decimals0: dec0,
				decimals1: dec1,
				Token0:    t0,
				Token1:    t1,
				Vendor:    factory.Name,
				Stable:    st,
				InitBlock: blocknumber,
			},
			Reserve0: r0,
			Reserve1: r1,
		}
		if st {
			pl.Fee = uint(factory.StableFee) // nolint
		} else {
			pl.Fee = uint(factory.Fee) // nolint
		}

		pairs = append(pairs, pl)
	}

	return pairs, nil
}

func getPairLength(provider *Provider, addr string, typ common.PoolType) (uint, error) {
	var pairsLen uint

	target := defitypes.MustAddressFromHex(addr)

	if typ == common.PoolTypeAeroAMM {
		call := defitypes.NewCall().
			SetTo(target).
			SetInput(allPoolsLengthMethodABI.MustEncodeArgs())

		// Call allPairsLength.
		b, _, err := provider.Call(context.Background(), call)
		if err != nil {
			return 0, err
		}

		allPoolsLengthMethodABI.MustDecodeValues(b, &pairsLen)
	} else {
		call := defitypes.NewCall().
			SetTo(target).
			SetInput(allPairsLengthMethodABI.MustEncodeArgs())

		// Call allPairsLength.
		b, _, err := provider.Call(context.Background(), call)
		if err != nil {
			return 0, err
		}

		allPairsLengthMethodABI.MustDecodeValues(b, &pairsLen)
	}

	return pairsLen, nil
}

// GetAeroV3Pools get aerodrome v3 pools.
func (pp *ProviderPool) GetAeroV3Pools(pairQuery, factory string,
	step int,
	parallel bool,
) ([]string, error) {
	provider := pp.GetV3()

	poolLen, err := getPairLength(provider, factory, common.PoolTypeAeroAMM) // aero v2, v3 allPoolsLength is same
	if err != nil {
		logger.Error().Err(err).Str("factory", factory).Msg("GetAeroV3Pools: get pool length failed")
		return nil, err
	}

	logger.Info().Msgf("GetAeroV3Pools: pools: %d", poolLen)

	call := defitypes.NewCall().
		SetTo(defitypes.MustAddressFromHex(pairQuery)).
		SetInput(getAeroV3PoolsMethodABI.MustEncodeArgs(factory, 0, 0))

	// Call getAeroV3Pools.
	b, _, err := provider.Call(context.Background(), call)
	if err != nil {
		return nil, err
	}

	var pools []string

	getAeroV3PoolsMethodABI.MustDecodeValues(b, &pools)

	return pools, nil
}

// GetV3Pools get uniswap v3 and variant swap pools.
func (pp *ProviderPool) GetV3Pools(pairQuery, factory, positionManager, name string,
	step int,
	parallel bool,
) ([]string, error) {
	call := defitypes.NewCall().
		SetTo(defitypes.MustAddressFromHex(pairQuery)).
		SetInput(getV3PositionsMethodABI.MustEncodeArgs(positionManager))

	// Call .
	b, _, err := pp.Get().Call(context.Background(), call)
	if err != nil {
		return nil, err
	}

	var howmany uint64

	getV3PositionsMethodABI.MustDecodeValues(b, &howmany)
	logger.Info().Str("factory", factory).Str("name", name).Msgf("v3 positions: %d", howmany)

	poolsMap := map[string]bool{}
	pools := []string{}
	lk := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1000)
	maxRoutines := 100
	routines := 0
	// Call method.
	callFn := func(cl *defitypes.Call, from, to int) {
		provider := pp.Get(true)
		// Call allPairs.
		b, _, err = provider.Call(context.Background(), cl)
		if err != nil {
			errCh <- err
			return
		}

		var items []string

		getV3PoolsMethodABI.MustDecodeValues(b, &items)
		log.Debug().Str("factory", factory).Msgf("position tokens from %d to %d: %d", from, to, len(items))

		if parallel {
			lk.Lock()
		}

		for _, item := range items {
			addr := strings.ToLower(item)
			if !poolsMap[addr] {
				pools = append(pools, addr)
				poolsMap[addr] = true
			}
		}

		if parallel {
			lk.Unlock()
			wg.Done()
		}
	}

	if step == 0 {
		step = 10000
	}

	for i := 0; i < int(howmany); i += step { // nolint
		end := i + step
		if end > int(howmany) { // nolint
			end = int(howmany) // nolint
		}

		call = defitypes.NewCall().
			SetTo(defitypes.MustAddressFromHex(pairQuery)).
			SetInput(getV3PoolsMethodABI.MustEncodeArgs(factory, positionManager, i, end))

		if parallel {
			wg.Add(1)

			go callFn(call, i, end)
		} else {
			callFn(call, i, end)
		}

		routines++
		if routines >= maxRoutines {
			logger.Info().Msg("GetV3Pools: reach max routines, wait routines complete")
			wg.Wait()

			routines = 0
			// sleep 1 second whatever
			time.Sleep(time.Second)
		}

		if len(errCh) > 0 {
			err := <-errCh
			return nil, err
		}
	}

	// if not parallel, wait is also ok
	wg.Wait()

	return pools, nil
}

// GetV2PoolAndLiquids get uniswap v2 pools and every pool liquidity.
func (pp *ProviderPool) GetV2PoolReserves(query string,
	name string,
	typ common.PoolType,
	pairs []*Pool,
	step int,
	parallel bool,
) error {
	if !typ.IsAMMVariety() {
		return fmt.Errorf("NOT v2 or variant factory")
	}

	ts1 := time.Now().UnixMilli()
	provider := pp.Get()
	lk := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 100)
	routines := 0
	pairsLen := len(pairs)
	queryAddr := defitypes.MustAddressFromHex(query)
	maxRoutines := provider.Tps
	maxRetry := 3
	// Call method.
	callFn := func(start int, howmany int, cl *defitypes.Call) {
		var (
			b   []byte
			err error
		)

		for retry := 0; retry < maxRetry; retry++ {
			b, _, err = provider.Call(context.Background(), cl)
			if err != nil {
				logger.Warn().Str("vendor", name).
					Err(err).
					Msgf("GetV2PoolReserves failed: from: %d howmany: %d retry: %d", start, howmany, retry)
				time.Sleep(time.Second)

				continue
			}

			var (
				reserve0s []*big.Int
				reserve1s []*big.Int
			)

			if typ == common.PoolTypeAeroAMM {
				getAeroV2PoolsLiquidityMethodABI.MustDecodeValues(b, &reserve0s, &reserve1s)
			} else {
				getV2PairsLiquidityMethodABI.MustDecodeValues(b, &reserve0s, &reserve1s)
			}

			if parallel {
				lk.Lock()
			}

			for k := 0; k < howmany; k++ {
				pairs[start+k].Reserve0 = reserve0s[k]
				pairs[start+k].Reserve1 = reserve1s[k]
			}

			if parallel {
				lk.Unlock()
				wg.Done()
			}

			return
		}

		wg.Done()
		logger.Error().Str("vendor", name).
			Err(err).
			Msg("GetV2PoolReserves failed reach maxRetry times")
		errCh <- err
	}

	if maxRoutines > 100 {
		maxRoutines = 100
	}

	if step == 0 {
		step = 10000
	}

	for i := 0; i < pairsLen; i += step {
		howmany := step
		if i+step >= pairsLen {
			howmany = pairsLen - i - 1
		}

		var call *defitypes.Call
		if typ == common.PoolTypeAeroAMM {
			call = defitypes.NewCall().
				SetTo(queryAddr).
				SetInput(getAeroPairsInfoMethodABI.MustEncodeArgs(pairs[i : i+howmany]))
		} else {
			call = defitypes.NewCall().
				SetTo(defitypes.MustAddressFromHex(query)).
				SetInput(getPairsInfoMethodABI.MustEncodeArgs(pairs[i : i+howmany]))
		}

		if parallel {
			wg.Add(1)

			go callFn(i, step, call)
		} else {
			callFn(i, step, call)
		}

		routines++
		if routines >= int(maxRoutines) { // nolint
			wg.Wait()

			routines = 0
			// sleep 1 second whatever
			time.Sleep(time.Second)
		}

		if len(errCh) > 0 {
			err := <-errCh
			return err
		}
	}

	wg.Wait()

	ts2 := time.Now().UnixMilli()
	logger.Info().Str("vendor", name).
		Msgf("load all v2 pairs reserves used: %d ms, pairs: %d", ts2-ts1, pairsLen)

	return nil
}

func (pp *ProviderPool) DecodeSetCustomFeeV2(evt *types.Log) (string, int, error) {
	var (
		addr string
		fee  int
	)

	err := setCustomeFeeV2EventABI.DecodeValues(convertTopics(evt.Topics), evt.Data, &addr, &fee)

	return addr, fee, err
}

func (pp *ProviderPool) DecodeSetCustomFeeV3(evt *types.Log) (string, int, error) {
	var (
		addr string
		fee  int
	)

	err := setCustomeFeeV3EventABI.DecodeValues(convertTopics(evt.Topics), evt.Data, &addr, &fee)

	return addr, fee, err
}

func convertTopics(hashes []ethcomm.Hash) (hs []defitypes.Hash) {
	for _, h := range hashes {
		hs = append(hs, defitypes.Hash(h.Bytes()))
	}

	return
}

func DecodeInfusionPairCreatedEvent(fac *common.SwapFactory, evt *types.Log) (*Pool, error) {
	var (
		token0   string
		token1   string
		pairAddr string
		stable   bool
		total    uint
	)

	err := infusionPairCreatedEventABI.DecodeValues(convertTopics(evt.Topics),
		evt.Data, &token0, &token1, &stable, &pairAddr, &total)
	if err != nil {
		logger.Error().Err(err).Str("txhash", evt.TxHash.String()).Msg("decode infusion PairCreated event failed")
		return nil, err
	}

	pl := &Pool{
		PoolInfo: PoolInfo{
			Typ:       fac.Typ,
			Address:   pairAddr,
			Factory:   fac.Address,
			Vendor:    fac.Name,
			Stable:    stable,
			Token0:    token0,
			Token1:    token1,
			InitBlock: evt.BlockNumber,
		},
		Reserve0: big.NewInt(0),
		Reserve1: big.NewInt(0),
	}
	if stable {
		pl.Fee = uint(fac.StableFee) // nolint
	} else {
		pl.Fee = uint(fac.Fee) // nolint
	}

	return pl, nil
}
