package metrics

import "github.com/prometheus/client_golang/prometheus"

type FilterMetrics interface{}
type filterMetrics struct {
	registry *prometheus.Registry
}

func NewFilterMetrics(registry *prometheus.Registry) FilterMetrics {
	return &filterMetrics{registry: registry}
}

func (m *filterMetrics) DeviceInfo() {
}
