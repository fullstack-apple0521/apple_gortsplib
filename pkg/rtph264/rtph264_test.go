package rtph264

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gortsplib/pkg/codech264"
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

type readerFunc func(p []byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) {
	return f(p)
}

var cases = []struct {
	name string
	dec  []*NALUAndTimestamp
	enc  [][]byte
}{
	{
		"single",
		[]*NALUAndTimestamp{
			{
				Timestamp: 25 * time.Millisecond,
				NALU: mergeBytes(
					[]byte{0x05},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8),
				),
			},
		},
		[][]byte{
			mergeBytes(
				[]byte{
					0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6f, 0x1f,
					0x9d, 0xbb, 0x78, 0x12, 0x05,
				},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8),
			),
		},
	},
	{
		"negative timestamp",
		[]*NALUAndTimestamp{
			{
				Timestamp: -20 * time.Millisecond,
				NALU: mergeBytes(
					[]byte{0x05},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8),
				),
			},
		},
		[][]byte{
			mergeBytes(
				[]byte{
					0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x5f, 0x4d,
					0x9d, 0xbb, 0x78, 0x12, 0x05,
				},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8),
			),
		},
	},
	{
		"fragmented",
		[]*NALUAndTimestamp{
			{
				Timestamp: 55 * time.Millisecond,
				NALU: mergeBytes(
					[]byte{0x05},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 256),
				),
			},
		},
		[][]byte{
			mergeBytes(
				[]byte{
					0x80, 0x60, 0x44, 0xed, 0x88, 0x77, 0x79, 0xab,
					0x9d, 0xbb, 0x78, 0x12, 0x1c, 0x85,
				},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
				[]byte{0x00, 0x01},
			),
			mergeBytes(
				[]byte{
					0x80, 0xe0, 0x44, 0xee, 0x88, 0x77, 0x79, 0xab,
					0x9d, 0xbb, 0x78, 0x12, 0x1c, 0x45,
				},
				[]byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 73),
			),
		},
	},
	{
		"aggregated",
		[]*NALUAndTimestamp{
			{
				NALU: []byte{0x09, 0xF0},
			},
			{
				NALU: []byte{
					0x41, 0x9a, 0x24, 0x6c, 0x41, 0x4f, 0xfe, 0xd6,
					0x8c, 0xb0, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x6d, 0x40,
				},
			},
		},
		[][]byte{
			{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x66, 0x55,
				0x9d, 0xbb, 0x78, 0x12, 0x18, 0x00, 0x02, 0x09,
				0xf0, 0x00, 0x44, 0x41, 0x9a, 0x24, 0x6c, 0x41,
				0x4f, 0xfe, 0xd6, 0x8c, 0xb0, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x6d, 0x40,
			},
		},
	},
	{
		"aggregated followed by single",
		[]*NALUAndTimestamp{
			{
				NALU: []byte{0x09, 0xF0},
			},
			{
				NALU: []byte{
					0x41, 0x9a, 0x24, 0x6c, 0x41, 0x4f, 0xfe, 0xd6,
					0x8c, 0xb0, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x6d, 0x40,
				},
			},
			{
				NALU: mergeBytes(
					[]byte{0x08},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 175),
				),
			},
		},
		[][]byte{
			{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x66, 0x55,
				0x9d, 0xbb, 0x78, 0x12, 0x18, 0x00, 0x02, 0x09,
				0xf0, 0x00, 0x44, 0x41, 0x9a, 0x24, 0x6c, 0x41,
				0x4f, 0xfe, 0xd6, 0x8c, 0xb0, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x6d, 0x40,
			},
			mergeBytes(
				[]byte{
					0x80, 0xe0, 0x44, 0xee, 0x88, 0x77, 0x66, 0x55,
					0x9d, 0xbb, 0x78, 0x12, 0x08,
				},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 175),
			),
		},
	},
	{
		"fragmented followed by aggregated",
		[]*NALUAndTimestamp{
			{
				NALU: mergeBytes(
					[]byte{0x05},
					bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 256),
				),
			},
			{
				NALU: []byte{0x09, 0xF0},
			},
			{
				NALU: []byte{0x09, 0xF0},
			},
		},
		[][]byte{
			mergeBytes(
				[]byte{
					0x80, 0x60, 0x44, 0xed, 0x88, 0x77, 0x66, 0x55,
					0x9d, 0xbb, 0x78, 0x12, 0x1c, 0x85,
				},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 182),
				[]byte{0x00, 0x01},
			),
			mergeBytes(
				[]byte{
					0x80, 0xe0, 0x44, 0xee, 0x88, 0x77, 0x66, 0x55,
					0x9d, 0xbb, 0x78, 0x12, 0x1c, 0x45,
				},
				[]byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 73),
			),
			{
				0x80, 0xe0, 0x44, 0xef, 0x88, 0x77, 0x66, 0x55,
				0x9d, 0xbb, 0x78, 0x12, 0x18, 0x00, 0x02, 0x09,
				0xf0, 0x00, 0x02, 0x09, 0xf0,
			},
		},
	},
}

func TestEncode(t *testing.T) {
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			sequenceNumber := uint16(0x44ed)
			ssrc := uint32(0x9dbb7812)
			initialTs := uint32(0x88776655)
			e := NewEncoder(96, &sequenceNumber, &ssrc, &initialTs)
			enc, err := e.Encode(ca.dec)
			require.NoError(t, err)
			require.Equal(t, ca.enc, enc)
		})
	}
}

func TestDecode(t *testing.T) {
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			i := 0
			r := readerFunc(func(p []byte) (int, error) {
				if i == len(ca.enc) {
					return 0, io.EOF
				}

				i++
				return copy(p, ca.enc[i-1]), nil
			})

			d := NewDecoder()

			// send an initial packet downstream
			// in order to compute the timestamp,
			// which is relative to the initial packet
			_, err := d.Decode([]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x66, 0x55,
				0x9d, 0xbb, 0x78, 0x12, 0x06, 0x00,
			})
			require.NoError(t, err)

			for _, dec0 := range ca.dec {
				dec, err := d.Read(r)
				require.NoError(t, err)
				require.Equal(t, dec0, dec)
			}

			_, err = d.Read(r)
			require.Equal(t, io.EOF, err)
		})
	}
}

func TestDecodeErrors(t *testing.T) {
	for _, ca := range []struct {
		name string
		byts []byte
		err  string
	}{
		{
			"missing payload",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12,
			},
			"payload is too short",
		},
		{
			"STAP-A without NALUs",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12, byte(codech264.NALUTypeStapA),
			},
			"STAP-A packet doesn't contain any NALU",
		},
		{
			"STAP-A without size",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12, byte(codech264.NALUTypeStapA), 0x01,
			},
			"Invalid STAP-A packet",
		},
		{
			"STAP-A with invalid size",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12, byte(codech264.NALUTypeStapA), 0x00, 0x15,
			},
			"Invalid STAP-A packet",
		},
		{
			"FU-A without payload",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12, byte(codech264.NALUTypeFuA),
			},
			"Invalid FU-A packet",
		},
		{
			"FU-A without start bit",
			[]byte{
				0x80, 0xe0, 0x44, 0xed, 0x88, 0x77, 0x6a, 0x15,
				0x9d, 0xbb, 0x78, 0x12, byte(codech264.NALUTypeFuA), 0x00,
			},
			"first NALU does not contain the start bit",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			d := NewDecoder()
			_, err := d.Decode(ca.byts)
			require.NotEqual(t, ErrMorePacketsNeeded, err)
			require.Equal(t, ca.err, err.Error())
		})
	}
}
