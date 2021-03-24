package rtpaac

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/pion/rtp"
)

const (
	rtpVersion        = 0x02
	rtpPayloadMaxSize = 1460 // 1500 (mtu) - 20 (ip header) - 8 (udp header) - 12 (rtp header)
)

// Encoder is a RTP/AAC encoder.
type Encoder struct {
	payloadType    uint8
	clockRate      float64
	sequenceNumber uint16
	ssrc           uint32
	initialTs      uint32
}

// NewEncoder allocates an Encoder.
func NewEncoder(payloadType uint8,
	clockRate int,
	sequenceNumber *uint16,
	ssrc *uint32,
	initialTs *uint32) *Encoder {
	return &Encoder{
		payloadType: payloadType,
		clockRate:   float64(clockRate),
		sequenceNumber: func() uint16 {
			if sequenceNumber != nil {
				return *sequenceNumber
			}
			return uint16(rand.Uint32())
		}(),
		ssrc: func() uint32 {
			if ssrc != nil {
				return *ssrc
			}
			return rand.Uint32()
		}(),
		initialTs: func() uint32 {
			if initialTs != nil {
				return *initialTs
			}
			return rand.Uint32()
		}(),
	}
}

func (e *Encoder) encodeTimestamp(ts time.Duration) uint32 {
	return e.initialTs + uint32(ts.Seconds()*e.clockRate)
}

// Encode encodes AUs into a RTP/AAC packet.
func (e *Encoder) Encode(ats []*AUAndTimestamp) ([]byte, error) {
	le := 2 // AU-headers-length
	for _, at := range ats {
		le += 2          // AU-header
		le += len(at.AU) // AU
	}

	if le > rtpPayloadMaxSize {
		return nil, fmt.Errorf("data is too big for a single packet")
	}

	payload := make([]byte, le)

	// AU-headers-length
	binary.BigEndian.PutUint16(payload, uint16(len(ats)*16))
	pos := 2

	// AU-headers
	for _, at := range ats {
		binary.BigEndian.PutUint16(payload[pos:], uint16(len(at.AU))<<3)
		pos += 2
	}

	// AUs
	for _, at := range ats {
		auLen := copy(payload[pos:], at.AU)
		pos += auLen
	}

	rpkt := rtp.Packet{
		Header: rtp.Header{
			Version:        rtpVersion,
			PayloadType:    e.payloadType,
			SequenceNumber: e.sequenceNumber,
			Timestamp:      e.encodeTimestamp(ats[0].Timestamp),
			SSRC:           e.ssrc,
			Marker:         true,
		},
		Payload: payload,
	}
	e.sequenceNumber++

	frame, err := rpkt.Marshal()
	if err != nil {
		return nil, err
	}

	return frame, nil
}
