package h264

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var casesAnnexB = []struct {
	name   string
	encin  []byte
	encout []byte
	dec    [][]byte
}{
	{
		"2 zeros, single",
		[]byte{0x00, 0x00, 0x01, 0xaa, 0xbb},
		[]byte{0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb},
		[][]byte{
			{0xaa, 0xbb},
		},
	},
	{
		"2 zeros, multiple",
		[]byte{
			0x00, 0x00, 0x01, 0xaa, 0xbb, 0x00, 0x00, 0x01,
			0xcc, 0xdd, 0x00, 0x00, 0x01, 0xee, 0xff,
		},
		[]byte{
			0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb, 0x00, 0x00,
			0x00, 0x01, 0xcc, 0xdd, 0x00, 0x00, 0x00, 0x01,
			0xee, 0xff,
		},
		[][]byte{
			{0xaa, 0xbb},
			{0xcc, 0xdd},
			{0xee, 0xff},
		},
	},
	{
		"3 zeros, single",
		[]byte{0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb},
		[]byte{0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb},
		[][]byte{
			{0xaa, 0xbb},
		},
	},
	{
		"3 zeros, multiple",
		[]byte{
			0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb, 0x00, 0x00,
			0x00, 0x01, 0xcc, 0xdd, 0x00, 0x00, 0x00, 0x01,
			0xee, 0xff,
		},
		[]byte{
			0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb, 0x00, 0x00,
			0x00, 0x01, 0xcc, 0xdd, 0x00, 0x00, 0x00, 0x01,
			0xee, 0xff,
		},
		[][]byte{
			{0xaa, 0xbb},
			{0xcc, 0xdd},
			{0xee, 0xff},
		},
	},
}

func TestAnnexBUnmarshal(t *testing.T) {
	for _, ca := range casesAnnexB {
		t.Run(ca.name, func(t *testing.T) {
			dec, err := AnnexBUnmarshal(ca.encin)
			require.NoError(t, err)
			require.Equal(t, ca.dec, dec)
		})
	}
}

func TestAnnexBMarshal(t *testing.T) {
	for _, ca := range casesAnnexB {
		t.Run(ca.name, func(t *testing.T) {
			enc, err := AnnexBMarshal(ca.dec)
			require.NoError(t, err)
			require.Equal(t, ca.encout, enc)
		})
	}
}

func TestAnnexBUnmarshalError(t *testing.T) {
	for _, ca := range []struct {
		name string
		enc  []byte
		err  string
	}{
		{
			"empty",
			[]byte{},
			"initial delimiter not found",
		},
		{
			"invalid initial delimiter 1",
			[]byte{0xaa, 0xbb},
			"unexpected byte: 170",
		},
		{
			"invalid initial delimiter 2",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x01},
			"initial delimiter not found",
		},
		{
			"empty NALU 1",
			[]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x01},
			"empty NALU",
		},
		{
			"empty NALU 2",
			[]byte{0x00, 0x00, 0x01, 0xaa, 0x00, 0x00, 0x01},
			"empty NALU",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			_, err := AnnexBUnmarshal(ca.enc)
			require.EqualError(t, err, ca.err)
		})
	}
}

func BenchmarkAnnexBUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AnnexBUnmarshal([]byte{
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
			0x00, 0x00, 0x00, 0x01,
			0x01, 0x02, 0x03, 0x04,
		})
	}
}
