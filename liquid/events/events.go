package events

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	lastRPCTS        int64
	errCanceled      = fmt.Errorf("canceled")
	errReachMaxRetry = fmt.Errorf("reach max retry times")
)

// EventSubscriber event subscriber.
type EventSubscriber struct {
	clients       []*ethclient.Client
	wssClient     *ethclient.Client
	RPC           string // http endpoint
	WSS           string // wss endpoint
	RateInterval  uint64 // milli second
	BlockInterval uint64 // milli second
	MaxRetry      uint32
	MaxStep       uint32
	SubEvents     bool
}

// NewEventSubscirber create EventSubscriber.
func NewEventSubscirber(rpcs, wss string,
	sub bool,
	rateInterval, blockInterval uint64,
	maxRetry, maxStep uint32,
) (*EventSubscriber, error) {
	var (
		err       error
		wssClient *ethclient.Client
		clients   = []*ethclient.Client{}
	)

	rpcList := strings.Split(rpcs, ",")
	for _, rpc := range rpcList {
		client, err := ethclient.Dial(strings.TrimSpace(rpc))
		if err != nil {
			logger.Error().Err(err).Str("rpc", rpc).Msg("create http eth client failed")
			return nil, err
		}

		clients = append(clients, client)
	}

	if sub || wss != "" {
		wssClient, err = ethclient.Dial(strings.TrimSpace(wss))
		if err != nil {
			logger.Error().Err(err).Str("wss", wss).Msg("create wss eth client failed")
			return nil, err
		}
	}

	if rateInterval == 0 {
		// 500 millisecond
		rateInterval = 500
	}

	if maxRetry == 0 {
		maxRetry = 5
	}

	if maxStep == 0 {
		maxStep = 100
	}

	if blockInterval == 0 {
		blockInterval = 2000 // default for base, op
	}

	es := &EventSubscriber{
		clients:       clients,
		wssClient:     wssClient,
		RPC:           rpcs,
		RateInterval:  rateInterval,
		BlockInterval: blockInterval,
		MaxRetry:      maxRetry,
		MaxStep:       maxStep,
		SubEvents:     sub,
	}

	return es, nil
}

// MustNewEventSubscriber NewEventSubscirber, panic if failed.
func MustNewEventSubscriber(rpc, wss string,
	sub bool,
	rateInterval, blockInterval uint64,
	maxRetry, maxStep uint32,
) *EventSubscriber {
	es, err := NewEventSubscirber(rpc, wss, sub, rateInterval, blockInterval, maxRetry, maxStep)
	if err != nil {
		logger.Fatal().Msg("create EventSubscriber failed: " + err.Error())
	}

	return es
}

func (es *EventSubscriber) SubscribeBlockHTTP(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			block, err := es.clients[0].BlockNumber(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("get blocknumber failed")
			} else {
				logger.Info().Uint64("block", block).Msg("new block")
			}
		}
	}
}

// SubscribeBlock subscribe new block. when error occurs, auto re-sub. blocked.
func (es *EventSubscriber) SubscribeBlock(ctx context.Context) {
	if !es.SubEvents {
		logger.Warn().Msg("SubscribeBlock: not subscribe mode")
		// es.SubscribeBlockHTTP(ctx)
		return
	}

	blockCh := make(chan *types.Header, 5)

	go func() {
		for {
			b := <-blockCh
			logger.Info().Uint64("block", b.Number.Uint64()).Msg("new block")
		}
	}()

	for {
		sub, err := es.wssClient.SubscribeNewHead(ctx, blockCh)
		if err != nil {
			log.Err(err).Msg("subscribe new head failed")
			continue
		}

		logger.Info().Msg("Subscribe New Head success")

	L:
		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("SubscribeNewHead was canceled, exit")
				sub.Unsubscribe()

				return

			case err = <-sub.Err():
				logger.Warn().Err(err).Msg("SubscribeNewHead error, Re-subscribe")
				sub.Unsubscribe()
				break L
			}
		}
	}
}

// SubscribeEventsFrom subscribe, events. This is a blocked call.
func (es *EventSubscriber) SubscribeEventsFromByHTTP(
	ctx context.Context,
	filterAddress []string,
	from uint64,
	logIndex uint,
	logCh chan types.Log,
	errCh chan error,
	wg *sync.WaitGroup,
) {
	filterAddr := []common.Address{}
	for _, addr := range filterAddress {
		filterAddr = append(filterAddr, common.HexToAddress(addr))
	}

	// first, drain logs between [from, end]
	start := from
	end := es.GetLatestBlockNumber()
	firstLoop := true

	logger.Info().Msgf("SubscribeEventsFromByHTTP: from=%d fromIndex=%d latestBlock=%d", from, logIndex, end)

	for {
		// events should be filtered by logIndex
		if firstLoop {
			startEvts, err := es.GetLogsFromToReturn(ctx, filterAddr, start, start, es.MaxStep)
			if err != nil {
				wg.Done()

				if err != errCanceled {
					errCh <- err
				}

				return
			}

			for _, evt := range startEvts {
				if (evt.BlockNumber == from && evt.Index > logIndex) || (evt.BlockNumber > from) {
					logCh <- evt
				}
			}

			start++
			firstLoop = false
		}

		_, err := es.GetLogsFromTo2(ctx, filterAddr, start, end, es.MaxStep, logCh)
		if err != nil {
			if err == errCanceled {
				err = nil
			}

			wg.Done()
			errCh <- err

			return
		}

		start = end + 1
		end = es.GetLatestBlockNumber()

		if start >= end || end-start < 20 {
			logger.Info().Msgf("SubscribeEventsFromByHTTP: catch up latest block, end=%v latest=%v, switch to subscribe", start, end)
			break
		} else {
			logger.Info().Msgf("SubscribeEventsFromByHTTP: GetLogsFromTo2 complete: start=%d latest=%d, continue", start, end)
		}
	}

	var (
		tmr               = time.NewTimer(time.Duration(es.BlockInterval/2) * time.Millisecond) // nolint
		wssCtx, wssCancel = context.WithCancel(ctx)
		blockCh           = make(chan *types.Header, 5)
	)

	if es.wssClient != nil {
		go func() {
			for {
				sub, err := es.wssClient.SubscribeNewHead(wssCtx, blockCh)
				if err != nil {
					if err == context.Canceled {
						logger.Warn().Msg("SubscribeEventsFromByHTTP: subscribe was canceled")
						return
					}

					logger.Warn().Err(err).Msg("SubscribeEventsFromByHTTP: subscribe new head failed")

					continue
				}

				logger.Info().Msg("SubscribeEventsFromByHTTP: subscribe New Head success")

			L:
				for {
					select {
					case <-wssCtx.Done():
						logger.Info().Msg("SubscribeEventsFromByHTTP: SubscribeNewHead was canceled, exit")
						sub.Unsubscribe()

						return

					case err = <-sub.Err():
						logger.Warn().Err(err).Msg("SubscribeEventsFromByHTTP: SubscribeNewHead error, Re-subscribe")
						sub.Unsubscribe()
						break L
					}
				}
			}
		}()
	}

	fetchLogs := func(endBlock uint64) {
		if endBlock >= start {
			count, err := es.GetLogsFromTo2(ctx, filterAddr, start, endBlock, es.MaxStep, logCh)
			if err != nil {
				if err == errCanceled {
					err = nil
				}

				wg.Done()
				wssCancel()
				errCh <- err

				return
			}

			if count == 0 {
				logger.Warn().Msgf("not found events: [%d, %d], not inc start", start, endBlock)
			} else {
				start = endBlock + 1
			}
		}
	}

	for {
		select {
		case <-tmr.C:
			end = es.GetLatestBlockNumber()
			logger.Info().Msgf("timer end, latest block: %d", end)

			fetchLogs(end)
			tmr.Reset(time.Duration(es.BlockInterval/2) * time.Millisecond) // nolint

		case b := <-blockCh:
			// cannot use new block fetch logs, in some case, different provider has NOT same latest block
			// end = b.Number.Uint64()
			// if end >= start {
			// 	err := es.GetLogsFromTo2(ctx, filterAddr, start, end, es.MaxStep, logCh)
			// 	if err != nil {
			// 		if err == errCanceled {
			// 			err = nil
			// 		}
			// 		wg.Done()
			// 		wssCancel()
			// 		errCh <- err
			// 		return
			// 	}
			// 	start = end + 1
			// }
			// tmr.Reset(time.Duration(es.BlockInterval) * time.Millisecond)
			logger.Info().Msgf("new block: %d %d", b.Number.Uint64(), b.Time)
			tmr.Reset(time.Duration(es.BlockInterval) * time.Millisecond) // nolint
			fetchLogs(b.Number.Uint64())

		case <-ctx.Done():
			logger.Info().Msg("SubscribeEventsFromByHTTP: canceled")
			wg.Done()
			wssCancel()

			return
		}
	}
}

// SubscribeEventsFrom subscribe, events. This is a blocked call.
func (es *EventSubscriber) SubscribeEventsFrom(
	ctx context.Context,
	filterAddress []string,
	from uint64,
	logIndex uint,
	logCh chan types.Log,
	errCh chan error,
	wg *sync.WaitGroup,
) {
	if !es.SubEvents {
		es.SubscribeEventsFromByHTTP(ctx, filterAddress, from, logIndex, logCh, errCh, wg)
		return
	}

	filterAddr := []common.Address{}
	for _, addr := range filterAddress {
		filterAddr = append(filterAddr, common.HexToAddress(addr))
	}

	// first, drain logs between [from, end]
	start := from
	end := es.GetLatestBlockNumber()

	logger.Info().Msgf("SubscribeEventsFrom: from=%d fromIndex=%d latestBlock=%d", from, logIndex, end)

	for {
		// events should be filtered by logIndex
		startEvts, err := es.GetLogsFromToReturn(ctx, filterAddr, start, start, es.MaxStep)
		if err != nil {
			wg.Done()

			if err != errCanceled {
				errCh <- err
			}

			return
		}

		for _, evt := range startEvts {
			if evt.Index > logIndex {
				logCh <- evt
			}
		}

		start++

		_, err = es.GetLogsFromTo2(ctx, filterAddr, start, end, es.MaxStep, logCh)
		if err != nil {
			if err == errCanceled {
				err = nil
			}

			wg.Done()
			errCh <- err

			return
		}

		start = end + 1
		end = es.GetLatestBlockNumber()

		if start >= end || end-start < 20 {
			logger.Info().Msgf("SubscribeEventsFrom: catch up latest block, end=%v latest=%v, switch to subscribe", start, end)
			break
		}
	}

	query := ethereum.FilterQuery{Addresses: filterAddr}
	ch := make(chan types.Log, 1000)
	cctx, cancel := context.WithCancel(ctx)

	wg.Add(1)

	first := true

	go func() {
		for {
			select {
			case e := <-ch:
				if first {
					logger.Info().Msgf("first event subscribed: block=%d, logIndex=%d, last fetch block=%d", e.BlockNumber, e.Index, start-1)
					// fetch events between gap: [start, e.blockNumber]
					// very important !!!
					if start <= e.BlockNumber {
						events, err := es.GetLogsFromToReturn(ctx, filterAddr, start, e.BlockNumber, es.MaxStep)
						if err != nil {
							logger.Error().Err(err).Msgf("SubscribeEventsFrom: got logs [%d, %d] failed", start, e.BlockNumber)
							close(ch)
							errCh <- err

							wg.Done()

							return
						}

						logger.Info().Msgf("fetch logs between [%d, %d]: %d", start, e.BlockNumber, len(events))

						for _, ele := range events {
							if ele.BlockNumber < e.BlockNumber ||
								(ele.BlockNumber == e.BlockNumber && ele.Index < e.Index) {
								logCh <- ele
							}
						}
					}

					first = false
				}
				logCh <- e

			case <-cctx.Done():
				left := 0

				for len(ch) > 0 {
					e := <-ch
					logCh <- e

					left++
				}

				wg.Done()
				close(ch)
				logger.Info().Msgf("event switch channel routine exit, drain events: %d", left)

				return
			}
		}
	}()

	logSub, err := es.wssClient.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		logger.Error().Err(err).Msg("subscribe filter logs failed")
		// logger.Fatal().Err(err).Msg("subscribe filter logs failed")
		wg.Done()
		cancel()
		errCh <- err

		return
	}

	go func() {
		subErr := <-logSub.Err()

		logger.Warn().Err(subErr).Msg("logs subscribe error occurs")
		errCh <- subErr
	}()

	logger.Info().Msg("Subscribe filter logs success")

	<-ctx.Done()
	logSub.Unsubscribe()
	cancel()

	wg.Done()

	logger.Info().Msg("logs subscription exited")
}

// / SubscribeEvents subscribe ETH events
// / @param ctx context
// / @param from The height start
// / @param to The end height. if 0, subscribe events
// / @param filterAddress address to query filter
// / @param eventCh event chan to receive []event.
func (es *EventSubscriber) SubscribeEvents(
	ctx context.Context,
	from uint64,
	filterAddress []string,
	eventsCh chan []types.Log,
) (err error) {
	filterAddr := []common.Address{}
	for _, addr := range filterAddress {
		filterAddr = append(filterAddr, common.HexToAddress(addr))
	}

	// first, drain logs between [from, end]
	start := from
	end := es.GetLatestBlockNumber()

	for {
		err := es.GetLogsFromTo(ctx, filterAddr, start, end, es.MaxStep, eventsCh)
		if err != nil {
			if err == errCanceled {
				err = nil
			}

			return err
		}

		logger.Info().Msgf("get events from [%d, %d]", start, end)
		start = end + 1
		end = es.GetLatestBlockNumber()

		if start >= end || end-start < 20 {
			logger.Info().Msgf("catch up latest block, end=%v latest=%v, switch to subscribe", start, end)
			break
		}
	}

	// then, subscribe logs
	query := ethereum.FilterQuery{Addresses: filterAddr}
	ch := make(chan types.Log, 1000)

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		_, err := es.wssClient.SubscribeFilterLogs(cctx, query, ch)
		if err != nil {
			log.Err(err).Msg("subscribe filter logs failed")
			logger.Fatal().Err(err).Msg("subscribe filter logs failed")
		}
	}()

	// indicate first receive events
	first := true

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("subscribe events was canceled")
			close(eventsCh)

			return

		case e := <-ch:
			if first {
				// fetch events between [start, latest]
				logger.Info().Msgf("got first subscribe event, blocknumber=%v, index=%d, latest fetched block %d", e.BlockNumber, e.Index, start)

				// very important !!!
				if start <= e.BlockNumber {
					logger.Info().Msgf("fetch logs between [%d, %d]", start, e.BlockNumber)

					events, err := es.GetLogsFromToReturn(ctx, filterAddr, start, e.BlockNumber, es.MaxStep)
					if err != nil {
						logger.Error().Err(err).Msgf("got logs [%d, %d] failed", start, e.BlockNumber)
						close(eventsCh)

						return err
					}

					idx := len(events)

					for i, ele := range events {
						if ele.BlockNumber < e.BlockNumber || (ele.BlockNumber == e.BlockNumber && ele.Index < e.Index) {
							continue
						} else {
							idx = i
							break
						}
					}

					if idx > 0 {
						eventsCh <- events[0:idx]
						logger.Info().Msgf("fill up %d events, start: [%d, %d], end: [%d, %d]",
							idx, events[0].BlockNumber, events[0].Index,
							events[idx-1].BlockNumber, events[idx-1].Index)
					} else {
						logger.Warn().Msgf("no events between last http fetch and wss subscribe")
					}
				}

				first = false
			}
			eventsCh <- []types.Log{e}
		}
	}
}

func (es *EventSubscriber) waitForRPCRateLimit(ms int64) {
	ts := time.Now().UnixMilli()
	if lastRPCTS == 0 {
		lastRPCTS = ts
		return
	}

	if ms == 0 {
		ms = int64(es.RateInterval) // nolint
	}

	if ts-lastRPCTS < ms {
		time.Sleep(time.Millisecond * time.Duration(ms+lastRPCTS-ts))
	}

	lastRPCTS = time.Now().UnixMilli()
}

// GetLogsFromTo get event logs base on filterAddress and block between [from, to], with context.
func (es *EventSubscriber) GetLogsFromTo(
	ctx context.Context,
	filterAddr []common.Address,
	from, to uint64,
	step uint32,
	eventsCh chan []types.Log,
) error {
	if from > to {
		logger.Warn().Msgf("GetLogsFromTo: invalid param: from(%d) >= to(%d)", from, to)
		return nil // fmt.Errorf("invalid param: from(%d) >= to(%d)", from, to)
	}

	start := from
	lastPrint := from

	end := from + uint64(step)
	if end > to {
		end = to
	}

	logger.Info().Msgf("get logs: [%d, %d]", from, to)

	var (
		err  error
		logs []types.Log
		// stepAdjusted = step
	)

	for end <= to {
		select {
		case <-ctx.Done():
			logger.Info().Msgf("GetLogsFromTo was canceled. from=%d current=%d to=%d", from, end, to)
			return errCanceled

		default:
		}

		logs, _, err = es.getLogsFromTo(filterAddr, start, end, step)
		if err != nil {
			return err
		}

		if end-lastPrint >= 100 {
			logger.Info().Msgf("get events from %d to %d: %d", start, end, len(logs))
			lastPrint = end
		}

		if len(logs) > 0 {
			eventsCh <- logs
		}

		start = end + 1
		if start > to {
			break
		}

		end = end + uint64(step)
		if end > to {
			end = to
		}
	}

	return nil
}

// GetLogsFromTo get event logs base on filterAddress and block between [from, to], with context.
func (es *EventSubscriber) GetLogsFromTo2(
	ctx context.Context,
	filterAddr []common.Address,
	from, to uint64,
	step uint32,
	eventCh chan types.Log,
) (int, error) {
	if from > to {
		logger.Warn().Msgf("GetLogsFromTo2: invalid param: from(%d) >= to(%d)", from, to)
		return 0, nil
	}

	tm0 := time.Now()
	start := from
	lastPrint := from

	end := from + uint64(step)
	if end > to {
		end = to
	}

	// logger.Info().Msgf("GetLogsFromTo2: [%d, %d]", from, to)

	var (
		err   error
		logs  []types.Log
		total int
		// stepAdjusted = step
	)

	for end <= to {
		select {
		case <-ctx.Done():
			logger.Info().Msgf("GetLogsFromTo2 was canceled. from=%d current=%d to=%d", from, end, to)
			return 0, errCanceled

		default:
		}

		logs, _, err = es.getLogsFromTo(filterAddr, start, end, step)
		if err != nil {
			return total, err
		}

		if end-lastPrint >= 100 {
			logger.Info().Msgf("GetLogsFromTo2: get events from %d to %d: %d", start, end, len(logs))
			lastPrint = end
		}

		for _, log := range logs {
			eventCh <- log
		}

		total += len(logs)

		start = end + 1
		if start > to {
			break
		}

		end = end + uint64(step)
		if end > to {
			end = to
		}
	}

	logger.Info().Msgf("GetLogsFromTo2: [%d, %d] complete, total: %d used: %v", from, to, total, time.Since(tm0))

	return total, nil
}

// GetLogsFromToParallel get event logs base on filterAddress and block between [from, to], return events.
func (es *EventSubscriber) GetLogsFromToParallel(
	ctx context.Context,
	filterAddr []common.Address,
	from, to uint64,
	step uint32,
) (evts []types.Log, err error) {
	wg := &sync.WaitGroup{}
	lk := sync.Mutex{}
	start := from
	end := from + uint64(step)
	routines := uint64(0)
	errorsCh := make(chan error, 1000)
	ts := time.Now().UnixMilli()

	maxRoutines := 1000 / es.RateInterval
	if maxRoutines > 100 {
		maxRoutines = 100
	}

	for end <= to {
		wg.Add(1)

		go func(s, e uint64) {
			res, _, err := es.getLogsFromTo(filterAddr, s, e, step)
			if err != nil {
				errorsCh <- err

				wg.Done()

				return
			}

			lk.Lock()
			evts = append(evts, res...)
			lk.Unlock()
			wg.Done()
		}(start, end)

		routines++
		start = end + 1
		end += uint64(step)

		if routines >= maxRoutines {
			wg.Wait()

			now := time.Now().UnixMilli()
			time.Sleep(time.Second)

			logger.Info().Msgf("wait routines complete: %d, endBlock: %d used: %d ms", routines, end, now-ts)
			// reset
			ts = now
			routines = 0
		}

		if len(errorsCh) > 0 {
			err = <-errorsCh

			wg.Wait()

			return nil, err
		}
	}

	wg.Wait()

	return
}

// GetLogsFromToReturn get event logs base on filterAddress and block between [from, to], return events.
func (es *EventSubscriber) GetLogsFromToReturn(
	ctx context.Context,
	filterAddr []common.Address,
	from, to uint64,
	step uint32,
) ([]types.Log, error) {
	if from > to {
		logger.Warn().Msgf("GetLogsFromToReturn: invalid param: from(%d) >= to(%d)", from, to)
		return []types.Log{}, nil // fmt.Errorf("invalid param: from(%d) >= to(%d)", from, to)
	}

	start := from
	lastPrint := from

	end := from + uint64(step)
	if end > to {
		end = to
	}

	events := []types.Log{}

	for end <= to {
		select {
		case <-ctx.Done():
			logger.Info().Msgf("GetLogsFromToReturn was canceled. from=%d current=%d to=%d", from, end, to)
			return nil, errCanceled

		default:
		}

		logs, _, err := es.getLogsFromTo(filterAddr, start, end, step)
		if err != nil {
			return nil, err
		}

		if end-lastPrint >= 100 {
			logger.Info().Msgf("GetLogsFromToReturn: get events from %d to %d: %d", start, end, len(logs))
			lastPrint = end
		}

		if len(logs) > 0 {
			events = append(events, logs...)
		} // else {
		//logger.Warn().Msgf("not found any events: blockRange=[%d, %d]", start, end)
		//}

		start = end + 1
		if start > to {
			break
		}

		end = end + uint64(step)
		if end > to {
			end = to
		}
	}

	logger.Info().Msgf("GetLogsFromToReturn: [%d, %d] events: %d", from, to, len(events))

	return events, nil
}

// getLogsFromTo get event logs base on filterAddress and block between [from, to].
func (es *EventSubscriber) getLogsFromTo(
	filterAddr []common.Address,
	from, to uint64,
	step uint32,
) (logs []types.Log, stepAdjusted uint32, err error) {
	var (
		endBlock uint64
		query    ethereum.FilterQuery
		retry    = uint32(0)
		ms       = int64(es.RateInterval) // nolint
	)

	// if filterAddr != nil {
	query.Addresses = filterAddr
	// }

	if step <= 0 {
		step = (es.MaxStep)
	}

	clientIdx := 0
	client := es.clients[clientIdx]

	for {
		endBlock = from + uint64(step)
		if endBlock > to {
			endBlock = to
		}

		query.FromBlock = big.NewInt(int64(from))   // nolint
		query.ToBlock = big.NewInt(int64(endBlock)) // nolint

		es.waitForRPCRateLimit(ms)

		result, err1 := client.FilterLogs(context.Background(), query)
		if err1 != nil {
			logger.Warn().Err(err1).Msg(fmt.Sprintf("getLogsFromTo failed: from=%d to=%d step=%d", from, endBlock, step))

			if retry >= es.MaxRetry {
				// sentry
				// utils.logger.Fatal().Msg("getLogsFromTo failed reach max retry: " + err.Error())
				stepAdjusted = step
				err = errReachMaxRetry

				logger.Error().Err(err1).Msg(fmt.Sprintf("getLogsFromTo failed: from=%d to=%d step=%d", from, endBlock, step))

				return
			}

			retry++

			if strings.Contains(err1.Error(), "rate limit") {
				ms *= 2
				es.waitForRPCRateLimit(ms)

				clientIdx++
				client = es.clients[clientIdx%len(es.clients)]
			} else {
				step = step / 2

				time.Sleep(200 * time.Millisecond)
			}

			continue
		}

		logs = append(logs, result...)

		if endBlock >= to {
			break
		}

		from = endBlock + 1
	}

	stepAdjusted = step

	return
}

func (es *EventSubscriber) GetLatestBlockNumber() uint64 {
	ms := int64(es.RateInterval) // nolint
	client := es.clients[0]

	for i := uint32(0); i < es.MaxRetry; i++ {
		es.waitForRPCRateLimit(ms)

		bn, err := client.BlockNumber(context.Background())
		if err == nil {
			return bn
		}

		logger.Warn().Msgf("get eth blockNumber failed: %v", err.Error())

		if strings.Contains(err.Error(), "rate limit") {
			ms *= 2
		}

		client = es.clients[int(i+1)%len(es.clients)]
	}

	// if any error, panic
	logger.Fatal().Msg("getLatestBlockNumber failed reach max retry")

	return 0
}
