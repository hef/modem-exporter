package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"modem-exporter/client"
	"modem-exporter/collector"
	"net/http"
	"os"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	options := []client.Option{
		client.WithLogger(logger.Named("client")),
	}

	if password, ok := os.LookupEnv("PASSWORD"); ok {
		options = append(options, client.WithPassword(password))
	}

	c, err := client.New(options...)
	if err != nil {
		logger.Error("failed to create client", zap.Error(err))
		return
	}

	reg := prometheus.NewRegistry()

	modemCollector, err := collector.New(
		collector.WithLogger(logger.Named("collector")),
		collector.WithClient(c),
	)
	if err != nil {
		logger.Error("failed to create collector",
			zap.Error(err))
		return
	}

	reg.MustRegister(modemCollector)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":2112", nil)

}
