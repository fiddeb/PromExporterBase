package exporter

import (
	//"fmt"

	"github.com/go-kit/log"
	//"github.com/go-kit/log/level"
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
  val int = 0
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


// Describe describes all the metrics ever exported by the Consul exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- rndgen
}

// Collect fetches the stats from configured Consul location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	ok := e.collectRandomGenMetric(ch)
	ok = e.collectTestCoubterMetric(ch) && ok



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

func (e *Exporter) collectTestCoubterMetric(ch chan<- prometheus.Metric) bool {
	value := float64(getCounterValue())
	ch <- prometheus.MustNewConstMetric(testCounter, prometheus.CounterValue, value , "asdf")
	ch <- prometheus.MustNewConstMetric(testCounter, prometheus.CounterValue, 1.0 , "asdf2")
	return true
}


func getCounterValue() int{
	counter := val
	if val >= 10 {
		val = 0
	} else {
		val ++
	}

	return counter
}