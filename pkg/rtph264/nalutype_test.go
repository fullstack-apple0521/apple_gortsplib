package rtph264

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNALUType(t *testing.T) {
	require.NotEqual(t, true, strings.HasPrefix(naluType(10).String(), "unknown"))
	require.NotEqual(t, true, strings.HasPrefix(naluType(26).String(), "unknown"))
	require.Equal(t, true, strings.HasPrefix(naluType(50).String(), "unknown"))
}
