package enrichmentprocessor

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/utils"
)

func TestEnrichmentprocessor(t *testing.T) {
	log.SetGlobalLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enrichmentprocessor Suite")
}

func assertAttributesEqual(act, exp pcommon.Map) {
	Expect(act.Len()).To(Equal(exp.Len()))
	Expect(act.Sort()).To(Equal(exp.Sort()))
}

func populateAttrsFromLabels(attr pcommon.Map, labels map[string]interface{}) {
	for k, v := range labels {
		// cast v to string or []string
		if str, ok := v.(string); ok {
			attr.InsertString(k, str)
		} else if slice, ok := v.([]string); ok {
			val := pcommon.NewValueSlice()
			sliceVal := val.SliceVal()
			for _, s := range slice {
				sliceVal.AppendEmpty().SetStringVal(s)
			}
			attr.Insert(k, val)
		}
	}
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})

var (
	hardCodedIPAddress  = "192.0.2.0"
	hardCodedEntityName = "test-entity"
	hardCodedServices   = []string{"svc1", "svc2"}
)
