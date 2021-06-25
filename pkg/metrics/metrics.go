package metrics

import (
	"strconv"
	"strings"

	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
)

type FilterMetrics interface {
	FilterData(filterData data.FilterData)
	UserData(userData data.UserData)
	NetworkClient(st data.NetworkDevice)
	NetworkAccessPoint(ap data.AccessPoint)
}

type filterMetrics struct {
	rotationSpeedGauge *prometheus.GaugeVec
	dfsFactorGauge     *prometheus.GaugeVec
	dfsGauge           *prometheus.GaugeVec
	frequency          *prometheus.GaugeVec
	pumpMode           *prometheus.GaugeVec
	networkClient      *prometheus.GaugeVec
	accessPoint        *prometheus.GaugeVec
}

func (m *filterMetrics) UserData(userData data.UserData) {

}

func (m *filterMetrics) NetworkClient(st data.NetworkDevice) {
	var power int
	if st.StPower != "" {
		power = 1
	}
	m.networkClient.WithLabelValues(st.From, st.SSID, ip(st.IP), ip(st.Gateway)).Set(float64(power))
}

func ip(input []int) string {
	parts := make([]string, len(input))
	for i, p := range input {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, ".")
}

func (m *filterMetrics) NetworkAccessPoint(ap data.AccessPoint) {
	var power int
	if ap.Power != "" {
		power = 1
	}
	m.accessPoint.WithLabelValues(ap.From, ap.SSID).Set(float64(power))
}

func NewFilterMetrics(registry *prometheus.Registry) FilterMetrics {
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "pump_mode", Help: "The pump mode"}, []string{"name", "mode"})
	registry.MustRegister(g)
	networkClient := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "network_client", Help: "Network Client Information"},
		[]string{"name", "ssid", "ip", "gateway"},
	)
	registry.MustRegister(networkClient)
	accessPoint := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "access_point", Help: "Access Point Information",
	}, []string{"name", "ssid"})
	registry.MustRegister(accessPoint)
	return &filterMetrics{
		rotationSpeedGauge: newGauge("rotation_speed", "The rotation speed of the filter motor", registry),
		dfsGauge:           newGauge("dfs", "unknown", registry),
		dfsFactorGauge:     newGauge("dfs_factor", "unknown factor", registry),
		frequency:          newGauge("frequency", "motor frequency", registry),
		pumpMode:           g,
		networkClient:      networkClient,
		accessPoint:        accessPoint,
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

func (m *filterMetrics) FilterData(input data.FilterData) {
	name := input.From
	m.dfsGauge.WithLabelValues(name).Set(float64(input.DFS))
	m.dfsFactorGauge.WithLabelValues(name).Set(float64(input.DFSFactor))
	m.rotationSpeedGauge.WithLabelValues(name).Set(float64(input.RotationSpeed))
	m.frequency.WithLabelValues(name).Set(float64(input.Frequency))
	m.pumpMode.WithLabelValues(name, input.PumpMode.String()).Set(1)
}
