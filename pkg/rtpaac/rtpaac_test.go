package rtpaac

import (
	"bytes"
	"testing"
	"time"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

func mergeBytes(vals ...[]byte) []byte {
	size := 0
	for _, v := range vals {
		size += len(v)
	}
	res := make([]byte, size)

	pos := 0
	for _, v := range vals {
		n := copy(res[pos:], v)
		pos += n
	}

	return res
}

var cases = []struct {
	name string
	aus  [][]byte
	pts  time.Duration
	pkts []*rtp.Packet
}{
	{
		"single",
		[][]byte{
			{
				0x21, 0x1a, 0xd4, 0xf5, 0x9e, 0x20, 0xc5, 0x42,
				0x89, 0x40, 0xa2, 0x9b, 0x3c, 0x94, 0xdd, 0x28,
				0x94, 0x48, 0xd5, 0x8b, 0xb0, 0x2, 0xdb, 0x1b,
				0xeb, 0xe0, 0xfa, 0x9f, 0xea, 0x91, 0xa7, 0x3,
				0xe8, 0x6b, 0xe5, 0x5, 0x95, 0x6, 0x62, 0x88,
				0x13, 0xa, 0x15, 0xa0, 0xeb, 0xef, 0x40, 0x82,
				0xdf, 0x49, 0xf2, 0xe0, 0x26, 0xfc, 0x52, 0x5b,
				0x6c, 0x2a, 0x2d, 0xe8, 0xa5, 0x70, 0xc5, 0xaf,
				0xfc, 0x98, 0x9a, 0x2f, 0x1f, 0xbb, 0xa2, 0xcb,
				0xb8, 0x26, 0xb6, 0x6e, 0x4c, 0x15, 0x6c, 0x21,
				0x3d, 0x35, 0xf6, 0xcf, 0xa4, 0x3b, 0x72, 0x26,
				0xe1, 0x3a, 0x3a, 0x99, 0xd8, 0x2d, 0x6a, 0x22,
				0xcd, 0x97, 0xa, 0xef, 0x52, 0x9c, 0x5f, 0xcd,
				0x5c, 0xd9, 0xd3, 0x12, 0x7e, 0x45, 0x45, 0xb3,
				0x24, 0xef, 0xd3, 0x4f, 0x2f, 0x96, 0xd9, 0x8b,
				0x9c, 0xc2, 0xcd, 0x54, 0xb, 0x6e, 0x19, 0x84,
				0x56, 0xeb, 0x85, 0x52, 0x63, 0x64, 0x28, 0xb2,
				0xf2, 0xcf, 0xb8, 0xa8, 0x71, 0x53, 0x6, 0x82,
				0x88, 0xf2, 0xc4, 0xe1, 0x7d, 0x65, 0x54, 0xe0,
				0x5e, 0xc8, 0x38, 0x75, 0x9d, 0xb0, 0x58, 0x65,
				0x41, 0xa2, 0xcd, 0xdb, 0x1b, 0x9e, 0xac, 0xd1,
				0xbe, 0xc9, 0x22, 0xf5, 0xe9, 0xc6, 0x6f, 0xaf,
				0xf8, 0xb1, 0x4c, 0xcb, 0xa2, 0x56, 0x11, 0xa4,
				0xd7, 0xfd, 0xe5, 0xef, 0x8e, 0xbf, 0xce, 0x4b,
				0xef, 0xe1, 0xd, 0xc0, 0x27, 0x18, 0xe2, 0x64,
				0x63, 0x5, 0x16, 0x6, 0xc, 0x34, 0xe, 0xf3, 0x62,
				0xc2, 0xd6, 0x42, 0x5d, 0x66, 0x81, 0x4, 0x65,
				0x76, 0xaa, 0xe7, 0x39, 0xdd, 0x8e, 0xfe, 0x48,
				0x23, 0x3a, 0x1, 0xc4, 0xd3, 0x65, 0x80, 0x28,
				0x6f, 0x9b, 0xc9, 0xb7, 0x4e, 0x44, 0x4c, 0x98,
				0x6a, 0x5f, 0x3b, 0x97, 0x81, 0x9b, 0xa9, 0xab,
				0xfd, 0xcf, 0x8e, 0x78, 0xbd, 0x4d, 0x70, 0x81,
				0x9b, 0x2d, 0x85, 0x94, 0x74, 0x2a, 0x3a, 0xb4,
				0xff, 0x4a, 0x13, 0x70, 0x76, 0x2c, 0x2f, 0x13,
				0x5b, 0x43, 0xf9, 0x17, 0xee, 0x26, 0x37, 0x1,
				0xbc, 0x9f, 0xb, 0xe, 0x68, 0xcb, 0x87, 0x65,
				0x86, 0xcc, 0x4c, 0x2f, 0x7a, 0x14, 0xd, 0xd1,
				0xb9, 0x57, 0xbd, 0x50, 0xb6, 0x95, 0x44, 0x1a,
				0xd, 0xc0, 0x15, 0xf, 0xd2, 0xc3, 0x72, 0x4d,
				0x6e, 0x4f, 0x8e, 0x6d, 0x64, 0xdc, 0x64, 0x1f,
				0x33, 0x53, 0x4e, 0xd8, 0xa4, 0x74, 0xf3, 0x33,
				0x4, 0x68, 0xd9, 0x92, 0xf3, 0x6e, 0xb7, 0x5b,
				0xe6, 0xf6, 0xc3, 0x55, 0x14, 0x54, 0x87, 0x0,
				0xaf, 0x7,
			},
		},
		20 * time.Millisecond,
		[]*rtp.Packet{
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17645,
					Timestamp:      2289527317,
					SSRC:           0x9dbb7812,
				},
				Payload: []byte{
					0x00, 0x10, 0x0a, 0xd8,
					0x21, 0x1a, 0xd4, 0xf5, 0x9e, 0x20, 0xc5, 0x42,
					0x89, 0x40, 0xa2, 0x9b, 0x3c, 0x94, 0xdd, 0x28,
					0x94, 0x48, 0xd5, 0x8b, 0xb0, 0x02, 0xdb, 0x1b,
					0xeb, 0xe0, 0xfa, 0x9f, 0xea, 0x91, 0xa7, 0x03,
					0xe8, 0x6b, 0xe5, 0x05, 0x95, 0x06, 0x62, 0x88,
					0x13, 0x0a, 0x15, 0xa0, 0xeb, 0xef, 0x40, 0x82,
					0xdf, 0x49, 0xf2, 0xe0, 0x26, 0xfc, 0x52, 0x5b,
					0x6c, 0x2a, 0x2d, 0xe8, 0xa5, 0x70, 0xc5, 0xaf,
					0xfc, 0x98, 0x9a, 0x2f, 0x1f, 0xbb, 0xa2, 0xcb,
					0xb8, 0x26, 0xb6, 0x6e, 0x4c, 0x15, 0x6c, 0x21,
					0x3d, 0x35, 0xf6, 0xcf, 0xa4, 0x3b, 0x72, 0x26,
					0xe1, 0x3a, 0x3a, 0x99, 0xd8, 0x2d, 0x6a, 0x22,
					0xcd, 0x97, 0x0a, 0xef, 0x52, 0x9c, 0x5f, 0xcd,
					0x5c, 0xd9, 0xd3, 0x12, 0x7e, 0x45, 0x45, 0xb3,
					0x24, 0xef, 0xd3, 0x4f, 0x2f, 0x96, 0xd9, 0x8b,
					0x9c, 0xc2, 0xcd, 0x54, 0x0b, 0x6e, 0x19, 0x84,
					0x56, 0xeb, 0x85, 0x52, 0x63, 0x64, 0x28, 0xb2,
					0xf2, 0xcf, 0xb8, 0xa8, 0x71, 0x53, 0x06, 0x82,
					0x88, 0xf2, 0xc4, 0xe1, 0x7d, 0x65, 0x54, 0xe0,
					0x5e, 0xc8, 0x38, 0x75, 0x9d, 0xb0, 0x58, 0x65,
					0x41, 0xa2, 0xcd, 0xdb, 0x1b, 0x9e, 0xac, 0xd1,
					0xbe, 0xc9, 0x22, 0xf5, 0xe9, 0xc6, 0x6f, 0xaf,
					0xf8, 0xb1, 0x4c, 0xcb, 0xa2, 0x56, 0x11, 0xa4,
					0xd7, 0xfd, 0xe5, 0xef, 0x8e, 0xbf, 0xce, 0x4b,
					0xef, 0xe1, 0x0d, 0xc0, 0x27, 0x18, 0xe2, 0x64,
					0x63, 0x05, 0x16, 0x06, 0x0c, 0x34, 0x0e, 0xf3,
					0x62, 0xc2, 0xd6, 0x42, 0x5d, 0x66, 0x81, 0x04,
					0x65, 0x76, 0xaa, 0xe7, 0x39, 0xdd, 0x8e, 0xfe,
					0x48, 0x23, 0x3a, 0x01, 0xc4, 0xd3, 0x65, 0x80,
					0x28, 0x6f, 0x9b, 0xc9, 0xb7, 0x4e, 0x44, 0x4c,
					0x98, 0x6a, 0x5f, 0x3b, 0x97, 0x81, 0x9b, 0xa9,
					0xab, 0xfd, 0xcf, 0x8e, 0x78, 0xbd, 0x4d, 0x70,
					0x81, 0x9b, 0x2d, 0x85, 0x94, 0x74, 0x2a, 0x3a,
					0xb4, 0xff, 0x4a, 0x13, 0x70, 0x76, 0x2c, 0x2f,
					0x13, 0x5b, 0x43, 0xf9, 0x17, 0xee, 0x26, 0x37,
					0x01, 0xbc, 0x9f, 0x0b, 0x0e, 0x68, 0xcb, 0x87,
					0x65, 0x86, 0xcc, 0x4c, 0x2f, 0x7a, 0x14, 0x0d,
					0xd1, 0xb9, 0x57, 0xbd, 0x50, 0xb6, 0x95, 0x44,
					0x1a, 0x0d, 0xc0, 0x15, 0x0f, 0xd2, 0xc3, 0x72,
					0x4d, 0x6e, 0x4f, 0x8e, 0x6d, 0x64, 0xdc, 0x64,
					0x1f, 0x33, 0x53, 0x4e, 0xd8, 0xa4, 0x74, 0xf3,
					0x33, 0x04, 0x68, 0xd9, 0x92, 0xf3, 0x6e, 0xb7,
					0x5b, 0xe6, 0xf6, 0xc3, 0x55, 0x14, 0x54, 0x87,
					0x00, 0xaf, 0x07,
				},
			},
		},
	},
	{
		"aggregated",
		[][]byte{
			{0x00, 0x01, 0x02, 0x03},
			{0x04, 0x05, 0x06, 0x07},
			{0x08, 0x09, 0x0A, 0x0B},
		},
		0,
		[]*rtp.Packet{
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17645,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: []byte{
					0x0, 0x30, 0x0, 0x20,
					0x0, 0x20, 0x0, 0x20, 0x0, 0x1, 0x2, 0x3,
					0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb,
				},
			},
		},
	},
	{
		"fragmented",
		[][]byte{
			bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 512),
		},
		0,
		[]*rtp.Packet{
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         false,
					PayloadType:    96,
					SequenceNumber: 17645,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: mergeBytes(
					[]byte{0x0, 0x10, 0x2d, 0x80},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
				),
			},
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         false,
					PayloadType:    96,
					SequenceNumber: 17646,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: mergeBytes(
					[]byte{0x00, 0x10, 0x2d, 0x80},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
				),
			},
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17647,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: mergeBytes(
					[]byte{0x00, 0x10, 0x25, 0x00},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 148),
				),
			},
		},
	},
	{
		"aggregated followed by fragmented",
		[][]byte{
			{0x00, 0x01, 0x02, 0x03},
			{0x04, 0x05, 0x06, 0x07},
			{0x08, 0x09, 0x0A, 0x0B},
			bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 256),
		},
		0,
		[]*rtp.Packet{
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17645,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: []byte{
					0x0, 0x30, 0x0, 0x20,
					0x0, 0x20, 0x0, 0x20, 0x0, 0x1, 0x2, 0x3,
					0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb,
				},
			},
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         false,
					PayloadType:    96,
					SequenceNumber: 17646,
					Timestamp:      2289529357,
					SSRC:           0x9dbb7812,
				},
				Payload: mergeBytes(
					[]byte{0x0, 0x10, 0x2d, 0x80},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
				),
			},
			{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17647,
					Timestamp:      2289529357,
					SSRC:           0x9dbb7812,
				},
				Payload: mergeBytes(
					[]byte{0x00, 0x10, 0x12, 0x80},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 74),
				),
			},
		},
	},
}

func TestDecode(t *testing.T) {
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			d := &Decoder{
				SampleRate:       48000,
				IndexLength:      3,
				SizeLength:       13,
				IndexDeltaLength: 3,
			}
			d.Init()

			// send an initial packet downstream
			// in order to compute the right timestamp,
			// that is relative to the initial packet
			pkt := rtp.Packet{
				Header: rtp.Header{
					Version:        2,
					Marker:         true,
					PayloadType:    96,
					SequenceNumber: 17645,
					Timestamp:      2289526357,
					SSRC:           0x9dbb7812,
				},
				Payload: []byte{0x00, 0x10, 0x00, 0x08, 0x0},
			}
			_, _, err := d.Decode(&pkt)
			require.NoError(t, err)

			var aus [][]byte
			expPTS := ca.pts

			for _, pkt := range ca.pkts {
				clone := pkt.Clone()

				addAUs, pts, err := d.Decode(pkt)
				if err == ErrMorePacketsNeeded {
					continue
				}

				require.NoError(t, err)
				require.Equal(t, expPTS, pts)
				aus = append(aus, addAUs...)
				expPTS += time.Duration(len(aus)) * 1000 * time.Second / 48000

				// test input integrity
				require.Equal(t, clone, pkt)
			}

			require.Equal(t, ca.aus, aus)
		})
	}
}

func TestDecodeErrors(t *testing.T) {
	for _, ca := range []struct {
		name string
		pkts []*rtp.Packet
		err  string
	}{
		{
			"missing payload",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
				},
			},
			"payload is too short",
		},
		{
			"missing au header",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x10},
				},
			},
			"EOF",
		},
		{
			"missing au",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x10, 0x0a, 0xd8},
				},
			},
			"payload is too short",
		},
		{
			"invalid au headers length 1",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x00},
				},
			},
			"invalid AU-headers-length (0)",
		},
		{
			"invalid au headers length 2",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x09},
				},
			},
			"invalid AU-headers-length (9)",
		},
		{
			"au index not zero",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x10, 0x0a, 0xd9},
				},
			},
			"AU-index different than zero is not supported",
		},
		{
			"au index delta not zero",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x20, 0x00, 0x08, 0x0a, 0xd9},
				},
			},
			"AU-index-delta different than zero is not supported",
		},
		{
			"fragmented with multiple AUs in 1st packet",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         false,
						PayloadType:    0x60,
						SequenceNumber: 0xea2,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x20, 0x00, 0x08, 0x00, 0x08},
				},
			},
			"a fragmented packet can only contain one AU",
		},
		{
			"fragmented with no payload in 1st packet",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         false,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x10, 0x0a, 0xd8},
				},
			},
			"payload is too short",
		},
		{
			"fragmented with multiple AUs in 2nd packet",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         false,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: mergeBytes(
						[]byte{0x0, 0x10, 0x2d, 0x80},
						bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
					),
				},
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ee,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: mergeBytes(
						[]byte{0x0, 0x20, 0x00, 0x08, 0x00, 0x08},
					),
				},
			},
			"a fragmented packet can only contain one AU",
		},
		{
			"fragmented with no payload in 2nd packet",
			[]*rtp.Packet{
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         false,
						PayloadType:    0x60,
						SequenceNumber: 0x44ed,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: mergeBytes(
						[]byte{0x0, 0x10, 0x2d, 0x80},
						bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
					),
				},
				{
					Header: rtp.Header{
						Version:        2,
						Marker:         true,
						PayloadType:    0x60,
						SequenceNumber: 0x44ee,
						Timestamp:      0x88776a15,
						SSRC:           0x9dbb7812,
					},
					Payload: []byte{0x00, 0x10, 0x0a, 0xd8},
				},
			},
			"payload is too short",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			d := &Decoder{
				SampleRate:  48000,
				IndexLength: 3,
				SizeLength:  13,
			}
			d.Init()

			var lastErr error
			for _, pkt := range ca.pkts {
				_, _, lastErr = d.Decode(pkt)
			}
			require.EqualError(t, lastErr, ca.err)
		})
	}
}

func TestEncode(t *testing.T) {
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			e := &Encoder{
				PayloadType: 96,
				SampleRate:  48000,
				SSRC: func() *uint32 {
					v := uint32(0x9dbb7812)
					return &v
				}(),
				InitialSequenceNumber: func() *uint16 {
					v := uint16(0x44ed)
					return &v
				}(),
				InitialTimestamp: func() *uint32 {
					v := uint32(0x88776655)
					return &v
				}(),
			}
			e.Init()

			pkts, err := e.Encode(ca.aus, ca.pts)
			require.NoError(t, err)
			require.Equal(t, ca.pkts, pkts)
		})
	}
}

func TestEncodeRandomInitialState(t *testing.T) {
	e := &Encoder{
		PayloadType: 96,
		SampleRate:  48000,
	}
	e.Init()
	require.NotEqual(t, nil, e.SSRC)
	require.NotEqual(t, nil, e.InitialSequenceNumber)
	require.NotEqual(t, nil, e.InitialTimestamp)
}
