# content-transcoder
This is a shell script based "service" that pulls transcoding & downscaling jobs from a work queue, and executes them.

Docker image will contain logic to do both kinds of jobs, job queue will be configured with environmental variable.

## Environmental variables
### Required
`REDIS_USERNAME`: username for redis
`REDIS_PASSWORD`: password for redis
`REDIS_HOST`: hostname for redis
`REDIS_PORT`: port for redis
`REDIS_JOB_QUEUE`: redis list name where to pop jobs from

`VIDEO_DOWNLOAD_PREFIX`: URL for video downloads, file ID will be appended to it.
`VIDEO_DETAILS_PREFIX`: URL for getting video details, file ID will be appended to it.

### Optional
???

## Jobs
Supported job types are `transcode` and `downscale`.

Examples:

[`transcode`](./json/examples/transcoding-job-example.json)

[`downscale`](./json/examples/downscaling-job-example.json)

### Transcoding job
Pull original file
Pull original metadata

Transcode from `mp4` => `webm` or vice versa

Copy metadata, modify name by appending new format to the end.
If original file was called `video.mp4`, and it was transcoded to `webm`, then it should be called `video.mp4.webm`.
Update encoding field in metadata.

Push new file with modified metadata. `fileID` field should be left blank.
`originalFileID` field should contain `fileID` of the original copy.

Example file details from original file: [](json/examples/video-details-example.json)

Example file details from transcoded file: [](json/examples/video-details-example-transcoded.json)
