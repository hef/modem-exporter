package client

import (
	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/html"
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

	return parseStatusPage(c.logger, doc)
}

func parseStatusPage(logger *zap.Logger, doc *html.Node) (downstreamBondendChannels []DownstreamBoundedChannel, err error) {
	rows := htmlquery.Find(doc, `//table[.//th[.="Downstream Bonded Channels"]]//tr`)
	if rows == nil {
		logger.Debug("couldn't find downstream bonded channels table.",
			zap.String("doc", debugPrint(doc)),
		)
		logger.Error("couldn't find 'Downstream Bonded Channels' table")
		return nil, errors.New("couldn't find downstream bonded channels table")
	}
	if len(rows) <= 2 {
		logger.Error("couldn't find enough rows in downstream bonded channels table", zap.Int("rows", len(rows)))
		return nil, errors.New("downstream bonded channels table didn't have enoigh rows")
	}

	var x []DownstreamBoundedChannel

	for _, row := range rows[2:] {
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
