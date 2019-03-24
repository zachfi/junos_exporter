package exporter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthStatus struct{}

func statusHandler(w http.ResponseWriter, r *http.Request) {

	health := &HealthStatus{}

	bytes, err := json.MarshalIndent(health, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func StartMetricsServer(bindAddr string) {
	d := http.NewServeMux()
	d.Handle("/metrics", promhttp.Handler())
	d.HandleFunc("/status/check", statusHandler)

	err := http.ListenAndServe(bindAddr, d)
	if err != nil {
		log.Fatal("Failed to start metrics server, error is:", err)
	}
}
