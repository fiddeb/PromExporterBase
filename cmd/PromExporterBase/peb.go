
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/fiddeb/PromExporterBase/pkg/exporter"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"github.com/fiddeb/PromExporterBase/pkg/build"
)

const exportername = "PromExporterBase"

type promHTTPLogger struct {
	logger *slog.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	l.logger.Error(fmt.Sprint(v...))
}

func init() {
	prometheus.MustRegister(versioncollector.NewCollector(exportername))
}

func main() {
	var (
		webConfig   = webflag.AddFlags(kingpin.CommandLine, ":9119")
		metricsPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		opts        = exporter.ExporterOpts{}
	)
	version.Version = build.Version
	version.Branch = build.Branch

	promlogConfig := &promslog.Config{}

	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promslog.New(promlogConfig)

	logger.Info("Starting "+exportername, "version", version.Info())

	exporter, err := exporter.New(opts, logger)
	if err != nil {
		logger.Error("Error creating the exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{
						logger: logger,
					},
				},
			),
		),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>` + exportername + `</title></head>
             <body>
             <h1>` + exportername + `</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </dl>
             <h2>Build</h2>
             <pre>` + version.Info() + `</pre>
             </body>
             </html>`))
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		logger.Error("Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
