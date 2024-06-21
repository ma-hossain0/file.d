package metric

import (
	"strings"

	"github.com/ozontech/file.d/fd"
	"github.com/ozontech/file.d/metric"
	"github.com/ozontech/file.d/pipeline"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

/*{ introduction
Metric plugin.

**Example:**
```yaml
pipelines:
  example_pipeline:
    ...
    actions:
    - type: metric
	  metric_name: errors_total
	  metric_labels:
	  	- level

    ...
```

}*/

type Plugin struct {
	config Config

	logger *zap.Logger

	//  plugin metrics
	appliedMetric *prometheus.CounterVec
}

// ! config-params
// ^ config-params
type Config struct {
	// > @3@4@5@6
	// >
	// > The metric name.
	MetricName string `json:"metric_name" default:"total"` // *

	// > @3@4@5@6
	// >
	// > Lists the event fields to add to the metric. Blank list means no labels.
	// > Important note: labels metrics are not currently being cleared.
	MetricLabels []string `json:"metric_labels"` // *
}

func init() {
	fd.DefaultPluginRegistry.RegisterAction(&pipeline.PluginStaticInfo{
		Type:    "metric",
		Factory: factory,
	})
}

func factory() (pipeline.AnyPlugin, pipeline.AnyConfig) {
	return &Plugin{}, &Config{}
}

func (p *Plugin) makeMetric(ctl *metric.Ctl, name, help string, labels ...string) *prometheus.CounterVec {
	if name == "" {
		return nil
	}

	uniq := make(map[string]struct{})
	labelNames := make([]string, 0, len(labels))
	for _, label := range labels {
		if label == "" {
			p.logger.Fatal("empty label name")
		}
		if _, ok := uniq[label]; ok {
			p.logger.Fatal("metric labels must be unique")
		}
		uniq[label] = struct{}{}

		labelNames = append(labelNames, label)
	}

	return ctl.RegisterCounterVec(name, help, labelNames...)
}

func (p *Plugin) Start(config pipeline.AnyConfig, params *pipeline.ActionPluginParams) {
	p.config = *config.(*Config) // copy shared config
	p.logger = params.Logger.Desugar()
	p.registerMetrics(params.MetricCtl)
}

func (p *Plugin) registerMetrics(ctl *metric.Ctl) {
	p.appliedMetric = p.makeMetric(ctl,
		p.config.MetricName,
		"Number of events",
		p.config.MetricLabels...,
	)
}

func (p *Plugin) Stop() {
}

func (p *Plugin) Do(event *pipeline.Event) pipeline.ActionResult {
	if p.config.MetricName != "" {
		labelValues := make([]string, 0, len(p.config.MetricLabels))
		for _, labelValuePath := range p.config.MetricLabels {
			value := "not_set"
			if node := event.Root.Dig(labelValuePath); node != nil {
				value = strings.Clone(node.AsString())
			}

			labelValues = append(labelValues, value)
		}

		p.appliedMetric.WithLabelValues(labelValues...).Inc()
	}

	return pipeline.ActionPass
}
