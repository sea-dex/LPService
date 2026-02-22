package arb

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/liquid/utils"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	minTokenETHReserve = utils.ToBigIntMust("10000000000000")      // 0.00001 ETH, about 0.02-0.03USD
	minPoolETH         = utils.ToBigIntMust("30000000000000000")   // 0.03 ETH
	minMidPoolETH      = utils.ToBigIntMust("1000000000000000000") // 1 ETH
	//lint:ignore U1000 Ignore unused function temporarily for debugging
	q96     = decimal.NewFromBigInt(pool.Q96, 0)
	d10     = decimal.NewFromInt(10)
	bigZero = big.NewInt(0)
	//lint:ignore U1000 Ignore unused function temporarily for debugging
	bigTwo = big.NewInt(2)
	bigTen = big.NewInt(10)
	//lint:ignore U1000 Ignore unused function temporarily for debugging
	big10000         = big.NewInt(10000)
	minMidTokenPools = uint(10)
	tokenMinTVLPools = int(2)
	p5               = decimal.NewFromFloat(0.05)
)

type PriceUpdated struct {
	LastBlock uint64
	Price     decimal.Decimal // price to ETH
}

type LPService struct {
	rc                   redis.Cmdable
	tokens               map[string]*common.Token // all tokens
	midTokens            map[string]*common.Token // middle tokens
	stableTokens         map[string]*common.Token // stable tokens
	pools                map[string]*pool.Pool
	midPools             map[string]*pool.Pool
	poolsGreatThanMinETH map[string]*pool.Pool
	flashPools           map[string]string
	blacklistTokens      map[string]bool
	nativeStablePair     string
	ethPrice             decimal.Decimal
	cfg                  *config.Config
	arbPairs             map[string]*ArbPairList
	midTokenPrice        map[string]*PriceUpdated
}

func CreateLPService(cfg *config.Config) *LPService {
	setDecimalPrecision(40)

	rc, err := common.CreatePoolRedisStore(false, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Fatal().Err(err).Msg("create redis connect failed")
		return nil
	}

	lpservice := &LPService{
		rc:              rc,
		stableTokens:    map[string]*common.Token{},
		tokens:          map[string]*common.Token{},
		midTokens:       map[string]*common.Token{},
		pools:           map[string]*pool.Pool{},
		flashPools:      map[string]string{},
		cfg:             cfg,
		blacklistTokens: map[string]bool{},
		midTokenPrice:   map[string]*PriceUpdated{},
	}

	return lpservice
}

func (lpservice *LPService) InitArb(arbCfg *config.ArbConfig) {
	minPoolETH = convertToBigInt(arbCfg.MinPoolETH, "0.03", "minPoolETH")
	minMidPoolETH = convertToBigInt(arbCfg.MinMidPoolETH, "10", "minMidPoolETH")

	for _, addr := range arbCfg.StablePools {
		addr = strings.ToLower(strings.TrimSpace(addr))
		lpservice.stableTokens[addr] = lpservice.tokens[addr]

		if lpservice.stableTokens[addr] == nil {
			logger.Fatal().Msgf("stable token %s not found", addr)
		}

		logger.Info().Msgf("set token %s %s as stable coin", lpservice.tokens[addr].Symbol, addr)
	}

	lpservice.nativeStablePair = strings.ToLower(strings.TrimSpace(arbCfg.NativeStablePair))
	if lpservice.nativeStablePair == "" {
		panic("invalid native stable pair")
	}

	logger.Info().Msgf("set native token pool to %s", lpservice.nativeStablePair)

	for _, addr := range arbCfg.BlackListTokens {
		addr = strings.TrimSpace(strings.ToLower(addr))
		lpservice.blacklistTokens[addr] = true
	}
}

func (lpservice *LPService) SetTokenPools(
	tokens map[string]*common.Token,
	pools map[string]*pool.Pool,
) {
	lpservice.tokens = tokens
	lpservice.pools = pools
}

func (lpservice *LPService) InitPools() {
	if len(lpservice.tokens) == 0 || len(lpservice.pools) == 0 {
		logger.Fatal().Msg("no tokens or no pools")
	}

	lpservice.InitArb(&lpservice.cfg.Arb)
	lpservice.refreshETHPrice()
	lpservice.recalculateMidTokenPools()
	lpservice.recalculateMidTokenPools()
}

func (lpservice *LPService) UpdateArbPairs(allTokens bool) {
	var arbPairs map[string]*ArbPairList

	if allTokens {
		arbPairs = lpservice.DiscoverArbPairList()
	} else {
		lpservice.tokens = lpservice.midTokens
		arbPairs = lpservice.DiscoverArbPairList()
	}

	lpservice.arbPairs = arbPairs
	logger.Info().Msgf("found arb pairs: allTokens=%v arbPairs=%d", allTokens, len(arbPairs))
}

func (lpservice *LPService) OnPoolUpdated(pls map[string]*pool.Pool, ratio decimal.Decimal, blocknumber uint64) {
	keys := map[string][]*pool.Pool{}

	tm0 := time.Now()
	for _, pl := range pls {
		if pl.Address == lpservice.nativeStablePair {
			lpservice.refreshETHPrice()
		}

		keys[pl.GetKey()] = append(keys[pl.GetKey()], pl)

		lpservice.updatePool(pl)
	}
	logger.Info().Msgf("update %d pools used: %v arb pairs: %d", len(pls), time.Since(tm0), len(keys))

	maxProfitUSD := decimal.Zero
	symbol := ""
	for key, pls := range keys {
		arb := lpservice.arbPairs[key]
		if arb != nil {
			start := time.Now()
			profitable, param, _ := lpservice.CalcProfitable(arb, pls, ratio, false)
			if profitable {
				profitUSD := lpservice.toUSD(param.tokenIn, param.bestProfit)
				logger.Info().Msgf("Arb list(%d/%d) was profitable: %s amountIn=%v profit=%v used: %v",
					len(pls), len(arb.Pairs),
					lpservice.getPoolSymbol(arb.Pairs[0]), param.bestAmtIn,
					profitUSD, time.Since(start))
				if profitUSD.Cmp(maxProfitUSD) > 0 {
					maxProfitUSD = profitUSD
					symbol = lpservice.getPoolSymbol(arb.Pairs[0])
				}
			}
		}
	}

	if maxProfitUSD.Cmp(decimal.NewFromFloat(0.02)) > 0 {
		logger.Info().Msgf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! %v max profit: %v block=%d used=%v",
			symbol, maxProfitUSD, blocknumber, time.Since(tm0))
	} else {
		logger.Info().Msgf("------- max profit: %v symbol=%v block=%d used=%v",
			maxProfitUSD, symbol, blocknumber, time.Since(tm0))
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (LPService) getArbGas() decimal.Decimal {
	// gasLimit = 1500000
	// gasPrice =

	return decimal.NewFromFloat(0.000)
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (lpservice *LPService) toETH(token string, amt *big.Int) decimal.Decimal {
	return decimal.NewFromFloat(0.002)
}

func (lpservice *LPService) toUSD(token string, amt *big.Int) decimal.Decimal {
	if isNativeOrWrapperNativeToken(token) {
		return decimal.NewFromBigInt(amt, 0).Mul(lpservice.ethPrice).Div(d_10_exp_18)
	}

	if lpservice.isStableToken(token) {
		tok := lpservice.tokens[token]
		return decimal.NewFromBigInt(amt, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tok.Decimals)))) // nolint
	}

	// middle token
	price := lpservice.midTokenPrice[token]
	if price == nil {
		logger.Warn().Msgf("not found token price: %v", token)
		return decimal.Zero
	}

	return decimal.NewFromBigInt(amt, 0).Mul(price.Price).Mul(lpservice.ethPrice).Div(d_10_exp_18)
}

// LoadTokenPools load token, pools from redis.
func (lpservice *LPService) LoadTokenPools() error {
	tokens, err := pool.LoadTokens(lpservice.rc)
	if err != nil {
		return err
	}

	pools, err := pool.LoadAllPools(lpservice.rc)
	if err != nil {
		return err
	}

	lpservice.tokens = tokens
	lpservice.pools = pools

	logger.Info().Msgf("load tokens: %d, load pools: %d", len(tokens), len(pools))

	return nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (lpservice *LPService) recalculateMidTokenPools() {
	logger.Info().Msg("recalulate middle tokens and pools ....")

	lpservice.setMidTokens()
	lpservice.refreshPools(true)
	// must be last, after TVL and pool price
	lpservice.setMidPools()
	lpservice.setMinETHPools()
}

func (lpservice *LPService) setMidTokens() {
	tokenQuoted := lpservice.getPoolTokensMap()
	_midTokensMap := map[string]*common.Token{}

	for addr, token := range tokenQuoted {
		if token.PoolCount > minMidTokenPools {
			if lpservice.setMidTokenPricePools(token) > tokenMinTVLPools {
				_midTokensMap[addr] = token
			}
		}
	}

	lpservice.midTokens = _midTokensMap
	logger.Info().Msgf("middle tokens: %d", len(_midTokensMap))
}

func (lpservice *LPService) setMidPools() {
	_midPoolsMap := map[string]*pool.Pool{}

	for addr, item := range lpservice.pools {
		_, ok0 := lpservice.midTokens[item.Token0]
		_, ok1 := lpservice.midTokens[item.Token1]

		if ok0 && ok1 && item.GetTVL() != nil && item.GetTVL().Cmp(minMidPoolETH) >= 0 {
			// logger.Info().Msgf("add pool %v %v %v to midPools, TVL: %v ETH",
			// 	item.Vendor, lpservice.getPoolName(item), item.Address,
			// 	readableETHAmount(item.GetTVL()))
			_midPoolsMap[addr] = item
		}
	}

	lpservice.midPools = _midPoolsMap
	logger.Info().Msgf("middle pools: %d", len(_midPoolsMap))
}

func (lpservice *LPService) getPoolSymbol(p *pool.Pool) string {
	token0 := lpservice.tokens[p.Token0]
	token1 := lpservice.tokens[p.Token1]

	name := ""
	if token0 != nil {
		name += token0.Symbol
	}

	name += "/"
	if token1 != nil {
		name += token1.Symbol
	}

	return name
}

func (lpservice *LPService) getPoolName(p *pool.Pool) string {
	if p.GetName() != "" {
		return p.GetName()
	}

	token0 := lpservice.tokens[p.Token0]
	token1 := lpservice.tokens[p.Token1]

	name := ""
	if token0 != nil {
		name += token0.Symbol
	}

	name += "/"
	if token1 != nil {
		name += token1.Symbol
	}

	name += "-" + fmt.Sprintf("%.3f%%", float64(p.Fee)/10000.0)

	return name
}

func (lpservice *LPService) refreshPools(updateTVL bool) {
	for _, p := range lpservice.pools {
		lpservice.updatePool(p)
	}

	if updateTVL {
		lpservice.updatePoolsTVL()
	}
}

func (lpservice *LPService) setMinETHPools() {
	m := map[string]*pool.Pool{}

	for addr, p := range lpservice.pools {
		if p.GetTVL() != nil && p.GetTVL().Cmp(minPoolETH) >= 0 {
			if p.Typ.IsCAMMVariety() {
				if p.Liquidity == nil || p.Liquidity.Cmp(bigZero) == 0 {
					continue
				}
			}

			m[addr] = p
		}
	}

	lpservice.poolsGreatThanMinETH = m
	logger.Info().Msgf("pools TVL great than minPoolETH: %d", len(m))
}

// token's MaxTVLPools is either ETH quoted or stable quoted.
func (lpservice *LPService) setMidTokenPricePools(token *common.Token) int {
	tokenAddr := token.Address
	if isNativeOrWrapperNativeToken(tokenAddr) || lpservice.isStableToken(tokenAddr) {
		return 100
	}

	pools := []*pool.Pool{}

	for _, item := range lpservice.poolsGreatThanMinETH {
		token0 := item.Token0
		token1 := item.Token1

		if (tokenAddr == token0 && (isNativeOrWrapperNativeToken(token1) || lpservice.isStableToken(token1))) ||
			(tokenAddr == token1 && (isNativeOrWrapperNativeToken(token0) || lpservice.isStableToken(token0))) {
			if minMidPoolETH.Cmp(item.GetTVL()) <= 0 {
				pools = append(pools, item)
			}
		}
	}

	if len(pools) == 0 {
		return 0
	}

	sort.Sort(sort.Reverse(PoolTVLSlice(pools)))

	poolAddrs := []string{}

	for i := 0; i < len(pools) && i <= tokenMinTVLPools; i++ {
		pl := pools[i]
		// logger.Info().Msgf("token: %v %v price pool: %v %v %v TVL: %v ETH",
		// 	token.Symbol, tokenAddr, pl.Vendor, pl.GetName(), pl.Address,
		// 	readableETHAmount(pl.GetTVL()))

		poolAddrs = append(poolAddrs, pl.Address)
	}

	token.MaxTVLPools = poolAddrs

	return len(poolAddrs)
}

func (lpservice *LPService) isStableToken(tokenAddr string) bool {
	return isStableToken(tokenAddr, lpservice.stableTokens)
}

func (lpservice *LPService) getPoolTokensMap() map[string]*common.Token {
	tokenQuoted := map[string]*common.Token{}

	for addr, pl := range lpservice.pools {
		if _, ok := lpservice.tokens[pl.Token0]; !ok {
			logger.Warn().Msgf("not found pool %v token0 %v in tokensMap", addr, pl.Token0)
			continue
		}

		if _, ok := lpservice.tokens[pl.Token1]; !ok {
			logger.Warn().Msgf("not found pool %v token1 %v in tokensMap", addr, pl.Token1)
			continue
		}

		token0Quoted, ok := tokenQuoted[pl.Token0]
		if !ok {
			token0Quoted = lpservice.tokens[pl.Token0]
			token0Quoted.PoolCount = 0
			tokenQuoted[pl.Token0] = token0Quoted
		}

		token0Quoted.PoolCount++

		token1Quoted, ok := tokenQuoted[pl.Token1]
		if !ok {
			token1Quoted = lpservice.tokens[pl.Token1]
			token1Quoted.PoolCount = 0
			tokenQuoted[pl.Token1] = token1Quoted
		}

		token1Quoted.PoolCount++
	}

	return tokenQuoted
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (lpservice *LPService) updatePool(pl *pool.Pool) {
	if pl.GetName() == "" {
		pl.UpdateName(lpservice.getPoolName(pl))
	}

	if pl.Stable && pl.Typ.IsAMMVariety() {
		if token0, ok := lpservice.tokens[pl.Token0]; ok {
			pl.SetDecimals0(new(big.Int).Exp(bigTen, big.NewInt(int64(token0.Decimals)), nil)) // nolint
		}

		if token1, ok := lpservice.tokens[pl.Token1]; ok {
			pl.SetDecimals1(new(big.Int).Exp(bigTen, big.NewInt(int64(token1.Decimals)), nil)) // nolint
		}
	}

	pl.UpdateBestPrice()

	if lpservice.midTokens[pl.Token0] != nil &&
		!isNativeOrWrapperNativeToken(pl.Token0) &&
		!lpservice.isStableToken(pl.Token0) {
		// update token0 eth price
		// update token1 eth price
		priceUpdated, ok := lpservice.midTokenPrice[pl.Token0]
		if !ok {
			logger.Warn().Msgf("not found pool token0 middle token price: pool=%v token=%v", pl.Address, pl.Token0)
		} else {
			if pl.LastBlockUpdated > priceUpdated.LastBlock {
				tok := lpservice.tokens[pl.Token0]
				lpservice.updateMidTokenETHPrice(tok)
			}
		}
	}

	if lpservice.midTokens[pl.Token1] != nil &&
		!isNativeOrWrapperNativeToken(pl.Token1) &&
		!lpservice.isStableToken(pl.Token1) {
		// update token1 eth price
		priceUpdated, ok := lpservice.midTokenPrice[pl.Token1]
		if !ok {
			logger.Warn().Msgf("not found pool token1 middle token price: pool=%v token=%v", pl.Address, pl.Token1)
		} else {
			if pl.LastBlockUpdated > priceUpdated.LastBlock {
				tok := lpservice.tokens[pl.Token1]
				lpservice.updateMidTokenETHPrice(tok)
			}
		}
	}

	pl.SetTVL(lpservice.calcPoolTVL(pl))
}

func (lpservice *LPService) updatePoolsTVL() {
	for _, pl := range lpservice.pools {
		pl.SetTVL(lpservice.calcPoolTVL(pl))
	}
}

func (lpservice *LPService) poolReservesGreatThan(pl *pool.Pool, minETH *big.Int) bool {
	r0ETH := lpservice.getReserveETH(pl.Token0, pl.Reserve0)
	r1ETH := lpservice.getReserveETH(pl.Token1, pl.Reserve1)

	return r0ETH.Cmp(minETH) >= 0 && r1ETH.Cmp(minETH) >= 0
}

func (lpservice *LPService) getReserveETH(token string, reserve *big.Int) *big.Int {
	if isNativeOrWrapperNativeToken(token) {
		return reserve
	}

	tok := lpservice.tokens[token]
	if tok == nil {
		// logger.Warn().Msgf("not found token: %v", token)
		return big.NewInt(0)
	}
	if lpservice.isStableToken(token) {
		return toETHAmountWithPrecision(reserve, tok, lpservice.ethPrice, true)
	}
	// middle token
	price, ok := lpservice.midTokenPrice[token]
	if !ok {
		logger.Warn().Msgf("not found middle token price: %v %v", tok.Symbol, tok.Address)
		return big.NewInt(0)
	}

	return toETHAmount(reserve, price.Price, false)
}

func (lpservice *LPService) calcPoolTVL(pl *pool.Pool) *big.Int {
	token0 := pl.Token0
	token1 := pl.Token1

	if pl.Typ.IsAMMVariety() && (pl.Reserve0.Cmp(bigZero) <= 0 || pl.Reserve1.Cmp(bigZero) <= 0) {
		return big.NewInt(0)
	}

	if pl.Typ.IsCAMMVariety() && (pl.Liquidity != nil && pl.Liquidity.Cmp(bigZero) <= 0) {
		return big.NewInt(0)
	}

	switch {
	case isNativeOrWrapperNativeToken(token0):
		return new(big.Int).Set(pl.Reserve0)

	case isNativeOrWrapperNativeToken(token1):
		return new(big.Int).Set(pl.Reserve1)

	case lpservice.isStableToken(token0):
		token0 := lpservice.tokens[token0]
		return toETHAmountWithPrecision(pl.Reserve0, token0, lpservice.ethPrice, true)

	case lpservice.isStableToken(token1):
		token1 := lpservice.tokens[token1]
		return toETHAmountWithPrecision(pl.Reserve1, token1, lpservice.ethPrice, true)

	default:
		var price decimal.Decimal

		if tok, ok := lpservice.midTokens[token0]; ok {
			priceUpdated, exist := lpservice.midTokenPrice[token0]
			if !exist {
				price = lpservice.updateMidTokenETHPrice(tok)
			} else {
				if priceUpdated.LastBlock < pl.LastBlockUpdated {
					price = lpservice.updateMidTokenETHPrice(tok)
				} else {
					price = priceUpdated.Price
				}
			}

			return toETHAmount(pl.Reserve0, price, false)
		}

		if tok, ok := lpservice.midTokens[token1]; ok {
			priceUpdated, exist := lpservice.midTokenPrice[token1]
			if !exist {
				price = lpservice.updateMidTokenETHPrice(tok)
			} else {
				if priceUpdated.LastBlock < pl.LastBlockUpdated {
					price = lpservice.updateMidTokenETHPrice(tok)
				} else {
					price = priceUpdated.Price
				}
			}

			return toETHAmount(pl.Reserve1, price, false)
		}
	}

	return big.NewInt(0)
}

func (lpservice *LPService) refreshETHPrice() {
	var (
		price    decimal.Decimal
		usdToken *common.Token
	)

	pl := lpservice.pools[lpservice.nativeStablePair]
	pl.UpdateBestPrice()

	if pl.Token0 == WETHAddress {
		price = pl.GetPrice0()
		usdToken = lpservice.tokens[pl.Token1]
	} else {
		price = pl.GetPrice1()
		usdToken = lpservice.tokens[pl.Token0]
	}

	if usdToken.Decimals == 18 {
		lpservice.ethPrice = price
	} else {
		prec := d10.Pow(decimal.NewFromInt(18 - int64(usdToken.Decimals))) // nolint
		lpservice.ethPrice = price.Mul(prec)
	}

	logger.Info().Msgf("refresh ETH price: %v", lpservice.ethPrice)
}

func (lpservice *LPService) updateMidTokenETHPrice(tok *common.Token) decimal.Decimal {
	total := decimal.NewFromInt(0)
	count := 0

	var (
		prevPrice *decimal.Decimal
		block     uint64
	)

	prevPool := ""

	for _, addr := range tok.MaxTVLPools {
		p := lpservice.pools[addr]

		var price decimal.Decimal
		if tok.Address == p.Token0 {
			price = p.GetPrice0()
		} else {
			price = p.GetPrice1()
		}

		if price.Cmp(decimal.Zero) == 0 {
			continue
		}

		if lpservice.isStableToken(p.Token0) {
			token0 := lpservice.tokens[p.Token0]
			prec := d10.Pow(decimal.NewFromInt(18 - int64(token0.Decimals))) // nolint
			price = price.Mul(prec).Div(lpservice.ethPrice)
		} else if lpservice.isStableToken(p.Token1) {
			token1 := lpservice.tokens[p.Token1]
			prec := d10.Pow(decimal.NewFromInt(18 - int64(token1.Decimals))) // nolint
			price = price.Mul(prec).Div(lpservice.ethPrice)
		}

		if price.Cmp(decimal.Zero) > 0 {
			total = total.Add(price)
			count++

			if prevPrice != nil {
				bias := prevPrice.Sub(price).Abs().Div(price)
				if bias.Cmp(p5) > 0 {
					logger.Warn().Msgf("token price bias too large: %s %s prev: %v %v current: %v %v",
						tok.Symbol, tok.Address, prevPrice.String(), prevPool, price, p.Address)
				}
			}

			prevPrice = &price
			prevPool = p.Address
		}

		if p.LastBlockUpdated > block {
			block = p.LastBlockUpdated
		}
	}

	price := total.Div(decimal.NewFromInt(int64(count)))
	lpservice.midTokenPrice[tok.Address] = &PriceUpdated{
		Price:     price,
		LastBlock: block,
	}

	logger.Info().Msgf("update middle token %v %v price to %v, block: %v", tok.Symbol, tok.Address, price, block)

	return price
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (lpservice *LPService) tokenQuotable(token0 string, token1 string) string {
	// return isNativeOrWrapperNativeToken(token) ||
	// 	isStableToken(token, lpservice.stableTokens) ||
	// 	lpservice.midTokens[token] != nil
	if isNativeOrWrapperNativeToken(token0) {
		return token0
	}

	if isNativeOrWrapperNativeToken(token1) {
		return token1
	}

	if lpservice.isStableToken(token0) {
		return token0
	}

	if lpservice.isStableToken(token1) {
		return token1
	}

	if lpservice.midTokens[token0] != nil {
		return token0
	}

	if lpservice.midPools[token1] != nil {
		return token1
	}

	logger.Warn().Msgf("both token %s %s is NOT quotable", token0, token1)

	return token0
}

/*
func (lpservice *LPService) calculatePriceByPool(pl *pool.Pool, token string, ethQuote bool) decimal.Decimal {
	isToken0 := token == pl.Token0
	token0 := lpservice.tokens[pl.Token0]
	token1 := lpservice.tokens[pl.Token1]

	if token0 == nil || token1 == nil {
		logger.Warn().
			Str("pool", pl.Address).
			Str("token0", pl.Token0).
			Str("token1", pl.Token1).
			Msg("not found pool token")

		return decimal.Zero
	}

	if !(lpservice.tokenQuotable(pl.Token0) ||
		lpservice.tokenQuotable(pl.Token1)) {
		logger.Warn().Msgf("no quotable token in pool: %s %s/%s", pl.Address, token0.Symbol, token1.Symbol)
		return decimal.Zero
	}

	if (!ethQuote) && WETHAddress != pl.Token0 && WETHAddress != pl.Token1 {
		logger.Fatal().Str("pool", pl.Address).Msg("not ETH quote, but neither token0 nor token1 is ETH")
	}

	//
	if pl.Typ.IsAMMVariety() {
		return lpservice.calculatePoolV2Price(pl, isToken0, token0, token1, ethQuote)
	} else if pl.Typ.IsCAMMVariety() {
		return lpservice.calculatePoolV3Price(pl, isToken0, token0, token1, ethQuote)
	} else {
		logger.Fatal().Str("poolType", pl.Typ.String()).
			Str("pool", pl.Address).Msg("invalid pool type")
	}

	return decimal.Zero
}

func (lpservice *LPService) calcPrice(
	pl *pool.Pool,
	isToken0 bool,
	token0 *common.Token,
	token1 *common.Token,
	ethQuote bool,
	price decimal.Decimal,
	precision decimal.Decimal,
) decimal.Decimal {
	if ethQuote {
		if isNativeOrWrapperNativeTokens(token0.Address, token1.Address) {
			if isToken0 {
				return price.Mul(precision)
			} else {
				return decimal.NewFromInt(1).Div(price).Div(precision)
			}
		} else if isStableTokens(token0.Address, token1.Address, lpservice.tokens) {
			if isToken0 {
				return price.Mul(precision).Div(lpservice.ethPrice)
			} else {
				return lpservice.ethPrice.Div(price).Div(precision)
			}
		}
	}

	// one token MUST be ETH
	if isToken0 {
		return price.Mul(precision)
	} else {
		return decimal.NewFromInt(1).Div(price).Div(precision)
	}
}

// todo aero drome v2 stable pool.
func (lpservice *LPService) calculatePoolV2Price(
	pl *pool.Pool,
	isToken0 bool,
	token0 *common.Token,
	token1 *common.Token,
	ethQuote bool,
) decimal.Decimal {
	price := decimal.NewFromBigInt(pl.Reserve1, 0).Div(decimal.NewFromBigInt(pl.Reserve0, 0))
	precision := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(token0.Decimals - token1.Decimals)))

	return lpservice.calcPrice(pl, isToken0, token0, token1, ethQuote, price, precision)
}

func (lpservice *LPService) calculatePoolV3Price(
	pl *pool.Pool,
	isToken0 bool,
	token0 *common.Token,
	token1 *common.Token,
	ethQuote bool,
) decimal.Decimal {
	sqrtPrice := decimal.NewFromBigInt(pl.SqrtPriceX96, 0).Div(q96)
	price := sqrtPrice.Mul(sqrtPrice)
	precision := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(token0.Decimals - token1.Decimals)))

	return lpservice.calcPrice(pl, isToken0, token0, token1, ethQuote, price, precision)
}
*/
