package collector

import (
	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
)

type Collector struct {
	client *client.Moodle
	log    *logrus.Logger
}
