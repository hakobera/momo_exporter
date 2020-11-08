package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"sort"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/koron/go-dproxy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "momo"
)

// https://www.w3.org/TR/webrtc-stats/#rtctatstype-*
var (
	codecLabelNames = []string{"codec"}
/*
"codec",
"inbound-rtp",
"outbound-rtp",
"remote-inbound-rtp",
"remote-outbound-rtp",
"media-source",
"csrc",
"peer-connection",
"data-channel",
"stream",
"track",
"transceiver",
"sender",
"receiver",
"transport",
"sctp-transport",
"candidate-pair",
"local-candidate",
"remote-candidate",
"certificate",
"ice-server"
*/
)

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

func newCodecMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "codec", metricName),
			docString,
			codecLabelNames,
			constLabels,
		),
		Type: t,
	}
}

type metrics map[int]metricInfo

func (m metrics) String() string {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = strconv.Itoa(k)
	}
	return strings.Join(s, ",")
}

var (
	momoInfo = prometheus.NewDesc(prometheus.BuildFQName(namespace, "version", "info"), "WebRTC Native Client Momo version info.", []string{"version", "environment", "libwebrtc"}, nil)
	momoUp = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of WebRTC Native Client Momo successful.", nil, nil)
)

// Exporter collects momo stats from given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	URI       string
	mutex     sync.RWMutex
	fetchStat func() (io.ReadCloser, error)

	up            prometheus.Gauge
	totalScrapes  prometheus.Counter
	jsonParseFailures prometheus.Counter
	serverMetrics map[int]metricInfo
	logger        log.Logger
}

// NewExporter returns an intialized Exporter.
func NewExporter(uri string, sslVerify bool, timeout time.Duration, logger log.Logger) (*Exporter, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	var fetchStat func() (io.ReadCloser, error)
	switch u.Scheme {
	case "http", "https":
		fetchStat = fetchHTTP(uri, sslVerify, timeout)
	default:
		return nil, fmt.Errorf("unsupported scheme: %q", u.Scheme)
	}

	return &Exporter{
		URI:       uri,
		fetchStat: fetchStat,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of WebRTC Native Client Momo successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total momo scrapse.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_json_parse_failures_total",
			Help:      "Number of failures while parsing JSON.",
		}),
		logger: logger,
	}, nil
}

// Describe describes all the metrics ever exported by the Momo exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- momoInfo
	ch <- momoUp
	ch <- e.totalScrapes.Desc()
	ch <- e.jsonParseFailures.Desc()
}

// Collect fetches the stats from configured WebRTC Native Client Momo location
// and delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	up := e.scrape(ch)

	ch <- prometheus.MustNewConstMetric(momoUp, prometheus.GaugeValue, up)
	ch <- e.totalScrapes
	ch <- e.jsonParseFailures
}

func fetchHTTP(uri string, sslVerify bool, timeout time.Duration) func() (io.ReadCloser, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerify}}
	client := http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	return func() (io.ReadCloser, error) {
		resp, err := client.Get(uri)
		if err != nil {
			return nil, err
		}
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
		}
		return resp.Body, nil
	}
}

type MomoMetrics struct {
	Version string `json:"version"`
	Environment string `json:"environment"`
	Libwebrtc string `json:"libwebrtc"`
	Stats string `json:"stats"`
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) (up float64) {
	e.totalScrapes.Inc()

	body, err := e.fetchStat()
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't scrape WebRTC Native Client Momo", "err", err)
		return 0
	}
	defer body.Close()

	var metrics MomoMetrics
	err = json.NewDecoder(body).Decode(&metrics)
	if err != nil {
		level.Error(e.logger).Log("msg", "Failed to parse response from WebRTC Native Client Momo", "err", err)
		e.jsonParseFailures.Inc()	
		return 0
	}

	ch <- prometheus.MustNewConstMetric(momoInfo, prometheus.GaugeValue, 1, metrics.Version, metrics.Environment, metrics.Libwebrtc)

	level.Debug(e.logger).Log("msg", metrics.Stats)

	var v interface{}
	json.Unmarshal([]byte(metrics.Stats), &v)

	stats, err := dproxy.New(v).Array()
	if err != nil {
		level.Error(e.logger).Log("msg", "Failed to parse WebRTC stats", "err", err)
		e.jsonParseFailures.Inc()	
		return 0
	}

	for _, s := range stats {
		e.parseStats(s, ch)
	}

	return 1
}

func (e *Exporter) parseStats(stats interface{}, ch chan<- prometheus.Metric) {
	s := dproxy.New(stats)
	t, err := s.M("type").String()
	if err != nil {
		level.Error(e.logger).Log("msg", "stats must have 'type' field", "err", err)
		e.jsonParseFailures.Inc()
		return
	}
	level.Info(e.logger).Log("type", t)

	switch t {
	case "codec":
		e.exportCodecMetrics(s, ch)
	}
}

func (e *Exporter) exportCodecMetrics(m dproxy.Proxy, ch chan<- prometheus.Metric) {
	id, _ := m.M("id").String()
	desc:= prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "codec", id),
			"",
			[]string{"id", "payload_type", "mime_type", "clock_rate"},
			nil,
		)
	
	payloadTypeValue, _ := m.M("payloadType").Float64()
	payloadType, _ := m.M("payloadType").Int64()
	mimeType, _ := m.M("mimeType").String()
	clockRate, _ := m.M("clockRate").Int64()
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, payloadTypeValue, id, strconv.FormatInt(payloadType, 10), mimeType, strconv.FormatInt(clockRate, 10))
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9801").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		momoScrapeURI = kingpin.Flag("momo.scrape-uri", "URI on which to scrape WebRTC Native Client Momo.").Default("http://localhost:8081/metrics").String()
		momoSSLVerify = kingpin.Flag("momo.ssl-verify", "Flag that enables SSL certificate verification for the scrape URI.").Default("true").Bool()
		momoTimeout   = kingpin.Flag("momo.timeout", "Timeout for trying to get stats from WebRTC Native Client Momo.").Default("5s").Duration()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("momo_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting momo_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	exporter, err := NewExporter(*momoScrapeURI, *momoSSLVerify, *momoTimeout, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating an exorter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("momo_exporter"))

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Momo Exporter</title></head>
		<body>
		<h1>WebRTC Native Client Momo Exporter</h1>
		<p><a href=` + *metricsPath + `>Metrics</a></p>
		</body>
		</html>`))
	})
	if err = http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
