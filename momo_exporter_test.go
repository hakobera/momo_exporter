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

func handlerStale(exit chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		<-exit
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
	e, err := NewExporter(h.URL, true, 5*time.Second, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
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

func TestInboundRTP(t *testing.T) {
	resp := `{
		"environment": "[x86_64] macOS Version 10.15.7 (Build 19H15)",
		"libwebrtc": "Shiguredo-Build M88.4324@{#3} (88.4324.3.0 b15b2915)",
		"stats": [
			{
				"bytesReceived": 10278549,
				"codecId": "RTCCodec_video_qDqHgY_Inbound_120",
				"decoderImplementation": "libvpx",
				"firCount": 0,
				"frameHeight": 720,
				"frameWidth": 1280,
				"framesDecoded": 2111,
				"framesDropped": 0,
				"framesPerSecond": 14,
				"framesReceived": 2112,
				"headerBytesReceived": 156448,
				"id": "RTCInboundRTPVideoStream_2189915641",
				"isRemote": false,
				"keyFramesDecoded": 1,
				"kind": "video",
				"lastPacketReceivedTimestamp": 4270.359,
				"mediaType": "video",
				"nackCount": 67,
				"packetsLost": 0,
				"packetsReceived": 9778,
				"pliCount": 0,
				"qpSum": 291917,
				"ssrc": 2189915641,
				"timestamp": 1609585297509136,
				"totalDecodeTime": 3.831,
				"totalInterFrameDelay": 160.8540000000006,
				"totalSquaredInterFrameDelay": 26.77570400000011,
				"trackId": "RTCMediaStreamTrack_receiver_3",
				"transportId": "RTCTransport_video_Nl0AIE_1",
				"type": "inbound-rtp"
			}
		]
	}`
	compare(t, resp, "inbound_rtp")
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

func TestTransport(t *testing.T) {
	resp := `{
		"version": "WebRTC Native Client Momo 2020.11 (db9d97e)",
		"libwebrtc": "Shiguredo-Build M88.4324@{#2} (88.4324.2.0 54bd8488)",
		"environment": "[aarch64] Ubuntu 18.04.5 LTS (nvidia-l4t-core 32.4.4-20201016123640)",
		"stats": [
			{
				"bytesReceived": 21186,
				"bytesSent": 5335226,
				"dtlsCipher": "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
				"dtlsState": "connected",
				"id": "RTCTransport_0_1",
				"localCertificateId": "RTCCertificate_D0:FF:4F:85:E1:67:31:83:33:F1:47:6E:08:65:FC:25:78:09:51:DF:04:51:7F:7B:E4:EE:CF:DA:D7:5C:94:41",
				"packetsReceived": 382,
				"packetsSent": 4904,
				"remoteCertificateId": "RTCCertificate_C0:92:CF:1A:63:65:5B:93:91:5A:81:F9:57:C0:0E:66:59:EC:47:BD:04:C8:8E:21:65:D3:F9:C9:F8:06:59:99",
				"selectedCandidatePairChanges": 2,
				"selectedCandidatePairId": "RTCIceCandidatePair_vpgjsoAn_zQSEz4UN",
				"srtpCipher": "AES_CM_128_HMAC_SHA1_80",
				"timestamp": 1608309189926189,
				"tlsVersion": "FEFD",
				"type": "transport"
			}
		]
	}`
	compare(t, resp, "transport")
}

func TestDeadline(t *testing.T) {
	exit := make(chan bool)
	h := httptest.NewServer(handlerStale(exit))
	defer func() {
		// s.Close() will block until the handler
		// returns, so we need to make it exit.
		exit <- true
		h.Close()
	}()

	e, err := NewExporter(h.URL, true, 1*time.Second, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}

	expectMetrics(t, e, "deadline")
}

func TestNotFound(t *testing.T) {
	h := httptest.NewServer(http.NotFoundHandler())
	defer h.Close()

	e, err := NewExporter(h.URL, true, 5*time.Second, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}

	expectMetrics(t, e, "not_found")
}
