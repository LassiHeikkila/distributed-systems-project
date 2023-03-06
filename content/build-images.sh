#!/bin/bash

version=${1:-"dev"}

docker build -t "flmnchll/content-provider:${version}" --file Dockerfile.contentprovider .
docker build -t "flmnchll/content-transcoder:${version}" contenttranscoder
docker build -t "flmnchll/content-manager:${version}" --file Dockerfile.contentmanager .
