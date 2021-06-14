package metrics

import (
	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
)

type FilterMetrics interface {
	RotationSpeed(string, int)
	DFS(string, int)
	DFSFactor(string, int)
	Frequency(string, int)
	PumpMode(string, data.PumpMode)
}

type filterMetrics struct {
	rotationSpeedGauge *prometheus.GaugeVec
	dfsFactorGauge     *prometheus.GaugeVec
	dfsGauge           *prometheus.GaugeVec
	frequency          *prometheus.GaugeVec
	pumpMode           *prometheus.GaugeVec
}

func NewFilterMetrics(registry *prometheus.Registry) FilterMetrics {
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "pump_mode", Help: "The pump mode"}, []string{"name", "mode"})
	registry.MustRegister(g)
	return &filterMetrics{
		rotationSpeedGauge: newGauge("rotation_speed", "The rotation speed of the filter motor", registry),
		dfsGauge:           newGauge("dfs", "unknown", registry),
		dfsFactorGauge:     newGauge("dfs_factor", "unknown factor", registry),
		frequency:          newGauge("frequency", "motor frequency", registry),
		pumpMode:           g,
	}
}

func newGauge(name string, help string, registry *prometheus.Registry) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, []string{"filter"})
	registry.MustRegister(gauge)
	return gauge
}

func (m *filterMetrics) DFS(name string, dfs int) {
	m.dfsGauge.WithLabelValues(name).Set(float64(dfs))
}

func (m *filterMetrics) DFSFactor(name string, dfsFactor int) {
	m.dfsFactorGauge.WithLabelValues(name).Set(float64(dfsFactor))
}

func (m *filterMetrics) RotationSpeed(name string, speed int) {
	m.rotationSpeedGauge.WithLabelValues(name).Set(float64(speed))
}

func (m *filterMetrics) Frequency(name string, frequency int) {
	m.frequency.WithLabelValues(name).Set(float64(frequency))
}

func (m *filterMetrics) PumpMode(name string, mode data.PumpMode) {
	pumpModeVal := "unknown"
	switch mode {
	case data.ConstantFlowMode:
		pumpModeVal = "constant_flow"
	case data.BioMode:
		pumpModeVal = "bio"
	case data.PulseFlowMode:
		pumpModeVal = "pulse"
	case data.ManualMode:
		pumpModeVal = "manual"
	}
	m.pumpMode.WithLabelValues(name, pumpModeVal).Set(1)
}
