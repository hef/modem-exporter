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

type UpstreamBoundedChannel struct {
	Channel       int
	ChannelId     int
	LockStatus    string
	UsChannelType string
	Frequency     int
	Width         int
	Power         float64
}

func (c *Client) Status() (downstreamBondendChannels []DownstreamBoundedChannel, upstreamBondedChannel []UpstreamBoundedChannel, err error) {

	req, err := c.newRequest("https://192.168.100.1/cmconnectionstatus.html")
	if err != nil {
		return nil, nil, err
	}

	doc, err := c.do(req)
	if err != nil {
		return nil, nil, err
	}

	return parseStatusPage(c.logger, doc)
}

func parseStatusPage(logger *zap.Logger, doc *html.Node) (downstreamBondendChannels []DownstreamBoundedChannel, upstreamBondedChannels []UpstreamBoundedChannel, err error) {
	rows := htmlquery.Find(doc, `//table[.//th[.="Downstream Bonded Channels"]]//tr`)
	if rows == nil {
		logger.Debug("couldn't find downstream bonded channels table.",
			zap.String("doc", debugPrint(doc)),
		)
		logger.Error("couldn't find 'Downstream Bonded Channels' table")
		return nil, nil, errors.New("couldn't find downstream bonded channels table")
	}
	if len(rows) <= 2 {
		logger.Error("couldn't find enough rows in downstream bonded channels table", zap.Int("rows", len(rows)))
		return nil, nil, errors.New("downstream bonded channels table didn't have enoigh rows")
	}

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
		if channelId == nil {
			logger.Error("failed to extract channel id from row in Downstream table, skipping row")
			continue
		}
		data.ChannelId, err = strconv.Atoi(channelId.Data)
		if err != nil {
			logger.Error("failed to parse channel id from row in Downstream table, skipping row")
			continue
		}

		if lockStatus != nil {
			data.LockStatus = lockStatus.Data
		} else {
			logger.Debug("failed to find lockstatus in Downstream table")
		}

		if modulation != nil {
			data.Modulation = modulation.Data
		} else {
			logger.Debug("failed to find modulation in Downstream table")
		}
		if frequency != nil {
			data.Frequency, err = strconv.Atoi(strings.TrimSuffix(frequency.Data, " Hz"))
			if err != nil {
				logger.Debug("failed to parse frequency data in Downstream table")
			}
		}

		if power != nil {
			data.Power, err = strconv.ParseFloat(strings.TrimSuffix(power.Data, " dBmV"), 64)
			if err != nil {
				logger.Debug("failed to parse power data in Downstream table")
			}
		}

		if snr != nil {
			data.SnrSmr, err = strconv.ParseFloat(strings.TrimSuffix(snr.Data, " dB"), 64)
			if err != nil {
				logger.Debug("failed to parse SNR in downstream table")
			}
		}

		if corrected != nil {
			data.Corrected, _ = strconv.Atoi(corrected.Data)
		}

		if Uncorrectables != nil {
			data.Uncorrectables, _ = strconv.Atoi(Uncorrectables.Data)
		}
		downstreamBondendChannels = append(downstreamBondendChannels, data)
	}

	rows = htmlquery.Find(doc, `//table[.//th[.="Upstream Bonded Channels"]]//tr`)

	for _, row := range rows[2:] {
		channel := htmlquery.FindOne(row, `//td[1]/text()`)
		channelId := htmlquery.FindOne(row, `//td[2]/text()`)
		lockStatus := htmlquery.FindOne(row, `//td[3]/text()`)
		usChannelType := htmlquery.FindOne(row, `//td[4]/text()`)
		frequency := htmlquery.FindOne(row, `//td[5]/text()`)
		width := htmlquery.FindOne(row, `//td[6]/text()`)
		power := htmlquery.FindOne(row, `//td[7]/text()`)

		data := UpstreamBoundedChannel{}

		data.Channel, err = strconv.Atoi(channel.Data)
		if err != nil {
			logger.Debug("failed to extract channel from upstream channel")
		}

		data.ChannelId, err = strconv.Atoi(channelId.Data)
		if err != nil {
			logger.Debug("failed to extract channel id from upstream channel")
		}

		data.LockStatus = lockStatus.Data
		data.UsChannelType = usChannelType.Data

		data.Frequency, err = strconv.Atoi(strings.TrimSuffix(frequency.Data, " Hz"))
		if err != nil {
			logger.Debug("failed to extract frequency from upstream channel")
		}

		data.Width, err = strconv.Atoi(strings.TrimSuffix(width.Data, " Hz"))
		if err != nil {
			logger.Debug("failed to parse width from upstream channel")
		}

		data.Power, err = strconv.ParseFloat(strings.TrimSuffix(power.Data, " dBmV"), 64)
		if err != nil {
			logger.Debug("failed to parse power from upstream channel")
		}
		upstreamBondedChannels = append(upstreamBondedChannels, data)
	}
	return downstreamBondendChannels, upstreamBondedChannels, nil
}
