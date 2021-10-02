package client

import (
	"github.com/antchfx/htmlquery"
	"strconv"
	"strings"
)

type DownstreamBoundedChannel struct {
	ChannelId      int
	LockStatus     string
	Modulation     string
	Frequency      int
	Power          float64
	SnrSmr         float64
	Corrected      int
	Uncorrectables int
}

func (c *Client) Status() (downstreamBondendChannels []DownstreamBoundedChannel, err error) {

	req, err := c.newRequest("https://192.168.100.1/cmconnectionstatus.html")
	if err != nil {
		return nil, err
	}

	doc, err := c.do(req)
	if err != nil {
		return nil, err
	}

	//log.Printf("original: %s", debugPrint(doc))

	rows := htmlquery.Find(doc, `//table[.//th[.="Downstream Bonded Channels"]]//tr`)[2:]
	var x []DownstreamBoundedChannel

	for _, row := range rows {
		//log.Printf("result: %s", debugPrint(row))
		channelId := htmlquery.FindOne(row, `//td[1]/text()`)
		lockStatus := htmlquery.FindOne(row, `//td[2]/text()`)
		modulation := htmlquery.FindOne(row, `//td[3]/text()`)
		frequency := htmlquery.FindOne(row, `//td[4]/text()`)
		power := htmlquery.FindOne(row, `//td[5]/text()`)
		snr := htmlquery.FindOne(row, `//td[6]/text()`)
		corrected := htmlquery.FindOne(row, `//td[7]/text()`)
		Uncorrectables := htmlquery.FindOne(row, `//td[8]/text()`)

		data := DownstreamBoundedChannel{}
		data.ChannelId, err = strconv.Atoi(channelId.Data)
		if err != nil {
			continue
		}
		data.LockStatus = lockStatus.Data
		data.Modulation = modulation.Data
		data.Frequency, err = strconv.Atoi(strings.TrimSuffix(frequency.Data, " Hz"))
		if err != nil {
			continue
		}
		data.Power, err = strconv.ParseFloat(strings.TrimSuffix(power.Data, " dBmV"), 64)
		if err != nil {
			continue
		}
		data.SnrSmr, err = strconv.ParseFloat(strings.TrimSuffix(snr.Data, " dB"), 64)
		if err != nil {
			continue
		}
		data.Corrected, err = strconv.Atoi(corrected.Data)
		if err != nil {
			continue
		}
		data.Uncorrectables, err = strconv.Atoi(Uncorrectables.Data)
		if err != nil {
			continue
		}
		x = append(x, data)
	}

	return x, nil
}
