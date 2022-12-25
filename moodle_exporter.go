package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
	"github.com/sysbind/moodle_exporter/collector"
)

var (
	log = logrus.New()
)

func main() {
	moodles, err := client.NewMoodleList(os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"))

	if err != nil {
		if os.Getenv("DEBUG") == "1" {
			log.Warning("Database Error, DEBUG=1, sleeping for 1 hour..")
			time.Sleep(time.Hour)
		}
		log.Fatal("Database Error, Exiting")
	}

	prometheus.MustRegister(collector.NewUserCollector(&moodles, log))
	prometheus.MustRegister(collector.NewStorageCollector(&moodles, log))
	log.Info("Starting moodle_exporter for ", moodles)

	http.Handle("/metrics", promhttp.Handler())

	// Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Moodle Exporter</title></head>
			<body>
			<h1>Moodle Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	http.ListenAndServe(":2345", nil)
}
