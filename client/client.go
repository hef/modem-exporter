package client

import (
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func (c *Client) newRequest(path string) (*http.Request, error) {

	u, err := url.Parse(path)
	if err != nil {
		c.logger.Error("failed to parse path",
			zap.String("path", path),
			zap.Error(err),
		)
	}

	csrfToken := c.csrfToken()
	if csrfToken != "" {
		q := u.Query()
		q.Set("ct_"+csrfToken, "")
		u.RawQuery = q.Encode()
		// Encode adds a trailing "=", but we need to remove it
		u.RawQuery = u.RawQuery[:len(u.RawQuery)-1]
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		c.logger.Error("error creating request",
			zap.String("url", u.String()),
			zap.Error(err),
		)
		return nil, err
	}
	return req, nil
}

func (c *Client) do(req *http.Request) (*goquery.Document, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("error sending request",
			zap.Error(err),
		)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.logger.Error("failed to parse response body",
			zap.Error(err),
		)
	}

	titleNode := doc.Find("title")
	if titleNode.Text() == "Login" {
		err = c.login()
		if err != nil {
			c.logger.Error("error logging in",
				zap.Error(err),
			)
			return nil, err
		}
	}

	req, err = c.newRequest(req.URL.String())
	if err != nil {
		c.logger.Error("error sending request",
			zap.Error(err),
		)
		return nil, err
	}
	resp, err = c.client.Do(req)
	if err != nil {
		c.logger.Error("error sending request",
			zap.Error(err),
		)
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	return doc, err
}
