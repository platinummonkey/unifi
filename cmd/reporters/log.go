package reporters

import (
	"fmt"

	"github.com/gobuffalo/velvet"
	"github.com/platinummonkey/unifi/cmd/log"
	"go.uber.org/zap"
)

func init() {
	RegisterReporterBuilder("log", NewLogWithConfig)
}

// Log is a simple log-based reporter to STDOUT.
type Log struct {
	logger         *zap.Logger
	metricTemplate *velvet.Template
	eventTemplate  *velvet.Template
}

// default templates are json structured
const (
	defaultLogMetricTemplate = `{"type": {{ json type }}, "metric": {{ json metric }}, "value": {{ value }}, "tags": {{ json tags }} }`
	defaultLogEventTemplate  = `{"type": {{ json type }}, "title": {{ json title }}, "message": {{ json message }}, "tags": {{ json tags }} }`
)

// NewLogWithConfig will create a new log reporter with config
func NewLogWithConfig(cfg map[string]interface{}) (Reporter, error) {
	var metricTemplate *velvet.Template
	var err error
	if tpl, ok := cfg["metricTemplate"]; ok {
		if tplStr, ok := tpl.(string); !ok {
			return nil, fmt.Errorf("invalid mustache metric template for log reporter: must be a string")
		} else {
			metricTemplate, err = velvet.NewTemplate(tplStr)
			if err != nil {
				return nil, fmt.Errorf("invalid mustache metric template for log reporter: %s", err)
			}
		}
	} else {
		metricTemplate, _ = velvet.NewTemplate(defaultLogMetricTemplate)
	}

	var eventTemplate *velvet.Template
	if tpl, ok := cfg["eventTemplate"]; ok {
		if tplStr, ok := tpl.(string); !ok {
			return nil, fmt.Errorf("invalid mustache event template for log reporter: must be a string")
		} else {
			eventTemplate, err = velvet.NewTemplate(tplStr)
			if err != nil {
				return nil, fmt.Errorf("invalid mustache event template for log reporter: %s", err)
			}
		}
	} else {
		eventTemplate, _ = velvet.NewTemplate(defaultLogEventTemplate)
	}

	return &Log{logger: log.Get(), metricTemplate: metricTemplate, eventTemplate: eventTemplate}, nil
}

type handlebarFloat64 float64

func (f handlebarFloat64) Interface() interface{} {
	return float64(f)
}

// ReportMetric will log a new metric
func (r *Log) ReportMetric(metricType MetricType, metric string, value float64, tags ...string) {
	ctx := velvet.NewContext()
	ctx.Set("type", string(metricType))
	ctx.Set("metric", metric)
	ctx.Set("value", handlebarFloat64(value))
	ctx.Set("tags", tags)
	l, err := r.metricTemplate.Exec(ctx)
	if err == nil {
		r.logger.Info(l)
	}
}

// ReportEvent will log a new event
func (r *Log) ReportEvent(title string, message string, tags ...string) {
	ctx := velvet.NewContext()
	ctx.Set("type", "event")
	ctx.Set("title", title)
	ctx.Set("message", message)
	ctx.Set("tags", tags)
	l, err := r.eventTemplate.Exec(ctx)
	if err == nil {
		r.logger.Info(l)
	}
}
