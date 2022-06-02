package h264

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDTSExtractor(t *testing.T) {
	sequence := []struct {
		nalus [][]byte
		pts   time.Duration
		dts   time.Duration
	}{
		{
			[][]byte{
				{
					0x67, 0x64, 0x00, 0x28, 0xac, 0xd9, 0x40, 0x78,
					0x02, 0x27, 0xe5, 0xc0, 0x44, 0x00, 0x00, 0x03,
					0x00, 0x04, 0x00, 0x00, 0x03, 0x00, 0x28, 0x3c,
					0x60, 0xc6, 0x58,
				},
				{0x68, 0xeb, 0xe3, 0xcb, 0x22, 0xc0},
				{
					0x65, 0x88, 0x82, 0x00, 0x05, 0xbf, 0xfe, 0xf7,
					0xd3, 0x3f, 0xcc, 0xb2, 0xec, 0x9a, 0x24, 0xb5,
					0xe3, 0xa8, 0xf7, 0xa2, 0x9e, 0x26, 0x5f, 0x43,
					0x75, 0x25, 0x01, 0x9b, 0x96, 0xc4, 0xed, 0x3a,
					0x80, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x55,
					0xda, 0xf7, 0x10, 0xe5, 0xc4, 0x70, 0xe1, 0xfe,
					0x83, 0xc0, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x1f, 0xa0, 0x00, 0x00, 0x05, 0x68, 0x00,
					0x00, 0x03, 0x01, 0xc6, 0x00, 0x00, 0x03, 0x01,
					0x0c, 0x00, 0x00, 0x03, 0x00, 0xb1, 0x00, 0x00,
					0x03, 0x00, 0x8f, 0x80, 0x00, 0x00, 0x8a, 0x80,
					0x00, 0x00, 0x9d, 0x00, 0x00, 0x03, 0x00, 0xb2,
					0x00, 0x00, 0x03, 0x01, 0x1c, 0x00, 0x00, 0x03,
					0x01, 0x7c, 0x00, 0x00, 0x03, 0x02, 0xf0, 0x00,
					0x00, 0x04, 0x40, 0x00, 0x00, 0x08, 0x80, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x0b,
					0x78,
				},
			},
			0,
			0,
		},
		{
			[][]byte{
				{
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
			800 * time.Millisecond,
			200 * time.Millisecond,
		},
		{
			[][]byte{
				{
					0x41, 0x9e, 0x42, 0x78, 0x82, 0x1f, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x02,
					0x0f,
				},
			},
			400 * time.Millisecond,
			400 * time.Millisecond,
		},
		{
			[][]byte{
				{
					0x01, 0x9e, 0x61, 0x74, 0x43, 0xff, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00,
					0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x04, 0x9c,
				},
			},
			200 * time.Millisecond,
			600 * time.Millisecond,
		},
	}

	ex := NewDTSExtractor()
	sps := &SPS{}

	for _, sample := range sequence {
		idrPresent := IDRPresent(sample.nalus)

		for _, nalu := range sample.nalus {
			if NALUType(nalu[0]&0x1F) == NALUTypeSPS {
				err := sps.Unmarshal(nalu)
				require.NoError(t, err)
				break
			}
		}

		dts, err := ex.Extract(sample.nalus, idrPresent, sample.pts, sps)
		require.NoError(t, err)
		require.Equal(t, sample.dts, dts)
	}
}
