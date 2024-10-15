package log

import (
	api "github.com/MSkrzypietz/proglog/api/v1"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestSegment(t *testing.T) {
	dir, _ := os.MkdirTemp("", "segment_test")
	defer os.RemoveAll(dir)

	want := &api.Record{Value: []byte("hello world")}

	cfg := Config{}
	cfg.Segment.MaxStoreBytes = 1024
	cfg.Segment.MaxIndexBytes = entWidth * 3

	s, err := newSegment(dir, 16, cfg)
	require.NoError(t, err)
	require.Equal(t, uint64(16), s.nextOffset)
	require.False(t, s.isMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, 16+i, off)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = s.Append(want)
	require.Equal(t, io.EOF, err)
	require.True(t, s.isMaxed())

	cfg.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	cfg.Segment.MaxIndexBytes = 1024

	s, err = newSegment(dir, 16, cfg)
	require.NoError(t, err)
	require.True(t, s.isMaxed())

	err = s.Remove()
	require.NoError(t, err)
	s, err = newSegment(dir, 16, cfg)
	require.NoError(t, err)
	require.False(t, s.isMaxed())
}
