package events

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"starbase.ag/liquidity/pkg/logger"
)

// CreateKafkaProducer new kafka producer.
func CreateKafkaProducer(topic, brokers string) *kafka.Writer {
	// Producer.RequiredAcks = sarama.WaitForAll
	addrs := strings.SplitSeq(brokers, ",")
	for addr := range addrs {
		conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, 0)
		if err == nil {
			// close the connection because we won't be using it
			conn.Close()
			break
		}
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(brokers, ",")...),
		Topic:        topic,
		BatchSize:    1000,
		BatchBytes:   1024 * 1024 * 500,
		RequiredAcks: kafka.RequireAll,
		// Balancer: &kafka.LeastBytes{},
	}
}

// ProduceEvents Listen event chan, send event to kafka.
func ProduceEvents(ctx context.Context, topic, brokers string, eventsCh chan []types.Log) error {
	p := &kafka.Writer{
		Addr:  kafka.TCP(strings.Split(brokers, ",")...),
		Topic: topic,
		// Balancer: &kafka.LeastBytes{},
	}

	maxRetry := uint(5)
	maxEvents := 5000

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("process was canceled, exit")
			p.Close()

			return nil

		case evts := <-eventsCh:
			// produce event
			if len(evts) > maxEvents {
				for i := 0; i < len(evts); i += maxEvents {
					end := i + maxEvents
					if end > len(evts) {
						end = len(evts)
					}

					part := evts[i:end]
					if err := produceMessage(p, part, maxRetry); err != nil {
						logger.Fatal().Msgf("produce event failed: %v", err)
					}
				}
			} else {
				if err := produceMessage(p, evts, maxRetry); err != nil {
					logger.Fatal().Msgf("produce event failed: %v", err)
				}
			}
		}
	}
}

// ConsumeEvents consume events from kafka.
func ConsumeEvents(ctx context.Context, topic, brokers, group string, fn func([]types.Log) error) error {
	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		GroupID: group,
		Topic:   topic,
	})

	defer c.Close()

	smpLog := log.Sample(&zerolog.BasicSampler{N: 100})

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("consumer was canceled")
			return nil

		default:
			evts, err := c.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil // Return nil if the context was canceled
				}

				logger.Error().Err(err).Msg("fetch message failed")

				return err
			}

			var logs []types.Log
			if err := json.Unmarshal(evts.Value, &logs); err != nil {
				logger.Error().Err(err).Msgf("unmarshal events failed")
				return err
			}

			smpLog.Info().Msgf("recv events, startBlock=%d offset=%d", logs[0].BlockNumber, evts.Offset)

			if fn != nil {
				if err := fn(logs); err != nil {
					logger.Error().Err(err).Msg("handle events failed")
					return err
				}
			}

			err = c.CommitMessages(ctx, evts)
			if err != nil {
				logger.Error().Err(err).Msgf("commit message failed: topic=%s offset=%v startBlock=%d endBlock=%d",
					topic, evts.Offset, logs[0].BlockNumber, logs[len(logs)-1].BlockNumber)
			}
		}
	}
}

func produceMessage(p *kafka.Writer, evts []types.Log, maxRetry uint) error {
	if len(evts) == 0 {
		return nil
	}

	buf, err := json.Marshal(evts)
	if err != nil {
		logger.Fatal().Err(err).Msg("marshal events failed")
		return err
	}

	retry := uint(0)
	ctx := context.Background()

	for retry < maxRetry {
		err = p.WriteMessages(ctx, kafka.Message{
			Value: buf,
		})
		if err == nil {
			return nil
		}

		retry++
	}

	return errors.New("produce message failed with retry: " + err.Error())
}

// FetchLatestBlock fetch watermark.
func FetchLatestBlock(topic, brokers, group string) (uint64, *types.Log, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers, topic, 0)
	if err != nil {
		return 0, nil, err
	}

	first, last, err := conn.ReadOffsets()
	if err != nil {
		return 0, nil, err
	}

	println("first", first, "last", last)
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   strings.Split(brokers, ","),
		Topic:     topic,
		Partition: 0,
	})
	defer r.Close() // Ensure the reader is closed when the function exits

	if err := r.SetOffset(last - 1); err != nil {
		return 0, nil, err
	}

	m, err := r.ReadMessage(context.Background())
	if err != nil {
		return 0, nil, err
	}

	var evts []types.Log

	err = json.Unmarshal(m.Value, &evts)
	if err != nil {
		logger.Error().Err(err).Msg("unmarshal data failed")
		return 0, nil, err
	}

	evt := evts[len(evts)-1]

	return evt.BlockNumber, &evt, nil
}

// func commitOffset(c *kafka.Consumer, topic string, partition int, offset int) error {
// 	res, err := c.CommitOffsets([]kafka.TopicPartition{{
// 		Topic:     &topic,
// 		Partition: int32(partition),
// 		Offset:    kafka.Offset(offset),
// 	}})
// 	if err != nil {
// 		return err
// 	}

// 	logger.Info().Msgf("Partition %d offset %d committed successfully", res[0].Partition, offset)

// 	return nil
// }
