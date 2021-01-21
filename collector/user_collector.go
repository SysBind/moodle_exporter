package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
)

type UserCollector struct {
	Collector
	up                               *prometheus.Desc
	liveUsers                        *prometheus.Desc
	expectedUpcomingExamParticipants *prometheus.Desc
}

func NewUserCollector(client *client.Moodle, log *logrus.Logger) *UserCollector {
	return &UserCollector{
		Collector:                        Collector{client: client, log: log},
		up:                               prometheus.NewDesc("moodle_up", "Whether the Moodle scrape was successful", nil, nil),
		liveUsers:                        prometheus.NewDesc("moodle_live_users", "Active users in last 5 minutes", nil, nil),
		expectedUpcomingExamParticipants: prometheus.NewDesc("moodle_expected_upcoming_partipicants", "users which have activity scheduled (or not yet finished) in next 20 minutes and are not currently active", nil, nil),
	}
}

func (c *UserCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.liveUsers
	ch <- c.expectedUpcomingExamParticipants
}

func (c *UserCollector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape: User")

	if stats, err := c.client.GetUserStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during User scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
		ch <- prometheus.MustNewConstMetric(c.liveUsers, prometheus.GaugeValue, float64(stats.LiveUsers))
		ch <- prometheus.MustNewConstMetric(c.expectedUpcomingExamParticipants, prometheus.GaugeValue, float64(stats.ExpectedUpcomingExamParticipants))
		c.log.Info("Scrape completed: User")
	}
}
