package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
)

type StorageCollector struct {
	Collector
	bytesAssignSubmission *prometheus.Desc
	bytesBackup           *prometheus.Desc
	bytesBackupAuto       *prometheus.Desc
	bytesAll              *prometheus.Desc
}

func NewStorageCollector(client *client.MoodleList, log *logrus.Logger) *StorageCollector {
	return &StorageCollector{
		Collector: Collector{client: client, log: log},

		bytesAssignSubmission: prometheus.NewDesc("moodle_bytes_assign_submission",
			"bytes used by assign submissions", []string{"moodle", "course"}, nil),

		bytesBackup: prometheus.NewDesc("moodle_bytes_backup",
			"bytes used by backups excluding automatic", []string{"moodle", "course"}, nil),
		bytesBackupAuto: prometheus.NewDesc("moodle_bytes_backup_auto",
			"bytes used by automatic backups", []string{"moodle", "course"}, nil),
		bytesAll: prometheus.NewDesc("moodle_bytes_all",
			"bytes used by all files", []string{"moodle"}, nil),
	}
}

func (c *StorageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bytesAssignSubmission
	ch <- c.bytesBackup
	ch <- c.bytesBackupAuto
	ch <- c.bytesAll
}

func (c *StorageCollector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape: Storage")

	var statsList []*client.StorageStats
	var err error
	if statsList, err = c.client.GetStorageStats();  err != nil {
		c.log.WithError(err).Error("Error during User scrape")
		return
	}

	for _, stats := range statsList {
		for course, bytes := range stats.BytesAssignSubmission {
			ch <- prometheus.MustNewConstMetric(c.bytesAssignSubmission,
				prometheus.GaugeValue,
				float64(bytes),
				fmt.Sprintf("%s", stats.MoodleShortName()),
				fmt.Sprintf("%d", course))
		}

		for course, bytes := range stats.BytesBackup {
			ch <- prometheus.MustNewConstMetric(c.bytesBackup,
				prometheus.GaugeValue,
				float64(bytes),
				fmt.Sprintf("%s", stats.MoodleShortName()),
				fmt.Sprintf("%d", course))
		}

		for course, bytes := range stats.BytesBackupAuto {
			ch <- prometheus.MustNewConstMetric(c.bytesBackupAuto,
				prometheus.GaugeValue,
				float64(bytes),
				fmt.Sprintf("%s", stats.MoodleShortName()),
				fmt.Sprintf("%d", course))
		}

		ch <- prometheus.MustNewConstMetric(c.bytesAll,
			prometheus.GaugeValue,
			float64(stats.BytesAll),
			fmt.Sprintf("%s", stats.MoodleShortName()))
	}
	c.log.Info("Scrape completed: Storage")
}
