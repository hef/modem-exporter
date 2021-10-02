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

	frequencyMetric      *prometheus.Desc
	lockStatusMetric     *prometheus.Desc
	powerMetric          *prometheus.Desc
	snrMetric            *prometheus.Desc
	correctedMetric      *prometheus.Desc
	uncorrectablesMetric *prometheus.Desc
}

func New(options ...Option) (*Collector, error) {
	c := Collector{
		frequencyMetric: prometheus.NewDesc(
			"modem_downstream_channel_frequency",
			"Downstream channel frequency.",
			[]string{"channel"},
			nil,
		),
		lockStatusMetric: prometheus.NewDesc(
			"modem_downstream_channel_locked",
			"Downstream channel Lock Status, 0 = Not Locked, 1 = Locked",
			[]string{"channel"},
			nil,
		),
		powerMetric: prometheus.NewDesc(
			"modem_downstream_channel_power",
			"Downstream channel power (dBmv)",
			[]string{"channel"},
			nil,
		),
		snrMetric: prometheus.NewDesc(
			"modem_downstream_channel_snr",
			"Downstream Channel Signal to Noise Ratio (dB)",
			[]string{"channel"},
			nil,
		),
		correctedMetric: prometheus.NewDesc(
			"modem_downstream_channel_corrected",
			"Downstream Channel corrected. (I don't know what this is)",
			[]string{"channel"},
			nil,
		),
		uncorrectablesMetric: prometheus.NewDesc(
			"modem_downstream_channel_uncorrectables",
			"Downstream Channel uncorrectables. (I don't know what this is)",
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
	descs <- c.frequencyMetric
	descs <- c.lockStatusMetric
	descs <- c.powerMetric
	descs <- c.snrMetric
	descs <- c.correctedMetric
	descs <- c.uncorrectablesMetric
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	channels, err := c.client.Status()
	if err != nil {
		c.logger.Error("error getting status", zap.Error(err))
		return
	}

	for _, channel := range channels {
		channelName := strconv.Itoa(channel.ChannelId)

		locked := 0
		if channel.LockStatus == "Locked" {
			locked = 1
		}

		metrics <- prometheus.MustNewConstMetric(c.frequencyMetric, prometheus.GaugeValue, float64(channel.Frequency), channelName)
		metrics <- prometheus.MustNewConstMetric(c.lockStatusMetric, prometheus.GaugeValue, float64(locked), channelName)
		metrics <- prometheus.MustNewConstMetric(c.powerMetric, prometheus.GaugeValue, channel.Power, channelName)
		metrics <- prometheus.MustNewConstMetric(c.snrMetric, prometheus.GaugeValue, channel.SnrSmr, channelName)
		metrics <- prometheus.MustNewConstMetric(c.correctedMetric, prometheus.GaugeValue, float64(channel.Corrected), channelName)
		metrics <- prometheus.MustNewConstMetric(c.uncorrectablesMetric, prometheus.GaugeValue, float64(channel.Uncorrectables), channelName)
	}
}
