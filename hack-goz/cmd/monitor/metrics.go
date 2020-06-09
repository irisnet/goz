package monitor

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace  = "hack"
	Subsystem  = "robot"
	TimeFormat = "2006-01-02 15:04:05.000Z"
)

type Metrics struct {
	Success *prom.GaugeVec
}

func PrometheusMetrics() *Metrics {
	metrics := &Metrics{
		Success: prom.NewGaugeVec(prom.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "success",
			Help:      "packet sent successfully",
		}, []string{}),
	}

	prom.MustRegister(metrics.Success)

	return metrics
}
