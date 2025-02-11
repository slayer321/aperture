package utils

import (
	"runtime"
	"time"

	"go.uber.org/goleak"
)

var ignoreFuncs = []string{
	"github.com/onsi/ginkgo/v2/internal/interrupt_handler.(*InterruptHandler).registerForInterrupts.func2",
	"github.com/onsi/ginkgo/v2/internal.(*Suite).runNode",
	"net/http.(*persistConn).readLoop",
	"net/http.(*persistConn).writeLoop",
	"go.opencensus.io/stats/view.(*worker).start",
	"github.com/klauspost/compress/zstd.(*blockDec).startDecoder",
	"k8s.io/klog/v2.(*loggingT).flushDaemon",
	"sync.runtime_notifyListWait",
	"google.golang.org/grpc.(*ccBalancerWrapper).watcher",
	"google.golang.org/grpc.(*addrConn).resetTransport",
	"go.opentelemetry.io/collector/service/internal/telemetry.(*ProcessMetricsViews).StartCollection.func1",
	"time.Sleep",
}

// GoLeakDetector holds options for the goleak detector.
type GoLeakDetector struct {
	goleakOptions []goleak.Option
}

// NewGoLeakDetector creates a new GoLeakDetector with goleakOptions.
func NewGoLeakDetector() *GoLeakDetector {
	options := []goleak.Option{}
	for _, ignoreFunc := range ignoreFuncs {
		options = append(options, goleak.IgnoreTopFunction(ignoreFunc))
	}
	return &GoLeakDetector{
		goleakOptions: options,
	}
}

// AddIgnoreTopFunctions adds functions to ignore in the leak detector.
func (l *GoLeakDetector) AddIgnoreTopFunctions(fs ...string) {
	for _, f := range fs {
		l.goleakOptions = append(l.goleakOptions, goleak.IgnoreTopFunction(f))
	}
}

// FindLeaks finds memory leaks in the current process.
func (l *GoLeakDetector) FindLeaks() error {
	time.Sleep(time.Second * 5)
	runtime.GC()
	return goleak.Find(l.goleakOptions...)
}
