package reporters

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/platinummonkey/unifi/cmd/log"
	"go.uber.org/zap"
)

func init() {
	RegisterReporterBuilder("dogstatsd", NewDogStatsdWithConfig)
}

type DogStatsd struct {
	client *statsd.Client
	logger *zap.Logger
}

func NewDogStatsdWithConfig(cfg map[string]interface{}) (Reporter, error) {
	var client *statsd.Client
	var err error
	endpoint, ok := cfg["endpoint"]
	if ok && endpoint.(string) != "" {
		client, err = statsd.New(endpoint.(string))
	} else {
		// use default
		client, err = statsd.New("127.0.0.1:8125")
	}
	if err != nil {
		return nil, err
	}
	return &DogStatsd{client: client, logger: log.Get()}, nil
}

func (d *DogStatsd) ReportMetric(metricType MetricType, metric string, value float64, tags ...string) {
	var err error
	switch metricType {
	case CountMetricType:
		err = d.client.Count(metric, int64(value), tags, 1.0)
	case DistributionMetricType:
		err = d.client.Distribution(metric, value, tags, 1.0)
	case GaugeMetricType:
		err = d.client.Gauge(metric, value, tags, 1.0)
	case HistogramMetricType:
		err = d.client.Histogram(metric, value, tags, 1.0)
	case TimingMetricType:
		err = d.client.TimeInMilliseconds(metric, value, tags, 1.0)
	default:
		err = fmt.Errorf("ignoring unhandled type: %s", metricType)
	}

	if err != nil {
		d.logger.Warn("unable to post dogstatsd metrics, check configuration")
	}
}

func (d DogStatsd) ReportEvent(title string, message string, tags ...string) {
	err := d.client.Event(&statsd.Event{
		Title: title,
		Text:  message,
		Tags:  tags,
	})
	if err != nil {
		d.logger.Warn("unable to post dogstatsd event, check configuration")
	}
}
