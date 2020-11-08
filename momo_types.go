package main

// MomoMetrics is metrics respose type of WebRTC Native Client Momo
type MomoMetrics struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Libwebrtc   string `json:"libwebrtc"`
	Stats       string `json:"stats"`
}

// WebRTC Stats Types
// https://www.w3.org/TR/webrtc-stats

// RTCStats is base type of all WebRTC Stats
type RTCStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`
}

type RTCRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`
}

type RTCCodecStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	PayloadType uint64 `json:"payloadType"`
	CodecType   string `json:"codecType"` // "encode" or "decode"
	TransportID string `json:"transportId"`
	MimeType    string `json:"mimeType"`
	ClockRate   uint64 `json:"clockRate"`
	Channels    uint64 `json:"chennels"`
	SDPFmtpLine string `json:"sdpFmtpLine"`
}

type RTCReceivedRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsReceived       uint64  `json:"packetsReceived"`
	PacketsLost           int64   `json:"packetsLost"`
	Jitter                float64 `json:"jitter"`
	PacketsDiscarded      uint64  `json:"packetsDiscarded"`
	PacketsReparied       uint64  `json:"packetsRepaired"`
	BurstPacketsLost      uint64  `json:"burstPacketsLost"`
	BurstPacketsDiscarded uint64  `json:"burstPacketsDiscarded"`
	BurstLossCount        uint64  `json:"burstLossCount"`
	BurstDiscardCount     uint64  `json:"burstDiscardCount"`
	BurstLoassRate        float64 `json:"burstLossRate"`
	BurstDiscardRate      float64 `json:"burstDiscardRate"`
	GapLossRate           float64 `json:"gapLossRate"`
	GapDiscardRate        float64 `json:"gapDiscardRate"`
	FrameDropped          uint64  `json:"framesDropped"`
	PartialFramesLost     uint64  `json:"partialFramesLost"`
	FullFramesLost        uint64  `json:"fullFramesLost"`
}

type RTCInboundRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsReceived       uint64  `json:"packetsReceived"`
	PacketsLost           int64   `json:"packetsLost"`
	Jitter                float64 `json:"jitter"`
	PacketsDiscarded      uint64  `json:"packetsDiscarded"`
	PacketsReparied       uint64  `json:"packetsRepaired"`
	BurstPacketsLost      uint64  `json:"burstPacketsLost"`
	BurstPacketsDiscarded uint64  `json:"burstPacketsDiscarded"`
	BurstLossCount        uint64  `json:"burstLossCount"`
	BurstDiscardCount     uint64  `json:"burstDiscardCount"`
	BurstLoassRate        float64 `json:"burstLossRate"`
	BurstDiscardRate      float64 `json:"burstDiscardRate"`
	GapLossRate           float64 `json:"gapLossRate"`
	GapDiscardRate        float64 `json:"gapDiscardRate"`
	FrameDropped          uint64  `json:"framesDropped"`
	PartialFramesLost     uint64  `json:"partialFramesLost"`
	FullFramesLost        uint64  `json:"fullFramesLost"`

	ReceiverID                  string  `json:"receiverId"`
	RemoteID                    string  `json:"remoteId"`
	FramesDecoded               uint64  `json:"framesDecoded"`
	KeyFramesDecoded            uint64  `json:"keyFramesDecoded"`
	FrameWidth                  uint64  `json:"frameWidth"`
	FrameHeight                 uint64  `json:"frameHeight"`
	FrameBitDepth               uint64  `json:"frameBitDepth"`
	FramesPerSecond             float64 `json:"framesPerSecond"`
	QPSum                       uint64  `json:"qpSum"`
	TotalDecodeTime             float64 `json:"totalDecodeTime"`
	TotalInterFrameDelay        float64 `json:"totalInterFrameDelay"`
	TotalSquaredInterFrameDelay float64 `json:"totalSquaredInterFrameDelay"`
	VoiceActivityFlag           bool    `json:"voiceActivityFlag"`
	LastPacketReceivedTimestamp float64 `json:"lastPacketReceivedTimestamp"`
	AverageRtcpInterval         float64 `json:"averageRtcpInterval"`
	HeaderBytesReceived         uint64  `json:"headerBytesReceived"`
	FecPacketsReceived          uint64  `json:"fecPacketsReceived"`
	FecPacketsDiscarded         uint64  `json:"fecPacketsDiscarded"`
	BytesReceived               uint64  `json:"bytesReceived"`
	PacketsFailedDecryption     uint64  `json:"packetsFailedDecryption"`
	PacketsDuplicated           uint64  `json:"packetsDuplicated"`
	//record<USVString, uint64> perDscpPacketsReceived"`
	NackCount                      uint64  `json:"nackCount"`
	FIRCount                       uint64  `json:"firCount"`
	PLICount                       uint64  `json:"pliCount"`
	SLICount                       uint64  `json:"sliCount"`
	TotalProcessingDelay           float64 `json:"totalProcessingDelay"`
	EstimatedPlayoutTimestamp      float64 `json:"estimatedPlayoutTimestamp"`
	JitterBufferDelay              float64 `json:"jitterBufferDelay"`
	JitterBufferEmittedCount       uint64  `json:"jitterBufferEmittedCount"`
	TotalSamplesReceived           uint64  `json:"totalSamplesReceived"`
	TotalSamplesDecoded            uint64  `json:"totalSamplesDecoded"`
	TotalSamplesDecodedWithSilk    uint64  `json:"samplesDecodedWithSilk"`
	SamplesDecodedWithCelt         uint64  `json:"samplesDecodedWithCelt"`
	ConcealedSamples               uint64  `json:"concealedSamples"`
	SilentConcealedSamples         uint64  `json:"silentConcealedSamples"`
	ConcealmentEvents              uint64  `json:"concealmentEvents"`
	InsertedSamplesForDeceleration uint64  `json:"insertedSamplesForDeceleration"`
	RemovedSamplesForAcceleration  uint64  `json:"removedSamplesForAcceleration"`
	AudioLevel                     float64 `json:"audioLevel"`
	TotalAudioEnergy               float64 `json:"totalAudioEnergy"`
	TotalSamplesDuration           float64 `json:"totalSamplesDuration"`
	FramesReceived                 uint64  `json:"framesReceived"`
	DecoderImplementation          string  `json:"decoderImplementation"`
}

type RTCRemoteInboundRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsReceived       uint64  `json:"packetsReceived"`
	PacketsLost           int64   `json:"packetsLost"`
	Jitter                float64 `json:"jitter"`
	PacketsDiscarded      uint64  `json:"packetsDiscarded"`
	PacketsReparied       uint64  `json:"packetsRepaired"`
	BurstPacketsLost      uint64  `json:"burstPacketsLost"`
	BurstPacketsDiscarded uint64  `json:"burstPacketsDiscarded"`
	BurstLossCount        uint64  `json:"burstLossCount"`
	BurstDiscardCount     uint64  `json:"burstDiscardCount"`
	BurstLoassRate        float64 `json:"burstLossRate"`
	BurstDiscardRate      float64 `json:"burstDiscardRate"`
	GapLossRate           float64 `json:"gapLossRate"`
	GapDiscardRate        float64 `json:"gapDiscardRate"`
	FrameDropped          uint64  `json:"framesDropped"`
	PartialFramesLost     uint64  `json:"partialFramesLost"`
	FullFramesLost        uint64  `json:"fullFramesLost"`

	LocalID                   string  `json:"localId"`
	RoundTripTime             float64 `json:"roundTripTime"`
	TotalRoundTripTime        float64 `json:"totalRoundTripTime"`
	FractionLost              float64 `json:"fractionLost"`
	ReportsReceived           uint64  `json:"reportsReceived"`
	RoundTripTimeMeasurements uint64  `json:"roundTripTimeMeasurements"`
}

type RTCSentRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsSent uint64 `json:"packetsSent"`
	BytesSent   uint64 `json:"bytesSent"`
}

type RTCOUtboundRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsSent uint64 `json:"packetsSent"`
	BytesSent   uint64 `json:"bytesSent"`

	RtxSsrc                  uint64  `json:"rtxSsrc"`
	MediaSourceID            string  `json:"mediaSourceId"`
	SenderID                 string  `json:"senderId"`
	RemoteID                 string  `json:"remoteId"`
	Rid                      string  `json:"rid"`
	LastPacketSentTimestamp  float64 `json:"lastPacketSentTimestamp"`
	HeaderBytesSent          uint64  `json:"headerBytesSent"`
	PacketsDiscardedOnSend   uint64  `json:"packetsDiscardedOnSend"`
	BytesDiscardedOnSend     uint64  `json:"bytesDiscardedOnSend"`
	FecPacketsSent           uint64  `json:"fecPacketsSent"`
	RetransmittedPacketsSent uint64  `json:"retransmittedPacketsSent"`
	RetransmittedBytesSent   uint64  `json:"retransmittedBytesSent"`
	TargetBitrate            float64 `json:"targetBitrate"`
	TotalEncodedBytesTarget  uint64  `json:"totalEncodedBytesTarget"`
	FrameWidth               uint64  `json:"frameWidth"`
	FrameHeight              uint64  `json:"frameHeight"`
	FramesBitDepth           uint64  `json:"frameBitDepth"`
	FramesPerSecond          float64 `json:"framesPerSecond"`
	FramesSent               uint64  `json:"framesSent"`
	HugeFramesSent           uint64  `json:"hugeFramesSent"`
	FramesEncoded            uint64  `json:"framesEncoded"`
	KeyFramesEncoded         uint64  `json:"keyFramesEncoded"`
	FramesDiscardedOnSend    uint64  `json:"framesDiscardedOnSend"`
	QPSum                    uint64  `json:"qpSum"`
	TotaoSamplesSend         uint64  `json:"totalSamplesSent"`
	SamplesEncodedWithSilk   uint64  `json:"samplesEncodedWithSilk"`
	SamplesEncodedWithCelt   uint64  `json:"samplesEncodedWithCelt"`
	VoiceActivityFlag        bool    `json:"voiceActivityFlag"`
	TotalEncodeTime          float64 `json:"totalEncodeTime"`
	TotalPacketSendDelay     float64 `json:"totalPacketSendDelay"`
	AverageRtcpInterval      float64 `json:"averageRtcpInterval"`
	QualityLimitationReason  string  `json:"qualityLimitationReason"`
	//record<string, float64> qualityLimitationDurations"`
	QualityLimitationResolutionChanges uint64 `json:"qualityLimitationResolutionChanges"`
	//record<USVString, uint64> perDscpPacketsSent"`
	NackCount             uint64 `json:"nackCount"`
	FIRCount              uint64 `json:"firCount"`
	PLICount              uint64 `json:"pliCount"`
	SLICount              uint64 `json:"sliCount"`
	EncoderImplementation string `json:"encoderImplementation"`
}

type RTCRemoteOutboundRtpStreamStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	SSRC        string `json:"ssrc"`
	Kind        string `json:"kind"`
	TransportID string `json:"transportId"`
	CodecID     string `json:"codecId"`

	PacketsSent uint64 `json:"packetsSent"`
	BytesSent   uint64 `json:"bytesSent"`

	LocalID         string  `json:"localId"`
	RemoteTimestamp float64 `json:"remoteTimestamp"`
	ReportsSent     uint64  `json:"reportsSent"`
}

type RTCMediaSourceStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Kind            string `json:"kind"`
	RelayedSource   bool   `json:"relayedSource"`
}

type RTCAudioSourceStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Kind            string `json:"kind"`
	RelayedSource   bool   `json:"relayedSource"`

	AudioLovel                float64 `json:"audioLevel"`
	TotalAudioEnergy          float64 `json:"totalAudioEnergy"`
	TotalSamplesDuration      float64 `json:"totalSamplesDuration"`
	EchoReturnLoss            float64 `json:"echoReturnLoss"`
	EchoReturnLossEnhancement float64 `json:"echoReturnLossEnhancement"`
}

type RTCVideoSourceStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Kind            string `json:"kind"`
	RelayedSource   bool   `json:"relayedSource"`

	Width           uint64  `json:"width"`
	Height          uint64  `json:"height"`
	BitDepth        uint64  `json:"bitDepth"`
	Frames          uint64  `json:"frames"`
	FramesPerSecond float64 `json:"framesPerSecond"`
}

type RTCPeerConnectionStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	DataChannelsOpened    uint64 `json:"dataChannelsOpened"`
	DataChannelsClosed    uint64 `json:"dataChannelsClosed"`
	DataChannelsRequested uint64 `json:"dataChannelsRequested"`
	DataChannelsAccepted  uint64 `json:"dataChannelsAccepted"`
}

type RTCMediaHandlerStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`
}

type RTCVideoHandlerStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`
}

type RTCVideoSenderStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`

	MediaSourceID string `json:"mediaSourceId"`
}

type RTCVideoReceiverStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`
}

type RTCAudioHandlerStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`
}

type RTCAudioSenderStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`

	MediaSourceID string `json:"mediaSourceId"`
}

type RTCAudioReceiverStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TrackIdentifier string `json:"trackIdentifier"`
	Ended           bool   `json:"ended"`
	Kind            string `json:"kind"`
}

type RTCDataChannelStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	Label                 string `json:"label"`
	Protocol              string `json:"protocol"`
	DataChannelIdentifier uint16 `json:"dataChannelIdentifier"`
	State                 string `json:"state"`
	MessagesSent          uint64 `json:"messagesSent"`
	BytesSent             uint64 `json:"bytesSent"`
	MessagesReceived      uint64 `json:"messagesReceived"`
	BytesReceived         uint64 `json:"bytesReceived"`
}

type RTCTransportStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	PacketsSent                  uint64 `json:"packetsSent"`
	PacketsReceived              uint64 `json:"packetsReceived"`
	BytesSent                    uint64 `json:"bytesSent"`
	BytesReceived                uint64 `json:"bytesReceived"`
	RTCPTransportStatsID         string `json:"rtcpTransportStatsId"`
	ICERole                      string `json:"iceRole"`
	ICELocalUsernameFragment     string `json:"iceLocalUsernameFragment"`
	DTLSState                    string `json:"dtlsState"`
	ICEState                     string `json:"iceState"`
	SelectedCandidatePairID      string `json:"selectedCandidatePairId"`
	LocalCertificateID           string `json:"localCertificateId"`
	RemoteCertificateID          string `json:"remoteCertificateId"`
	TLSVersion                   string `json:"tlsVersion"`
	DTLSCipher                   string `json:"dtlsCipher"`
	SRTPCipher                   string `json:"srtpCipher"`
	TLSGroup                     string `json:"tlsGroup"`
	SelectedCandidatePairChanges uint64 `json:"selectedCandidatePairChanges"`
}

type RTCSctpTransportStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TransportID           string  `json:"transportId"`
	SmoothedRoundTripTime float64 `json:"smoothedRoundTripTime"`
	CongestionWindow      uint64  `json:"congestionWindow"`
	ReceiverWindow        uint64  `json:"receiverWindow"`
	MTU                   uint64  `json:"mtu"`
	UnackData             uint64  `json:"unackData"`
}

type RTCIceCandidateStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TransportID   string `json:"transportId"`
	Address       string `json:"address"`
	Port          int64  `json:"port"`
	Protocol      string `json:"protocol"`
	CandidateType string `json:"candidateType"`
	Priority      int64  `json:"priority"`
	URL           string `json:"url"`
	RelayProtocol string `json:"relayProtocol"`
}

type RTCIceCandidatePairStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	TransportID                 string  `json:"transportId"`
	LocalCandidateID            string  `json:"localCandidateId"`
	RemoteCandidateID           string  `json:"remoteCandidateId"`
	State                       string  `json:"state"`
	Nominated                   bool    `json:"nominated"`
	PacketsSent                 uint64  `json:"packetsSent"`
	PacketsReceived             uint64  `json:"packetsReceived"`
	BytesSent                   uint64  `json:"bytesSent"`
	BytesReceived               uint64  `json:"bytesReceived"`
	LastPacketSentTimestamp     float64 `json:"lastPacketSentTimestamp"`
	LastPacketReceivedTimestamp float64 `json:"lastPacketReceivedTimestamp"`
	FirstRequestTimestamp       float64 `json:"firstRequestTimestamp"`
	LastRequestTimestamp        float64 `json:"lastRequestTimestamp"`
	LastResponseTimestamp       float64 `json:"lastResponseTimestamp"`
	TotalRoundTripTime          float64 `json:"totalRoundTripTime"`
	CurrentRoundTripTime        float64 `json:"currentRoundTripTime"`
	AvailableOutgoingBitrate    float64 `json:"availableOutgoingBitrate"`
	AvailableIncomingBitrate    float64 `json:"availableIncomingBitrate"`

	CircuitBreakerTriggerCount uint64  `json:"circuitBreakerTriggerCount"`
	RequestsReceived           uint64  `json:"requestsReceived"`
	RequestsSent               uint64  `json:"requestsSent"`
	ResponsesReceived          uint64  `json:"responsesReceived"`
	ResponsesSent              uint64  `json:"responsesSent"`
	TetransmissionsReceived    uint64  `json:"retransmissionsReceived"`
	RetransmissionsSent        uint64  `json:"retransmissionsSent"`
	ConsentRequestsSent        uint64  `json:"consentRequestsSent"`
	ConsentExpiredTimestamp    float64 `json:"consentExpiredTimestamp"`
	PacketsDiscardedOnSend     uint64  `json:"packetsDiscardedOnSend"`
	BytesDiscardedOnSend       uint64  `json:"bytesDiscardedOnSend"`
	RequestBytesSent           uint64  `json:"requestBytesSent"`
	ConsentRequestBytesSent    uint64  `json:"consentRequestBytesSent"`
	ResponseBytesSent          uint64  `json:"responseBytesSent"`
}

type RTCCertificateStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	Fingerprint          string `json:"fingerprint"`
	FingerprintAlgorithm string `json:"fingerprintAlgorithm"`
	Base64Certificate    string `json:"base64Certificate"`
	IssuerCertificateID  string `json:"issuerCertificateId"`
}

type RTCIceServerStats struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`

	URL                    string  `json:"url"`
	Port                   int64   `json:"port"`
	RelayProtocol          string  `json:"relayProtocol"`
	TotalRequestsSent      uint64  `json:"totalRequestsSent"`
	TotalResponsesReceived uint64  `json:"totalResponsesReceived"`
	TotalRoundTripTime     float64 `json:"totalRoundTripTime"`
}
