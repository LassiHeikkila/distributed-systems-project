FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y \
        curl \
        ffmpeg  \
        redis-tools \
        jq \
    && apt clean

COPY ./scripts/downscale-video /bin/downscale-video
COPY ./scripts/get-resolution /bin/get-resolution
COPY ./scripts/mp4-to-webm /bin/mp4-to-webm
COPY ./scripts/webm-to-mp4 /bin/webm-to-mp4
COPY ./scripts/download-video /bin/download-video
COPY ./scripts/upload-video /bin/upload-video
COPY ./scripts/get-job-from-redis /bin/get-job-from-redis
COPY ./scripts/get-file-details /bin/get-file-details
COPY ./scripts/run-downscaling-job /bin/run-downscaling-job
COPY ./scripts/run-transcoding-job /bin/run-transcoding-job

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

WORKDIR /work
