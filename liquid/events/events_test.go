package events

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/pkg/utils"
)

func createES(rpc, wss string) *EventSubscriber {
	if rpc == "" {
		rpc = os.Getenv("BASE_RPC")
	}

	if wss == "" {
		wss = os.Getenv("BASE_WSS")
	}

	if rpc == "" {
		rpc = "https://mainnet.base.org"
	}

	if wss == "" {
		wss = "wss://base-rpc.publicnode.com"
	}

	return MustNewEventSubscriber(rpc, wss, true, 500, 0, 0, 0)
}

func TestGetLogsFromTo(t *testing.T) {
	utils.SkipCI(t)

	es := createES("", "")
	start := uint64(16543038)
	end := uint64(16543377)

	// steps := []uint32{1, 2, 3, 5, 10, 24, 25, 50, 73, 100, 111, 338, 339, 340, 500}
	steps := []uint32{1, 2, 41, 100, 338, 339, 340, 500}
	counts := make([]uint32, len(steps))

	for i, step := range steps {
		stop := make(chan bool)
		stopped := make(chan bool)
		eventCh := make(chan []types.Log, 100)

		go func(idx int) {
			for {
				select {
				case <-stop:
				L:
					for {
						_, ok := <-eventCh
						if !ok {
							break L
						} else {
							counts[idx]++
						}
					}

					close(stopped)

					return

				case <-eventCh:
					counts[idx]++
				}
			}
		}(i)

		err := es.GetLogsFromTo(context.Background(), []common.Address{}, start, end, step, eventCh)
		assert.Nil(t, err)
		// t.Log("GetLogsFromTo success")

		close(stop)
		close(eventCh)
		<-stopped
	}

	for i := 0; i < len(counts)-1; i++ {
		assert.Equal(t, counts[i], counts[i+1])
	}
}

func TestSubscribeLogs(t *testing.T) {
	utils.SkipCI(t)

	es := createES("", "")
	start := es.GetLatestBlockNumber() - 300

	stopped := make(chan bool)
	eventCh := make(chan []types.Log, 1000)
	blocknumber := uint64(0)
	count := 0

	go func() {
		for {
			e, ok := <-eventCh
			if !ok {
				close(stopped)
				return
			}

			if len(e) == 0 {
				continue
			}

			blocknumber = e[len(e)-1].BlockNumber
			count += len(e)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(60 * time.Second)
		cancel()
	}()

	err := es.SubscribeEvents(ctx, start, []string{}, eventCh)
	assert.Nil(t, err)
	<-stopped

	t.Logf("count=%d blocknumber=%d", count, blocknumber)
}

func TestSubscribeBlockAndEvents(t *testing.T) {
	utils.SkipCI(t)

	es := createES("", "")
	ctx, cancel := context.WithCancel(context.Background())
	blockCh := make(chan *types.Header, 10)
	eventCh := make(chan types.Log, 1000)

	wg := &sync.WaitGroup{}
	errch := make(chan error)
	es.SubscribeEventsFrom(ctx, nil, 0, 0, eventCh, errch, wg)

	tmr := time.NewTimer(time.Second * 10)
L:
	for {
		select {
		case b := <-blockCh:
			t.Logf("new block: %d %v", time.Now().UnixMilli(), b.Number)

		case e := <-eventCh:
			t.Logf("new event: %d %v %d %s", time.Now().UnixMilli(), e.BlockNumber, e.Index, e.TxHash.String())

		case <-tmr.C:
			tmr.Stop()
			cancel()
			t.Logf("time end")
			break L
		}
	}

	wg.Wait()
}

func TestGetLogsFromTo2(t *testing.T) {
	utils.SkipCI(t)

	es := createES("", "")
	start := uint64(18101616)
	end := uint64(18101619)
	step := uint32(100)
	addr := []common.Address{common.HexToAddress("0xb2cc224c1c9fee385f8ad6a55b4d94e92359dc59")}

	evts, _, err := es.getLogsFromTo(addr, start, end, step)
	assert.Nil(t, err)

	t.Logf("events: %d", len(evts))
	t.Logf("first event: %d", evts[0].BlockNumber)
	t.Logf("last event: %d", evts[len(evts)-1].BlockNumber)

	for _, evt := range evts {
		if len(evt.Topics) == 0 {
			t.Logf("%v", evt)
		} else {
			t.Logf("event: block=%d topic=%s txhash=%s", evt.BlockNumber, evt.Topics[0].String(), evt.TxHash.String())
		}
	}
}

func TestGetLogsFromToParallel(t *testing.T) {
	utils.SkipCI(t)

	start := uint64(0)
	end := uint64(17804820)
	es := createES("https://xxx.base-mainnet.quiknode.pro/", "")
	es.RateInterval = 5
	step := uint32(5000)
	addr := []common.Address{
		common.HexToAddress("0x33128a8fC17869897dcE68Ed026d694621f6FDfD"),
		common.HexToAddress("0xc35DADB65012eC5796536bD9864eD8773aBc74C4"),
		common.HexToAddress("0x41ff9AA7e16B8B1a8a8dc4f0eFacd93D02d071c9"),
	}
	evts, err := es.GetLogsFromToParallel(context.Background(), addr, start, end, step)
	assert.Nil(t, err)

	t.Logf("events: %d", len(evts))
}

func TestWithCancel(t *testing.T) {
	ctx := context.Background()
	cctx, cancel := context.WithCancel(ctx)
	_ = cctx

	time.Sleep(time.Second)
	cancel()
	t.Logf("after cancel")
}

func TestPublicWSS(t *testing.T) {
	utils.SkipCI(t)

	wss := "wss://base-rpc.publicnode.com"
	wssClient, err := ethclient.Dial(wss)
	assert.Nil(t, err)

	ctx := context.Background()
	blockCh := make(chan *types.Header, 5)
	sub, err := wssClient.SubscribeNewHead(ctx, blockCh)
	assert.Nil(t, err)

	go func() {
		for {
			b := <-blockCh
			t.Logf("new heade: %d", b.Number.Uint64())
		}
	}()

	time.Sleep(5 * time.Second)
	sub.Unsubscribe()
}
