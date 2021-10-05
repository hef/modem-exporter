package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"modem-exporter/client"
	"strconv"
)

type Collector struct {
	logger *zap.Logger
	client *client.Client

	downstreamFrequencyMetric      *prometheus.Desc
	downstreamLockStatusMetric     *prometheus.Desc
	downstreamPowerMetric          *prometheus.Desc
	downstreamSnrMetric            *prometheus.Desc
	downstreamCorrectedMetric      *prometheus.Desc
	downstreamUncorrectablesMetric *prometheus.Desc

	upstreamLockStatusMetric *prometheus.Desc
	upstreamFrequencyMetric  *prometheus.Desc
	upstreamWidthMetric      *prometheus.Desc
	upstreamPowerMetric      *prometheus.Desc
}

func New(options ...Option) (*Collector, error) {
	c := Collector{
		downstreamFrequencyMetric: prometheus.NewDesc(
			"modem_downstream_channel_frequency",
			"Downstream channel frequency.",
			[]string{"channel"},
			nil,
		),
		downstreamLockStatusMetric: prometheus.NewDesc(
			"modem_downstream_channel_locked",
			"Downstream channel Lock Status, 0 = Not Locked, 1 = Locked",
			[]string{"channel"},
			nil,
		),
		downstreamPowerMetric: prometheus.NewDesc(
			"modem_downstream_channel_power",
			"Downstream channel power (dBmv)",
			[]string{"channel"},
			nil,
		),
		downstreamSnrMetric: prometheus.NewDesc(
			"modem_downstream_channel_snr",
			"Downstream Channel Signal to Noise Ratio (dB)",
			[]string{"channel"},
			nil,
		),
		downstreamCorrectedMetric: prometheus.NewDesc(
			"modem_downstream_channel_corrected",
			"Downstream Channel corrected. (I don't know what this is)",
			[]string{"channel"},
			nil,
		),
		downstreamUncorrectablesMetric: prometheus.NewDesc(
			"modem_downstream_channel_uncorrectables",
			"Downstream Channel uncorrectables. (I don't know what this is)",
			[]string{"channel"},
			nil,
		),
		upstreamLockStatusMetric: prometheus.NewDesc(
			"modem_upstream_channel_locked",
			"Upstream channel Lock Status, 0 = Not Locked, 1 = Locked",
			[]string{"channel"},
			nil,
		),
		upstreamFrequencyMetric: prometheus.NewDesc(
			"modem_upstream_channel_frequency",
			"Upstream channel frequency.",
			[]string{"channel"},
			nil,
		),
		upstreamWidthMetric: prometheus.NewDesc(
			"modem_upstream_channel_width",
			"Upstream channel width.",
			[]string{"channel"},
			nil,
		),
		upstreamPowerMetric: prometheus.NewDesc(
			"modem_upstream_channel_power",
			"Upstream channel power (dBmv)",
			[]string{"channel"},
			nil,
		),
	}

	for _, option := range options {
		option(&c)
	}

	if c.logger == nil {
		c.logger = zap.NewNop()
	}

	if c.client == nil {
		modemClient, err := client.New()
		if err != nil {
			return nil, err
		}
		c.client = modemClient
	}
	return &c, nil
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.downstreamFrequencyMetric
	descs <- c.downstreamLockStatusMetric
	descs <- c.downstreamPowerMetric
	descs <- c.downstreamSnrMetric
	descs <- c.downstreamCorrectedMetric
	descs <- c.downstreamUncorrectablesMetric
	descs <- c.upstreamLockStatusMetric
	descs <- c.upstreamFrequencyMetric
	descs <- c.upstreamWidthMetric
	descs <- c.upstreamPowerMetric
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	downstreamChannels, upstreamChannels, err := c.client.Status()
	if err != nil {
		c.logger.Error("error getting status", zap.Error(err))
		return
	}

	for _, channel := range downstreamChannels {
		channelName := strconv.Itoa(channel.ChannelId)

		locked := 0
		if channel.LockStatus == "Locked" {
			locked = 1
		}

		metrics <- prometheus.MustNewConstMetric(c.downstreamFrequencyMetric, prometheus.GaugeValue, float64(channel.Frequency), channelName)
		metrics <- prometheus.MustNewConstMetric(c.downstreamLockStatusMetric, prometheus.GaugeValue, float64(locked), channelName)
		metrics <- prometheus.MustNewConstMetric(c.downstreamPowerMetric, prometheus.GaugeValue, channel.Power, channelName)
		metrics <- prometheus.MustNewConstMetric(c.downstreamSnrMetric, prometheus.GaugeValue, channel.SnrSmr, channelName)
		metrics <- prometheus.MustNewConstMetric(c.downstreamCorrectedMetric, prometheus.GaugeValue, float64(channel.Corrected), channelName)
		metrics <- prometheus.MustNewConstMetric(c.downstreamUncorrectablesMetric, prometheus.GaugeValue, float64(channel.Uncorrectables), channelName)

	}

	for _, channel := range upstreamChannels {
		channelName := strconv.Itoa(channel.ChannelId)

		locked := 0
		if channel.LockStatus == "Locked" {
			locked = 1
		}

		metrics <- prometheus.MustNewConstMetric(c.upstreamLockStatusMetric, prometheus.GaugeValue, float64(locked), channelName)
		metrics <- prometheus.MustNewConstMetric(c.upstreamFrequencyMetric, prometheus.GaugeValue, float64(channel.Frequency), channelName)
		metrics <- prometheus.MustNewConstMetric(c.upstreamWidthMetric, prometheus.GaugeValue, float64(channel.Width), channelName)
		metrics <- prometheus.MustNewConstMetric(c.upstreamPowerMetric, prometheus.GaugeValue, channel.Power, channelName)
	}
}
