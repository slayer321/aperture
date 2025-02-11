package enrichmentprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

type enrichmentProcessor struct {
	cache *entitycache.EntityCache
}

func newProcessor(cache *entitycache.EntityCache) *enrichmentProcessor {
	return &enrichmentProcessor{
		cache: cache,
	}
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (ep *enrichmentProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

// Start indicates and logs the start of the enrichment processor.
func (ep *enrichmentProcessor) Start(_ context.Context, _ component.Host) error {
	log.Debug().Msg("enrichment processor start")
	return nil
}

// Shutdown indicates and logs the shutdown of the enrichment processor.
func (ep *enrichmentProcessor) Shutdown(context.Context) error {
	log.Debug().Msg("enrichment processor shutdown")
	return nil
}

// ConsumeMetrics receives pmetric.Metrics for consumption then returns the enriched metrics.
func (ep *enrichmentProcessor) ConsumeMetrics(ctx context.Context, origMd pmetric.Metrics) (pmetric.Metrics, error) {
	if ep.cache == nil {
		return pmetric.Metrics{}, fmt.Errorf("cache not provided")
	}
	md := origMd.Clone()
	err := otelcollector.IterateMetrics(md, func(metric pmetric.Metric) error {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dataPoints := metric.Gauge().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeSum:
			dataPoints := metric.Sum().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeSummary:
			dataPoints := metric.Summary().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeHistogram:
			dataPoints := metric.Histogram().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeExponentialHistogram:
			dataPoints := metric.ExponentialHistogram().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		}
		return nil
	})
	return md, err
}

func (ep *enrichmentProcessor) enrichMetrics(attributes pcommon.Map) {
	hostNamex, ok := attributes.Get(otelcollector.EntityNameLabel)
	if !ok {
		return
	}
	hostName := hostNamex.StringVal()
	hostEntity, err := ep.cache.GetByName(hostName)
	attributes.Remove(otelcollector.EntityNameLabel)
	if err != nil {
		log.Trace().Str("name", hostName).Msg("Skipping because entity not found in cache")
		return
	}
	servicesValue := pcommon.NewValueSlice()
	for _, service := range hostEntity.Services {
		servicesValue.SliceVal().AppendEmpty().SetStringVal(service)
	}
	servicesValue.CopyTo(attributes.PutEmpty(otelcollector.ApertureServicesLabel))
}
