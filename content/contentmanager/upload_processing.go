package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

const (
	resolution1080p = "1920x1080"
	resolution720p  = "1280x720"
	resolution480p  = "640x480"

	encodingMp4  = "mp4"
	encodingWebm = "webm"

	JobTypeDownscale = "downscale"
	JobTypeTranscode = "transcode"
)

var supportedResolutions = []string{
	resolution1080p,
	resolution720p,
	resolution480p,
}

var DownscalingNeeded = map[string][]string{
	resolution1080p: supportedResolutions[1:],
	resolution720p:  supportedResolutions[2:],
	resolution480p:  nil,
}

var TranscodingNeeded = map[string][]string{
	encodingMp4:  {encodingWebm},
	encodingWebm: {encodingMp4},
}

type GenericJob struct {
	JobID   string `json:"jobID"`
	JobType string `json:"jobType"`
}

type TranscodingJob struct {
	GenericJob
	Job TranscodingDetails `json:"job"`
}

type DownscalingJob struct {
	GenericJob
	Job DownscalingDetails `json:"job"`
}

type TranscodingDetails struct {
	SourceEncoding string `json:"sourceEncoding"`
	TargetEncoding string `json:"targetEncoding"`
	FileID         string `json:"fileID"`
}

type DownscalingDetails struct {
	SourceResolution string `json:"sourceResolution"`
	TargetResolution string `json:"targetResolution"`
	FileID           string `json:"fileID"`
}

func ProcessNewUpload(v contentdb.Video) error {
	log.Println("processing uploaded video with id", v.ContentID)

	SubmitDownscalingJobs(v.Files[0])
	SubmitTranscodingJobs(v.Files[0])

	return nil
}

func SubmitDownscalingJobs(vf contentdb.VideoFile) {
	downscalingNeeded := DownscalingNeeded[vf.Resolution]
	for _, ds := range downscalingNeeded {
		job := DownscalingJob{
			GenericJob: GenericJob{
				JobID:   contentdb.GenerateUUID(),
				JobType: JobTypeDownscale,
			},
			Job: DownscalingDetails{
				SourceResolution: vf.Resolution,
				TargetResolution: ds,
				FileID:           vf.FileID,
			},
		}
		log.Println("submitting downscaling job with id:", job.JobID)
		go SubmitToRedisAsJSON(context.Background(), redisQueue, &job)
	}
}
func SubmitTranscodingJobs(vf contentdb.VideoFile) {
	transcodingNeeded := TranscodingNeeded[vf.Encoding]
	for _, tc := range transcodingNeeded {
		job := TranscodingJob{
			GenericJob: GenericJob{
				JobID:   contentdb.GenerateUUID(),
				JobType: JobTypeTranscode,
			},
			Job: TranscodingDetails{
				SourceEncoding: vf.Encoding,
				TargetEncoding: tc,
				FileID:         vf.FileID,
			},
		}
		log.Println("submitting transcoding job with id:", job.JobID)
		go SubmitToRedisAsJSON(context.Background(), redisQueue, &job)
	}
}

func GetResolution(file string) (string, error) {
	res, err := ffmpeg.Probe(
		file,
		ffmpeg.KwArgs{
			"select_streams": "v:0",
			"show_entries":   "stream=width,height",
		},
	)

	if err != nil {
		return "", err
	}

	type result struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	var r result
	_ = json.Unmarshal([]byte(res), &r)

	return fmt.Sprintf("%dx%d", r.Streams[0].Width, r.Streams[0].Height), nil
}

func GetDuration(file string) (time.Duration, error) {
	res, err := ffmpeg.Probe(file)

	if err != nil {
		return 0, err
	}

	type result struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	var r result
	_ = json.Unmarshal([]byte(res), &r)

	if r.Format.Duration != "" {
		return time.ParseDuration(r.Format.Duration + "s")
	}

	return 0, errors.New("unable to determine duration")
}

func RoundDurationToSecondsCeil(d time.Duration) int {
	ms := d.Milliseconds()
	r := ms % 1000
	if r == 0 {
		return int(ms / 1000)
	}

	return int((ms + 1000 - r) / 1000)
}
