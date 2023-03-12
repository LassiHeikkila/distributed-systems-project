package main

import (
	"crypto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

const (
	kB = 1024
	MB = 1024 * kB
)

func VideoUploadHandler(w http.ResponseWriter, req *http.Request) {
	e := json.NewEncoder(w)

	err := req.ParseMultipartForm(64 * MB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "failed to parse multipart form",
		})
		return
	}

	// form should contain
	// - title
	// - license
	// - attribution
	// - category
	// - original content id
	// - file with key video_upload
	title := req.FormValue("title")
	license := req.FormValue("license")
	attribution := req.FormValue("attribution")
	category := req.FormValue("category")
	origContentID := req.FormValue("originalContentID")
	if origContentID != "" {
		VideoUploadFileAdditionHandler(w, req, origContentID)
		return
	}

	f, fh, err := req.FormFile("video_upload")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "could not get file with key \"video_upload\" from form",
		})
		return
	}
	wf, err := os.CreateTemp(videoFileTempDirectory, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "failed to create temporary file to hold video",
		})
		return
	}
	_, fileSha256Sum, err := CopyAndHash(wf, f, crypto.SHA256)
	if err != nil {
		os.Remove(wf.Name())

		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "failed to store video onto disk",
		})
		return
		// big problem
	}

	contentID := contentdb.GenerateUUID()
	enc := strings.TrimPrefix(path.Ext(fh.Filename), ".")
	fileID := contentdb.GenerateUUID() + "." + enc
	// this could also fail but unlikely if temp directory is
	// on same storage device as final directory
	_ = os.Rename(wf.Name(), path.Join(videoFileDirectory, fileID))

	resolution, err := GetResolution(path.Join(videoFileDirectory, fileID))
	if err != nil {
		log.Println("failed to get resolution:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	duration, err := GetDuration(path.Join(videoFileDirectory, fileID))
	if err != nil {
		log.Println("failed to get duration:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := contentdb.Video{
		ContentID:       contentID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Name:            title,
		License:         license,
		Attribution:     attribution,
		DurationSeconds: RoundDurationToSecondsCeil(duration),
		Category:        category,
		Files: []contentdb.VideoFile{
			{
				FileID:        fileID,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
				ContentID:     contentID,
				Uploaded:      time.Now().UTC(),
				Encoding:      enc,
				Resolution:    resolution,
				FileSizeBytes: uint(fh.Size),
				Hash:          fmt.Sprintf("%x", fileSha256Sum),
			},
		},
	}

	videoID, err := contentdb.AddVideo(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = e.Encode(map[string]any{
			"ok":    false,
			"msg":   "inserting video to db failed",
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
	_ = e.Encode(map[string]any{
		"ok":        true,
		"msg":       "video added to database successfully",
		"contentID": videoID,
	})

	go ProcessNewUpload(v)
}

func VideoUploadFileAdditionHandler(w http.ResponseWriter, req *http.Request, origContentID string) {
	e := json.NewEncoder(w)

	f, fh, err := req.FormFile("video_upload")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "could not get file with key \"video_upload\" from form",
		})
		return
	}

	wf, err := os.CreateTemp(videoFileTempDirectory, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "failed to create temporary file to hold video",
		})
		return
	}
	_, fileSha256Sum, err := CopyAndHash(wf, f, crypto.SHA256)
	if err != nil {
		os.Remove(wf.Name())

		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(map[string]any{
			"ok":  false,
			"msg": "failed to store video onto disk",
		})
		return
		// big problem
	}
	enc := strings.TrimPrefix(path.Ext(fh.Filename), ".")
	fileID := contentdb.GenerateUUID() + "." + enc
	// this could also fail but unlikely if temp directory is
	// on same storage device as final directory
	_ = os.Rename(wf.Name(), path.Join(videoFileDirectory, fileID))

	resolution, err := GetResolution(path.Join(videoFileDirectory, fileID))
	if err != nil {
		log.Println("failed to get resolution:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vf := contentdb.VideoFile{
		FileID:        fileID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		ContentID:     origContentID,
		Uploaded:      time.Now().UTC(),
		Encoding:      enc,
		Resolution:    resolution,
		FileSizeBytes: uint(fh.Size),
		Hash:          fmt.Sprintf("%x", fileSha256Sum),
	}

	err = contentdb.AddVideoFile(vf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = e.Encode(map[string]any{
			"ok":    false,
			"msg":   "adding video file to db failed",
			"error": err.Error(),
		})
		return
	}

	_ = e.Encode(map[string]any{
		"ok":        true,
		"msg":       "video file added to database successfully",
		"contentID": origContentID,
		"fileID":    fileID,
	})

	// if this was the first addition with current encoding, trigger downscale jobs
	SubmitDownscalingJobsIfNeeded(vf.ContentID, vf.Encoding)
}

func SubmitDownscalingJobsIfNeeded(contentID string, enc string) {
	v, err := contentdb.GetVideo(contentID)
	if err != nil {
		return
	}

	count := 0
	idx := 0
	for i, vf := range v.Files {
		if vf.Encoding == enc {
			count++
			idx = i
		}
	}

	if count == 1 {
		SubmitDownscalingJobs(v.Files[idx])
	}
}
