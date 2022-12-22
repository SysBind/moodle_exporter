package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sysbind/moodle_exporter/client"
)

var (
	log = logrus.New()
)

func main() {
	_, err := client.NewMoodleList(os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"))

	if err != nil {
		log.Fatal("Database Error, Exiting")
	}
	// prometheus.MustRegister(collector.NewUserCollector(moodle, log))
	// prometheus.MustRegister(collector.NewStorageCollector(moodle, log))
	// log.Info("Starting moodle_exporter for ", moodle)

	// http.Handle("/metrics", promhttp.Handler())

	// // Home Page
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(`<html>
	// 		<head><title>Moodle Exporter</title></head>
	// 		<body>
	// 		<h1>Moodle Exporter</h1>
	// 		<p><a href="/metrics">Metrics</a></p>
	// 		</body>
	// 		</html>`))
	// })
	// http.ListenAndServe(":2345", nil);
}
