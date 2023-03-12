package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetResolution(t *testing.T) {
	tests := map[string]struct {
		inputFile          string
		expectedResolution string
	}{
		"SUV-Iceland.webm": {
			inputFile:          "./video-samples/SUV-Iceland.webm",
			expectedResolution: "1280x720",
		},
		"file_example_MP4_1920_18MG.mp4": {
			inputFile:          "./video-samples/file_example_MP4_1920_18MG.mp4",
			expectedResolution: "1920x1080",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := GetResolution(tc.inputFile)
			require.Nil(t, err)
			require.Equal(t, tc.expectedResolution, res)
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := map[string]struct {
		inputFile               string
		expectedDuration        time.Duration
		expectedDurationSeconds int
	}{
		"SUV-Iceland.webm": {
			inputFile:               "./video-samples/SUV-Iceland.webm",
			expectedDuration:        (13 * time.Second) + (81 * time.Millisecond),
			expectedDurationSeconds: 14,
		},
		"file_example_MP4_1920_18MG.mp4": {
			inputFile:               "./video-samples/file_example_MP4_1920_18MG.mp4",
			expectedDuration:        (30 * time.Second) + (526 * time.Millisecond) + (667 * time.Microsecond),
			expectedDurationSeconds: 31,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := GetDuration(tc.inputFile)
			require.Nil(t, err)
			require.Equal(t, tc.expectedDuration, d)
			require.Equal(t, tc.expectedDurationSeconds, RoundDurationToSecondsCeil(d))
		})
	}
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("failed to parse duration: " + err.Error())
	}

	return d
}

func TestRoundDurationToSecondsCeil(t *testing.T) {
	tests := map[time.Duration]int{
		mustParseDuration("12.25s"): 13,
		mustParseDuration("0.012s"): 1,
		mustParseDuration("1h"):     3600,
	}

	for d, want := range tests {
		t.Run(d.String(), func(t *testing.T) {
			got := RoundDurationToSecondsCeil(d)
			require.Equal(t, want, got)
		})
	}
}
