package client

import (
	b64 "encoding/base64"
	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
)

func isLoginPage(doc *html.Node) bool {
	titleNode := htmlquery.FindOne(doc, "//title[text() = 'Login']")
	if titleNode != nil {
		return true
	}
	return false
}

func (c *Client) login() error {
	auth := b64.StdEncoding.EncodeToString([]byte("admin" + ":" + c.password))

	req, err := http.NewRequest(http.MethodGet, "https://192.168.100.1/cmconnectionstatus.html?login_"+auth, nil)
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
			Name:     "csrftoken",
			Value:    csrfToken,
			Secure:   true,
			HttpOnly: true,
		},
	})
	return nil
}

func (c *Client) setCsrfTokenOnUrl(u *url.URL) {
	csrfToken := c.csrfToken()
	if csrfToken != "" {
		q := u.Query()
		q.Set("ct_"+csrfToken, "")
		u.RawQuery = q.Encode()
		// Encode adds a trailing "=", but we need to remove it
		u.RawQuery = u.RawQuery[:len(u.RawQuery)-1]
	}
}

func (c *Client) csrfToken() string {
	u, _ := url.Parse("https://192.168.100.1/")
	cookies := c.client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "csrftoken" {
			return cookie.Value
		}
	}
	return ""
}
