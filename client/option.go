package client

import (
	"go.uber.org/zap"
	"net/http"
)

type Option func(*Client)

func WithPassword(password string) Option {
	return func(c *Client) {
		c.password = password
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}
