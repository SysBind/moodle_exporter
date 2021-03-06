package main

import (
	"net/http"
	"os"

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
	moodle, err := client.New(os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"))

	if err != nil {
		log.Fatal("Database Error, Exiting")
	}
	prometheus.MustRegister(collector.NewUserCollector(moodle, log))
	prometheus.MustRegister(collector.NewStorageCollector(moodle, log))
	log.Info("Starting moodle_exporter for ", moodle.Shortname)

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
