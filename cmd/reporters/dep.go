package reporters

import (
	"fmt"
)

// MetricType indicates a valid metric type
type MetricType string

// most common metric types
const (
	CountMetricType MetricType = "count"
	GaugeMetricType MetricType = "gauge"
	HistogramMetricType MetricType = "histogram"
	DistributionMetricType MetricType = "distribution"
	TimingMetricType MetricType = "timing"
)

// Reporter is an interface on which to report events and metrics
type Reporter interface {
	ReportMetric(metricType MetricType, metric string, value float64, tags ...string)
	ReportEvent(title string, message string, tags ...string)
}

// ReporterFromConfigBuilder defines the common builder interface
type ReporterFromConfigBuilder func(map[string]interface{}) (Reporter, error)

var registeredReporters = map[string]ReporterFromConfigBuilder{}

// RegisterReporterBuilder can be used in init() methods on reporters to dynamically register
func RegisterReporterBuilder(name string, builder ReporterFromConfigBuilder) {
	registeredReporters[name] = builder
}

// ReporterFromTypeAndConfig will build a reporter from a config
func ReporterFromTypeAndConfig(reporterType string, cfg map[string]interface{}) (Reporter, error) {
	builder, ok := registeredReporters[reporterType]
	if !ok {
		return nil, fmt.Errorf("invalid reporterType: %s", reporterType)
	}
	return builder(cfg)
}

// Reporters are many reporters convenience function
type Reporters map[string]Reporter

func (r Reporters) ReportMetric(metricType MetricType, metric string, value float64, tags ...string) {
	for _, reporter := range r {
		reporter.ReportMetric(metricType, metric, value, tags...)
	}
}

func (r Reporters) ReportEvent(title string, message string, tags ...string) {
	for _, reporter := range r {
		reporter.ReportEvent(title, message, tags...)
	}
}