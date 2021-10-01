package main

import (
	"go.uber.org/zap"
	"modem-exporter/client"
	"os"
)

func main() {

	logger, err := zap.NewDevelopment()
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
	c.Status()
}
