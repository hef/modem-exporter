package client

import "log"

type DownstreamBoundedChannels struct {
	ChannelId      string
	LockStatus     string
	Modulation     string
	Frequency      string
	Power          string
	SnrSmr         string
	Corrected      string
	Uncorrectables string
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


	// I really need xpath
	//table := doc.Find("table >")

	//Downstream Bonded Channels

	html, _ := doc.Html()
	log.Printf("%s", html)

	return nil
}
