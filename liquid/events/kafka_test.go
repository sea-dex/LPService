package events

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/pkg/utils"
)

func TestFetchLatestBlock(t *testing.T) {
	utils.SkipCI(t)

	topic := "sep-testnet"
	brokers := "localhost:9092"
	block, _, err := FetchLatestBlock(topic, brokers, "producer")
	assert.Nil(t, err)

	_ = block
}

func TestProducer(t *testing.T) {
	utils.SkipCI(t)

	ctx, cancel := context.WithCancel(context.Background())
	topic := "sep-testnet111"
	brokers := "localhost:9092"
	eventsCh := make(chan []types.Log, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err := ProduceEvents(ctx, topic, brokers, eventsCh)
		assert.Nil(t, err)
		wg.Done()
	}()

	ts := time.Now().Unix()
	ms := big.NewInt(time.Now().UnixMilli())
	eventsCh <- []types.Log{{BlockNumber: uint64(ts), Address: common.BigToAddress(ms), Topics: []common.Hash{common.BigToHash(ms)}}} // nolint

	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}

func TestProduceMessageSync(t *testing.T) {
	utils.SkipCI(t)

	brokers := "localhost:9092"
	topic := "sep-testnet"
	p := CreateKafkaProducer(topic, brokers)

	tmr := time.NewTimer(30 * time.Second)
	ctx := context.Background()
	count := 0
	batch := 10000
L:
	for {
		select {
		case <-tmr.C:
			t.Log("time end")
			break L

		default:
		}

		// produce message
		msgs := []kafka.Message{}
		for i := 0; i < batch; i++ {
			msgs = append(msgs, kafka.Message{
				Value: []byte(fmt.Sprint(count + i)),
			})
		}
		err := p.WriteMessages(ctx, msgs...)
		count += batch
		if err == nil {
			continue
		}

		t.Log(err.Error())
		break
	}
}

func TestProduceMessage(t *testing.T) {
	utils.SkipCI(t)

	brokers := "localhost:9092"
	topic := "sep-testnet"
	p := CreateKafkaProducer(topic, brokers)

	evts := []types.Log{{BlockNumber: 100, Address: common.HexToAddress("0x0123456"), Topics: []common.Hash{}}}
	err := produceMessage(p, evts, 5)
	assert.Nil(t, err)
}

func TestConsumer1(t *testing.T) {
	utils.SkipCI(t)

	brokers := "localhost:9092"
	topic := "base-mainnet"
	group := "consumer"

	ctx, cancel := context.WithCancel(context.Background())
	// FetchLatestBlock(topic, brokers, group)
	go func() {
		err := ConsumeEvents(ctx, topic, brokers, group, nil)
		if err != nil {
			panic(err.Error())
		}
	}()
	// 48675 48857 48993
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

func TestConsumerFrom(t *testing.T) {
	utils.SkipCI(t)

	brokers := "localhost:9092"
	topic := "base-mainnet"
	group := "consumer"

	ctx, cancel := context.WithCancel(context.Background())
	// FetchLatestBlock(topic, brokers, group)
	go func() {
		err := ConsumeEvents(ctx, topic, brokers, group, nil)
		if err != nil {
			panic(err.Error())
		}
	}()
	// 48675 48857 48993
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}
