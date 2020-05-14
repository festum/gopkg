package kafka

import (
	"crypto/tls"

	"github.com/caarlos0/env"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type (
	config struct {
		Endpoint      string `env:"KAFKA_URL" envDefault:"localhost:9093"`
		Network       string `env:"KAFKA_NETWORK" envDefault:"tcp"`
		Username      string `env:"KAFKA_USERNAME" envDefault:""`
		Password      string `env:"KAFKA_PASSWORD"`
		ClientTimeout int64  `env:"KAFKA_CLIENT_TIMEOUT_SECONDS" envDefault:"60"`
		MinBytes      int
		MaxBytes      int
	}
)

func NewConfig() *config {
	var k config
	if err := env.Parse(&k); err != nil {
		_logger.Fatalw("failed on creating kafka config", "error", err)
		return nil
	}
	k.MinBytes = 1
	k.MaxBytes = 10e6 // 10MB
	return &k
}

func (k *config) getDialer() (dialer *kafka.Dialer) {
	dialer = kafka.DefaultDialer
	if k.Username == "" {
		return
	}
	dialer.SASLMechanism = plain.Mechanism{
		Username: k.Username,
		Password: k.Password,
	}
	// Use TLS config to enable SSL
	dialer.TLS = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
	}
	return
}
