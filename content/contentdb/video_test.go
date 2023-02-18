package contentdb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func cleanup() {
	dbHandle = nil
}

func deference[T any](v T) *T {
	return &v
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("error parsing time string as time.RFC3339: " + err.Error())
	}
	return t
}

func compareSlices[T any](a []T, b []T, binaryPredicate func(T, T) bool) bool {
	l := len(a)
	if l != len(b) {
		return false
	}

	for i := 0; i < l; i++ {
		if !binaryPredicate(a[i], b[i]) {
			return false
		}
	}

	return true
}

func TestVideoDB(t *testing.T) {
	defer cleanup()

	f, err := os.CreateTemp("", "contentdb-test-db")
	if err != nil {
		t.Fatal("failed to create temporary file for DB:", err)
	}
	defer os.Remove(f.Name())

	require.Nil(t, Connect(f.Name()))
	require.NotNil(t, dbHandle)

	require.Nil(t, Init())

	v1 := Video{
		FileID: "file1",
		Name:   "file1.mp4",
	}

	id, err := AddVideo(v1)
	require.Equal(t, "file1", id)
	require.Nil(t, err)

	v1_got, err := GetVideo("file1")
	require.Nil(t, err)
	// only check the fields we've set, stuff like ID and CreatedAt gets autopopulated
	require.Equal(t, v1.FileID, v1_got.FileID)
	require.Equal(t, v1.Name, v1_got.Name)

	v, err := GetVideo("file2")
	require.Nil(t, v)
	require.ErrorIs(t, err, ErrVideoNotFound)

	require.Nil(t, DeleteVideo("file1"))
	v, err = GetVideo("file1")
	require.Nil(t, v)
	require.ErrorIs(t, err, ErrVideoNotFound)

	// no problem deleting non-existent video
	require.Nil(t, DeleteVideo("file2"))

	require.Nil(t, Disconnect())
	require.Nil(t, dbHandle)
}

func TestVideoAddDuplicateFileID(t *testing.T) {
	defer cleanup()

	f, err := os.CreateTemp("", "contentdb-test-db")
	if err != nil {
		t.Fatal("failed to create temporary file for DB:", err)
	}
	defer os.Remove(f.Name())

	require.Nil(t, Connect(f.Name()))
	require.NotNil(t, dbHandle)

	require.Nil(t, Init())

	v1 := Video{
		FileID: "file1",
		Name:   "file1.mp4",
	}

	id, err := AddVideo(v1)
	require.Equal(t, "file1", id)
	require.Nil(t, err)

	v2 := Video{
		FileID: "file1",
		Name:   "file2.mp4",
	}

	id2, err2 := AddVideo(v2)
	require.Empty(t, id2)
	require.Error(t, err2)
}

func TestVideoDBSearch(t *testing.T) {
	defer cleanup()

	f, err := os.CreateTemp("", "contentdb-test-db")
	if err != nil {
		t.Fatal("failed to create temporary file for DB:", err)
	}
	defer os.Remove(f.Name())

	require.Nil(t, Connect(f.Name()))
	require.NotNil(t, dbHandle)

	require.Nil(t, Init())

	videos := []Video{
		{
			FileID:          "file1",
			Name:            "movie with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer",
			Uploaded:        mustParseTime("2023-01-31T07:28:00Z"),
			Encoding:        "mp4",
			DurationSeconds: 90 * 60,
			Resolution:      "1920x1080",
			FileSizeBytes:   3 * 1024 * 1024 * 1024,
			Category:        "movie",
		},
		{
			FileID:          "file2",
			Name:            "documentary with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer 2",
			Uploaded:        mustParseTime("2023-01-21T07:28:00Z"),
			Encoding:        "webm",
			DurationSeconds: 60 * 60,
			Resolution:      "1920x1080",
			FileSizeBytes:   1024 * 1024 * 1024,
			Category:        "documentary",
		},
		{
			FileID:          "file3",
			Name:            "short film with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer",
			Uploaded:        mustParseTime("2023-01-11T07:28:00Z"),
			Encoding:        "mp4",
			DurationSeconds: 15 * 60,
			Resolution:      "1280x720",
			FileSizeBytes:   256 * 1024 * 1024,
			Category:        "movie",
		},
	}

	for _, v := range videos {
		_, err := AddVideo(v)
		require.Nil(t, err)
	}

	tests := map[string]struct {
		SearchOptions []SearchOption
		Want          []Video
	}{
		"search by name": {
			SearchOptions: []SearchOption{SearchVideoByName("%documentary%")},
			Want:          []Video{videos[1]},
		},
		"search by license (+)": {
			SearchOptions: []SearchOption{SearchVideoByLicense("CC-BY 3.0")},
			Want:          []Video{videos[0], videos[1], videos[2]},
		},
		"search by license (-)": {
			SearchOptions: []SearchOption{SearchVideoByLicense("Copyright")},
			Want:          nil,
		},
		"search by attribution": {
			SearchOptions: []SearchOption{SearchVideoByAttribution("Famous producer")},
			Want:          []Video{videos[0], videos[2]},
		},
		"search by uploaded before": {
			SearchOptions: []SearchOption{SearchVideoByUploadedBeforeOrAfterDate(mustParseTime("2023-01-25T00:00:00Z"), true)},
			Want:          []Video{videos[1], videos[2]},
		},
		"search by uploaded after": {
			SearchOptions: []SearchOption{SearchVideoByUploadedBeforeOrAfterDate(mustParseTime("2023-01-25T00:00:00Z"), false)},
			Want:          []Video{videos[0]},
		},
		"search by encoding": {
			SearchOptions: []SearchOption{SearchVideoByEncoding("mp4")},
			Want:          []Video{videos[0], videos[2]},
		},
		"search by duration (>=)": {
			SearchOptions: []SearchOption{SearchVideoByDuration(45*time.Minute, false)},
			Want:          []Video{videos[0], videos[1]},
		},
		"search by duration (<=)": {
			SearchOptions: []SearchOption{SearchVideoByDuration(45*time.Minute, true)},
			Want:          []Video{videos[2]},
		},
		"search by resolution": {
			SearchOptions: []SearchOption{SearchVideoByResolution("1280x720")},
			Want:          []Video{videos[2]},
		},
		"search by size (>= 1GB)": {
			SearchOptions: []SearchOption{SearchVideoByFileSize(1024*1024*1024, false)},
			Want:          []Video{videos[0], videos[1]},
		},
		"search by size (<= 1GB)": {
			SearchOptions: []SearchOption{SearchVideoByFileSize(1024*1024*1024, true)},
			Want:          []Video{videos[1], videos[2]},
		},
		"search by category": {
			SearchOptions: []SearchOption{SearchVideoByCategory("documentary")},
			Want:          []Video{videos[1]},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := SearchVideos(tc.SearchOptions...)
			require.Nil(t, err)
			isMatch := compareSlices(tc.Want, got, func(a, b Video) bool {
				return a.FileID == b.FileID
			})

			require.True(t, isMatch, "expected: %v, got: %v", tc.Want, got)
		})
	}
}
