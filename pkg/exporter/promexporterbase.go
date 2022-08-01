package exporter

import (
	"github.com/go-kit/log"
	"math/rand"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "asdf"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of exporter successful.",
		nil, nil,
	)
	rndgen = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "randomgauge"),
		"A random number as gaue for test.",
		nil, nil,
	)
	testCounter = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "testcounter_total"),
		"a test counter",
		[]string{"keyname"}, nil,
	)

)

var (
  val1 int = 0
  val2 int = 0
)

// ExporterOpts configures options for connecting to target.
type ExporterOpts struct {
	RequestLimit int
}

type Exporter struct {

	logger           log.Logger
	requestLimitChan chan struct{}
}

// New returns an initialized Exporter.
func New(opts ExporterOpts, logger log.Logger) (*Exporter, error) {

	var requestLimitChan chan struct{}
	if opts.RequestLimit > 0 {
		requestLimitChan = make(chan struct{}, opts.RequestLimit)
	}

	// Init our exporter.
	return &Exporter{
		logger:           logger,
		requestLimitChan: requestLimitChan,
	}, nil
}


// Describe describes all the metrics ever exported by the exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- rndgen
	ch <- testCounter
}

// Collect fetches the stats and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	ok := e.collectRandomGenMetric(ch)
	ok = e.collectTestCounterMetric(ch) && ok



	if ok {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}
}


func (e *Exporter) collectRandomGenMetric(ch chan<- prometheus.Metric) bool {
	var value float64
	rnd := rand.Intn(100)
	value =  float64(rnd)
	ch <- prometheus.MustNewConstMetric(
		rndgen, prometheus.GaugeValue, value,
	)
	return true
}

func (e *Exporter) collectTestCounterMetric(ch chan<- prometheus.Metric) bool {
	value1, value2 := getCounterValue()
	ch <- prometheus.MustNewConstMetric(testCounter, prometheus.CounterValue,  float64(value1) , "counter1")
	ch <- prometheus.MustNewConstMetric(testCounter, prometheus.CounterValue,  float64(value2) , "counter2")
	return true
}


func getCounterValue() (int, int){
	counter1 := val1
	counter2 := val2
	if val1 >= 10 {
		val1 = 0
	} else {
		val1 ++
	}
	val2 ++
	return counter1, counter2
}