package collector

import (
	"go.uber.org/zap"
	"modem-exporter/client"
)

type Option func(collector *Collector)

func WithLogger(logger *zap.Logger) Option {
	return func(collector *Collector) {
		collector.logger = logger
	}
}

func WithClient(client *client.Client) Option {
	return func(collector *Collector) {
		collector.client = client
	}
}
