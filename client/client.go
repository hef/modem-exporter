package client

import (
	"crypto/tls"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	client   *http.Client
	logger   *zap.Logger
	password string
}

func New(options ...Option) (*Client, error) {
	c := Client{}

	for _, option := range options {
		option(&c)
	}
	if c.client == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			panic(err)
		}

		c.client = &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	if c.logger == nil {
		c.logger = zap.NewNop()
	}

	return &c, nil
}
