package kafka

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type (
	Writer interface {
		Topic() string
		Write(key []byte, value []byte) error
		Shutdown()
	}
	writer struct {
		*kafka.Writer
	}
)

func NewWriter(topic string) (Writer, error) {
	conf := NewConfig()
	wtr := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{conf.Endpoint},
		Topic:    topic,
		Balancer: &kafka.Murmur2Balancer{},
		Dialer:   conf.getDialer(),
	})

	out := writer{wtr}

	return out, nil
}

func (w writer) Topic() string {
	return w.Stats().Topic
}

func (w writer) Write(key []byte, value []byte) error {
	ctx, cancel := getTimeLimitedContext(NewConfig().ClientTimeout)
	defer cancel()
	return w.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

func (w writer) Shutdown() {
	if err := w.Close(); err != nil {
		return
	}
}

func WriteOnce(topic string, key, value []byte) error {
	w, err := NewWriter(topic)
	if err != nil {
		return err
	}
	if err := w.Write(key, value); err != nil {
		return err
	}
	w.Shutdown()
	return nil
}

func getTimeLimitedContext(timeout int64) (context.Context, context.CancelFunc) {
	// If app fail to cancel the context, the goroutine that WithCancel
	//  or WithTimeout created will be retained in memory indefinitely
	//  (until the program shuts down), causing a memory leak.
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
