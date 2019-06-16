package reporters

import (
	"fmt"
	"time"

	"github.com/platinummonkey/unifi/cmd/log"
	datadog "github.com/zorkian/go-datadog-api"
	"go.uber.org/zap"
)

func init() {
	RegisterReporterBuilder("datadog", NewDataDogWithConfig)
}

type DataDog struct {
	client *datadog.Client
	logger *zap.Logger
}

func NewDataDogWithConfig(cfg map[string]interface{}) (Reporter, error) {
	apiKey, ok := cfg["api_key"]
	if !ok || apiKey.(string) == "" {
		return nil, fmt.Errorf("DataDog reporter requires a valid api_key")
	}
	client := datadog.NewClient(apiKey.(string), "")
	baseUrl, ok := cfg["baseurl"]
	if ok && baseUrl.(string) != "" {
		client.SetBaseUrl(baseUrl.(string))
	}
	return &DataDog{client: client, logger: log.Get()}, nil
}

func (d *DataDog) ReportMetric(metricType MetricType, metric string, value float64, tags ...string) {
	mType := string(metricType)
	now := float64(time.Now().UTC().Unix())
	err := d.client.PostMetrics([]datadog.Metric{
		{
			Metric: &metric,
			Type:   &mType,
			Points: []datadog.DataPoint{
				{
					&now,
					&value,
				},
			},
			Tags: tags,
		},
	})
	if err != nil {
		d.logger.Warn("unable to post datadog metrics, check configuration")
	}
}

func (d *DataDog) ReportEvent(title string, message string, tags ...string) {
	_, err := d.client.PostEvent(&datadog.Event{
		Title: &title,
		Text:  &message,
		Tags:  tags,
	})
	if err != nil {
		d.logger.Warn("unable to post datadog event, check configuration")
	}
}
