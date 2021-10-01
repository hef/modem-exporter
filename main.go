package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

func doLogin(c *http.Client) (token string) {
	auth := b64.StdEncoding.EncodeToString([]byte("admin" + ":" + ""))

	req, err := http.NewRequest(http.MethodGet, "http://192.168.100.1/cmconnectionstatus.html?login_"+auth, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("admin", "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	csrfToken := string(body)
	//csrf token and potentially auth token
	// todo, save token as csrftoken in cookie
	c.Jar.SetCookies(req.URL, []*http.Cookie{
		&http.Cookie{
			Name:  "csrftoken",
			Value: csrfToken,
		},
	})
	log.Printf("token: %s", string(body))
	return string(body)
}

func main() {

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	c := http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	token := doLogin(&c)

	req, err := http.NewRequest(http.MethodGet, "http://192.168.100.1/cmconnectionstatus.html?ct_"+token, nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("dump: %s", string(body))
}
