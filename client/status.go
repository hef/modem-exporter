package client

import (
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
)

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
		q.Set("ct_" + csrfToken, "")
		u.RawQuery = q.Encode()
		// Encode adds a trailing "=", but we need to remove it
		u.RawQuery = u.RawQuery[:len(u.RawQuery)-1]
	}

	c.logger.Debug("csrf token",
		zap.String("csrftoken", csrfToken),
		zap.String("url", u.String()),
	)

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

func (c *Client) Status() (err error) {

	req, err := c.newRequest("http://192.168.100.1/cmconnectionstatus.html")
	if err != nil {
		return err
	}

	doc, err := c.do(req)
	if err != nil {
		return err
	}

	html, _ := doc.Html()
	log.Printf("%s", html)

	return nil
}
