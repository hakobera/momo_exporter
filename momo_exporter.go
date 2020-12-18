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

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/iancoleman/strcase"
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

// MomoMetrics is metrics respose type of WebRTC Native Client Momo
type MomoMetrics struct {
	Version     string      `json:"version"`
	Environment string      `json:"environment"`
	Libwebrtc   string      `json:"libwebrtc"`
	Stats       interface{} `json:"stats"`
}

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

var (
	momoInfo = prometheus.NewDesc(prometheus.BuildFQName(namespace, "version", "info"), "WebRTC Native Client Momo version info.", []string{"version", "environment", "libwebrtc"}, nil)
	momoUp   = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of WebRTC Native Client Momo successful.", nil, nil)
)

// Exporter collects momo stats from given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	URI       string
	mutex     sync.RWMutex
	fetchStat func() (io.ReadCloser, error)

	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter
	serverMetrics     map[int]metricInfo
	logger            log.Logger
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
			Help:      "Current total momo scrapes.",
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
	for _, m := range outboundRTPMetrics {
		ch <- m.Desc
	}
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

	stats, err := dproxy.New(metrics.Stats).Array()
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
	level.Debug(e.logger).Log("msg", "Metrics type", "type", t)

	// https://www.w3.org/TR/webrtc-stats/#summary
	switch t {
	case "data-channel":
		e.exportDataChannelMetrics(s, ch)
	case "outbound-rtp":
		e.exportOutboundRTPMetrics(s, ch)
	case "peer-connection":
		e.exportPeerConnectionMetrics(s, ch)
	case "transport":
		e.exportTransportMetrics(s, ch)
	}
}

func (e *Exporter) exportDataChannelMetrics(m dproxy.Proxy, ch chan<- prometheus.Metric) {
	id, _ := m.M("id").String()
	label, _ := m.M("label").String()

	for key, metric := range dataChannelMetrics {
		val, _ := m.M(strcase.ToLowerCamel(key)).Float64()
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, val, id, label)
	}
}

func (e *Exporter) exportOutboundRTPMetrics(m dproxy.Proxy, ch chan<- prometheus.Metric) {
	id, _ := m.M("id").String()
	codecID, _ := m.M("codecId").String()
	encoderImplementation, _ := m.M("encoderImplementation").String()
	kind, _ := m.M("kind").String()
	mediaSourceID, _ := m.M("mediaSourceId").String()

	for key, metric := range outboundRTPMetrics {
		val, _ := m.M(strcase.ToLowerCamel(key)).Float64()
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, val, id, codecID, encoderImplementation, kind, mediaSourceID)
	}
}

func (e *Exporter) exportPeerConnectionMetrics(m dproxy.Proxy, ch chan<- prometheus.Metric) {
	id, _ := m.M("id").String()

	for key, metric := range peerConnectionMetrics {
		val, _ := m.M(strcase.ToLowerCamel(key)).Float64()
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, val, id)
	}
}

func (e *Exporter) exportTransportMetrics(m dproxy.Proxy, ch chan<- prometheus.Metric) {
	id, _ := m.M("id").String()

	for key, metric := range transportMetrics {
		val, _ := m.M(strcase.ToLowerCamel(key)).Float64()
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, val, id)
	}
}

type metrics map[string]metricInfo

var (
	// https://www.w3.org/TR/webrtc-stats/#dom-rtcdatachannelstats
	dataChannelLabelNames = []string{"id", "label"}
	dataChannelMetrics    = metrics{
		"bytesSent":        newDataChannelMetric("bytes_sent_total", "Total number of payload bytes sent on this RTCDataChannel", prometheus.CounterValue, nil),
		"bytesReceived":    newDataChannelMetric("bytes_received_total", "Total number of payload bytes sent on this RTCDataChannel.", prometheus.CounterValue, nil),
		"messagesSent":     newDataChannelMetric("messages_sent_total", "Total number of API \"message\" events sent.", prometheus.CounterValue, nil),
		"messagesReceived": newDataChannelMetric("messages_received_total", "Total number of API \"message\" events received.", prometheus.CounterValue, nil),
	}

	// https://www.w3.org/TR/webrtc-stats/#dom-rtcoutboundrtpstreamstats
	outboundRTPLabelNames = []string{"id", "codecId", "encoderImplementation", "kind", "mediaSourceId"}
	outboundRTPMetrics    = metrics{
		"bytesSent":                          newOutboundRTPMetric("bytes_sent_total", "Total number of bytes sent for this SSRC.", prometheus.CounterValue, nil),
		"headerBytesSent":                    newOutboundRTPMetric("header_bytes_sent_total", "Total number of RTP header and padding bytes sent for this SSRC.", prometheus.CounterValue, nil),
		"retransmittedBytesSent":             newOutboundRTPMetric("retransmitted_bytes_sent_total", "Total number of bytes that were retransmitted for this SSRC.", prometheus.CounterValue, nil),
		"packetsSent":                        newOutboundRTPMetric("packets_sent_total", "Total number of RTP packets sent for this SSRC.", prometheus.CounterValue, nil),
		"retransmittedPacketsSent":           newOutboundRTPMetric("retransmitted_packets_sent_total", "Total number of RTP packets sent for this SSRC.", prometheus.CounterValue, nil),
		"framesSent":                         newOutboundRTPMetric("frames_sent_total", "Total number of frames sent on this RTP stream.", prometheus.CounterValue, nil),
		"firCount":                           newOutboundRTPMetric("fir_count_total", "Total number of Full Intra Request (FIR) packets received by this sender.", prometheus.CounterValue, nil),
		"pliCount":                           newOutboundRTPMetric("pli_count_total", "Total number of Picture Loss Indication (PLI) packets received by this sender.", prometheus.CounterValue, nil),
		"sliCount":                           newOutboundRTPMetric("sli_count_total", "Total number of Slice Loss Indication (SLI) packets received by this sender.", prometheus.CounterValue, nil),
		"nackCount":                          newOutboundRTPMetric("nack_count_total", "Total number of Negative ACKnowledgement (NACK) packets received by this sender.", prometheus.CounterValue, nil),
		"qpSum":                              newOutboundRTPMetric("qp_sum", "Sum of the QP values of frames encoded by this sender.", prometheus.CounterValue, nil),
		"framesEncoded":                      newOutboundRTPMetric("frames_encoded_total", "Total number of frames successfully encoded for this RTP media stream.", prometheus.CounterValue, nil),
		"keyFramesEncoded":                   newOutboundRTPMetric("key_frames_encoded_total", "Total number of key frames successfully encoded for this RTP media stream.", prometheus.CounterValue, nil),
		"totalEncodeTime":                    newOutboundRTPMetric("encode_time_total", "Total number of seconds that has been spent encoding the framesEncoded frames of this stream.", prometheus.CounterValue, nil),
		"frameWidth":                         newOutboundRTPMetric("frame_width", "Width of the last encoded frame.", prometheus.GaugeValue, nil),
		"frameHeight":                        newOutboundRTPMetric("frame_height", "Height of the last encoded frame.", prometheus.GaugeValue, nil),
		"framesPerSecond":                    newOutboundRTPMetric("frames_per_second", "Number of encoded frames during the last second.", prometheus.GaugeValue, nil),
		"totalPacketSendDelay":               newOutboundRTPMetric("packet_send_delay_total", "Total number of seconds that packets have spent buffered locally before being transmitted onto the network.", prometheus.CounterValue, nil),
		"totalSamplesSent":                   newOutboundRTPMetric("samples_sent_total", "Total number of samples that have been sent over this RTP stream.", prometheus.CounterValue, nil),
		"qualityLimitationResolutionChanges": newOutboundRTPMetric("quality_limitation_resolution_changes_total", "Number of times that the resolution has changed because we are quality limited (qualityLimitationReason has a value other than \"none\").", prometheus.CounterValue, nil),
	}

	// https://www.w3.org/TR/webrtc-stats/#dom-rtcpeerconnectionstats
	peerConnectionLabelNames = []string{"id"}
	peerConnectionMetrics    = metrics{
		"dataChannelsOpened": newPeerConnectionMetric("data_channels_opened_total", "Number of unique RTCDataChannels that have entered the \"open\" state during their lifetime.", prometheus.CounterValue, nil),
		"dataChannelsClosed": newPeerConnectionMetric("data_chennels_closed_total", "Number of unique RTCDataChannels that have left the \"open\" state during their lifetime.", prometheus.CounterValue, nil),
	}

	// https://www.w3.org/TR/webrtc-stats/#transportstats-dict*
	transportLabelNames = []string{"id"}
	transportMetrics    = metrics{
		"bytesSent":                    newTransportMetric("bytes_sent_total", "Total number of payload bytes sent on this RTCIceTransport.", prometheus.CounterValue, nil),
		"bytesReceived":                newTransportMetric("bytes_received_total", "Total number of payload bytes received on this RTCIceTransport.", prometheus.CounterValue, nil),
		"packetsSent":                  newTransportMetric("packets_sent_total", "Total number of packets sent over this transport.", prometheus.CounterValue, nil),
		"packetsReceived":              newTransportMetric("packets_received_total", "Total number of packets received on this transport.", prometheus.CounterValue, nil),
		"selectedCandidatePairChanges": newTransportMetric("selected_candidate_pair_changes_total", "Number of times that the selected candidate pair of this transport has changed.", prometheus.CounterValue, nil),
	}
)

func newMetric(category string, metricName string, docString string, t prometheus.ValueType, variableLabels []string, constLabels prometheus.Labels) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, category, metricName),
			docString,
			variableLabels,
			constLabels,
		),
		Type: t,
	}
}

func newDataChannelMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return newMetric("datachannel", metricName, docString, t, dataChannelLabelNames, constLabels)
}

func newOutboundRTPMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return newMetric("outbound_rtp", metricName, docString, t, outboundRTPLabelNames, constLabels)
}

func newPeerConnectionMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return newMetric("peerconnection", metricName, docString, t, peerConnectionLabelNames, constLabels)
}

func newTransportMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return newMetric("transport", metricName, docString, t, transportLabelNames, constLabels)
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
