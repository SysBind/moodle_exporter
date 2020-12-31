
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Collector struct {
	client *sb.ServiceBusClient
	log    *logrus.Logger

	up *prometheus.Desc

	queueMessageCounts messageCountMetrics
	queueSizes         sizeMetrics

	topicMessageCounts messageCountMetrics
	topicSizes         sizeMetrics

	subscriptionMessageCounts messageCountMetrics
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	describeMessageCountMetrics(ch, &c.queueMessageCounts)
	describeSizeMetrics(ch, &c.queueSizes)

	describeMessageCountMetrics(ch, &c.topicMessageCounts)
	describeSizeMetrics(ch, &c.topicSizes)

	describeMessageCountMetrics(ch, &c.subscriptionMessageCounts)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape")

	if stats, err := c.client.GetServiceBusStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectQueueMetrics(c, ch, stats)
		collectTopicAndSubscriptionMetrics(c, ch, stats)

		c.log.Info("Scrape completed")
	}

}
