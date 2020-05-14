package kafka

import (
	"github.com/festum/gopkg/logger"
	kafka "github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

var _logger = logger.NewLogger()

type (
	Reader interface {
		Topic() string
		ReadLag() (totalLag int64, err error)
		ReadAll(stoppable bool, fn func([]byte))
		Shutdown()
	}

	reader struct {
		topic   string
		config  *config
		readers []*kafka.Reader
	}
)

func NewReader(topic string) (Reader, error) {
	conf := NewConfig()
	conn, err := NewConnection(conf)
	if err != nil {
		return nil, err
	}
	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return nil, err
	}
	// For topics with multiple partitions we need to create readers, one per partition
	r := reader{
		topic:   topic,
		config:  conf,
		readers: make([]*kafka.Reader, len(partitions)),
	}
	for i, partition := range partitions {
		r.readers[i] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{conf.Endpoint},
			Topic:          topic,
			Partition:      partition.ID,
			MinBytes:       conf.MinBytes,
			MaxBytes:       conf.MaxBytes,
			SessionTimeout: 10,
			Dialer:         conf.getDialer(),
		})
	}

	return r, nil
}

func NewConnection(c *config) (*kafka.Conn, error) {
	ctx, cancel := getTimeLimitedContext(NewConfig().ClientTimeout)
	defer cancel()
	return c.getDialer().DialContext(ctx, c.Network, c.Endpoint)
}

func (r reader) Topic() string {
	return r.topic
}

func (r reader) ReadLag() (totalLag int64, err error) {
	for _, rr := range r.readers {
		ctx, cancel := getTimeLimitedContext(NewConfig().ClientTimeout)
		defer cancel()
		lag, err := rr.ReadLag(ctx)
		if err != nil {
			return totalLag, err
		}
		totalLag += lag
	}
	return
}

func (r reader) ReadAll(stoppable bool, fn func([]byte)) {
	for _, rr := range r.readers {
		go func(reader *kafka.Reader) {
			for {
				ctx := context.Background()
				if stoppable {
					var cancel context.CancelFunc
					ctx, cancel = getTimeLimitedContext(NewConfig().ClientTimeout)
					defer cancel()
					// allows single ReadLag to fail silently and return 0, the 'ReadLag()' verifies progress for the whole topic.
					lag, _ := reader.ReadLag(ctx)
					if lag <= 0 {
						_logger.Debugf("reached the end of %s in partition %d", r.Topic(), reader.Config().Partition)
						break
					}
				}
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					_logger.Errorw("failed on reading message", "topic", reader.Config().Topic, "partition", reader.Config().Partition, "error", err)
					break
				}
				fn(msg.Value)
			}
		}(rr)
	}
}

func (r reader) Shutdown() {
	go func() {
		for _, rr := range r.readers {
			config := rr.Config()
			if err := rr.Close(); err != nil {
				_logger.Infow("failed on closing reader", "topic", config.Topic, "partition", config.Partition, "error", err)
			}
		}
	}()

}
