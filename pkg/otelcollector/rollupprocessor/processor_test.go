package rollupprocessor_test

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/fluxninja/datasketches-go/sketches"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/pkg/otelcollector"
	. "github.com/fluxninja/aperture/pkg/otelcollector/rollupprocessor"
)

var _ = Describe("Rollup processor", func() {
	var (
		config       *Config
		testConsumer *fakeConsumer
	)

	BeforeEach(func() {
		config = &Config{}
		testConsumer = &fakeConsumer{
			receivedLogs:    []plog.Logs{},
			receivedMetrics: []pmetric.Metrics{},
			receivedTraces:  []ptrace.Traces{},
		}
	})

	Describe("Logs", func() {
		var logsProcessor component.LogsProcessor

		JustBeforeEach(func() {
			var err error
			logsProcessor, err = CreateLogsProcessor(
				context.TODO(), component.ProcessorCreateSettings{}, config, testConsumer)
			Expect(err).NotTo(HaveOccurred())
		})

		It("works for single log record", func() {
			attributeValues := []int{5}
			expectedSerializedDatasketch, err := serializedDatasketchFromAttributeValues(attributeValues)
			Expect(err).NotTo(HaveOccurred())

			input := plog.NewLogs()
			logs := input.ResourceLogs().AppendEmpty().
				ScopeLogs().AppendEmpty().
				LogRecords()
			logRecord := logs.AppendEmpty()
			logRecord.Attributes().InsertString("fizz", "buzz")
			logRecord.Attributes().InsertString(otelcollector.WorkloadDurationLabel, strconv.Itoa(attributeValues[0]))

			err = logsProcessor.ConsumeLogs(context.TODO(), input)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConsumer.receivedLogs).To(HaveLen(1))
			attributes := testConsumer.receivedLogs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes().AsRaw()
			Expect(attributes).To(HaveLen(7))
			Expect(attributes).To(HaveKeyWithValue(RollupCountKey, int64(1)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupDatasketch), expectedSerializedDatasketch))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupSum), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupMin), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupMax), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupSumOfSquares), float64(25)))
			Expect(attributes).To(HaveKeyWithValue("fizz", "buzz"))
		})

		It("works for multiple log records", func() {
			attributeValues := []int{5, 6, 7}
			expectedSerializedDatasketch, err := serializedDatasketchFromAttributeValues(attributeValues)
			Expect(err).NotTo(HaveOccurred())

			input := plog.NewLogs()
			logs := input.ResourceLogs().AppendEmpty().
				ScopeLogs().AppendEmpty().
				LogRecords()
			logRecord := logs.AppendEmpty()
			logRecord.Attributes().InsertString(otelcollector.WorkloadDurationLabel, strconv.Itoa(attributeValues[0]))
			logRecord = logs.AppendEmpty()
			logRecord.Attributes().InsertString(otelcollector.WorkloadDurationLabel, strconv.Itoa(attributeValues[1]))
			logRecord = logs.AppendEmpty()
			logRecord.Attributes().InsertString(otelcollector.WorkloadDurationLabel, strconv.Itoa(attributeValues[2]))

			err = logsProcessor.ConsumeLogs(context.TODO(), input)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConsumer.receivedLogs).To(HaveLen(1))
			attributes := testConsumer.receivedLogs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes().AsRaw()
			Expect(attributes).To(HaveLen(6))
			Expect(attributes).To(HaveKeyWithValue(RollupCountKey, int64(3)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupDatasketch), expectedSerializedDatasketch))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupSum), float64(18)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupMin), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupMax), float64(7)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelcollector.WorkloadDurationLabel, RollupSumOfSquares), float64(110)))
		})
	})
})

type fakeConsumer struct {
	receivedLogs    []plog.Logs
	receivedMetrics []pmetric.Metrics
	receivedTraces  []ptrace.Traces
}

func (fc *fakeConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (fc *fakeConsumer) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	fc.receivedLogs = append(fc.receivedLogs, ld)
	return nil
}

func (fc *fakeConsumer) ConsumeMetrics(_ context.Context, ld pmetric.Metrics) error {
	fc.receivedMetrics = append(fc.receivedMetrics, ld)
	return nil
}

func (fc *fakeConsumer) ConsumeTraces(_ context.Context, ld ptrace.Traces) error {
	fc.receivedTraces = append(fc.receivedTraces, ld)
	return nil
}

func serializedDatasketchFromAttributeValues(values []int) (string, error) {
	sketch, err := sketches.NewDoublesSketch(128)
	if err != nil {
		return "", err
	}
	for _, v := range values {
		err = sketch.Update(float64(v))
		if err != nil {
			return "", err
		}
	}
	sketchBytes, _ := sketch.Compact().Serialize()
	return base64.StdEncoding.EncodeToString(sketchBytes), nil
}
