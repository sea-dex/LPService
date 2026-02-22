package swapor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/arb"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

type EventHandler struct {
	startBlock   uint64
	lastBlock    uint64
	totalEvents  uint64
	parsedEvents uint64
	blockNumber  uint64
	logIndex     uint

	// reseting bool
	// resetCh  chan any

	pairsQuery     string
	factory        map[string]*common.SwapFactory // known factory
	v3poolQuery    map[common.PoolType]string
	pools          map[string]*pool.Pool
	tokens         map[string]*common.Token
	unknownFactory map[string]bool                // unknown factory
	unknownPools   map[string]*pool.PoolBasicInfo // pools which is already unknown, and will NOT fetch it onchain
	pp             *pool.ProviderPool
	store          redis.Cmdable
	// producer       *KafkaProducer // remove kafka dependency
	// producer    *kafkago.Writer
	// feeProducer *kafkago.Writer
	subscriber *events.EventSubscriber
	cfg        *config.Config
	lpservice  *arb.LPService
}

var (
	ErrUnknownPool = errors.New("pool is already unknown")
	topicMap       = map[string]bool{
		pool.TopicSYNC:                true,
		pool.TopicSYNCAero:            true,
		pool.TopicInfusionPairCreated: true,
		// TopicInitialize:  true,
		pool.TopicAeroV2SetFee: true,
		pool.TopicAeroV3SetFee: true,
		pool.TopicMint:         true,
		pool.TopicBurn:         true,
		pool.TopicAeroSwapV2:   true,
		pool.TopicSwapV2:       true,
		pool.TopicSwap:         true,
		pool.TopicCollect:      true,
		pool.TopicPancakeSwap:  true,
	}
)

// CreateEventHandler create event handler.
func CreateEventHandler(cfg *config.Config) *EventHandler {
	var err error

	eh := &EventHandler{
		factory:        map[string]*common.SwapFactory{}, // known factory
		pools:          map[string]*pool.Pool{},
		unknownFactory: map[string]bool{},
		unknownPools:   map[string]*pool.PoolBasicInfo{},
		v3poolQuery:    map[common.PoolType]string{},
		pairsQuery:     cfg.Chain.PairsQuery,
		tokens:         map[string]*common.Token{},
		cfg:            cfg,
	}

	if eh.pairsQuery == "" {
		logger.Fatal().Msg("pairsQuery address is empty")
	}

	for _, item := range cfg.Chain.Factory {
		// every factory in config is Known
		eh.AddFactory(item.Address, item.PositionManager, item.Name, item.Typ, item.Fee, item.StableFee, true)
	}

	for k, v := range cfg.Chain.PoolQuery {
		eh.v3poolQuery[common.PoolType(v)] = k
		logger.Info().Uint("poolType", v).Str("address", k).Msg("add v3 pool query contract")
	}

	// validate every v3 pool type has query contract
	for _, f := range eh.factory {
		if f.Typ.IsCAMMVariety() {
			if _, ok := eh.v3poolQuery[f.Typ]; !ok {
				logger.Fatal().Uint("poolType", uint(f.Typ)).Msg("not found v3 pool query")
			}
		}
	}

	eh.pp = pool.CreateProviderPool(cfg.Chain.Providers)

	eh.store, err = common.CreatePoolRedisStore(cfg.IsProd(), cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Fatal().Err(err).Msg("create pool redis store failed")
	}

	// remove kafka dependency
	// eh.producer = events.CreateKafkaProducer(cfg.Kafka.Topic, cfg.Kafka.Brokers)
	// eh.feeProducer = events.CreateKafkaProducer(config.PoolFeeTopic, cfg.Kafka.Brokers)
	// eh.producer = NewKafkaProducer(strings.Split(cfg.Kafka.Brokers, ","))
	eh.lpservice = arb.CreateLPService(cfg)

	return eh
}

// SetSubscriber set subscriber.
func (eh *EventHandler) SetSubscriber(sub *events.EventSubscriber) {
	eh.subscriber = sub
}

func (eh *EventHandler) restartSubscribeEvents(ctx context.Context,
	block uint64,
	logIdx uint,
	logCh chan types.Log,
	errCh chan error,
	wg *sync.WaitGroup,
) {
	swg := &sync.WaitGroup{}
	swg.Add(1)

	cctx, cancel := context.WithCancel(ctx)
	subErrCh := make(chan error, 1)

	go eh.subscriber.SubscribeEventsFrom(cctx, nil, block, logIdx, logCh, subErrCh, swg)

	var err error
	select {
	case err = <-subErrCh:
		logger.Info().Msg("ws subscribe error occurs, should be restart")

	case <-ctx.Done():
		logger.Info().Msg("ws subscribe canceled")
	}

	cancel()
	swg.Wait()
	logger.Info().Msg("subscribe logs exited")

	if err != nil {
		errCh <- err
	}

	wg.Done()
}

func (eh *EventHandler) parseEvents(logCh chan types.Log) {
	if len(logCh) == 0 {
		return
	}

	addrEvents, feeEvents, blockNumber, logIndex, evtsCount := eh.drainEvents(logCh)
	results := []common.ParseResult{}

	poolAddrs := []string{}

	for addr := range addrEvents {
		if _, ok := eh.factory[addr]; ok {
			// factory events
			continue
		}

		if eh.pools[addr] == nil && eh.unknownPools[addr] == nil {
			poolAddrs = append(poolAddrs, addr)
		}
	}

	eh.fetchNewPoolAndTokens(poolAddrs)

	v2Updated := 0
	v3Updated := 0
	v3pools := []*pool.Pool{}

	for addr, evts := range addrEvents {
		if fac, ok := eh.factory[addr]; ok {
			for _, evt := range evts {
				result := eh.doFactoryEvents(fac, &evt)
				if result.Status == common.ParsePoolCreated {
					results = append(results, result)
				}
			}

			continue
		}

		result := eh.ParseAddrEvents(addr, evts)
		if result.Status == common.ParsePoolUpdated {
			results = append(results, result)

			if result.Data.(*pool.Pool).Typ.IsAMMVariety() {
				v2Updated++
			} else {
				v3Updated++

				v3pools = append(v3pools, result.Data.(*pool.Pool))
			}
		}
	}

	// swaps := eh.ParseAddrSwapEvents(addrEvents)

	if len(v3pools) > 0 {
		eh.pp.ReloadV3PoolReserves(v3pools)
	}

	if len(feeEvents) > 0 {
		eh.handleFeeEvents(feeEvents)
	}

	logger.Info().Msgf("---------- parse events to %d, %d, evts: %d pools updated: total=%d v2=%d v3=%d",
		blockNumber, logIndex, evtsCount, len(results), v2Updated, v3Updated)

	if err := eh.SyncPoolLiquids(results, blockNumber, logIndex); err != nil {
		logger.Fatal().Err(err).Msg("sync pool liquidity failed")
	}
}

// EventRoutine parse events.
func (eh *EventHandler) StartEventsRoutine(ctx context.Context, sub *events.EventSubscriber, wg *sync.WaitGroup) {
	eh.SetSubscriber(sub)
	mode := eh.getMode()
	// init event consumer
	err := eh.Init(mode, false)
	if err != nil {
		logger.Fatal().Msg("init LP tracker failed: " + err.Error())
	}

	logger.Info().Msgf("start LP tracker: mode=%s blockNumber=%d logIndex=%d", mode, eh.blockNumber, eh.logIndex)

	logCh := make(chan types.Log, 5000)
	errCh := make(chan error)

	cctx, cancel := context.WithCancel(ctx)

	swg := &sync.WaitGroup{}
	swg.Add(1)

	go eh.restartSubscribeEvents(cctx, eh.blockNumber, eh.logIndex, logCh, errCh, swg)

	tmr := time.NewTicker(time.Millisecond * 100)
	hourTmr := time.NewTicker(time.Second * 3600)

	for {
		select {
		case <-ctx.Done():
			cancel()

			for len(logCh) > 0 {
				eh.parseEvents(logCh)
			}

			wg.Done()
			swg.Wait()
			logger.Info().Msg("events routine exited")
			close(logCh)

			return

		case <-errCh:
			logger.Info().Err(err).
				Uint64("blockNumber", eh.blockNumber).
				Uint("logIndex", eh.logIndex).
				Msg("main thread channel recv ws error")
			time.Sleep(time.Second)

			for len(logCh) > 0 {
				logger.Info().
					Uint64("blockNumber", eh.blockNumber).
					Uint("logIndex", eh.logIndex).
					Msgf("ws subscribe error occurs, before restart, drain up channel events: %d", len(logCh))
				eh.parseEvents(logCh)
			}

			logger.Info().
				Uint64("blockNumber", eh.blockNumber).
				Uint("logIndex", eh.logIndex).
				Msg("restart events subscribe")

			swg.Add(1)

			go eh.restartSubscribeEvents(ctx, eh.blockNumber, eh.logIndex, logCh, errCh, wg)

		case <-tmr.C:
			eh.parseEvents(logCh)

		case <-hourTmr.C:
			eh.checkPoolsFee()
		}
	}
}

// SyncPoolLiquids send pool to redis and kafka.
func (eh *EventHandler) SyncPoolLiquids(updates []common.ParseResult, blockNumber uint64, logIndex uint) error {
	pipe := eh.store.Pipeline()
	ctx := context.Background()
	// msgs := []kafkago.Message{}
	msgs := [][]byte{}
	changed := map[string]*pool.Pool{}
	// totalBytes := 0

	if len(updates) > 0 {
		for _, update := range updates {
			buf, err := json.Marshal(update.Data)
			if err != nil {
				logger.Error().Err(err).Msg("marshal pool failed")
				return err
			}

			data := update.Data.(*pool.Pool)
			changed[data.Address] = data

			if err := pipe.HSet(ctx, common.PoolLiquidKey, data.Address, buf).Err(); err != nil {
				logger.Error().Err(err).Str("pool", data.Address).Msg("HSET pool update failed")
				return err
			}
			// msgs = append(msgs, kafkago.Message{
			// 	Value: buf,
			// })
			// totalBytes += len(buf)
			msgs = append(msgs, buf)
		}
	}

	if err := pipe.Set(ctx, eventsParsedAtKey, fmt.Sprintf("%d,%d", blockNumber, logIndex), 0).Err(); err != nil {
		logger.Error().Err(err).Msg("SET eventParsedAt failed")
		return err
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("sync pools updates to redis failed")
		return err
	}

	// Update cursor only if Redis operations were successful
	eh.blockNumber = blockNumber
	eh.logIndex = logIndex

	if len(msgs) > 0 {
		// logger.Info().Msgf("message bytes: %d", totalBytes)
		// only send 100 max per time
		// if err := eh.producer.WriteMessages(ctx, msgs...); err != nil {
		// 	logger.Error().Err(err).Msg("producer pools updates to kafka failed")
		// 	return err
		// }
		// remove kafka dependency
		// eh.producer.ProduceMessage(eh.cfg.Kafka.Topic, msgs)
		eh.lpservice.OnPoolUpdated(changed, decimal.NewFromFloat(1.00005), blockNumber)
	} else {
		logger.Warn().Msgf("no messages to produce: updates=%d block=%v logIndex=%v", len(updates), blockNumber, logIndex)
	}

	return nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
// func (eh *EventHandler) produceSwaps(swaps []*pool.SwapEvent) {
// 	msgs := [][]byte{}

// 	for _, item := range swaps {
// 		buf, err := json.Marshal(item)
// 		if err != nil {
// 			logger.Fatal().Err(err).Msgf("marshal swap event failed")
// 		}

// 		msgs = append(msgs, buf)
// 	}

// 	eh.producer.ProduceMessage("swap", msgs)
// }

func (eh *EventHandler) drainEvents(logCh chan types.Log) (map[string][]types.Log, []types.Log, uint64, uint, int) {
	var events []types.Log

	total := 0
	maxEvts := 50000
	blockNumber := eh.blockNumber
	logIndex := eh.logIndex
	feeEvts := []types.Log{}

	for len(logCh) > 0 {
		e := <-logCh
		total++

		if e.BlockNumber == blockNumber {
			if logIndex+1 != e.Index {
				logger.Fatal().Msgf("logIndex NOT continuous: blocknumber=%d current logIndex=%d prev logIndex=%d", blockNumber, e.Index, logIndex)
			}

			logIndex++
		} else if e.BlockNumber < blockNumber {
			logger.Fatal().Msgf("event duplicate: current event blocknumber: %d prev event blocknumber: %d", e.BlockNumber, blockNumber)
		} else {
			if blockNumber+1 != e.BlockNumber {
				logger.Warn().Msgf("blocknumber NOT continuous: event blocknumber: %d prev event blocknumber: %d", e.BlockNumber, blockNumber)
			}

			blockNumber = e.BlockNumber
			logIndex = 0
		}
		// blockNumber = e.BlockNumber
		// logIndex = e.Index

		if e.Removed {
			logger.Warn().
				Uint64("block", e.BlockNumber).
				Uint("logIdx", e.Index).
				Str("txhash", e.TxHash.String()).
				Msg("event was removed")

			continue
		}

		if eh.startBlock == 0 {
			eh.startBlock = e.BlockNumber
		}

		// there is event which has no topic
		// https://basescan.org/tx/0xd70778d2578a126df246be3d4b4135158896949fa767d018d2ce939016d4b51c#eventlog
		if len(e.Topics) == 0 {
			continue
		}

		if e.BlockNumber > eh.lastBlock {
			eh.lastBlock = e.BlockNumber
			logger.Info().Uint64("blocknumber", e.BlockNumber).Msg("new block events arrived")
		}
		//  if (e.BlockNumber < eh.lastBlock) || (e.BlockNumber == eh.lastBlock && e.Index ) {
		// 	logger.Error().
		// 		Uint64("eventBlockNumber", e.BlockNumber).
		// 		Uint64("lastBlockParsed", eh.lastBlock).
		// 		Msg("new event block less than latest block parsed")
		// 	// mostly occur in startup
		// 	// just skip events
		// 	continue
		// }

		if e.BlockNumber < eh.blockNumber || (e.BlockNumber == eh.blockNumber && e.Index <= eh.logIndex) {
			logger.Fatal().
				Uint64("eventBlockNumber", e.BlockNumber).
				Uint("eventLogIndex", e.Index).
				Msg("event has been parsed")
		}

		topic := strings.ToLower(e.Topics[0].String())
		if topicMap[topic] {
			if topic == pool.TopicAeroV2SetFee || topic == pool.TopicAeroV3SetFee {
				feeEvts = append(feeEvts, e)
			} else {
				events = append(events, e)
			}
		}

		if len(events) > maxEvts {
			break
		}
	}

	eh.totalEvents += uint64(total) // nolint
	eh.parsedEvents += uint64(len(events))

	addrEvents := map[string][]types.Log{}
	// ethusdc := strings.ToLower("0xd0b53D9277642d899DF5C87A3966A349A798F224")

	for _, evt := range events {
		addr := strings.ToLower(evt.Address.String())
		addrEvents[addr] = append(addrEvents[addr], evt)
	}

	return addrEvents, feeEvts, blockNumber, logIndex, len(events)
}

// ParseAddrEvents parse address events. addr is always lower case.
func (eh *EventHandler) ParseAddrEvents(addr string, events []types.Log) (result common.ParseResult) {
	pl, err := eh.getPool(addr)
	if err != nil {
		return
	}

	if pl.Typ.IsAMMVariety() {
		eh.doV2PoolEvents(pl, events)
	} else if pl.Typ.IsCAMMVariety() {
		eh.doV3PoolEvents(pl, events)
	} else {
		logger.Error().Str("pool", pl.Address).Uint("poolType", uint(pl.Typ)).Msg("unknown pool type")
		logger.Fatal().Msgf("unknown pool type: %s", pl.Typ)
	}

	pl.LastBlockUpdated = events[len(events)-1].BlockNumber
	result.Status = common.ParsePoolUpdated
	result.Data = pl

	return
}

func (eh *EventHandler) ParseAddrSwapEvents(events map[string][]types.Log) []*pool.SwapEvent {
	swaps := []*pool.SwapEvent{}

	for addr, evts := range events {
		pl, err := eh.getPool(addr)
		if err != nil {
			continue
		}

		tok0 := eh.tokens[pl.Token0]
		if tok0 == nil {
			logger.Warn().Str("pool", pl.Address).Str("token", pl.Token0).Msg("not found token0")
			continue
		}

		tok1 := eh.tokens[pl.Token1]
		if tok1 == nil {
			logger.Warn().Str("pool", pl.Address).Str("token", pl.Token1).Msg("not found token1")
			continue
		}

		decimal0 := tok0.Decimals
		decimal1 := tok1.Decimals

		// addr evts, only care swap events
		for _, evt := range evts {
			topic := evt.Topics[0].Hex()
			if !isSwapEvent(topic) {
				continue
			}

			zeroForOne, amt0, amt1 := pl.ParseSwapEvent(topic, &evt)
			swaps = append(swaps, &pool.SwapEvent{
				ZeroForOne:       zeroForOne,
				Amount0:          amt0,
				Amount1:          amt1,
				Decimals0:        decimal0,
				Decimals1:        decimal1,
				Txhash:           evt.TxHash.String(),
				Address:          pl.Address,
				Token0:           pl.Token0,
				Token1:           pl.Token1,
				Factory:          pl.Factory,
				Vendor:           pl.Vendor,
				Typ:              pl.Typ,
				Fee:              pl.Fee,
				Stable:           pl.Stable,
				TickSpacing:      pl.TickSpacing,
				LastBlockUpdated: pl.LastBlockUpdated,
				InitBlock:        pl.InitBlock,
			})
		}
	}

	return swaps
}

func isSwapEvent(topic string) bool {
	return topic == pool.TopicSwapV2 || topic == pool.TopicSwap || topic == pool.TopicAeroSwapV2
}
