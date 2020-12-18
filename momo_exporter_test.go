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

func TestEmptyStats(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
 		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
  		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": []
	}`
	compare(t, resp, "empty_stats")
}

func TestOutboundRTP(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
  		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
  		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": [
			{
				"bytesSent": 5157622,
				"codecId": "RTCCodec_0_Outbound_102",
				"encoderImplementation": "Jetson Video Encoder",
				"firCount": 0,
				"frameHeight": 720,
				"frameWidth": 1280,
				"framesEncoded": 603,
				"framesPerSecond": 30,
				"framesSent": 603,
				"headerBytesSent": 120652,
				"hugeFramesSent": 1,
				"id": "RTCOutboundRTPVideoStream_2372247626",
				"isRemote": false,
				"keyFramesEncoded": 6,
				"kind": "video",
				"mediaSourceId": "RTCVideoSource_1",
				"mediaType": "video",
				"nackCount": 0,
				"packetsSent": 4788,
				"pliCount": 0,
				"qpSum": 11409,
				"qualityLimitationReason": "none",
				"qualityLimitationResolutionChanges": 0,
				"remoteId": "RTCRemoteInboundRtpVideoStream_2372247626",
				"retransmittedBytesSent": 0,
				"retransmittedPacketsSent": 0,
				"ssrc": 2372247626,
				"timestamp": 1608309189926189,
				"totalEncodeTime": 10.865,
				"totalEncodedBytesTarget": 0,
				"totalPacketSendDelay": 127.646,
				"trackId": "RTCMediaStreamTrack_sender_1",
				"transportId": "RTCTransport_0_1",
				"type": "outbound-rtp"
    		}
		]
	}`
	compare(t, resp, "outbound_rtp")
}

func TestDataChannel(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
  		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
  		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": [
			{
				"bytesReceived": 10,
				"bytesSent": 20,
				"dataChannelIdentifier": 1,
				"id": "RTCDataChannel_1",
				"label": "serial",
				"messagesReceived": 1,
				"messagesSent": 2,
				"protocol": "",
				"state": "open",
				"timestamp": 1608309189926189,
				"type": "data-channel"
			}
		]
	}`
	compare(t, resp, "data_channel")
}

func TestPeerConnection(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": [
			{
				"dataChannelsClosed": 0,
				"dataChannelsOpened": 1,
				"id": "RTCPeerConnection",
				"timestamp": 1608309189926189,
				"type": "peer-connection"
			}
		]
	}`
	compare(t, resp, "peer_connection")
}
