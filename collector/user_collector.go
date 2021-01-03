package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
)

type Collector struct {
	client    *client.Moodle
	log       *logrus.Logger
	up        *prometheus.Desc
	liveUsers *prometheus.Desc
}

func New(client *client.Moodle, log *logrus.Logger) *Collector {
	return &Collector{
		client: client,
		log:    log,

		up:        prometheus.NewDesc("moodle_up", "Whether the Moodle scrape was successful", nil, nil),
		liveUsers: prometheus.NewDesc("moodle_live_users", "Active users in last 5 minutes", nil, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.liveUsers
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape")

	if stats, err := c.client.GetUserStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
		ch <- prometheus.MustNewConstMetric(c.liveUsers, prometheus.GaugeValue, float64(stats.LiveUsers))
		c.log.Info("Scrape completed")
	}

}
