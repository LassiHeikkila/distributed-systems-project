package contentdb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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
		ContentID: "file1",
		Name:      "file1.mp4",
	}

	id, err := AddVideo(v1)
	require.Equal(t, "file1", id)
	require.Nil(t, err)

	v1_got, err := GetVideo("file1")
	require.Nil(t, err)
	// only check the fields we've set, stuff like ID and CreatedAt gets autopopulated
	require.Equal(t, v1.ContentID, v1_got.ContentID)
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
		ContentID: "file1",
		Name:      "file1.mp4",
	}

	id, err := AddVideo(v1)
	require.Equal(t, "file1", id)
	require.Nil(t, err)

	v2 := Video{
		ContentID: "file1",
		Name:      "file2.mp4",
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
			ContentID:       "content1",
			CreatedAt:       mustParseTime("2023-01-31T07:28:00Z"),
			Name:            "movie with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer",
			DurationSeconds: 90 * 60,
			Category:        "movie",
			Files: []VideoFile{
				{
					FileID:        "file1",
					CreatedAt:     time.Time{},
					UpdatedAt:     time.Time{},
					DeletedAt:     gorm.DeletedAt{},
					ContentID:     "content1",
					Uploaded:      mustParseTime("2023-01-31T07:28:00Z"),
					Encoding:      "mp4",
					Resolution:    "1920x1080",
					FileSizeBytes: 3 * 1024 * 1024 * 1024,
				},
			},
		},
		{
			ContentID:       "content2",
			CreatedAt:       mustParseTime("2023-01-21T07:28:00Z"),
			Name:            "documentary with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer 2",
			DurationSeconds: 60 * 60,
			Category:        "documentary",
			Files: []VideoFile{
				{
					FileID:        "file2",
					CreatedAt:     time.Time{},
					UpdatedAt:     time.Time{},
					DeletedAt:     gorm.DeletedAt{},
					ContentID:     "content2",
					Uploaded:      mustParseTime("2023-01-21T07:28:00Z"),
					Encoding:      "webm",
					Resolution:    "1920x1080",
					FileSizeBytes: 1024 * 1024 * 1024,
				},
			},
		},
		{
			ContentID:       "content3",
			CreatedAt:       mustParseTime("2023-01-11T07:28:00Z"),
			Name:            "short film with title",
			License:         "CC-BY 3.0",
			Attribution:     "Famous producer",
			DurationSeconds: 15 * 60,
			Category:        "movie",
			Files: []VideoFile{
				{
					FileID:        "file3",
					CreatedAt:     time.Time{},
					UpdatedAt:     time.Time{},
					DeletedAt:     gorm.DeletedAt{},
					ContentID:     "content3",
					Uploaded:      mustParseTime("2023-01-11T07:28:00Z"),
					Encoding:      "mp4",
					Resolution:    "1280x720",
					FileSizeBytes: 256 * 1024 * 1024,
				},
			},
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
		"search by duration (>=)": {
			SearchOptions: []SearchOption{SearchVideoByDuration(45*time.Minute, false)},
			Want:          []Video{videos[0], videos[1]},
		},
		"search by duration (<=)": {
			SearchOptions: []SearchOption{SearchVideoByDuration(45*time.Minute, true)},
			Want:          []Video{videos[2]},
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
				return a.ContentID == b.ContentID
			})

			require.True(t, isMatch, "expected: %v, got: %v", tc.Want, got)
		})
	}
}

func TestVideoFileAdd(t *testing.T) {
	defer cleanup()

	f, err := os.CreateTemp("", "contentdb-test-db")
	if err != nil {
		t.Fatal("failed to create temporary file for DB:", err)
	}
	defer os.Remove(f.Name())

	require.Nil(t, Connect(f.Name()))
	require.NotNil(t, dbHandle)

	require.Nil(t, Init())

	file1 := VideoFile{
		FileID:        "video1.mp4",
		ContentID:     "content1",
		Encoding:      "mp4",
		Resolution:    "1920x1080",
		FileSizeBytes: 256 * 1024 * 1024,
		Hash:          "sdfsdfsdf",
	}
	file2 := VideoFile{
		FileID:        "video1.webm",
		ContentID:     "content1",
		Encoding:      "webm",
		Resolution:    "1920x1080",
		FileSizeBytes: 128 * 1024 * 1024,
		Hash:          "dwdwdsf",
	}

	v1 := Video{
		ContentID:       "content1",
		Files:           []VideoFile{file1},
		Name:            "some film",
		License:         "example license",
		Attribution:     "Carla the Creator",
		DurationSeconds: 3600,
		Category:        "movie",
	}

	id, err := AddVideo(v1)
	require.Equal(t, "content1", id)
	require.Nil(t, err)

	err = AddVideoFile(file2)
	require.Nil(t, err)

	v, _ := GetVideo("content1")
	require.Len(t, v.Files, 2)
}
