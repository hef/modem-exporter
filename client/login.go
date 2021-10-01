package client

import (
	b64 "encoding/base64"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *Client) login() error {
	auth := b64.StdEncoding.EncodeToString([]byte("admin" + ":" + c.password))

	req, err := http.NewRequest(http.MethodGet, "http://192.168.100.1/cmconnectionstatus.html?login_"+auth, nil)
	if err != nil {
		c.logger.Error("error create request",
			zap.Error(err),
		)
		return err
	}
	req.SetBasicAuth("admin", c.password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("error issuing login request",
			zap.Error(err),
		)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("error reading login response",
			zap.Error(err),
		)
		return err
	}
	csrfToken := string(body)

	c.client.Jar.SetCookies(req.URL, []*http.Cookie{
		&http.Cookie{
			Name:  "csrftoken",
			Value: csrfToken,
		},
	})
	return nil
}

func (c *Client) csrfToken() string {
	u, _ := url.Parse("http://192.168.100.1/")
	cookies := c.client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "csrftoken" {
			return cookie.Value
		}
	}
	return ""
}
