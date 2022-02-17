package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister()
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9400", nil); err != nil {
		panic(err)
	}
}
