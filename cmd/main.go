package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/IBM/prometheus-metrics-server-simulator/pkg/metrics"
)

func main() {
	prometheus.NewCounter(prometheus.CounterOpts{})
	rand.Int()
	r := prometheus.NewRegistry()
	g, err := metrics.NewGenerator("/etc/conf/config.yaml", r)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

	}
	g.Start()

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.HandleFunc("/metrics/value", metrics.SetValue)
	http.ListenAndServe(":8080", nil)
}
