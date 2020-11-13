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
	exp, err := os.Open(path.Join("test", fixture))
	if err != nil {
		t.Fatalf("Error opening fixture file %q: %v", fixture, err)
	}
	if err := testutil.CollectAndCompare(c, exp); err != nil {
		t.Fatal("Unexpected metrics returned:", err)
	}
}

func TestInvalidFormat(t *testing.T) {
	h := newMomo([]byte("{"))
	defer h.Close()

	e, _ := NewExporter(h.URL, true, 5*time.Second, log.NewNopLogger())

	expectMetrics(t, e, "invalid_format.metrics")
}
