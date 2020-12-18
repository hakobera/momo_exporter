package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const (
	testVersion     = "WebRTC Native Client Momo 20XX.Y (test)"
	testEnvironment = "Test Environment"
	testLibwebrtc   = "Test-Build MXX.YYYY@{#Z} (XX.YYYY.Z test)"
)

type momo struct {
	*httptest.Server
	response []byte
}

func newMomo(response []byte) *momo {
	h := &momo{response: response}
	h.Server = httptest.NewServer(handler(h))
	return h
}

func handler(h *momo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(h.response)
	}
}

func expectMetrics(t *testing.T, c prometheus.Collector, fixture string) {
	exp, err := os.Open(path.Join("test", fixture) + ".metrics")
	if err != nil {
		t.Fatalf("Error opening fixture file %q: %v", fixture, err)
	}
	if err := testutil.CollectAndCompare(c, exp); err != nil {
		t.Fatal("Unexpected metrics returned:", err)
	}
}

func compare(t *testing.T, response string, fixture string) {
	h := newMomo([]byte(response))
	defer h.Close()
	e, _ := NewExporter(h.URL, true, 5*time.Second, log.NewNopLogger())
	expectMetrics(t, e, fixture)
}

func TestInvalidFormat(t *testing.T) {
	compare(t, "{", "invalid_format")
}

func TestSuccessEmptyStats(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
 		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
  		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": []
	}`
	compare(t, resp, "success_empty_stats")
}
