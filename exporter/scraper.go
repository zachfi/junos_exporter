package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	junos "github.com/scottdware/go-junos"
)

var (
	macAddress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mac",
		Help: "A individual MAC Adrress",
	}, []string{"ip", "day"})
)

func init() {
	prometheus.MustRegister(
		macAddress,
	)
}

func ScrapeMetrics(auth *junos.AuthMethod, hosts []string) {
	for _, h := range hosts {
		session, err := junos.NewSession(h, auth)
		defer session.Close()
		if err != nil {
			log.Error(err)
			continue
		}

		views, err := session.View("arp")
		if err != nil {
			log.Error(err)
			continue
		}

		for _, arp := range views.Arp.Entries {
			macAddress.WithLabelValues(arp.MACAddress, arp.IPAddress).Set(1)
		}
	}
}
